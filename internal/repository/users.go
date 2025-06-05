package repository

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *sql.DB
}

type User struct {
	Id        int          `json:"id"`
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	Password  HashPassword `json:"-"`
	IsActive  bool         `json:"is_active"`
	RoleId    int          `json:"role_id"`
	CreatedAt string       `json:"created_at"`
	UpdateAt  string       `json:"update_at"`
}

type HashPassword struct {
	Text *string
	Hash []byte
}

func (h *HashPassword) Set(text string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	h.Text = &text
	h.Hash = hashed

	return nil
}

func (h *HashPassword) Compare(password string) error {
	err := bcrypt.CompareHashAndPassword(h.Hash, []byte(password))
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) Create(ctx context.Context, payload *User) error {

	query := `INSERT INTO users(username,email,password,role_id) VALUES(?,?,?,?)`

	if payload.RoleId == 0 {
		payload.RoleId = 1
	}

	_, err := u.db.ExecContext(ctx, query, payload.Username, payload.Email, payload.Password.Hash, payload.RoleId)
	if err != nil {
		return err
	}

	return nil
}
