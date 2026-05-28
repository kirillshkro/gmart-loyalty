package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	RunAddress           string `env:"RUN_ADDRESS" env-default:"localhost:8080"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

type AuthConfig struct {
	SecretKey string `env:"AUTH_SECRET_KEY"`
}

func GetAppConfig() *AppConfig {
	var (
		cfg  AppConfig
		once sync.Once
	)

	once.Do(func() {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			log.Fatalln(err)
		}
	})
	return &cfg
}

func GetAuthConfig() *AuthConfig {
	var (
		cfg  AuthConfig
		once sync.Once
	)
	once.Do(func() {
		if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
			log.Fatal(err)
		}
	})
	return &cfg
}
