package services

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

//Middleware, осуществляющее авторизацию пользователя.

func (u *UserService) AuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Реализация логики проверки
		//Извлекаем токен из заголовка или куки
		userCookie, err := r.Cookie("user_token")
		if err != nil {
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
