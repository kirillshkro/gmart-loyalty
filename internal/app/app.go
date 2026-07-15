package app

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

type App struct {
	cfg       *config.AppConfig
	router    *mux.Router
	server    *http.Server
	interrupt chan os.Signal
}

func NewApp(cfg *config.AppConfig) (*App, error) {
	app := &App{
		cfg:       cfg,
		interrupt: make(chan os.Signal, 1),
	}
	return app, nil
}

func (a *App) setupDB() (*gorm.DB, error) {
	dbLog := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	db, err := gorm.Open(postgres.Open(a.cfg.DatabaseURI), &gorm.Config{
		Logger: logger.NewSlogLogger(dbLog, logger.Config{
			Colorful:             true,
			SlowThreshold:        200 * time.Millisecond,
			LogLevel:             logger.Info,
			ParameterizedQueries: true,
		}),
		PrepareStmt:    true,
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)
	err = db.AutoMigrate(&model.UserProfile{}, &model.Order{}, &model.UserBalance{}, &model.Withdrawal{})
	return db, err
}

func (a *App) setupService(db *gorm.DB) (*services.Service, error) {
	service := services.NewService()
	service.Repo = repository.NewRepository(db)
	return service, nil
}

func (a *App) setupRouter(service *services.Service) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/user/register", service.Register).Methods(http.MethodPost)
	router.HandleFunc("/api/user/login", service.Login).Methods(http.MethodPost)

	authRouter := router.PathPrefix("/api/user").Subrouter()
	authRouter.HandleFunc("/orders", service.SetOrderUser).Methods(http.MethodPost)
	authRouter.HandleFunc("/orders", service.OrdersByUser).Methods(http.MethodGet)
	authRouter.HandleFunc("/balance", service.UserBalance).Methods(http.MethodGet)
	authRouter.HandleFunc("/balance/withdraw", service.SetUserWithdraw).Methods(http.MethodPost)
	authRouter.HandleFunc("/withdrawals", service.UserWithdrawals).Methods(http.MethodGet)
	router.Use(middleware.LoggerMiddleware)
	authRouter.Use(middleware.LoggerMiddleware)
	return router
}

func (a *App) runServer() error {
	server := &http.Server{
		Addr:    a.cfg.RunAddress,
		Handler: a.router,
	}
	go func() {
		log.Printf("server is listening on %s\n", a.cfg.RunAddress)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error listen server is %s\n", err.Error())
		}
	}()
	a.server = server
	return nil
}

func (a *App) gracefulShutdown() {
	signal.Notify(a.interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-a.interrupt
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}
	log.Println("Shutting down server...")
}

func (a *App) Run() error {

	db, err := a.setupDB()
	if err != nil {
		return err
	}

	service, err := a.setupService(db)
	if err != nil {
		return err
	}

	a.router = a.setupRouter(service)

	err = a.runServer()
	if err != nil {
		return err
	}

	a.gracefulShutdown()
	return nil
}
