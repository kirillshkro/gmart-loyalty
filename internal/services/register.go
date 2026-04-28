package services

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"github.com/kirillshkro/gmart-loyalty/internal/types"
)

type Registerer interface {
	Register(w http.ResponseWriter, r *http.Request)
}

// Регистрирует нового пользователя в системе
func (u UserService) Register(w http.ResponseWriter, r *http.Request) {
	var (
		regUser model.User
		profile model.UserProfile
		ok      bool
		err     error
	)
	if err = json.NewDecoder(r.Body).Decode(&regUser); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request", http.StatusInternalServerError)
	}

	profile = model.UserProfile{
		User: regUser,
	}

	if err = u.Repo.Create(profile); err != nil {
		log.Println(err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		if _, ok = errors.AsType[types.ErrDuplicateUser](err); ok {
			w.WriteHeader(http.StatusConflict)
		}
	}
	w.WriteHeader(http.StatusOK)
}
