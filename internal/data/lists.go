package data

import (
	"context"
	"database/sql"
	"time"
)

// ReadingList represents a reading list in the book club system.
type ReadingList struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   int       `json:"created_by"`
	Books       []int     `json:"books"` // List of book IDs
	Status      string    `json:"status"` // "currently reading" or "completed"
	CreatedAt   time.Time `json:"created_at"`
}

// ReadingListModel handles the database interactions for reading lists.
type ReadingListModel struct {
	DB *sql.DB
}

// Insert a new reading list
func (m *ReadingListModel) Insert(list *ReadingList) error {
	query := `
		INSERT INTO reading_lists (name, description, created_by, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	args := []interface{}{list.Name, list.Description, list.CreatedBy, list.Status, time.Now()}

	return m.DB.QueryRow(query, args...).Scan(&list.ID)
}

// Get a single reading list by ID
func (m *ReadingListModel) Get(id int) (*ReadingList, error) {
	query := `
		SELECT id, name, description, created_by, status, created_at
		FROM reading_lists
		WHERE id = $1`

	var list ReadingList
	err := m.DB.QueryRow(query, id).Scan(
		&list.ID, &list.Name, &list.Description, &list.CreatedBy, &list.Status, &list.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNoRecord
	}
	return &list, err
}

// Update an existing reading list
func (m *ReadingListModel) Update(list *ReadingList) error {
	query := `
		UPDATE reading_lists
		SET name = $1, description = $2, status = $3
		WHERE id = $4`
	args := []interface{}{list.Name, list.Description, list.Status, list.ID}

	_, err := m.DB.Exec(query, args...)
	return err
}

// Delete a reading list by ID
func (m *ReadingListModel) Delete(id int) error {
	query := `DELETE FROM reading_lists WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	return err
}
