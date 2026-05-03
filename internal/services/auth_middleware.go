package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kirillshkro/gmart-loyalty/internal/model/claims"
)

//Middleware, осуществляющее авторизацию пользователя.

func (u *UserService) AuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Реализация логики проверки
		//Извлекаем токен из заголовка или куки
		userCookie, err := r.Cookie("auth_cookie")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				userCookie, err = u.createCookie()
				if err != nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				http.SetCookie(w, userCookie)
				//Выдадим UserID из куки
				tk := userCookie.Value
				authClaims := claims.NewUserClaims()
				userID, err := authClaims.UserIDByToken(tk)
				if err != nil {
					log.Print(err)
					http.Error(w, "Invalid token", http.StatusInternalServerError)
					next.ServeHTTP(w, r)
					return
				}
				c := context.WithValue(r.Context(), UserID, userID)
				r = r.WithContext(c)
				next.ServeHTTP(w, r)
				return
			}
			//Выдадим UserID из куки
			tk := userCookie.Value
			if tk != "" {
				http.Error(w, "Invalid token", http.StatusInternalServerError)
				next.ServeHTTP(w, r)
				return
			}
			authClaims := claims.NewUserClaims()
			userID, err := authClaims.UserIDByToken(tk)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusInternalServerError)
				next.ServeHTTP(w, r)
				return
			}
			c := context.WithValue(r.Context(), UserID, userID)
			r = r.WithContext(c)
			next.ServeHTTP(w, r)
			return
		}
		//Проверяем валидность токена
		if !u.validateToken(userCookie.Value) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// Реализация логики проверки JWT токена
func (u UserService) validateToken(token string) bool {
	// Проверка на пустоту токена
	if token == "" {
		return false
	}
	// Проверка валидности токена
	claims, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(u.authCfg.SecretKey), nil
	})

	if err != nil || !claims.Valid {
		return false
	}
	// Если токен валиден
	return true
}

func (u UserService) createCookie() (*http.Cookie, error) {
	authUser := claims.NewUserClaims()
	tk, err := authUser.Token()
	if err != nil {
		return nil, fmt.Errorf("error creating token: %w", err)
	}
	return &http.Cookie{
		Name:     "auth_cookie",
		Value:    tk,
		Expires:  time.Now().Add(24 * 7 * time.Hour),
		HttpOnly: true,
		Secure:   false,
	}, nil
}
