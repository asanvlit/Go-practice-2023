package repository

import (
	myErrors "Golang-practice-2023/internal/domain/apperrors"
	"Golang-practice-2023/internal/domain/logger"
	"Golang-practice-2023/internal/domain/user"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db     *sqlx.DB
	logger logger.Logger
}

func New(db *sqlx.DB, logger logger.Logger) *Repository {
	return &Repository{db: db, logger: logger}
}

func (r *Repository) GetDbInstance() *sqlx.DB {
	return r.db
}

func (r *Repository) Create(ctx context.Context, user *user.User) error {
	query := "INSERT INTO account (email, passwordhash) VALUES ($1, $2) RETURNING id"

	row := r.db.QueryRowContext(ctx, query, user.Email, user.Passwordhash)
	err := row.Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetById(ctx context.Context, id uuid.UUID) (*user.User, error) {
	query := "SELECT id, email, passwordhash FROM account WHERE id=$1"

	var u user.User
	err := r.db.GetContext(ctx, &u, query, id)
	if err != nil {
		return nil, myErrors.ErrUserNotFound
	}

	return &u, nil
}

func (r *Repository) Update(ctx context.Context, user *user.User) error {
	query := "UPDATE account SET email=$1, passwordhash=$2 WHERE id=$3"

	result, err := r.db.ExecContext(ctx, query, user.Email, user.Passwordhash, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("account not found")
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM account WHERE id=$1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("account not found")
	}

	return nil
}
