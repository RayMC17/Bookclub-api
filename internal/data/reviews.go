package data

import (
	"context"
	"database/sql"
	"time"
)

// Review represents a review of a book.
type Review struct {
	ID         int       `json:"id"`
	BookID     int       `json:"book_id"`
	UserID     int       `json:"user_id"`
	Rating     int       `json:"rating"`
	Content    string    `json:"content"`
	ReviewDate time.Time `json:"review_date"`
}

// ReviewModel handles the database interactions for reviews.
type ReviewModel struct {
	DB *sql.DB
}

// Insert adds a new review to the database.
func (m *ReviewModel) Insert(review *Review) error {
	query := `
		INSERT INTO reviews (book_id, user_id, rating, content, review_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	args := []interface{}{review.BookID, review.UserID, review.Rating, review.Content, time.Now()}

	return m.DB.QueryRow(query, args...).Scan(&review.ID)
}

// Get retrieves a review by its ID.
func (m *ReviewModel) Get(id int) (*Review, error) {
	query := `
		SELECT id, book_id, user_id, rating, content, review_date
		FROM reviews
		WHERE id = $1`

	var review Review
	err := m.DB.QueryRow(query, id).Scan(
		&review.ID, &review.BookID, &review.UserID, &review.Rating, &review.Content, &review.ReviewDate,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNoRecord
	}
	return &review, err
}

// Update modifies an existing review in the database.
func (m *ReviewModel) Update(review *Review) error {
	query := `
		UPDATE reviews
		SET rating = $1, content = $2
		WHERE id = $3`
	args := []interface{}{review.Rating, review.Content, review.ID}

	_, err := m.DB.Exec(query, args...)
	return err
}

// Delete removes a review by its ID.
func (m *ReviewModel) Delete(id int) error {
	query := `DELETE FROM reviews WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	return err
}
