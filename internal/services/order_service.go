package services

import (
	"context"
	"errors"
	"io"
	"net/http"
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

func (o *Service) SetOrderUser(w http.ResponseWriter, r *http.Request) {
	//Получить данные пользователя из куки
	if !o.cookieExist(r, authCookie) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

func (o *Service) processingOrder(userID int, numOrder string, errCh chan<- error) error {
	//Проверить статус заказа
	order := model.Order{
		Number: string(numOrder),
		UserID: userID,
		Status: model.StatusNew,
	}
	//Проверить формат номера заказ
	if !utils.Valid(string(numOrder)) {
		errCh <- &types.ErrInvalidFormatOrder{Number: string(numOrder)}
		order.Status = model.StutusInvalid
		return nil
	}

	ctx := context.WithValue(context.TODO(), types.UserID, userID)
	ctx = context.WithValue(ctx, types.OrderNum, string(numOrder))
	//Сохранить заказ в базе данных
	//Если пользователь дублировал заказ, вернуть
	if err := o.Repo.CreateOrder(ctx, &order); err != nil {
		errCh <- err
		return nil
	}
	return nil
}
