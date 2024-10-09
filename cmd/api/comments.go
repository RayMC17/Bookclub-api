package main

import (
	//"encoding/json"
	"fmt"
	"net/http"
	//"github.com/RayMC17/comments/internal/data"
)

func (a *applicationDependencies)createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var incomingData struct {
	Content string `json:"content"`
	Author string `json:"author"`
	}

	//err := json.NewDecoder (r.Body) .Decode(&incomingData)
	err := a.readJSON(w, r, &incomingData)
if err != nil {
	a.badRequestResponse(w, r, err)
	//a.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
return
}
fmt.Fprintf(w, "%+v\n", incomingData)
}




