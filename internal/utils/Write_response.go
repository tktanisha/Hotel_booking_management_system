package utils

import (
	"encoding/json"
	"net/http"
)

type WriteResponse struct {
	Status     bool        `json:"status"` //false-error, true-success
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func WriteSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(WriteResponse{
		Status:     true,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}
