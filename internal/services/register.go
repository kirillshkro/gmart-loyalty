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

	if regUser.UserName == "" || regUser.Password == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if regUser.Password != regUser.Password2 {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	//Шифруем пароль пользователя
	cryptPass, err := bcrypt.GenerateFromPassword([]byte(regUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error encrypting password", http.StatusInternalServerError)
		return
	}
	regUser.Password = string(cryptPass)
	regUser.Password2 = regUser.Password
	//Создаем профиль пользователя

	profile = model.UserProfile{
		User: regUser,
	}

	if err = s.Repo.CreateUser(profile); err != nil {
		if _, ok = errors.AsType[*types.ErrDuplicateUser](err); ok {
			w.WriteHeader(http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
