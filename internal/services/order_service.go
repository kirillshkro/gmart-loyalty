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

func (o *Service) SetOrderUser(w http.ResponseWriter, r *http.Request) {
	var order model.Order
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
	if !utils.Valid(string(numOrder)) {
		http.Error(w, "Invalid format order", http.StatusUnprocessableEntity)
		return
	}
	order = model.Order{
		OrderNum:  string(numOrder),
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	ctx := context.WithValue(context.TODO(), types.UserID, userID)
	ctx = context.WithValue(ctx, types.OrderNum, string(numOrder))
	//Сохранить заказ в базе данных
	//Если пользователь дублировал заказ, вернуть
	if err = o.Repo.CreateOrder(ctx, &order); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			http.Error(w, err.Error(), http.StatusOK)
			return
		}
		//Проверить, что нет других пользователей с таким номером заказа
		if _, ok := errors.AsType[*types.ErrOwnAnotherUser](err); ok {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
