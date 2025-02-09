package responses

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

// sendResponse is a helper function to send JSON responses
func SendResponse(w http.ResponseWriter, status string, message string, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If we fail to encode the response, we've already set the header status,
		// so we'll just log the error and return
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
