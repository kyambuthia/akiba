package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"akiba/backend/internal/domain"
	"akiba/backend/internal/usecase"
)

type AuthHandler struct{ authService *usecase.AuthService }

func NewAuthHandler(authService *usecase.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type signupRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type userResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAt"`
	Status    string `json:"status,omitempty"`
}

func mapUser(u *domain.User) userResponse {
	return userResponse{ID: u.ID, Email: u.EmailLower, Phone: u.PhoneE164, Username: u.UsernameLower, CreatedAt: u.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"), Status: string(u.Status)}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload", nil)
		return
	}
	res, fields, err := h.authService.Signup(r.Context(), usecase.SignupInput(req))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			writeError(w, http.StatusBadRequest, "validation_error", "invalid signup payload", fields)
		case errors.Is(err, domain.ErrUserExists):
			writeError(w, http.StatusConflict, "user_exists", "user already exists", fields)
		default:
			writeError(w, http.StatusInternalServerError, "internal_error", "internal server error", nil)
		}
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"user": mapUser(res.User), "accessToken": res.AccessToken})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload", nil)
		return
	}
	res, fields, err := h.authService.Login(r.Context(), usecase.LoginInput(req))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			writeError(w, http.StatusBadRequest, "validation_error", "invalid login payload", fields)
		case errors.Is(err, domain.ErrInvalidCredentials):
			writeError(w, http.StatusUnauthorized, "invalid_credentials", "invalid login or password", nil)
		default:
			writeError(w, http.StatusInternalServerError, "internal_error", "internal server error", nil)
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": mapUser(res.User), "accessToken": res.AccessToken})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(ctxKeyUserID{}).(string)
	user, err := h.authService.Me(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrUnauthorized) {
			writeError(w, http.StatusUnauthorized, "unauthorized", "unauthorized", nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "internal server error", nil)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": mapUser(user)})
}

func userIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(ctxKeyUserID{}).(string)
	return v
}
