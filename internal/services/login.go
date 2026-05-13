package services

import (
	"net/http"

	"github.com/kirillshkro/gmart-loyalty/internal/model/claims"
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

	userCookie, err := r.Cookie(authCookie)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
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
