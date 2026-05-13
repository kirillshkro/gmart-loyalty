package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type LogWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (l LogWriter) WriteHeader(statusCode int) {
	l.status = statusCode
	l.ResponseWriter.WriteHeader(l.status)
}

func (l LogWriter) Write(b []byte) (int, error) {

	n, err := l.ResponseWriter.Write(b)
	l.size += n
	return n, err
}

func newLogWriter(resp http.ResponseWriter) *LogWriter {
	return &LogWriter{
		ResponseWriter: resp,
		size:           0,
		status:         http.StatusOK,
	}
}

func LoggerHandler(next http.Handler) http.Handler {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))
	fn := func(w http.ResponseWriter, r *http.Request) {
		writer := newLogWriter(w)
		method := r.Method
		uri := r.RequestURI
		startTime := time.Now()
		logger.Info(fmt.Sprint("Method: %s, uri: %s\n", method, uri))
		next.ServeHTTP(writer, r)
		duration := time.Since(startTime)
		logger.Info(fmt.Sprintf("Time: %d ms, req size: %d, req status: %d\n", duration.Milliseconds(), writer.size, writer.status))
	}
	return http.HandlerFunc(fn)
}
