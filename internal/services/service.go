package services

import (
	"log"
	"log/slog"
	"os"

	"github.com/kirillshkro/gmart-loyalty/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Service struct {
	logger  *log.Logger
	cfg     *config.AppConfig
	authCfg *config.AuthConfig
	db      *gorm.DB
}

func NewService() (*Service, error) {
	cfg := config.GetAppConfig()
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURI), &gorm.Config{
		PrepareStmt:    true,
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}
	return &Service{
		logger:  slog.NewLogLogger(slog.NewJSONHandler(os.Stderr, nil), slog.LevelError),
		cfg:     cfg,
		authCfg: config.GetAuthConfig(),
		db:      db,
	}, nil
}
