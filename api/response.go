package api

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrRequestBody = errors.New("invalid request body")
)

type ErrorResponse struct {
	StatusText string `json:"status_text"`
	Message    string `json:"message"`
}

func Success(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func Error(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	rsp := ErrorResponse {
		Message:    err.Error(),
		StatusText: http.StatusText(code),
	}
	json.NewEncoder(w).Encode(rsp)
}
