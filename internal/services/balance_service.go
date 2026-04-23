package services

import "github.com/kirillshkro/gmart-loyalty/internal/config"

type BalanceService struct {
	cfg *config.AppConfig
}

func NewBalanceService(cfg *config.AppConfig) *BalanceService {
	return &BalanceService{
		cfg: cfg,
	}
}
