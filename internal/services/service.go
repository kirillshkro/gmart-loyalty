package services

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/kirillshkro/gmart-loyalty/internal/config"
	"github.com/kirillshkro/gmart-loyalty/internal/repository"
	"github.com/kirillshkro/gmart-loyalty/internal/types"
)

type Service struct {
	logger  *log.Logger
	cfg     *config.AppConfig
	authCfg *config.AuthConfig
	Repo    repository.IRepository
}

func NewService() *Service {
	cfg := config.GetAppConfig()

	return &Service{
		logger:  slog.NewLogLogger(slog.NewJSONHandler(os.Stderr, nil), slog.LevelError),
		cfg:     cfg,
		authCfg: config.GetAuthConfig(),
	}
}

type IBalanceService interface {
	UserBalance(w http.ResponseWriter, r *http.Request)
	SetUserWithdraw(w http.ResponseWriter, r *http.Request)
	UserWithdrawals(w http.ResponseWriter, r *http.Request)
}

func (s Service) UserBalance(w http.ResponseWriter, r *http.Request) {
	if !s.cookieExist(r, authCookie) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID := r.Context().Value(types.UserID).(int)
	balance, err := s.Repo.BalanceByUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(balance); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
