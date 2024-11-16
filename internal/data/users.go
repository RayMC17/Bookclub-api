package data

import (
	//"context"
	 //"errors"
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
	UpdatedAt   time.Time `json:"updated_at"`
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
		return nil, ErrRecordNotFound
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

// // Get retrieves a user by ID.
// func (m *UserModel) Get(id int64) (*User, error) {
//     // SQL query to retrieve a user by ID
//     query := `
//         SELECT id, name, email, created_at, updated_at
//         FROM users
//         WHERE id = $1`

//     // Initialize an empty User object to hold the result
//     var user User

//     // Context with a timeout for the query execution
//     ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//     defer cancel()

//     // Execute the query and scan the result into the user struct
//     err := m.DB.QueryRowContext(ctx, query, id).Scan(
//         &user.ID,
//         &user.Username,
//         &user.Email,
//         &user.CreatedAt,
//         &user.UpdatedAt,
//     )

//     // Handle any errors that occur
//     if err != nil {
//         if errors.Is(err, sql.ErrNoRows) {
//             return nil, ErrNoRecord
//         }
//         return nil, err
//     }

//     return &user, nil
// }