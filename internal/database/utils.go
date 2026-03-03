package database

import (
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("Not found")

func CheckQueryResult(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
