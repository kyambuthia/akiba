package http

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"akiba/backend/internal/auth"

	"github.com/go-chi/chi/v5/middleware"
)

type ctxKeyUserID struct{}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) { r.status = code; r.ResponseWriter.WriteHeader(code) }

func Recoverer() func(http.Handler) http.Handler { return middleware.Recoverer }
func RequestID() func(http.Handler) http.Handler { return middleware.RequestID }

func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r)
			logger.Info("http_request", "method", r.Method, "path", r.URL.Path, "status", rec.status, "latency_ms", time.Since(start).Milliseconds(), "request_id", middleware.GetReqID(r.Context()))
		})
	}
}

func RequireAuth(jwtMgr *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				writeError(w, http.StatusUnauthorized, "unauthorized", "missing or invalid bearer token", nil)
				return
			}
			claims, err := jwtMgr.Verify(strings.TrimSpace(parts[1]))
			if err != nil || claims.Sub == "" {
				writeError(w, http.StatusUnauthorized, "unauthorized", "invalid token", nil)
				return
			}
			ctx := context.WithValue(r.Context(), ctxKeyUserID{}, claims.Sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
