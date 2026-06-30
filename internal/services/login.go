package services

import (
	"encoding/json"
	"net/http"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"github.com/kirillshkro/gmart-loyalty/internal/model/claims"
	"golang.org/x/crypto/bcrypt"
)

type Loginer interface {
	Login(w http.ResponseWriter, r *http.Request)
}

func (u Service) Login(w http.ResponseWriter, r *http.Request) {
	//Получить пользователя из запроса
	var (
		user    model.User
		err     error
		profile model.UserProfile
	)

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	//Получить пользователя из базы данных
	if profile, err = u.Repo.UserByName(user.Login); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(profile.Password), []byte(user.Password)) != nil {
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
