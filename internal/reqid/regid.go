package reqid

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// уникальный тип ключа для context, чтобы избежать конфликтов
type ctxKey struct{}

var key = &ctxKey{}

// Header - имя заголовка для request ID
const Header = "X-Request-ID"

// Middleware добавляет request id к запросу (если отсутствует — генерирует),
// кладёт его в контекст и устанавливает заголовок в ответ.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(Header)
		if rid == "" {
			rid = uuid.New().String()
			// также записываем в заголовок запроса, чтобы downstream видели его
			r.Header.Set(Header, rid)
		}
		// возвращаем id в ответе
		w.Header().Set(Header, rid)

		// кладем id в context запроса
		ctx := context.WithValue(r.Context(), key, rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FromContext извлекает request id из контекста или пустую строку
func FromContext(ctx context.Context) string {
	if v := ctx.Value(key); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
