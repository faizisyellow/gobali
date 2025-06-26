package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
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
	Bookings  []Booking    `json:"bookings"`
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
		switch err {
		case sql.ErrNoRows:
			return ErrNoRows
		default:
			return err
		}
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

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT users.id,users.email,users.password,r.name FROM users JOIN roles r ON users.role_id = r.id WHERE email = ? AND is_active = 1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	row := u.db.QueryRowContext(ctx, query, email)

	user := User{}
	err := row.Scan(&user.Id, &user.Email, &user.Password.Hash, &user.Role.Name)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNoRows
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (u *UserRepository) GetByID(ctx context.Context, userId int) (*User, error) {
	query := `
	SELECT users.id, username, email, password, role_id, roles.id, roles.name, roles.level, roles.description
	FROM users JOIN roles ON users.role_id = roles.id WHERE users.id = ? AND is_active = 1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res := u.db.QueryRowContext(ctx, query, userId)

	user := User{}
	err := res.Scan(
		&user.Id, &user.Username, &user.Email,
		&user.Password.Hash,
		&user.RoleId,
		&user.Role.Id,
		&user.Role.Name,
		&user.Role.Level,
		&user.Role.Description,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNoRows
		default:
			return nil, err
		}
	}

	return &user, nil
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

func (u *UserRepository) GetUserInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.is_active FROM users u
		JOIN user_invitation ui ON u.id = ui.user_id
		WHERE ui.token = ? AND ui.expire > ? 
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])
	user := &User{}

	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(&user.Id, &user.Username, &user.Email, &user.IsActive)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNoRows
		default:
			return nil, err
		}
	}

	return user, nil
}

func (u *UserRepository) UpdateWithTx(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `UPDATE users SET username = ?, is_active = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	// updating the user
	_, err := tx.ExecContext(ctx, query, &user.Username, &user.IsActive, &user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) deleteUserInvitation(ctx context.Context, tx *sql.Tx, userID int) error {
	query := `DELETE FROM user_invitation WHERE user_id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) DeleteUser(ctx context.Context, tx *sql.Tx, userId int) error {
	query := `DELETE FROM users WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) Delete(ctx context.Context, userId int) error {
	return withTx(u.db, ctx, func(tx *sql.Tx) error {
		err := u.DeleteUser(ctx, tx, userId)
		if err != nil {
			return err
		}

		err = u.deleteUserInvitation(ctx, tx, userId)
		if err != nil {
			return err
		}

		return nil
	})
}

func (u *UserRepository) Activate(ctx context.Context, token string) error {
	return withTx(u.db, ctx, func(tx *sql.Tx) error {
		user, err := u.GetUserInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		user.IsActive = true
		err = u.UpdateWithTx(ctx, tx, user)
		if err != nil {
			return err
		}

		err = u.deleteUserInvitation(ctx, tx, user.Id)
		if err != nil {
			return err
		}

		return nil
	})
}

// TODO: populate the user with their bookings villa
func (u *UserRepository) GetUserBookings(ctx context.Context, userId int, pq PaginatedUserBookingsQuery) (*User, error) {
	query := `
	SELECT u.id,u.email,b.villa_name,b.status b.total_price b.created_at FROM users u LEFT JOIN bookings b ON b.user_id = u.id
	ORDER BY b.created_at ` + pq.Sort + ` LIMIT ? OFFSET ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := u.db.ExecContext(ctx, query, pq.Limit, pq.Offset)
	if err != nil {
		return nil, err
	}

	return nil, nil

}
