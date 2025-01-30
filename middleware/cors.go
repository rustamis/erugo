package middleware

import (
	"net/http"

	"github.com/DeanWard/erugo/config"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		appUrl := config.AppConfig.AppUrl
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", appUrl)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Pass to the next handler
		next.ServeHTTP(w, r)
	})
}
