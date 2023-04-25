package repository

import (
	"Go-scheduler-service/internal/domain/apperrors"
	"Go-scheduler-service/internal/domain/logger"
	"Go-scheduler-service/internal/domain/user"
	"context"
	"database/sql"
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
	if u, _ := r.GetByEmail(ctx, user.Email); u != nil {
		return apperrors.ErrAlreadyRegisteredUserEmail
	}

	query := "INSERT INTO account (email, passwordhash) VALUES ($1, $2) RETURNING id, createdAt, updatedAt"

	row := r.db.QueryRowContext(ctx, query, user.Email, user.Passwordhash)
	err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		r.logger.Warning(err.Error())
		return apperrors.ErrDbQueryProcessing
	}

	return nil
}

func (r *Repository) Save(ctx context.Context, user *user.User) error {
	if u, _ := r.GetByEmail(ctx, user.Email); u != nil {
		return apperrors.ErrAlreadyRegisteredUserEmail
	}

	query := "INSERT INTO account (email, passwordhash, createdat, updatedat) VALUES ($1, $2, $3, $4) RETURNING id"

	row := r.db.QueryRowContext(ctx, query, user.Email, user.Passwordhash, user.CreatedAt, user.UpdatedAt)
	err := row.Scan(&user.ID)
	if err != nil {
		r.logger.Warning(err.Error())
		return apperrors.ErrDbQueryProcessing
	}

	return nil
}

func (r *Repository) GetById(ctx context.Context, id uuid.UUID) (*user.User, error) {
	query := "SELECT id, email, passwordhash, createdAt, updatedAt FROM account WHERE id=$1"

	var u user.User
	err := r.db.GetContext(ctx, &u, query, id)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}

	return &u, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := "SELECT id, email, passwordhash, createdAt, updatedAt FROM account WHERE email=$1"

	var u user.User
	err := r.db.GetContext(ctx, &u, query, email)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}

	return &u, nil
}

func (r *Repository) GetLastRegisteredUser(ctx context.Context) (*user.User, error) {
	query := "SELECT id, email, passwordhash, createdAt, updatedAt FROM account ORDER BY createdAt DESC LIMIT 1"

	var u user.User
	err := r.db.GetContext(ctx, &u, query)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}

	return &u, nil
}

func (r *Repository) GetRegisteredLaterThen(ctx context.Context, registerDate string, limit int) (*[]user.User, error) {
	var users []user.User

	query := "SELECT id, email, passwordhash, createdAt, updatedAt FROM account WHERE createdat > $1 LIMIT $2"

	err := r.db.SelectContext(ctx, &users, query, registerDate, limit)
	if err != nil {
		if err == sql.ErrNoRows {
			return &users, nil
		}
		return nil, apperrors.ErrDbQueryProcessing
	}

	return &users, nil
}

func (r *Repository) GetWithOffsetAndLimit(ctx context.Context, offset int, limit int) (*[]user.User, error) {
	var users []user.User

	query := "SELECT id, email, passwordhash, createdAt, updatedAt FROM account ORDER BY createdat OFFSET $1 LIMIT $2"

	err := r.db.SelectContext(ctx, &users, query, offset, limit)
	if err != nil {
		if err == sql.ErrNoRows {
			return &users, nil
		}
		return nil, apperrors.ErrDbQueryProcessing
	}

	return &users, nil
}

func (r *Repository) Update(ctx context.Context, user *user.User) error {
	if u, _ := r.GetById(ctx, user.ID); u == nil {
		return apperrors.ErrUserNotFound
	}

	query := "UPDATE account SET email=$1, passwordhash=$2, updatedAt=current_timestamp WHERE id=$3 RETURNING createdAt, updatedAt"

	row := r.db.QueryRowContext(ctx, query, user.Email, user.Passwordhash, user.ID)
	err := row.Scan(&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		r.logger.Warning(err.Error())
		return apperrors.ErrDbQueryProcessing
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM account WHERE id=$1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return apperrors.ErrDbQueryProcessing
	}

	if u, _ := r.GetById(ctx, id); u == nil {
		return apperrors.ErrUserNotFound
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperrors.ErrDbQueryProcessing
	}
	if rowsAffected == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}
