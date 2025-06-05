package repository

import (
	"context"
	"database/sql"
)

type Repository struct {
	Users interface {
		Create(context.Context, *User) error
	}
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Users: &UserRepository{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fnc func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = fnc(tx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
