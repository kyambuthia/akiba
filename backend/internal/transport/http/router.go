package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"akiba/backend/internal/auth"
	"akiba/backend/internal/usecase"

	"github.com/go-chi/chi/v5"
)

func NewRouter(logger *slog.Logger, authService *usecase.AuthService, jwtMgr *auth.JWTManager, readinessCheck func(context.Context) error) http.Handler {
	r := chi.NewRouter()
	r.Use(RequestID())
	r.Use(Recoverer())
	r.Use(Logging(logger))

	h := NewAuthHandler(authService)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/signup", h.Signup)
		r.Post("/auth/login", h.Login)
		r.With(RequireAuth(jwtMgr)).Get("/me", h.Me)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		if readinessCheck == nil {
			writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := readinessCheck(ctx); err != nil {
			writeError(w, http.StatusServiceUnavailable, "service_unavailable", "service not ready", nil)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})
	return r
}
