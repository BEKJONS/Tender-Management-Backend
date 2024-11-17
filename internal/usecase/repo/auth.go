package repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"tender_management/internal/entity"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

// CRUD for Users Table

func (u *UserRepo) CreateUser(user entity.User) (entity.User, error) {
	query := `INSERT INTO users (username, password, role, email) VALUES ($1, $2, $3, $4) RETURNING id`
	err := u.db.QueryRowx(query, user.Username, user.Password, user.Role, user.Email).Scan(&user.ID)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (u *UserRepo) IsEmailExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := u.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}

func (u *UserRepo) GetUserByUsername(username string) (entity.User, error) {
	query := `SELECT id, username, password, role, email FROM users WHERE username = $1`
	var user entity.User
	err := u.db.Get(&user, query, username)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (u *UserRepo) UpdateUser(user entity.User) (entity.User, error) {
	query := `UPDATE users SET username=$1, password=$2, role=$3, email=$4 WHERE id = $5 RETURNING id`
	err := u.db.QueryRowx(query, user.Username, user.Password, user.Role, user.Email, user.ID).Scan(&user.ID)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to update user: %w", err)
	}
	return user, nil
}

func (u *UserRepo) DeleteUser(userID string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := u.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
