package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/DeanWard/erugo/config"
	"github.com/golang-jwt/jwt/v5"
)

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("Attempting to authenticate request")
		token := r.Header.Get("Authorization")
		log.Printf("Token: %s", token)
		//strip the "Bearer " prefix
		token = strings.TrimPrefix(token, "Bearer ")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Println("No token provided")
			return
		}

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JwtSecret), nil
		})
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Printf("Failed to parse token: %v", err)
			return
		}

		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			r.Header.Set("X-User-ID", claims["sub"].(string))
			log.Printf("Authenticated user: %s", claims["sub"].(string))
		}

		next.ServeHTTP(w, r)
	})
}
