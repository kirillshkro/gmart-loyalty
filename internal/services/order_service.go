package services

import (
	"context"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"github.com/kirillshkro/gmart-loyalty/internal/types"
	"github.com/kirillshkro/gmart-loyalty/pkg/utils"
	"gorm.io/gorm"
)

type IOrderService interface {
	SetOrderUser(w http.ResponseWriter, r *http.Request)
	OrdersByUser(w http.ResponseWriter, r *http.Request)
}

const MAX_ORDERS_PER_USER = 100
const WORKER_POOL_SIZE = 10

func (o *Service) SetOrderUser(w http.ResponseWriter, r *http.Request) {
	//Получить данные пользователя из куки
	if !o.cookieExist(r, authCookie) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cookie, err := r.Cookie(authCookie)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID, err := o.userFromCookie(cookie)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//Получить номер заказа из тела запроса
	numOrder, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	task := func(w http.ResponseWriter) {
		errCh := make(chan error, MAX_ORDERS_PER_USER)
		go func() {
			err = o.processingOrder(userID, string(numOrder), errCh)
			close(errCh)
		}()
		select {
		case err := <-errCh:
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				w.WriteHeader(http.StatusOK)
				return
			}
			if _, ok := errors.AsType[*types.ErrOwnAnotherUser](err); ok {
				w.WriteHeader(http.StatusConflict)
				return
			}
			if _, ok := errors.AsType[*types.ErrInvalidFormatOrder](err); ok {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
		case <-time.After(5 * time.Second):
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusAccepted)
	}
	o.processTaskInPool(task, w)
}

func (o *Service) processingOrder(userID int, numOrder string, errCh chan<- error) error {
	order := model.Order{
		Number: string(numOrder),
		UserID: userID,
		Status: model.StatusNew,
	}
	//Проверить формат номера заказ
	if !utils.Valid(string(numOrder)) {
		errCh <- &types.ErrInvalidFormatOrder{Number: string(numOrder)}
		order.Status = model.StatusInvalid
		return nil
	}

	ctx := context.WithValue(context.TODO(), types.UserID, userID)
	ctx = context.WithValue(ctx, types.OrderNum, string(numOrder))
	//Сохранить заказ в базе данных
	if err := o.Repo.CreateOrder(ctx, &order); err != nil {
		errCh <- err
		return nil
	}
	return nil
}

func (o *Service) processTaskInPool(task func(w http.ResponseWriter), w http.ResponseWriter) {
	var wg sync.WaitGroup
	taskCh := make(chan func(w http.ResponseWriter))

	for range WORKER_POOL_SIZE {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				task(w)
			}
		}()
	}
	taskCh <- task
	close(taskCh)
	wg.Wait()
}
