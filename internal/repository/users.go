package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

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
	Role      Role         `json:"role"`
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

	query := `INSERT INTO users(username,email,password,role_id) VALUES(?,?,?,(SELECT id FROM roles WHERE name = ?))`

	if payload.Role.Name == "" {
		payload.Role.Name = "user"
	}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := u.db.ExecContext(ctx, query, payload.Username, payload.Email, payload.Password.Hash, payload.Role.Name)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, payload *User) error {

	query := `INSERT INTO users(username,email,password,role_id) VALUES(?,?,?,(SELECT id FROM roles WHERE name = ?))`

	if payload.Role.Name == "" {
		payload.Role.Name = "user"
	}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	row, err := tx.ExecContext(ctx, query, payload.Username, payload.Email, payload.Password.Hash, payload.Role.Name)
	if err != nil {
		duplicateKey := "Error 1062"
		switch {
		case strings.Contains(err.Error(), duplicateKey):
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	userId, err := row.LastInsertId()
	if err != nil {
		return err
	}

	payload.Id = int(userId)

	return nil
}

func (u *UserRepository) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userID int) error {
	query := `INSERT INTO user_invitation(token,user_id,expire) VALUES(?,?,?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(exp))
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {

	return withTx(u.db, ctx, func(tx *sql.Tx) error {

		if err := u.CreateWithTx(ctx, tx, user); err != nil {
			return err
		}

		if err := u.createUserInvitation(ctx, tx, token, invitationExp, user.Id); err != nil {
			return err
		}

		return nil
	})
}

func (u *UserRepository) Delete(ctx context.Context, userId int) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := u.db.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}

	return nil
}
