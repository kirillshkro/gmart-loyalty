package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/kirillshkro/gmart-loyalty/internal/config"
	"github.com/kirillshkro/gmart-loyalty/internal/middleware"
	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"github.com/kirillshkro/gmart-loyalty/internal/repository"
	"github.com/kirillshkro/gmart-loyalty/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := config.GetAppConfig()
	router := setupRouter(cfg)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func setupRouter(cfg *config.AppConfig) *mux.Router {
	db, err := setupDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	service := services.NewService()
	service.Repo = repository.NewRepository(db)
	r := mux.NewRouter()
	r.HandleFunc("/api/user/register", service.Register).Methods(http.MethodPost)
	r.HandleFunc("/api/user/login", service.Login).Methods(http.MethodPost)
	r.HandleFunc("/api/user/orders", service.SetOrderUser).Methods(http.MethodPost)
	r.HandleFunc("/api/user/orders", service.OrdersByUser).Methods(http.MethodGet)
	r.HandleFunc("api/user/balance", service.UserBalance).Methods(http.MethodGet)
	r.HandleFunc("/api/user/balance/withdraw", service.SetUserWithdraw).Methods(http.MethodPost)
	r.HandleFunc("/api/user/withdrawals", service.UserWithdrawals).Methods(http.MethodGet)
	r.Use(middleware.LoggerHandler)
	r.Use(service.AuthMiddleware)
	return r
}

func setupDB(cfg *config.AppConfig) (*gorm.DB, error) {
	dbLog := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURI), &gorm.Config{
		Logger:         logger.NewSlogLogger(dbLog, logger.Config{Colorful: true, SlowThreshold: 1000 * time.Millisecond}),
		PrepareStmt:    true,
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&model.UserProfile{}, &model.Order{}, &model.UserBalance{})
	return db, err
}
