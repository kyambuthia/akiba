package usecase

import (
	"context"
	"testing"
	"time"

	"akiba/backend/internal/auth"
	"akiba/backend/internal/domain"
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

func TestSignupValidation(t *testing.T) {
	repo := &memRepo{users: map[string]*domain.User{}}
	svc := NewAuthService(repo, auth.NewJWTManager("secret", "test"), time.Hour)
	_, fields, err := svc.Signup(context.Background(), SignupInput{Email: "bad-email", Phone: "123", Username: "ab", Password: "weak"})
	if err == nil {
		t.Fatalf("expected error")
	}
	if fields["phone"] == "" || fields["username"] == "" || fields["password"] == "" {
		t.Fatalf("expected validation fields, got %#v", fields)
	}
}

func TestSignupAndLoginHappyPath(t *testing.T) {
	repo := &memRepo{users: map[string]*domain.User{}}
	svc := NewAuthService(repo, auth.NewJWTManager("secret", "test"), time.Hour)
	res, fields, err := svc.Signup(context.Background(), SignupInput{Email: "USER@example.com", Phone: "+14155552671", Username: "User_Name", Password: "Password1"})
	if err != nil || len(fields) > 0 {
		t.Fatalf("signup failed: err=%v fields=%#v", err, fields)
	}
	if res.User.EmailLower != "user@example.com" || res.User.UsernameLower != "user_name" {
		t.Fatalf("normalization failed: %#v", res.User)
	}
	loginRes, _, err := svc.Login(context.Background(), LoginInput{Login: "USER_NAME", Password: "Password1"})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if loginRes.AccessToken == "" {
		t.Fatalf("missing token")
	}
}

func TestLoginValidation(t *testing.T) {
	repo := &memRepo{users: map[string]*domain.User{}}
	svc := NewAuthService(repo, auth.NewJWTManager("secret", "test"), time.Hour)
	_, fields, err := svc.Login(context.Background(), LoginInput{Login: "", Password: ""})
	if err == nil {
		t.Fatalf("expected error")
	}
	if fields["login"] == "" || fields["password"] == "" {
		t.Fatalf("expected login/password validation fields, got %#v", fields)
	}
}
