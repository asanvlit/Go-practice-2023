package user

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, user *User) error
	GetById(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetWithOffsetAndLimit(ctx context.Context, offset int, limit int) (*[]User, error)
	GetRegisteredLaterThenWithLimit(ctx context.Context, registerDate string, limit int) (*[]User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
