package main

import (
	"encoding/json"
	"net/http"
)

func (a *applicationDependencies) hwriteJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	jsResponse, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	jsResponse = append(jsResponse, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	_, err = w.Write(jsResponse)
	if err != nil {
		return err
	}
	return nil
}
