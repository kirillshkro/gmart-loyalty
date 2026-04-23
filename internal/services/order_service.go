package services

import "github.com/kirillshkro/gmart-loyalty/internal/config"

type OrderService struct {
	cfg *config.AppConfig
}

func NewOrderService(cfg *config.AppConfig) *OrderService {
	return &OrderService{
		cfg: cfg,
	}
}
