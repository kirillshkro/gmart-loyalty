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
	"github.com/kirillshkro/gmart-loyalty/pkg/utils"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestOrderSuite struct {
	suite.Suite
	service *Service
	uCookie *http.Cookie
	resp    *http.Response
}

func (s *TestOrderSuite) SetupSuite() {
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
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &opts)
	if err != nil {
		s.T().Error(err)
		return
	}
	s.service = NewService()
	s.service.Repo = repository.NewRepository(db)
	err = db.AutoMigrate(&model.UserProfile{}, &model.Order{}, &model.UserBalance{})
	if err != nil {
		s.T().Error(err)
		return
	}

	user := model.User{
		Login:    "testuser",
		Password: "testpass",
	}
	var reqBody bytes.Buffer
	if err := json.NewEncoder(&reqBody).Encode(user); err != nil {
		s.T().Fatal(err)
	}
	rr := httptest.NewRecorder()
	regReq := httptest.NewRequest(http.MethodPost, "/api/user/register", &reqBody)
	s.service.Register(rr, regReq)
	s.resp = rr.Result()
	s.uCookie = s.resp.Cookies()[0]
}

func (s *TestOrderSuite) TearDownSuite() {
	s.resp.Body.Close()
}

// Тест для проверки корректности работы метода SetOrder,
// Юзер существует в базе данных, другие юзеры не создавали заказ с этим же номером заказа.
func (s *TestOrderSuite) Test_NormalSetOrder() {
	firstOrder := utils.GenNumberOrder(8)
	secondOrder := firstOrder

	reqFirstOrder := httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString(firstOrder))
	rr := httptest.NewRecorder()
	reqFirstOrder.AddCookie(s.uCookie)
	s.service.SetOrderUser(rr, reqFirstOrder)
	if s.Assert().Equal(http.StatusAccepted, rr.Code) {
		reqSecOrder := httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString(secondOrder))
		rr = httptest.NewRecorder()
		reqSecOrder.AddCookie(s.uCookie)
		s.service.SetOrderUser(rr, reqSecOrder)
		s.Assert().Equal(http.StatusOK, rr.Code)
	}
}

// Тест если пользователь неавторизован
func (s *TestOrderSuite) Test_UnautorizedUser() {
	//Не сохраняем куки
	orderNum := utils.GenNumberOrder(8)
	reqOrder := httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString(orderNum))
	rr := httptest.NewRecorder()
	s.service.SetOrderUser(rr, reqOrder)

	s.Assert().Equal(http.StatusUnauthorized, rr.Code)
}

func (s *TestOrderSuite) Test_AnotherUser() {
	user2 := model.User{
		Login:    "another",
		Password: "dirtyharry",
	}
	var reqBody2 bytes.Buffer

	if err := json.NewEncoder(&reqBody2).Encode(user2); err != nil {
		s.T().Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/user/register", &reqBody2)
	rr := httptest.NewRecorder()
	s.service.Register(rr, req)

	s.resp = rr.Result()

	user2Cookie := s.resp.Cookies()[0]

	numOrder := utils.GenNumberOrder(8)

	reqOrder1 := httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString(numOrder))
	reqOrder1.AddCookie(s.uCookie)
	rr = httptest.NewRecorder()
	s.service.SetOrderUser(rr, reqOrder1)
	if s.Assert().Equal(http.StatusAccepted, rr.Code) {
		reqOrder2 := httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString(numOrder))
		reqOrder2.AddCookie(user2Cookie)
		rr = httptest.NewRecorder()
		s.service.SetOrderUser(rr, reqOrder2)
		s.Assert().Equal(http.StatusConflict, rr.Code)
	}
}

func (s *TestOrderSuite) Test_InvalidNumOrder() {
	numOrder := "10000111"

	reqOrder1 := httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString(numOrder))
	reqOrder1.AddCookie(s.uCookie)
	rr := httptest.NewRecorder()
	s.service.SetOrderUser(rr, reqOrder1)
	s.Assert().Equal(http.StatusUnprocessableEntity, rr.Code)
}

func TestOrder(t *testing.T) {
	suite.Run(t, new(TestOrderSuite))
}
