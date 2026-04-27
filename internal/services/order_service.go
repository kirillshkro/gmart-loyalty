package services

import (
	"net/http"

	"github.com/kirillshkro/gmart-loyalty/internal/config"
)

type OrderService struct {
	cfg *config.AppConfig
}

type IOrderService interface {
	SetOrderUser(w http.ResponseWriter, r *http.Request)
	OrdersByUser(w http.ResponseWriter, r *http.Request)
}

func NewOrderService(cfg *config.AppConfig) *OrderService {
	return &OrderService{
		cfg: cfg,
	}
}

func (o *OrderService) SetOrderUser(w http.ResponseWriter, r *http.Request) {
}

func (o OrderService) OrdersByUser(w http.ResponseWriter, r *http.Request) {
}
