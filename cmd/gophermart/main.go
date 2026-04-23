package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kirillshkro/gmart-loyalty/internal/config"
	"github.com/kirillshkro/gmart-loyalty/internal/handlers"
	"github.com/kirillshkro/gmart-loyalty/internal/services"
)

func main() {
	cfg := config.GetAppConfig()
	router := setupRouter(cfg)
	log.Fatal(http.ListenAndServe(cfg.RunAddress, router))
}

func setupRouter(cfg *config.AppConfig) *mux.Router {
	userService := services.NewUserService(cfg)
	orderService := services.NewOrderService(cfg)
	balanceService := services.NewBalanceService(cfg)
	r := mux.NewRouter()
	r.Handle("/api/user/register", handlers.RegisterUser(userService)).Methods(http.MethodPost)
	r.Handle("/api/user/login", handlers.LoginUser(userService))
	r.Handle("/api/user/orders", handlers.SetOrderUser(orderService)).Methods(http.MethodPost)
	r.Handle("/api/user/orders", handlers.OrdersByUser(orderService)).Methods(http.MethodGet)
	r.Handle("api/user/balance", handlers.UserBalance(balanceService)).Methods(http.MethodGet)
	r.Handle("/api/user/balance/withdraw", handlers.SetUserWithdraw(balanceService)).Methods(http.MethodPost)
	r.Handle("/api/user/withdrawals", handlers.UserWithdrawals(balanceService)).Methods(http.MethodGet)
	return r
}
