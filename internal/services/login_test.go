package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
)

func (s *TestUserSuite) Test_Login() {
	testCases := []struct {
		name         string
		username     string
		password     string
		expectedCode int
	}{
		{
			name:         "Match login",
			username:     "testusermatch",
			password:     "password123",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Mismatch login",
			username:     "testusermismatch",
			password:     "password123",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Empty login username",
			username:     "",
			password:     "password123",
			expectedCode: http.StatusBadRequest,
		},
	}

	//создать юзера для теста
	testUser := model.User{
		UserName:  "testusermatch",
		Password:  "password123",
		Password2: "password123",
	}
	var testUserBody bytes.Buffer
	err := json.NewEncoder(&testUserBody).Encode(testUser)
	s.Require().NoError(err)
	reqTestUser := httptest.NewRequest(http.MethodPost, "/api/user/register", &testUserBody)
	rr := httptest.NewRecorder()
	s.service.Register(rr, reqTestUser)
	authMiddleware := s.service.AuthMiddleware(http.HandlerFunc(s.service.Login)).(http.HandlerFunc)
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			user := model.User{
				UserName:  tc.username,
				Password:  tc.password,
				Password2: tc.password,
			}
			err := json.NewEncoder(&testUserBody).Encode(user)
			s.Require().NoError(err)
			req := httptest.NewRequest(http.MethodPost, "/api/user/login", &testUserBody)
			rr = httptest.NewRecorder()
			authMiddleware(rr, req)
			resp := rr.Result()
			defer resp.Body.Close()
			s.Equal(tc.expectedCode, resp.StatusCode)
		})
	}
}
