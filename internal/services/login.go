package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"github.com/kirillshkro/gmart-loyalty/internal/model/claims"
	"github.com/kirillshkro/gmart-loyalty/internal/types"
)

type Loginer interface {
	Login(w http.ResponseWriter, r *http.Request)
}

func (u Service) Login(w http.ResponseWriter, r *http.Request) {
	//Получить пользователя из запроса
	var (
		err error
	)

	userCookie, err := r.Cookie(authCookie)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//Проверить валидность логин/пароля
	if ok, err := u.validateUser(user); !ok {
		if _, ok := errors.AsType[*types.ErrInvalidLogin](err); ok {
			http.Error(w, "Empty login or password", http.StatusBadRequest)
			return
		}
	}

	userID, err := u.userFromCookie(userCookie)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	//Получить пользователя из базы данных
	if _, err = u.Repo.UserByID(userID); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (u Service) userFromCookie(c *http.Cookie) (int, error) {
	tk := c.Value
	tkClaims := claims.NewUserClaims()
	userID, err := tkClaims.UserIDByToken(tk)
	return userID, err
}
