package repository

import (
	"akiba/backend/internal/domain"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByLogin(ctx context.Context, login string) (*domain.User, error)
	EnsureIndexes(ctx context.Context) error
}
