package data

import (
	"context"
	"database/sql"
	"time"
)

type Book struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Authors        []string  `json:"authors"`
	ISBN           string    `json:"isbn"`
	PublicationDate time.Time `json:"publication_date"`
	Genre          string    `json:"genre"`
	Description    string    `json:"description"`
	AverageRating  float64   `json:"average_rating"`
}

type BookModel struct {
	DB *sql.DB
}

// Insert a new book
func (m *BookModel) Insert(book *Book) error {
	query := `
		INSERT INTO books (title, authors, isbn, publication_date, genre, description, average_rating)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
	args := []interface{}{book.Title, book.Authors, book.ISBN, book.PublicationDate, book.Genre, book.Description, book.AverageRating}

	return m.DB.QueryRow(query, args...).Scan(&book.ID)
}

// Get a single book by ID
func (m *BookModel) Get(id int) (*Book, error) {
	query := `
		SELECT id, title, authors, isbn, publication_date, genre, description, average_rating
		FROM books
		WHERE id = $1`

	var book Book
	err := m.DB.QueryRow(query, id).Scan(
		&book.ID, &book.Title, &book.Authors, &book.ISBN,
		&book.PublicationDate, &book.Genre, &book.Description, &book.AverageRating,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNoRecord
	}
	return &book, err
}

// Update a book
func (m *BookModel) Update(book *Book) error {
	query := `
		UPDATE books
		SET title = $1, authors = $2, isbn = $3, publication_date = $4, genre = $5, description = $6, average_rating = $7
		WHERE id = $8`
	args := []interface{}{book.Title, book.Authors, book.ISBN, book.PublicationDate, book.Genre, book.Description, book.AverageRating, book.ID}

	_, err := m.DB.Exec(query, args...)
	return err
}

// Delete a book by ID
func (m *BookModel) Delete(id int) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	return err
}
