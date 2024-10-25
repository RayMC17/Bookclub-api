package main

import (
	//"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/RayMC17/comments/internal/data"
	"github.com/RayMC17/comments/internal/validator"
)

func (a *applicationDependencies) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var incomingData struct {
		Content string `json:"content"`
		Author  string `json:"author"`
	}

	//err := json.NewDecoder (r.Body) .Decode(&incomingData)
	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		//a.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
		return
	}

	comment := &data.Comment{
		Content: incomingData.Content,
		Author:  incomingData.Author,
	}

	v := validator.New()

	data.ValidateComment(v, comment)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}
	a.logger.Info("OKAY")
	err = a.commentModel.Insert(comment)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	//fmt.Fprintf(w, "%+v\n", incomingData)
	// Set a Location header. The path to the newly created comment
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/comments/%d", comment.ID))

	data := envelope{
		"comment": comment,
	}
	err = a.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

}
func (a *applicationDependencies) displayCommentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return

	}
	comment, err := a.commentModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}
	data := envelope{
		"comment": comment,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

}

func (a *applicationDependencies) updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id from the URL
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}
	// Call Get() to retrieve the comment with the specified id
	comment, err := a.commentModel.Get(id)
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
		Content *string `json:"content"`
		Author  *string `json:"author"`
	}
	err = a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	// We need to now check the fields to see which ones need updating
	// if incomingData.Content is nil, no update was provided
	if incomingData.Content != nil {
		comment.Content = *incomingData.Content
	}
	// if incomingData.Author is nil, no update was provided
	if incomingData.Author != nil {
		comment.Author = *incomingData.Author
	}

	// Before we write the updates to the DB let's validate
	v := validator.New()
	data.ValidateComment(v, comment)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = a.commentModel.Update(comment)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	data := envelope{
		"comment": comment,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *applicationDependencies) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}
	err = a.commentModel.Delete(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}
	// display the comment
	data := envelope{
		"message": "comment successfully deleted",
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}

}

func (a *applicationDependencies) listCommentsHandler(w http.ResponseWriter, r *http.Request) {

	var queryParametersData struct {
		Content string
		Author  string
	}
	queryParameters := r.URL.Query()

	// Load the query parameters into our struct
	queryParametersData.Content = a.getSingleQueryParameter(
		queryParameters,
		"content",
		"")

	queryParametersData.Author = a.getSingleQueryParameter(
		queryParameters,
		"author",
		"",
	)
	comments, err := a.commentModel.GetAll(
		queryParametersData.Content,
		queryParametersData.Author,
	)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	data := envelope{
		"comments": comments,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
