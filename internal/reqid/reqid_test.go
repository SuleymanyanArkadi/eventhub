package reqid

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMiddleware_SetsHeaderAndContext(t *testing.T) {
	// handler проверяет, что FromContext возвращает ненулевой id
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := FromContext(r.Context())
		if id == "" {
			http.Error(w, "no id", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(id))
	})

	wrapped := Middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	// ответ должен содержать X-Request-ID
	rid := rec.Header().Get(Header)
	if rid == "" {
		t.Fatalf("expected response header %s to be set", Header)
	}

	// тело тоже содержит id (наш handler писал id в тело)
	body := strings.TrimSpace(rec.Body.String())
	if body != rid {
		t.Fatalf("expected body %q to equal header %q", body, rid)
	}
}

func TestMiddleware_PreservesClientHeader(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := FromContext(r.Context())
		_, _ = w.Write([]byte(id))
	})
	wrapped := Middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(Header, "client-id-123")
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	if rec.Header().Get(Header) != "client-id-123" {
		t.Fatalf("expected header preserved, got %q", rec.Header().Get(Header))
	}
}
