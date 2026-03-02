package logging

import (
	"log"
	"net/http"
	"time"

	"github.com/SuleymanyanArkadi/eventhub/internal/reqid"
)

// statusRecorder перехватывает код ответа, чтобы логировать его потом.
type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (r *statusRecorder) WriteHeader(code int) {
	if !r.wroteHeader {
		r.status = code
		r.wroteHeader = true
		r.ResponseWriter.WriteHeader(code)
	}
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	// Если заголовок еще не был записан, то по умолчанию статус 200 OK
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	return r.ResponseWriter.Write(b)
}

// Middleware логирует вход/выход запросы: method, path, status, duration, request-id.
// Он пытается взять request-id из контекста (reqid.FromContext), а если нет — из заголовка X-Request-ID.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		// Выполняем handler
		next.ServeHTTP(rec, r)

		// Получаем request-id для логирования
		id := reqid.FromContext(r.Context())
		if id == "" {
			id = r.Header.Get(reqid.Header)
		}

		duration := time.Since(start)
		log.Printf("method=%s path=%s status=%d duration=%s request_id=%s", r.Method, r.URL.Path, rec.status, duration.String(), id)
	})
}
