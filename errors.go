package main

import (
	"net/http"
)

// APIError is an error from the api
type APIError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// Render implements the go-chi render interface
func (e *APIError) Render(w http.ResponseWriter, req *http.Request) error {
	if e.Status != 0 {
		w.WriteHeader(e.Status)
	}
	return nil
}

// ErrNotFound returns not found
func ErrNotFound(err error) *APIError {
	if err != nil {
		return &APIError{Message: err.Error(), Status: 404}
	}
	return &APIError{Message: "Not Found", Status: 404}
}
