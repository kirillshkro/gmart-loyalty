package services

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kirillshkro/gmart-loyalty/internal/types"
)

// Реализация метода для получения заказов пользователя
func (o Service) OrdersByUser(w http.ResponseWriter, r *http.Request) {
	//Проверка авторицации пользователя
	if !o.cookieExist(r, authCookie) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	uc, err := r.Cookie(authCookie)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Получение ID пользователя из куки
	userID, err := o.userFromCookie(uc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Получение контекста запроса
	ctx := context.WithValue(r.Context(), types.UserID, userID)
	orders, err := o.Repo.GetAll(ctx)
	if err != nil {
		o.logger.Error("Ошибка при получении заказов")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		o.logger.Error("Orders not found")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
