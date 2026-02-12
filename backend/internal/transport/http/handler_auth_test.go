package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"akiba/backend/internal/auth"
	"akiba/backend/internal/domain"
	"akiba/backend/internal/usecase"
)

type memRepo struct{ users map[string]*domain.User }

func (m *memRepo) EnsureIndexes(ctx context.Context) error { return nil }
func (m *memRepo) Create(ctx context.Context, user *domain.User) error {
	for _, u := range m.users {
		if u.EmailLower == user.EmailLower || u.PhoneE164 == user.PhoneE164 || u.UsernameLower == user.UsernameLower {
			return domain.ErrUserExists
		}
	}
	user.ID = "u1"
	m.users[user.ID] = user
	return nil
}
func (m *memRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return u, nil
}
func (m *memRepo) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	for _, u := range m.users {
		if u.EmailLower == login || u.PhoneE164 == login || u.UsernameLower == login {
			return u, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func testRouter() http.Handler {
	repo := &memRepo{users: map[string]*domain.User{}}
	jwtMgr := auth.NewJWTManager("secret", "test")
	authSvc := usecase.NewAuthService(repo, jwtMgr, time.Hour)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewRouter(logger, authSvc, jwtMgr, func(ctx context.Context) error { return nil })
}

func TestSignupValidationError(t *testing.T) {
	r := testRouter()
	body := map[string]string{"email": "bad", "phone": "111", "username": "ab", "password": "123"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewReader(b))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSignupAndMeFlow(t *testing.T) {
	r := testRouter()
	body := map[string]string{"email": "user@example.com", "phone": "+14155552671", "username": "user_1", "password": "Password1"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewReader(b))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var out map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &out)
	tok, _ := out["accessToken"].(string)
	if tok == "" {
		t.Fatalf("expected token")
	}
	meReq := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	meReq.Header.Set("Authorization", "Bearer "+tok)
	meW := httptest.NewRecorder()
	r.ServeHTTP(meW, meReq)
	if meW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", meW.Code)
	}
}

func TestLoginValidationError(t *testing.T) {
	r := testRouter()
	body := map[string]string{"login": "", "password": ""}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(b))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSignupUnknownFieldError(t *testing.T) {
	r := testRouter()
	body := `{"email":"user@example.com","phone":"+14155552671","username":"user_1","password":"Password1","unexpected":"x"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestReadyEndpoint(t *testing.T) {
	r := testRouter()
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
