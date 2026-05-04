package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"github.com/kirillshkro/gmart-loyalty/internal/types"
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
		same bool
	)
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//Проверка логина и пароля
	if ok, err := u.validateUser(user); !ok {
		if _, same = errors.AsType[*types.ErrInvalidLogin](err); same {
			http.Error(w, "Bad request: empty login or password", http.StatusBadRequest)
			return
		}
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
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

func (u UserService) validateUser(user model.User) (bool, error) {
	//Проверка логина и пароля
	if user.UserName == "" || user.Password == "" {
		return false, &types.ErrInvalidLogin{}
	}
	//Проверка сущетвования пользователя в базе данных
	profile, err := u.Repo.GetByName(user.UserName)
	if err != nil {
		return false, err
	}
	//Проверка пароля с помощью bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(profile.Password), []byte(user.Password))
	if err != nil {
		return false, err
	}
	return true, nil
}
