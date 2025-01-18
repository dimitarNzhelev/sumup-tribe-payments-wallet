package user

import (
	"context"
	"database/sql"
)

type PostgresUsersRepo struct {
	db *sql.DB
}

func NewPostgresUsersRepo(db *sql.DB) *PostgresUsersRepo {
	return &PostgresUsersRepo{db: db}
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
}

func (r *PostgresUsersRepo) CreateUser(ctx context.Context, user *User) error {
	_, err := r.db.Exec("INSERT INTO users (first_name, last_name, email, password_hash) VALUES ($1, $2, $3, $4)", user.FirstName, user.LastName, user.Email, user.PasswordHash)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresUsersRepo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	err := r.db.QueryRow("SELECT id, first_name, last_name, email, password_hash FROM users WHERE email = $1", email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgresUsersRepo) GetUserByID(ctx context.Context, id string) (*User, error) {
	user := &User{}
	err := r.db.QueryRow("SELECT id, first_name, last_name, email, password_hash FROM users WHERE id = $1", id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return user, nil
}
