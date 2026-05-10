package services

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/kirillshkro/gmart-loyalty/internal/config"
	"github.com/kirillshkro/gmart-loyalty/internal/repository"
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
}

func (s Service) SetUserWithdraw(w http.ResponseWriter, r *http.Request) {
}

func (s Service) UserWithdrawals(w http.ResponseWriter, r *http.Request) {
}
