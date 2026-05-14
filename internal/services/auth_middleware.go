package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"github.com/kirillshkro/gmart-loyalty/internal/model/claims"
	"github.com/kirillshkro/gmart-loyalty/internal/types"
	"golang.org/x/crypto/bcrypt"
)

const authCookie = "auth_cookie"

//Middleware, осуществляющее авторизацию пользователя.

func (s *Service) AuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		//Проверка существования куки
		if !s.cookieExist(r, authCookie) {
			w.WriteHeader(http.StatusUnauthorized)
			next.ServeHTTP(w, r)
			return
		}
		var user model.User
		newReq := io.NopCloser(r.Body)
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			next.ServeHTTP(w, r)
			return
		}
		//Проверить валидность логин/пароля
		if ok, err := s.validateUser(user); !ok {
			if _, ok := errors.AsType[*types.ErrInvalidLogin](err); ok {
				w.WriteHeader(http.StatusBadRequest)
				next.ServeHTTP(w, r)
				return
			}
			s.logger.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			next.ServeHTTP(w, r)
			return
		}
		r.Body = newReq
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (s Service) createCookie(userID int) (*http.Cookie, error) {
	authUser := claims.NewUserClaims()
	authUser.UserID = userID
	tk, err := authUser.Token()
	if err != nil {
		return nil, fmt.Errorf("error creating token: %w", err)
	}
	cookie := &http.Cookie{
		Name:     "auth_cookie",
		Value:    tk,
		Path:     "/",
		Expires:  time.Now().Add(24 * 7 * time.Hour),
		HttpOnly: true,
		Secure:   false,
	}
	return cookie, nil
}

func (s Service) cookieExist(req *http.Request, cookieName string) bool {
	newReq := io.NopCloser(req.Body)
	_, err := req.Cookie(cookieName)
	if err != nil {
		return false
	}
	req.Body = newReq
	return true
}

func (s Service) refreshCookie(ctx context.Context) (*http.Cookie, error) {
	userID := ctx.Value(UserID).(int)
	return s.createCookie(userID)
}

func (s Service) validateUser(user model.User) (bool, error) {
	//Проверка логина и пароля
	if user.Login == "" || user.Password == "" {
		return false, &types.ErrInvalidLogin{}
	}
	//Проверка сущетвования пользователя в базе данных
	profile, err := s.Repo.UserByName(user.Login)
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
