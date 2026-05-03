package services

import (
	"encoding/json"
	"net/http"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserIDKey string

const UserID UserIDKey = "user_id"

type Loginer interface {
	Login(w http.ResponseWriter, r *http.Request)
}

func (u UserService) Login(w http.ResponseWriter, r *http.Request) {
	//Получить пользователя из запроса
	var (
		user model.User
		err  error
	)
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
	if _, err = u.Repo.GetByName(user.UserName); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	//Создвать куки для пользователя
	userCookie, err := u.createCookie()
	if err != nil {
		http.Error(w, "Cannot create cookie", http.StatusInternalServerError)
		return
	}
	//Установить куки в ответе
	http.SetCookie(w, userCookie)
	w.WriteHeader(http.StatusOK)
}

func (u UserService) validateUser(user model.User) bool {
	//Проверка логина и пароля
	if user.UserName == "" || user.Password == "" {
		return false
	}
	//Проверка сущетвования пользователя в базе данных
	profile, err := u.Repo.GetByName(user.UserName)
	if err != nil {
		return false
	}
	//Проверка пароля с помощью bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(profile.Password), []byte(user.Password))
	return err == nil
}
