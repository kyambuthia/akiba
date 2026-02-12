package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"akiba/backend/internal/auth"
	"akiba/backend/internal/domain"
	"akiba/backend/internal/repository"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type SignupInput struct {
	Email    string `validate:"required,email"`
	Phone    string `validate:"required"`
	Username string `validate:"required"`
	Password string `validate:"required"`
}
type LoginInput struct {
	Login    string `validate:"required"`
	Password string `validate:"required"`
}
type AuthResult struct {
	User        *domain.User
	AccessToken string
}

type AuthService struct {
	users          repository.UserRepository
	jwt            *auth.JWTManager
	validate       *validator.Validate
	accessTokenTTL time.Duration
}

func NewAuthService(users repository.UserRepository, jwtMgr *auth.JWTManager, accessTokenTTL time.Duration) *AuthService {
	return &AuthService{users: users, jwt: jwtMgr, validate: validator.New(), accessTokenTTL: accessTokenTTL}
}

func (s *AuthService) Signup(ctx context.Context, in SignupInput) (*AuthResult, domain.FieldErrors, error) {
	fields := domain.FieldErrors{}
	if err := s.validate.Var(in.Email, "required,email"); err != nil {
		fields["email"] = "must be a valid email address"
	}
	email := domain.NormalizeEmail(in.Email)
	phone := domain.NormalizePhone(in.Phone)
	username := domain.NormalizeUsername(in.Username)
	password := strings.TrimSpace(in.Password)
	if !domain.ValidatePhoneE164(phone) {
		fields["phone"] = "must be valid E.164 format"
	}
	if !domain.ValidateUsername(username) {
		fields["username"] = "must be 3-20 chars and only letters, numbers, underscore"
	}
	if !domain.ValidatePassword(password) {
		fields["password"] = "must be at least 8 chars and include a letter and number"
	}
	if len(fields) > 0 {
		return nil, fields, domain.ErrInvalidInput
	}
	now := time.Now().UTC()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}
	user := &domain.User{EmailLower: email, PhoneE164: phone, UsernameLower: username, PasswordHash: string(hash), Status: domain.UserStatusActive, CreatedAt: now, UpdatedAt: now}
	if err := s.users.Create(ctx, user); err != nil {
		if errors.Is(err, domain.ErrUserExists) {
			return nil, domain.FieldErrors{"login": "email, phone, or username already exists"}, domain.ErrUserExists
		}
		return nil, nil, err
	}
	token, err := s.jwt.IssueAccessToken(user.ID, s.accessTokenTTL)
	if err != nil {
		return nil, nil, err
	}
	return &AuthResult{User: user, AccessToken: token}, nil, nil
}

func (s *AuthService) Login(ctx context.Context, in LoginInput) (*AuthResult, domain.FieldErrors, error) {
	fields := domain.FieldErrors{}
	if strings.TrimSpace(in.Login) == "" {
		fields["login"] = "is required"
	}
	password := strings.TrimSpace(in.Password)
	if password == "" {
		fields["password"] = "is required"
	}
	if len(fields) > 0 {
		return nil, fields, domain.ErrInvalidInput
	}
	login := strings.TrimSpace(in.Login)
	if strings.Contains(login, "@") {
		login = domain.NormalizeEmail(login)
	} else if strings.HasPrefix(login, "+") {
		login = domain.NormalizePhone(login)
	} else {
		login = domain.NormalizeUsername(login)
	}
	user, err := s.users.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, nil, domain.ErrInvalidCredentials
		}
		return nil, nil, err
	}
	if user.Status != domain.UserStatusActive {
		return nil, nil, domain.ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)) != nil {
		return nil, nil, domain.ErrInvalidCredentials
	}
	token, err := s.jwt.IssueAccessToken(user.ID, s.accessTokenTTL)
	if err != nil {
		return nil, nil, err
	}
	return &AuthResult{User: user, AccessToken: token}, nil, nil
}

func (s *AuthService) Me(ctx context.Context, userID string) (*domain.User, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, domain.ErrUnauthorized
	}
	return s.users.GetByID(ctx, userID)
}
