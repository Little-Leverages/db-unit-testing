package users

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID   string
	Name string
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepo {
	return &userRepo{db: db}
}

func (u *userRepo) Create(user User) error {
	query := `INSERT INTO users (id, name) VALUES ($1, $2)`
	_, err := u.db.Exec(query, user.ID, user.Name)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (u *userRepo) Get(id string) (*User, error) {
	query := `SELECT id, name FROM users WHERE id = $1`
	row := u.db.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}
