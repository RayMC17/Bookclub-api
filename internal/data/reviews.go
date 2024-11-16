package data

import (
    "context"
    "database/sql"
    "errors"
	"fmt"
    "time"

    "github.com/RayMC17/bookclub-api/internal/validator"
)

var ErrNoRecord = errors.New("record not found")

// Review represents a review for a book.
type Review struct {
    ID        int64     `json:"id"`
    BookID    int64     `json:"book_id"`
    Author    string    `json:"author"`
    Rating    int       `json:"rating"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}

// ReviewModel wraps a SQL database connection pool.
type ReviewModel struct {
    DB *sql.DB
}

// ValidateReview validates the review data.
func ValidateReview(v *validator.Validator, review *Review) {
    v.Check(review.Author != "", "author", "must be provided")
    v.Check(len(review.Author) <= 100, "author", "must not be more than 100 characters long")

    v.Check(review.Rating >= 1 && review.Rating <= 5, "rating", "must be between 1 and 5")
    
    v.Check(review.Content != "", "content", "must be provided")
    v.Check(len(review.Content) <= 1000, "content", "must not be more than 1000 characters long")
}

// Insert adds a new review to the database.
func (m *ReviewModel) Insert(review *Review) error {
    query := `
        INSERT INTO reviews (book_id, author, rating, content, created_at)
        VALUES ($1, $2, $3, $4, NOW())
        RETURNING id, created_at`

    args := []interface{}{review.BookID, review.Author, review.Rating, review.Content}

    return m.DB.QueryRow(query, args...).Scan(&review.ID, &review.CreatedAt)
}

// Get retrieves a specific review by ID.
func (m *ReviewModel) Get(id int64) (*Review, error) {
    query := `
        SELECT id, book_id, author, rating, content, created_at
        FROM reviews
        WHERE id = $1`

    var review Review

    err := m.DB.QueryRow(query, id).Scan(
        &review.ID,
        &review.BookID,
        &review.Author,
        &review.Rating,
        &review.Content,
        &review.CreatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, ErrNoRecord
    } else if err != nil {
        return nil, err
    }

    return &review, nil
}

// Update modifies the data of a specific review.
func (m *ReviewModel) Update(review *Review) error {
    query := `
        UPDATE reviews
        SET author = $1, rating = $2, content = $3
        WHERE id = $4`

    args := []interface{}{review.Author, review.Rating, review.Content, review.ID}

    _, err := m.DB.Exec(query, args...)
    return err
}

// Delete removes a specific review from the database.
func (m *ReviewModel) Delete(id int64) error {
    query := `
        DELETE FROM reviews
        WHERE id = $1`

    _, err := m.DB.Exec(query, id)
    return err
}


// GetAll retrieves all reviews for a specific book with optional filters for pagination and sorting.
func (m *ReviewModel) GetAll(bookID int64, author string, filters Filters) ([]*Review, Metadata, error) {
    query := `
        SELECT COUNT(*) OVER(), id, book_id, author, rating, content, created_at
        FROM reviews
        WHERE (book_id = $1)
        AND (author ILIKE '%' || $2 || '%' OR $2 = '')
        ORDER BY %s %s, id ASC
        LIMIT $3 OFFSET $4`

    formattedQuery := formatQuery(query, filters.SortColumn(), filters.SortDirection())

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    rows, err := m.DB.QueryContext(ctx, formattedQuery, bookID, author, filters.Limit(), filters.Offset())
    if err != nil {
        return nil, Metadata{}, err
    }
    defer rows.Close()

    totalRecords := 0
    reviews := []*Review{}

    for rows.Next() {
        var review Review
        err := rows.Scan(
            &totalRecords,
            &review.ID,
            &review.BookID,
            &review.Author,
            &review.Rating,
            &review.Content,
            &review.CreatedAt,
        )
        if err != nil {
            return nil, Metadata{}, err
        }
        reviews = append(reviews, &review)
    }

    err = rows.Err()
    if err != nil {
        return nil, Metadata{}, err
    }

    metadata := CalculateMetadata(totalRecords, filters.Page, filters.PageSize)
    return reviews, metadata, nil
}

// Helper function to safely format the SQL query with sort options.
func formatQuery(query, sortColumn, sortDirection string) string {
    return fmt.Sprintf(query, sortColumn, sortDirection)
}


func (m *ReviewModel) GetAllByUser(userID int64, filters Filters) ([]*Review, Metadata, error) {
    query := fmt.Sprintf(`
        SELECT COUNT(*) OVER(), id, book_id, author, content, rating
        FROM reviews
        WHERE user_id = $1
        ORDER BY %s %s, id ASC
        LIMIT $2 OFFSET $3`, filters.SortColumn(), filters.SortDirection())

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    rows, err := m.DB.QueryContext(ctx, query, userID, filters.Limit(), filters.Offset())
    if err != nil {
        return nil, Metadata{}, err
    }
    defer rows.Close()

    totalRecords := 0
    reviews := []*Review{}

    for rows.Next() {
        var review Review
        err := rows.Scan(
            &totalRecords,
            &review.ID,
            &review.BookID,
            &review.Author,
            &review.Content,
            &review.Rating,
        )
        if err != nil {
            return nil, Metadata{}, err
        }
        reviews = append(reviews, &review)
    }

    err = rows.Err()
    if err != nil {
        return nil, Metadata{}, err
    }

    metadata := CalculateMetadata(totalRecords, filters.Page, filters.PageSize)
    return reviews, metadata, nil
}