package main

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]any

func (a *applicationDependencies) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	jsResponse, err := json.MarshalIndent(data, "", "\t")

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
