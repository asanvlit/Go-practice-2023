package user

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, user *User) error
	GetById(ctx context.Context, id uuid.UUID) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
