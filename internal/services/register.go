package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"github.com/kirillshkro/gmart-loyalty/internal/types"
	"golang.org/x/crypto/bcrypt"
)

type Registerer interface {
	Register(w http.ResponseWriter, r *http.Request)
}

// Регистрирует нового пользователя в системе
func (s Service) Register(w http.ResponseWriter, r *http.Request) {
	var (
		regUser model.User
		profile model.UserProfile
		ok      bool
		err     error
	)
	if err = json.NewDecoder(r.Body).Decode(&regUser); err != nil {
		http.Error(w, "Invalid request", http.StatusInternalServerError)
		return
	}

	if regUser.Login == "" || regUser.Password == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//Шифруем пароль пользователя
	cryptPass, err := bcrypt.GenerateFromPassword([]byte(regUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error encrypting password", http.StatusInternalServerError)
		return
	}
	regUser.Password = string(cryptPass)

	//Создаем профиль пользователя

	profile = model.UserProfile{
		User: regUser,
	}
	var userID int
	if userID, err = s.Repo.CreateUser(profile); err != nil {
		if _, ok = errors.AsType[*types.ErrDuplicateUser](err); ok {
			w.WriteHeader(http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user_cookie, err := s.createCookie(userID)
	if err != nil {
		http.Error(w, "Error creating cookie", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, user_cookie)
	w.Header().Set("Content-Type", "application/json")
	r.AddCookie(user_cookie)
	s.Login(w, r)
}
