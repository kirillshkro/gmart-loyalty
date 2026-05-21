package services

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"github.com/kirillshkro/gmart-loyalty/internal/repository"
	"github.com/kirillshkro/gmart-loyalty/internal/types"
	"github.com/kirillshkro/gmart-loyalty/pkg/utils"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type BalanceTestSuite struct {
	suite.Suite
	service *Service
	uCookie *http.Cookie
	resp    *http.Response
}

func (t *BalanceTestSuite) SetupSuite() {
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
		t.T().Error(err)
		return
	}
	t.service = NewService()
	t.service.Repo = repository.NewRepository(db)
	err = db.AutoMigrate(&model.UserProfile{}, &model.Order{}, &model.UserBalance{})
	if err != nil {
		t.T().Error(err)
		return
	}

	user := model.User{
		Login:     "balanceuser",
		Password:  "testpass",
		Password2: "testpass",
	}
	var reqBody bytes.Buffer
	if err := json.NewEncoder(&reqBody).Encode(user); err != nil {
		t.T().Fatal(err)
	}
	rr := httptest.NewRecorder()
	regReq := httptest.NewRequest(http.MethodPost, "/api/user/register", &reqBody)
	t.service.Register(rr, regReq)
	t.resp = rr.Result()
	t.uCookie = t.resp.Cookies()[0]
}

func (t *BalanceTestSuite) TearDownSuite() {
	t.resp.Body.Close()
}

func (t *BalanceTestSuite) Test_GetBalance() {
	userID, err := t.service.userFromCookie(t.uCookie)
	if err != nil {
		t.T().Fatal(err)
	}
	for range 5 {
		no := utils.GenNumberOrder(8)
		oReq := httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString(no))
		oReq.AddCookie(t.uCookie)
		rr := httptest.NewRecorder()
		t.service.SetOrderUser(rr, oReq)
	}
	ctx := context.WithValue(context.TODO(), types.UserID, userID)
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/user/balance", nil)
	req.AddCookie(t.uCookie)
	rr := httptest.NewRecorder()
	t.service.UserBalance(rr, req)

	t.resp = rr.Result()
	if t.Assert().Equal(http.StatusOK, t.resp.StatusCode) {
		var balance model.UserBalance
		if err = json.NewDecoder(t.resp.Body).Decode(&balance); err != nil {
			t.T().Fatal(err)
		}
		t.Assert().NotZero(balance.ID)
	}
}

func (t *BalanceTestSuite) Test_GetBalanceUnautorized() {
	userID, err := t.service.userFromCookie(t.uCookie)
	if err != nil {
		t.T().Fatal(err)
	}
	ctx := context.WithValue(context.TODO(), types.UserID, userID)
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/user/balance", nil)
	rr := httptest.NewRecorder()
	t.service.UserBalance(rr, req)

	t.resp = rr.Result()
	t.Assert().Equal(http.StatusUnauthorized, t.resp.StatusCode)
}

func TestBalance(t *testing.T) {
	suite.Run(t, new(BalanceTestSuite))
}
