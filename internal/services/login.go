package services

import (
	"encoding/json"
	"net/http"

	ctypes "github.com/kirillshkro/gmart-loyalty/internal/types"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"github.com/kirillshkro/gmart-loyalty/internal/model/claims"
)

type Loginer interface {
	Login(w http.ResponseWriter, r *http.Request)
}

func (s Service) Login(w http.ResponseWriter, r *http.Request) {
	//Получить пользователя из запроса
	var (
		user    model.User
		err     error
		ok      bool
		profile model.UserProfile
	)

	if user, ok = r.Context().Value(ctypes.UserID).(model.User); !ok {
		if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
			s.logger.Error("Bad request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	//Получить пользователя из базы данных
	if profile, err = s.Repo.UserByName(user.Login); err != nil {
		s.logger.Error("Invalid credentials")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userCookie, err := s.createCookie(profile.ID)
	if err != nil {
		s.logger.Error("Error create cookie")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, userCookie)

	w.WriteHeader(http.StatusOK)
}

func (u Service) userFromCookie(c *http.Cookie) (int, error) {
	tk := c.Value
	tkClaims := claims.NewUserClaims()
	userID, err := tkClaims.UserIDByToken(tk)
	return userID, err
}
