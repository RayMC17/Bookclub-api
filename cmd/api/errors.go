package main

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrNoRecord = errors.New("record not found")
)

// logError is a helper function for logging errors.
func (a *applicationDependencies) logError(r *http.Request, err error) {
	a.logger.Error(err.Error(), "method", r.Method, "url", r.URL.String())
}

// errorResponseJSON sends a JSON-formatted error message with the specified status code.
func (a *applicationDependencies) errorResponseJSON(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	// Create an envelope containing the error message.
	errorData := envelope{"error": message}

	// Write the JSON response.
	err := a.writeJSON(w, status, errorData, nil)
	if err != nil {
		a.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// serverErrorResponse sends a 500 Internal Server Error response.
func (a *applicationDependencies) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	a.errorResponseJSON(w, r, http.StatusInternalServerError, message)
}

// notFoundResponse sends a 404 Not Found response.
func (a *applicationDependencies) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	a.errorResponseJSON(w, r, http.StatusNotFound, message)
}

// methodNotAllowedResponse sends a 405 Method Not Allowed response.
func (a *applicationDependencies) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	a.errorResponseJSON(w, r, http.StatusMethodNotAllowed, message)
}

// badRequestResponse sends a 400 Bad Request response.
func (a *applicationDependencies) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
}

// failedValidationResponse sends a 422 Unprocessable Entity response with validation errors.
func (a *applicationDependencies) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	a.errorResponseJSON(w, r, http.StatusUnprocessableEntity, errors)
}

// rateLimitExceededResponse sends a 429 Too Many Requests response when rate limit is exceeded.
func (a *applicationDependencies) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	a.errorResponseJSON(w, r, http.StatusTooManyRequests, message)
}
