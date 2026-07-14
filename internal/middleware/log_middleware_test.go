package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggerMiddleware(t *testing.T) {
	// Создаем тестовый обработчик
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Применяем middleware
	middleware := LoggerMiddleware(handler)

	// Создаем тестовый запрос
	req := httptest.NewRequest("GET", "/test", strings.NewReader("test request body"))
	req.Header.Set("Content-Type", "application/json")

	// Создаем тестовый ResponseWriter
	rec := httptest.NewRecorder()

	// Выполняем запрос
	middleware.ServeHTTP(rec, req)

	// Проверяем, что ответ был обработан правильно
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Body.String() != "test response" {
		t.Errorf("Expected body 'test response', got '%s'", rec.Body.String())
	}
}

func TestLoggerMiddlewareWithDifferentMethods(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	middleware := LoggerMiddleware(handler)

	methods := []string{"GET", "POST", "PUT", "DELETE"}

	for _, method := range methods {
		req := httptest.NewRequest(method, "/test", nil)
		rec := httptest.NewRecorder()

		middleware.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status code %d for method %s, got %d", http.StatusCreated, method, rec.Code)
		}
	}
}

func TestLoggerMiddlewareWithNoContent(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	middleware := LoggerMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, rec.Code)
	}
}

func TestLoggerMiddlewareWithPostRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := LoggerMiddleware(handler)

	req := httptest.NewRequest("POST", "/test", strings.NewReader(`{"key": "value"}`))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestLoggerMiddlewareStructure(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := LoggerMiddleware(handler)

	// Проверяем, что middleware возвращает правильный тип
	if middleware == nil {
		t.Error("Middleware should not be nil")
	}

	// Создаем тестовый запрос
	req := httptest.NewRequest("GET", "/api/test", nil)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	// Проверяем, что запрос обработан
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}
