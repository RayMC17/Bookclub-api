package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/RayMC17/bookclub-api/internal/data"
	"github.com/RayMC17/bookclub-api/internal/validator"
)

func (a *applicationDependencies) createBookHandler(w http.ResponseWriter, r *http.Request) {
	var incomingData struct {
		Title           string   `json:"title"`
		Authors         []string `json:"authors"`
		ISBN            string   `json:"isbn"`
		PublicationDate string   `json:"publication_date"`
		Genre           string   `json:"genre"`
		Description     string   `json:"description"`
		AverageRating   float64  `json:"average_rating"`
	}
	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	book := &data.Book{
		Title:         incomingData.Title,
		Authors:       incomingData.Authors,
		ISBN:          incomingData.ISBN,
		Genre:         incomingData.Genre,
		Description:   incomingData.Description,
		AverageRating: incomingData.AverageRating,
	}

	v := validator.New()
	data.ValidateBook(v, book)
	if !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.bookModel.Insert(book)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/api/v1/books/%d", book.ID))
	data := envelope{"book": book}
	err = a.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) getBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	book, err := a.bookModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{"book": book}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	book, err := a.bookModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	var incomingData struct {
		Title         *string   `json:"title"`
		Authors       *[]string `json:"authors"`
		ISBN          *string   `json:"isbn"`
		Genre         *string   `json:"genre"`
		Description   *string   `json:"description"`
		AverageRating *float64  `json:"average_rating"`
	}
	err = a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if incomingData.Title != nil {
		book.Title = *incomingData.Title
	}
	if incomingData.Authors != nil {
		book.Authors = *incomingData.Authors
	}
	if incomingData.ISBN != nil {
		book.ISBN = *incomingData.ISBN
	}
	if incomingData.Genre != nil {
		book.Genre = *incomingData.Genre
	}
	if incomingData.Description != nil {
		book.Description = *incomingData.Description
	}
	if incomingData.AverageRating != nil {
		book.AverageRating = *incomingData.AverageRating
	}

	v := validator.New()
	data.ValidateBook(v, book)
	if !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.bookModel.Update(book)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{"book": book}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.bookModel.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{"message": "book successfully deleted"}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) listBooksHandler(w http.ResponseWriter, r *http.Request) {
	var queryParametersData struct {
		Title  string
		Author string
		data.Filters
	}

	queryParameters := r.URL.Query()

	// Load the query parameters into our struct
	queryParametersData.Title = a.getSingleQueryParameter(queryParameters, "title", "")
	queryParametersData.Author = a.getSingleQueryParameter(queryParameters, "author", "")
	v := validator.New()

	queryParametersData.Filters.Page = a.getSingleIntegerParameter(queryParameters, "page", 1, v)
	queryParametersData.Filters.PageSize = a.getSingleIntegerParameter(queryParameters, "page_size", 10, v)
	queryParametersData.Filters.Sort = a.getSingleQueryParameter(queryParameters, "sort", "id")
	queryParametersData.Filters.SortSafelist = []string{"id", "title", "author", "-id", "-title", "-author"}

	// Check if our filters are valid
	data.ValidateFilters(v, &queryParametersData.Filters)
	if !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	books, metadata, err := a.bookModel.GetAll(
		queryParametersData.Title,
		queryParametersData.Author,
		queryParametersData.Filters,
	)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	responseData := envelope{
		"books":     books,
		"@metadata": metadata,
	}
	err = a.writeJSON(w, http.StatusOK, responseData, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// searchBooksHandler handles searching for books based on query parameters.
func (a *applicationDependencies) searchBooksHandler(w http.ResponseWriter, r *http.Request) {
	var queryParams struct {
		Title  string
		Author string
		data.Filters
	}

	// Get the query parameters for title and author
	queryParams.Title = a.getSingleQueryParameter(r.URL.Query(), "title", "")
	queryParams.Author = a.getSingleQueryParameter(r.URL.Query(), "author", "")

	// Initialize the validator and set up filters
	v := validator.New()
	queryParams.Filters.Page = a.getSingleIntegerParameter(r.URL.Query(), "page", 1, v)
	queryParams.Filters.PageSize = a.getSingleIntegerParameter(r.URL.Query(), "page_size", 10, v)
	queryParams.Filters.Sort = a.getSingleQueryParameter(r.URL.Query(), "sort", "id")
	queryParams.Filters.SortSafelist = []string{"id", "title", "author", "-id", "-title", "-author"}

	// Validate filters
	data.ValidateFilters(v, &queryParams.Filters)
	if !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Search books in the database
	books, metadata, err := a.bookModel.GetAll(queryParams.Title, queryParams.Author, queryParams.Filters)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Create response with books and metadata
	response := envelope{
		"books":    books,
		"metadata": metadata,
	}

	// Write response as JSON
	err = a.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// listReadingListsHandler retrieves a list of reading lists.
func (a *applicationDependencies) listReadingListsHandler(w http.ResponseWriter, r *http.Request) {
	var filters data.Filters

	// Set up query parameters for pagination and sorting
	v := validator.New()
	filters.Page = a.getSingleIntegerParameter(r.URL.Query(), "page", 1, v)
	filters.PageSize = a.getSingleIntegerParameter(r.URL.Query(), "page_size", 10, v)
	filters.Sort = a.getSingleQueryParameter(r.URL.Query(), "sort", "id")
	filters.SortSafelist = []string{"id", "name", "-id", "-name"} // Define allowed sort fields

	// Validate filters
	data.ValidateFilters(v, &filters)
	if !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Retrieve reading lists from the database
	lists, metadata, err := a.readingListModel.GetAll(filters)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Send response with lists and metadata
	response := envelope{
		"reading_lists": lists,
		"metadata":      metadata,
	}
	err = a.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) getReadingListHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the reading list ID from the URL parameters
	idParam := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 1 {
		a.notFoundResponse(w, r)
		return
	}

	// Fetch the reading list from the database
	readingList, err := a.readingListModel.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			a.notFoundResponse(w, r)
		} else {
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// Send the reading list in the response
	err = a.writeJSON(w, http.StatusOK, envelope{"reading_list": readingList}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) createReadingListHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Books       []int  `json:"books"` // IDs of books in the list
		Status      string `json:"status"`
	}

	// Decode JSON body
	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	readingList := &data.ReadingList{
		Name:        input.Name,
		Description: input.Description,
		Books:       input.Books,
		Status:      input.Status,
	}

	// Validate the input
	v := validator.New()
	data.ValidateReadingList(v, readingList)
	if !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Insert the reading list into the database
	err = a.readingListModel.Insert(readingList)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", "/api/v1/lists/"+strconv.Itoa(readingList.ID))

	// Send the created reading list in the response
	err = a.writeJSON(w, http.StatusCreated, envelope{"reading_list": readingList}, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) updateReadingListHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the reading list ID from the URL parameters
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Fetch the existing reading list from the database
	readingList, err := a.readingListModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// Parse the JSON request body into an input struct
	var input struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Books       *[]int  `json:"books"` // IDs of books in the list
		Status      *string `json:"status"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Update the fields in the reading list based on the input
	if input.Name != nil {
		readingList.Name = *input.Name
	}
	if input.Description != nil {
		readingList.Description = *input.Description
	}
	if input.Books != nil {
		readingList.Books = *input.Books
	}
	if input.Status != nil {
		readingList.Status = *input.Status
	}

	// Validate the updated reading list
	v := validator.New()
	data.ValidateReadingList(v, readingList)
	if !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Save the updated reading list to the database
	err = a.readingListModel.Update(readingList)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Send the updated reading list in the response
	data := envelope{"reading_list": readingList}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) deleteReadingListHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the reading list ID from the URL.
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Delete the reading list from the database.
	err = a.readingListModel.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// Respond with a success message.
	response := envelope{"message": "reading list successfully deleted"}
	err = a.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) addBookToReadingListHandler(w http.ResponseWriter, r *http.Request) {
	// Get the reading list ID from the URL parameters
	readingListID, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Decode the request body to get the book ID
	var input struct {
		BookID int `json:"book_id"`
	}
	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Add the book to the reading list
	err = a.readingListModel.AddBook(readingListID, input.BookID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// Send a success response
	envelope := envelope{"message": "book successfully added to the reading list"}
	err = a.writeJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) removeBookFromReadingListHandler(w http.ResponseWriter, r *http.Request) {
	// Get the reading list ID from the URL parameters
	readingListID, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Decode the request body to get the book ID
	var input struct {
		BookID int `json:"book_id"`
	}
	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Remove the book from the reading list
	err = a.readingListModel.RemoveBook(readingListID, input.BookID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// Send a success response
	envelope := envelope{"message": "book successfully removed from the reading list"}
	err = a.writeJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) listReviewsHandler(w http.ResponseWriter, r *http.Request) {
	var queryParams struct {
		Rating int
		Author string
		data.Filters
	}

	queryParams.Rating = a.getSingleIntegerParameter(r.URL.Query(), "rating", 0, validator.New())
	queryParams.Author = a.getSingleQueryParameter(r.URL.Query(), "author", "")

	v := validator.New()
	queryParams.Filters.Page = a.getSingleIntegerParameter(r.URL.Query(), "page", 1, v)
	queryParams.Filters.PageSize = a.getSingleIntegerParameter(r.URL.Query(), "page_size", 10, v)
	queryParams.Filters.Sort = a.getSingleQueryParameter(r.URL.Query(), "sort", "id")
	queryParams.Filters.SortSafelist = []string{"id", "rating", "author", "-id", "-rating", "-author"}

	data.ValidateFilters(v, &queryParams.Filters)
	if !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	reviews, metadata, err := a.reviewModel.GetAll(int64(queryParams.Rating), queryParams.Author, queryParams.Filters)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	response := envelope{
		"reviews":  reviews,
		"metadata": metadata,
	}
	err = a.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// createReviewHandler handles the creation of a new review for a specific book.
func (a *applicationDependencies) createReviewHandler(w http.ResponseWriter, r *http.Request) {
	// Parse book ID from the URL
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Define a structure to hold the expected data from the request body
	var input struct {
		Author  string `json:"author"`
		Content string `json:"content"`
		Rating  int    `json:"rating"`
	}

	// Parse JSON request body
	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Create a Review instance with the parsed data
	review := &data.Review{
		BookID:  int64(id),
		Author:  input.Author,
		Content: input.Content,
		Rating:  input.Rating,
	}

	// Validate the review data
	v := validator.New()
	data.ValidateReview(v, review)
	if !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Insert the new review into the database
	err = a.reviewModel.Insert(review)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Send a response with the created review
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/api/v1/books/%d/reviews/%d", id, review.ID))
	err = a.writeJSON(w, http.StatusCreated, envelope{"review": review}, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) updateReviewHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the review ID from the URL
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}
	id64 := int64(id)
	// Fetch the existing review from the database
	review, err := a.reviewModel.Get(id64)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// Define a struct for holding the updated data
	var input struct {
		Content *string `json:"content"`
		Rating  *int    `json:"rating"`
	}

	// Parse the input from the request body
	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Update the review fields if new data is provided
	if input.Content != nil {
		review.Content = *input.Content
	}
	if input.Rating != nil {
		review.Rating = *input.Rating
	}

	// Validate the updated review
	v := validator.New()
	data.ValidateReview(v, review)
	if !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Save the updated review
	err = a.reviewModel.Update(review)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Send the updated review in the response
	response := envelope{"review": review}
	err = a.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) deleteReviewHandler(w http.ResponseWriter, r *http.Request) {
    // Parse the review ID from the URL and convert it to int64
    id, err := a.readIDParam(r)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    // Convert the id to int64 if it's not already
    id64 := int64(id)

    // Delete the review
    err = a.reviewModel.Delete(id64)
    if err != nil {
        switch {
        case errors.Is(err, data.ErrNoRecord):
            a.notFoundResponse(w, r)
        default:
            a.serverErrorResponse(w, r, err)
        }
        return
    }

    // Respond with a 204 No Content status code
    w.WriteHeader(http.StatusNoContent)
}

func (a *applicationDependencies) getUserProfileHandler(w http.ResponseWriter, r *http.Request) {
    // Extract user ID from the URL parameters
    id, err := a.readIDParam(r)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    // Get the user profile from the database using the user model
    profile, err := a.userModel.Get(id)
    if err != nil {
        switch {
        case errors.Is(err, data.ErrNoRecord):
            a.notFoundResponse(w, r)
        default:
            a.serverErrorResponse(w, r, err)
        }
        return
    }

    // Respond with the user profile data in JSON format
    err = a.writeJSON(w, http.StatusOK, envelope{"user_profile": profile}, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}

func (a *applicationDependencies) getUserReadingListsHandler(w http.ResponseWriter, r *http.Request) {
    // Get the user ID from the URL parameters
    id, err := a.readIDParam(r)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    // Initialize the filters for pagination and sorting
    var filters data.Filters
    v := validator.New()

    filters.Page = a.getSingleIntegerParameter(r.URL.Query(), "page", 1, v)
    filters.PageSize = a.getSingleIntegerParameter(r.URL.Query(), "page_size", 10, v)
    filters.Sort = a.getSingleQueryParameter(r.URL.Query(), "sort", "id")
    filters.SortSafelist = []string{"id", "name", "-id", "-name"}

    data.ValidateFilters(v, &filters)
    if !v.Valid() {
        a.failedValidationResponse(w, r, v.Errors)
        return
    }

    // Get the reading lists associated with the user from the model
    readingLists, metadata, err := a.readingListModel.GetAllByUser(int64(id), filters)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    // Respond with the reading lists and metadata in JSON format
    response := envelope{
        "reading_lists": readingLists,
        "metadata":      metadata,
    }
    err = a.writeJSON(w, http.StatusOK, response, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}


func (a *applicationDependencies) getUserReviewsHandler(w http.ResponseWriter, r *http.Request) {
    // Get the user ID from the URL parameters
    id, err := a.readIDParam(r)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    // Initialize the filters for pagination and sorting
    var filters data.Filters
    v := validator.New()

    filters.Page = a.getSingleIntegerParameter(r.URL.Query(), "page", 1, v)
    filters.PageSize = a.getSingleIntegerParameter(r.URL.Query(), "page_size", 10, v)
    filters.Sort = a.getSingleQueryParameter(r.URL.Query(), "sort", "id")
    filters.SortSafelist = []string{"id", "rating", "author", "-id", "-rating", "-author"}

    data.ValidateFilters(v, &filters)
    if !v.Valid() {
        a.failedValidationResponse(w, r, v.Errors)
        return
    }

    // Get the reviews associated with the user from the model
    reviews, metadata, err := a.reviewModel.GetAllByUser(int64(id), filters)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    // Respond with the reviews and metadata in JSON format
    response := envelope{
        "reviews":  reviews,
        "metadata": metadata,
    }
    err = a.writeJSON(w, http.StatusOK, response, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}
