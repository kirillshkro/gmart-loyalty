package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
)

// Тест AuthMiddleware для проверки авторизации пользователя.
func (s *TestUserSuite) Test_AuthMiddleware() {

	// Создаем токен для теста
	user := model.User{
		Login:    "testuser2",
		Password: "password123",
	}

	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(user); err != nil {
		s.T().Error(err)
	}

	//Передать куки

	reqTestUser := httptest.NewRequest(http.MethodPost, "/api/user/register", &body)
	rr := httptest.NewRecorder()
	s.service.Register(rr, reqTestUser)
	resp := rr.Result()
	defer resp.Body.Close()
	tcookie := resp.Cookies()
	if err := json.NewEncoder(&body).Encode(user); err != nil {
		s.T().Error(err)
	}
	wrapped := s.service.AuthMiddleware(http.HandlerFunc(s.service.Login)).(http.HandlerFunc)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", &body)
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	req.AddCookie(tcookie[0])
	wrapped(rr, req)
	resp = rr.Result()
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)
	for _, cookie := range resp.Cookies() {
		s.Assert().Equal("auth_cookie", cookie.Name)
	}
}
