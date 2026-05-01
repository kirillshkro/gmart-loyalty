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
	service *UserService
}

func (s *TestUserSuite) SetupTest() {
	l := logger.NewSlogLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)), logger.Config{
		Colorful:             true,
		ParameterizedQueries: true,
		SlowThreshold:        1000 * time.Millisecond,
		LogLevel:             logger.Info,
	})
	opts := gorm.Config{
		PrepareStmt:    false,
		TranslateError: true,
		Logger:         l,
	}
	s.service = NewUserService()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &opts)
	if err != nil {
		s.T().Error(err)
	}
	s.service.Repo = repository.NewUserRepository(db)
	if err = db.AutoMigrate(&model.UserProfile{}); err != nil {
		s.T().Error(err)
	}
}

func (s *TestUserSuite) TearDownSuite() {

}

// Тест регистрации пользователя
func (s *TestUserSuite) Test_RegisterUserEmptyPassword() {
	testCases := []struct {
		name         string
		username     string
		password     string
		expectedCode int
	}{
		{
			name:         "Normal register",
			username:     "John Doe",
			password:     "password123",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Empty Username",
			username:     "",
			password:     "password123",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Empty password",
			username:     "John Doe",
			password:     "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Empty Username and password",
			username:     "",
			password:     "",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			user := model.User{
				UserName: tc.username,
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
		UserName:  "existinguser",
		Password:  "password",
		Password2: "password",
	}
	user2 := model.User{
		UserName:  "existinguser",
		Password:  "password2",
		Password2: "password2",
	}

	reqBody1, _ := json.Marshal(user1)
	req1 := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(reqBody1))
	w := httptest.NewRecorder()
	s.service.Register(w, req1)
	s.Equal(http.StatusOK, w.Code)
	reqBody2, _ := json.Marshal(user2)
	req2 := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(reqBody2))
	w = httptest.NewRecorder()
	s.service.Register(w, req2)
	s.Equal(http.StatusConflict, w.Code)
}

func (s *TestUserSuite) Test_PasswordsNotEquals() {
	user := model.User{
		UserName:  "newuser",
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
