package main

import (
	// "fmt"
	// "net/http"
	"encoding/json"
	"net/http"
)

//	func (a *applicationDependencies) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
//		fmt.Fprintln(w, "status: available")
//		fmt.Fprintf(w, "environment: %s\n", a.config.environment)
//		fmt.Fprintf(w, "version: %s\n", appVersion)
//	}
func (a *applicationDependencies) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": a.config.environment,
		"version":     appVersion,
	}
	jsResponse, err := json.Marshal(data)
	if err != nil {
		a.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}
	jsResponse = append(jsResponse, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsResponse)

}
