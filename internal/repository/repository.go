package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicateCategory = errors.New("category already exists")
	ErrDuplicateLocation = errors.New("location already exist")
	ErrNoRows            = errors.New("records not found")
	QueryTimeoutDuration = 5 * time.Second
)

type Repository struct {
	Users interface {
		Create(context.Context, *User) error
		CreateWithTx(context.Context, *sql.Tx, *User) error
		CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error
		Delete(context.Context, int) error
		Activate(context.Context, string) error
		GetUserInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error)
		UpdateWithTx(ctx context.Context, tx *sql.Tx, user *User) error
	}
	Roles interface {
		Create(context.Context, *Role) error
		GetByName(context.Context, string) (*Role, error)
	}
	Categories interface {
		Create(ctx context.Context, name string) error
		GetByID(ctx context.Context, id int) (*Category, error)
		GetCategories(ctx context.Context) ([]*Category, error)
		Update(ctx context.Context, category *Category) error
		Delete(ctx context.Context, id int) error
	}
	Location interface {
		Create(ctx context.Context, area string) error
		GetByID(ctx context.Context, id int) (*Location, error)
		GetLocations(ctx context.Context) ([]*Location, error)
		Update(ctx context.Context, location *Location) error
		Delete(ctx context.Context, id int) error
	}
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Users:      &UserRepository{db},
		Roles:      &RolesRepository{db},
		Categories: &CategoriesRepository{db},
		Location:   &LocationsRepository{db},
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
