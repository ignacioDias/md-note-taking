package database

import (
	"context"
	"database/sql"
	"errors"
	"mdeditor/internal/domain"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

var ErrUserNotFound error = errors.New("User not found")

func (userRepo *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
	INSERT INTO users (email, name, hashed_password, profile_picture)
	VALUES (:email, :name, :hashed_password, :profile_picture)
	RETURNING id, created_at
	`
	stmt, err := userRepo.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	return stmt.GetContext(ctx, user, user)
}

func (userRepo *UserRepository) FindUserByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	query := "SELECT id, email, name, created_at, hashed_password, profile_picture FROM users WHERE id = $1"
	if err := userRepo.db.GetContext(ctx, &user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (userRepo *UserRepository) DeleteUserByID(ctx context.Context, id int64) error {
	query := "DELETE FROM users WHERE id = $1"
	result, err := userRepo.db.ExecContext(ctx, query, id)
	ret := CheckQueryResult(result, err)
	if ret == ErrNotFound {
		return ErrUserNotFound
	}
	return ret
}

func (userRepo *UserRepository) UpdateUserByID(ctx context.Context, id int64, newUser *domain.User) error {
	query := `UPDATE users SET email = $1, name = $2, hashed_password = $3, profile_picture = $4 WHERE id = $5`
	result, err := userRepo.db.ExecContext(ctx, query, newUser.Email, newUser.Name, newUser.HashedPassword, newUser.ProfilePicture, id)
	ret := CheckQueryResult(result, err)
	if ret == ErrNotFound {
		return ErrUserNotFound
	}
	return ret
}
