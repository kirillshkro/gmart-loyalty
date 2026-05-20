package services

import (
	"encoding/json"
	"net/http"
)

// Реализация метода для получения заказов пользователя
func (o Service) OrdersByUser(w http.ResponseWriter, r *http.Request) {
	//Проверка авторицации пользователя
	if !o.cookieExist(r, authCookie) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Получение контекста запроса
	ctx := r.Context()
	orders, err := o.Repo.GetAll(ctx)
	if err != nil {
		http.Error(w, "Ошибка при получении заказов", http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		http.Error(w, "Orders not found", http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, "Ошибка при кодировании JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
