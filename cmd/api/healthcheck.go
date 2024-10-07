package main

import (
	// "fmt"
	// "net/http"
	// "encoding/json"
	"net/http"
)

//	func (a *applicationDependencies) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
//		fmt.Fprintln(w, "status: available")
//		fmt.Fprintf(w, "environment: %s\n", a.config.environment)
//		fmt.Fprintf(w, "version: %s\n", appVersion)
//	}
func (a *applicationDependencies) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	//panic("Apples & Oranges") 
	data := envelope {
		"status": "available",
		"system_info": map[string]string{
			"environment": a.config.environment,
			"version":     appVersion,
		},
	}
	//jsResponse, err := json.Marshal(data)
	err := a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		// a.logger.Error(err.Error())
		// http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		// return
		a.serverErrorResponse(w, r, err)
	}
	// jsResponse = append(jsResponse, '\n')
	// w.Header().Set("Content-Type", "application/json")
	// w.Write(jsResponse)

}
