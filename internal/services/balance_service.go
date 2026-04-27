package services

import (
	"net/http"

	"github.com/kirillshkro/gmart-loyalty/internal/config"
)

type BalanceService struct {
	cfg *config.AppConfig
}

type IBalanceService interface {
	UserBalance(w http.ResponseWriter, r *http.Request)
	SetUserWithdraw(w http.ResponseWriter, r *http.Request)
	UserWithdrawals(w http.ResponseWriter, r *http.Request)
}

func NewBalanceService(cfg *config.AppConfig) *BalanceService {
	return &BalanceService{
		cfg: cfg,
	}
}

func (b BalanceService) UserBalance(w http.ResponseWriter, r *http.Request) {
}

func (b BalanceService) SetUserWithdraw(w http.ResponseWriter, r *http.Request) {
}

func (b BalanceService) UserWithdrawals(w http.ResponseWriter, r *http.Request) {
}
