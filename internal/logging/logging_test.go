package logging

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMiddleware_LogsAndPreservesStatus(t *testing.T) {
	// Перехватываем вывод логгера в буфер
	var buf bytes.Buffer
	old := log.Default().Writer()
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(old)
	}()

	// handler возвращает статус 201 и тело
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("created"))
	})

	wrapped := Middleware(handler)

	req := httptest.NewRequest(http.MethodPost, "/items", nil)
	req.Header.Set("X-Request-ID", "test-id-42")
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	// Проверяем, что статус дошёл до клиента
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	// Проверяем, что лог содержит request_id и путь
	logContent := buf.String()
	if !strings.Contains(logContent, "request_id=test-id-42") {
		t.Fatalf("expected log to contain request_id, got: %q", logContent)
	}
	if !strings.Contains(logContent, "path=/items") {
		t.Fatalf("expected log to contain path, got: %q", logContent)
	}
}
