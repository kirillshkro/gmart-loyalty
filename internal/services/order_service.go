package services

import (
	"errors"
	"io"
	"net/http"

	"github.com/kirillshkro/gmart-loyalty/internal/config"
	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"github.com/kirillshkro/gmart-loyalty/internal/repository"
	"github.com/kirillshkro/gmart-loyalty/internal/types"
	"github.com/kirillshkro/gmart-loyalty/pkg/utils"
	"gorm.io/gorm"
)

type OrderService struct {
	Repo repository.IOrderRepository
}

type IOrderService interface {
	SetOrderUser(w http.ResponseWriter, r *http.Request)
	OrdersByUser(w http.ResponseWriter, r *http.Request)
}

func NewOrderService(cfg *config.AppConfig) *OrderService {
	return &OrderService{}
}

func (o *OrderService) SetOrderUser(w http.ResponseWriter, r *http.Request) {
	//Извлечь модель из запроса
	var order model.Order
	//Получить данные пользователя из контекста
	userID := r.Context().Value(UserID).(int)
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
	order.OrderNum = string(numOrder)
	order.UserID = userID
	//Сохранить заказ в базе данных
	//Если пользователь дублировал заказ, вернуть
	if err = o.Repo.Create(order); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			http.Error(w, err.Error(), http.StatusOK)
			return
		}
		//Проверить, что нет других пользователей с таким номером заказа
		if _, ok := errors.AsType[*types.ErrOwnAnotherUser](err); ok {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
