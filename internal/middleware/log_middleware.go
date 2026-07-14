package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggerMiddleware is a middleware that logs request information
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom ResponseWriter to capture the status code and response body
		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Process the request
		next.ServeHTTP(wrappedWriter, r)

		// Calculate duration
		duration := time.Since(start).Milliseconds()

		// Log the information
		log.Printf("URI: %s, Method: %s, Request Size: %d, Duration: %d ms, Status Code: %d",
			r.URL.RequestURI(),
			r.Method,
			r.ContentLength,
			duration,
			wrappedWriter.statusCode,
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body = b
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
