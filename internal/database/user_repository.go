package database

import "github.com/jmoiron/sqlx"

type UserRepository struct {
	db *sqlx.DB
}
