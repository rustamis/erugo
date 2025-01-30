package handlers

import (
	"fmt"
	"net/http"

	"github.com/DeanWard/erugo/config"
	"github.com/DeanWard/erugo/responses/json_response"
)

func HealthCheckHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		healthCheck(w)
	})
}

func healthCheck(w http.ResponseWriter) {
	payload := map[string]string{
		"status":         "ok",
		"max_share_size": fmt.Sprintf("%d", config.GetMaxShareSize()),
	}
	json_response.New(json_response.SuccessStatus, "OK", payload, http.StatusOK).Send(w)
}
