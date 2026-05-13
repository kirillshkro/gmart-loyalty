package services

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"github.com/kirillshkro/gmart-loyalty/internal/repository"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestUserSuite struct {
	suite.Suite
	service *Service
}

func (s *TestUserSuite) SetupSuite() {
	var err error
	l := logger.NewSlogLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), logger.Config{
		Colorful:             true,
		ParameterizedQueries: false,
		SlowThreshold:        1000 * time.Millisecond,
		LogLevel:             logger.Info,
	})
	opts := gorm.Config{
		PrepareStmt:    false,
		TranslateError: true,
		Logger:         l,
	}
	s.service = NewService()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &opts)
	if err != nil {
		s.T().Error(err)
	}
	s.service.Repo = repository.NewRepository(db)
	if err = db.AutoMigrate(&model.UserProfile{}); err != nil {
		s.T().Error(err)
	}
}

func (s *TestUserSuite) TearDownSuite() {

}

// Тест регистрации пользователя
func (s *TestUserSuite) Test_RegisterUser() {
	testCases := []struct {
		name         string
		login        string
		password     string
		expectedCode int
	}{
		{
			name:         "Normal register",
			login:        "John Doe",
			password:     "password123",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Empty Username",
			login:        "",
			password:     "password123",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Empty password",
			login:        "John Doe",
			password:     "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Empty Username and password",
			login:        "",
			password:     "",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			user := model.User{
				Login:    tc.login,
				Password: tc.password,
			}
			userJSON, _ := json.Marshal(user)
			req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(userJSON))
			w := httptest.NewRecorder()
			s.service.Register(w, req)
			s.Equal(http.StatusBadRequest, w.Code)
		})
	}
}

// Проверяем обработку запроса с существующим именем пользователя
func (s *TestUserSuite) Test_RegisterDuplicateUsername() {
	user1 := model.User{
		Login:     "existinguser",
		Password:  "password",
		Password2: "password",
	}
	user2 := model.User{
		Login:     "existinguser",
		Password:  "password2",
		Password2: "password2",
	}

	reqBody1, _ := json.Marshal(user1)
	req1 := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(reqBody1))
	w := httptest.NewRecorder()
	s.service.Register(w, req1)
	resp := w.Result()
	defer resp.Body.Close()
	s.Assert().Greater(len(resp.Cookies()), 0)
	if s.Assert().Equal(http.StatusOK, resp.StatusCode) {
		reqBody2, _ := json.Marshal(user2)
		req2 := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(reqBody2))
		w = httptest.NewRecorder()
		s.service.Register(w, req2)
		resp = w.Result()
		defer resp.Body.Close()
		s.Assert().Equal(http.StatusConflict, resp.StatusCode)
	}
}

func (s *TestUserSuite) Test_PasswordsNotEquals() {
	user := model.User{
		Login:     "newuser",
		Password:  "password1",
		Password2: "password2",
	}
	reqBody, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()
	s.service.Register(w, req)
	s.Equal(http.StatusBadRequest, w.Code)
}

func TestMainUserSuite(t *testing.T) {
	suite.Run(t, new(TestUserSuite))
}
