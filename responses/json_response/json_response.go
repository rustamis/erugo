package json_response

import (
	"encoding/json"
	"net/http"
)

// StatusType defines allowed status values
type StatusType string

const (
	SuccessStatus StatusType = "success"
	ErrorStatus   StatusType = "error"
)

// JSONResponse defines the structure of the JSON response
type JSONResponse struct {
	StatusMessage StatusType  `json:"status_message"`
	Message       string      `json:"message,omitempty"`
	Data          interface{} `json:"data,omitempty"`
	StatusCode    int         `json:"status_code"`
}

func New(status StatusType, message string, data interface{}, statusCode int) *JSONResponse {
	return &JSONResponse{
		StatusMessage: status,
		Message:       message,
		Data:          data,
		StatusCode:    statusCode,
	}
}

func (r *JSONResponse) Send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.StatusCode)
	json.NewEncoder(w).Encode(r)
}
