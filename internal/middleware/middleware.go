package middleware

import (
	"net/http"
)

// Middleware interface defines the contract for middleware
type Middleware interface {
	// LoggerMiddleware logs request information
	LoggerMiddleware(next http.Handler) http.Handler
}

// Mux is a struct that holds all middleware implementations
type Mux struct{}

// LoggerMiddleware is a middleware that logs request information
func (m *Mux) LoggerMiddleware(next http.Handler) http.Handler {
	return LoggerMiddleware(next)
}
