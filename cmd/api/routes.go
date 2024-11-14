package main

import (
    "net/http"
    "github.com/julienschmidt/httprouter"
)

func (a *applicationDependencies) routes() http.Handler {
    router := httprouter.New()
    router.NotFound = http.HandlerFunc(a.notFoundResponse)
    router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)

    // Book endpoints
    router.HandlerFunc(http.MethodGet, "/api/v1/books", a.listBooksHandler)
    router.HandlerFunc(http.MethodGet, "/api/v1/books/:id", a.displayBookHandler)
    router.HandlerFunc(http.MethodPost, "/api/v1/books", a.createBookHandler)
    router.HandlerFunc(http.MethodPut, "/api/v1/books/:id", a.updateBookHandler)
    router.HandlerFunc(http.MethodDelete, "/api/v1/books/:id", a.deleteBookHandler)
    router.HandlerFunc(http.MethodGet, "/api/v1/books/search", a.searchBooksHandler)

    // Reading List endpoints
    router.HandlerFunc(http.MethodGet, "/api/v1/lists", a.listReadingListsHandler)
    router.HandlerFunc(http.MethodGet, "/api/v1/lists/:id", a.displayReadingListHandler)
    router.HandlerFunc(http.MethodPost, "/api/v1/lists", a.createReadingListHandler)
    router.HandlerFunc(http.MethodPut, "/api/v1/lists/:id", a.updateReadingListHandler)
    router.HandlerFunc(http.MethodDelete, "/api/v1/lists/:id", a.deleteReadingListHandler)
    router.HandlerFunc(http.MethodPost, "/api/v1/lists/:id/books", a.addBookToListHandler)
    router.HandlerFunc(http.MethodDelete, "/api/v1/lists/:id/books", a.removeBookFromListHandler)

    // Review endpoints
    router.HandlerFunc(http.MethodGet, "/api/v1/books/:id/reviews", a.listReviewsHandler)
    router.HandlerFunc(http.MethodPost, "/api/v1/books/:id/reviews", a.createReviewHandler)
    router.HandlerFunc(http.MethodPut, "/api/v1/reviews/:id", a.updateReviewHandler)
    router.HandlerFunc(http.MethodDelete, "/api/v1/reviews/:id", a.deleteReviewHandler)

    // User endpoints
    router.HandlerFunc(http.MethodGet, "/api/v1/users/:id", a.getUserProfileHandler)
    router.HandlerFunc(http.MethodGet, "/api/v1/users/:id/lists", a.getUserReadingListsHandler)
    router.HandlerFunc(http.MethodGet, "/api/v1/users/:id/reviews", a.getUserReviewsHandler)

    return a.recoverPanic(router)
}
