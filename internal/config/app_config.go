package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func GetAppConfig() *AppConfig {
	var (
		cfg  AppConfig
		once sync.Once
	)

	once.Do(func() {
		if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
			log.Fatalln(err)
		}
	})
	return &cfg
}
