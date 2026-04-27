package services

import (
	"github.com/kirillshkro/gmart-loyalty/internal/config"
)

type UserService struct {
	cfg     *config.AppConfig
	authCfg *config.AuthConfig
}

type IUserService interface {
	Registerer
	Loginer
}

func NewUserService() *UserService {
	userCfg := config.GetAppConfig()
	authCfg := config.GetAuthConfig()
	return &UserService{
		cfg:     userCfg,
		authCfg: authCfg,
	}
}
