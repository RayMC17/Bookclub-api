package data

import (
	"context"
	"database/sql"
	"time"
)

// User represents a user in the book club system.
type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserModel handles the database interactions for users.
type UserModel struct {
	DB *sql.DB
}

// Insert adds a new user to the database.
func (m *UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (username, email, password, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
	args := []interface{}{user.Username, user.Email, user.Password, time.Now()}

	return m.DB.QueryRow(query, args...).Scan(&user.ID)
}

// Get retrieves a user by their ID.
func (m *UserModel) Get(id int) (*User, error) {
	query := `
		SELECT id, username, email, created_at
		FROM users
		WHERE id = $1`

	var user User
	err := m.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNoRecord
	}
	return &user, err
}

// Update modifies an existing user in the database.
func (m *UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2
		WHERE id = $3`
	args := []interface{}{user.Username, user.Email, user.ID}

	_, err := m.DB.Exec(query, args...)
	return err
}

// Delete removes a user by their ID.
func (m *UserModel) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	return err
}
