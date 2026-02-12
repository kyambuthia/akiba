package http

import (
	"log/slog"
	"net/http"

	"akiba/backend/internal/auth"
	"akiba/backend/internal/usecase"

	"github.com/go-chi/chi/v5"
)

func NewRouter(logger *slog.Logger, authService *usecase.AuthService, jwtMgr *auth.JWTManager) http.Handler {
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
	return r
}
