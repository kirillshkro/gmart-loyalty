package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kirillshkro/gmart-loyalty/internal/config"
	"github.com/kirillshkro/gmart-loyalty/internal/services"
)

func main() {
	cfg := config.GetAppConfig()
	router := setupRouter(cfg)
	log.Fatal(http.ListenAndServe(cfg.RunAddress, router))
}

func setupRouter(cfg *config.AppConfig) *mux.Router {
	userService := services.NewUserService()
	orderService := services.NewOrderService(cfg)
	balanceService := services.NewBalanceService(cfg)
	r := mux.NewRouter()
	r.HandleFunc("/api/user/register", userService.Register).Methods(http.MethodPost)
	r.HandleFunc("/api/user/login", userService.Login).Methods(http.MethodPost)
	r.HandleFunc("/api/user/orders", orderService.SetOrderUser).Methods(http.MethodPost)
	r.HandleFunc("/api/user/orders", orderService.OrdersByUser).Methods(http.MethodGet)
	r.HandleFunc("api/user/balance", balanceService.UserBalance).Methods(http.MethodGet)
	r.HandleFunc("/api/user/balance/withdraw", balanceService.SetUserWithdraw).Methods(http.MethodPost)
	r.HandleFunc("/api/user/withdrawals", balanceService.UserWithdrawals).Methods(http.MethodGet)
	r.Use(userService.AuthMiddleware)
	return r
}
