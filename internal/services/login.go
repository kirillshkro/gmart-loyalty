package services

import (
	"net/http"
)

type UserIDKey string

const UserID UserIDKey = "user_id"

type Loginer interface {
	Login(w http.ResponseWriter, r *http.Request)
}

func (u Service) Login(w http.ResponseWriter, r *http.Request) {
	//Получить пользователя из запроса
	var (
		err error
	)

	ctx := r.Context()
	userID := ctx.Value(UserID).(int)

	//Получить пользователя из базы данных
	if _, err = u.Repo.UserByID(userID); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}
