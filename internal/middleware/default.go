package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/fx"
	"golang.org/x/net/context"
)

const RequestIDKey = "request_id"

type DefaultMiddleware struct{}

func NewDefaultMiddleware() *DefaultMiddleware {
	return &DefaultMiddleware{}
}

func AddRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()

		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r)
	})
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}

	return "unknown"
}

var MiddlewareModule = fx.Module("middlewares",
	fx.Provide(NewDefaultMiddleware),
)
