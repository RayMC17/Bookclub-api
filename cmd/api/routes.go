package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *applicationDependencies) routes() http.Handler {
	router := httprouter.New()

	// Health check route
	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", a.healthCheckHandler)

	// Books routes
	router.HandlerFunc(http.MethodGet, "/v1/books", a.listBooksHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/books/:id", a.getBookHandler)
	router.HandlerFunc(http.MethodPost, "/v1/books", a.createBookHandler)
	router.HandlerFunc(http.MethodPut, "/v1/books/:id", a.updateBookHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/books/:id", a.deleteBookHandler)

	router.HandlerFunc(http.MethodGet, "/api/v1/books", a.listBooksHandler)

	// Reading Lists routes
	router.HandlerFunc(http.MethodGet, "/v1/lists", a.listReadingListsHandler)//change to search
	router.HandlerFunc(http.MethodGet, "/v1/lists/:id", a.getReadingListHandler)
	router.HandlerFunc(http.MethodPost, "/v1/lists", a.createReadingListHandler)
	router.HandlerFunc(http.MethodPut, "/v1/lists/:id", a.updateReadingListHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/lists/:id", a.deleteReadingListHandler)
	router.HandlerFunc(http.MethodPost, "/v1/lists/:id/books", a.addBookToReadingListHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/lists/:id/books", a.removeBookFromReadingListHandler)

	// Reviews routes
	router.HandlerFunc(http.MethodGet, "/v1/books/:id/reviews", a.listReviewsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/books/:id/reviews", a.createReviewHandler)
	router.HandlerFunc(http.MethodPut, "/v1/reviews/:id", a.updateReviewHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/reviews/:id", a.deleteReviewHandler)

	// Users routes
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", a.getUserProfileHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/:id/lists", a.getUserReadingListsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/:id/reviews", a.getUserReviewsHandler)

	// Wrap the entire router with global middleware
	return a.logRequest(a.rateLimit(a.recoverPanic(router)))
}
