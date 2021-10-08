package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	StatusText string `json:"status_text"`
	Message    string `json:"message"`
}

func MarshalPayload(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func MarshalError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response := ErrorResponse{
		Message:    err.Error(),
		StatusText: http.StatusText(code),
	}
	json.NewEncoder(w).Encode(response)
}
