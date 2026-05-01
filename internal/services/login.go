package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
)

type Loginer interface {
	Login(w http.ResponseWriter, r *http.Request)
}

func (u UserService) Login(w http.ResponseWriter, r *http.Request) {
	//Получить пользователя из запроса
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//Проверка логина и пароля
	if ok := u.validateUser(user); !ok {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}
	//Получить пользователя из базы данных
	if _, err := u.Repo.GetByName(user.UserName); err != nil {
		http.Error(w, fmt.Errorf("User %s not found", user.UserName).Error(), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (u UserService) validateUser(user model.User) bool {
	//Проверка логина и пароля
	if user.UserName == "" || user.Password == "" {
		return false
	}
	return true
}
