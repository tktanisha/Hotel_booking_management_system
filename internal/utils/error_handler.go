package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Status     bool   `json:"status"` //false-error, true-success
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
	Details    string `json:"details,omitempty"`
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, errMsg string, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(ErrorResponse{
		Status:     false,
		StatusCode: statusCode,
		Error:      errMsg,
		Details:    details,
	})
}
