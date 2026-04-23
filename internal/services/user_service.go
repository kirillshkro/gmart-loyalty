package services

import "github.com/kirillshkro/gmart-loyalty/internal/config"

type UserService struct {
	cfg *config.AppConfig
}

func NewUserService(cfg *config.AppConfig) *UserService {
	return &UserService{
		cfg: cfg,
	}
}
