package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/DeanWard/erugo/config"
	"github.com/DeanWard/erugo/responses"
	"github.com/DeanWard/erugo/utils"
	"github.com/golang-jwt/jwt/v5"
)

// Define a custom type for context keys
type ContextKey string

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("Attempting to authenticate request")
		token := r.Header.Get("Authorization")
		log.Printf("Token: %s", token)
		//strip the "Bearer " prefix
		token = strings.TrimPrefix(token, "Bearer ")
		if token == "" {
			responses.SendResponse(w, responses.StatusError, "Unauthorized", nil, http.StatusUnauthorized)
			log.Println("No token provided")
			return
		}

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JwtSecret), nil
		})
		if err != nil {
			responses.SendResponse(w, responses.StatusError, "Unauthorized", nil, http.StatusUnauthorized)
			log.Printf("Failed to parse token: %v", err)
			return
		}

		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			userID, ok := claims["sub"].(float64)
			if !ok {
				responses.SendResponse(w, responses.StatusError, "Unauthorized", nil, http.StatusUnauthorized)
				log.Println("User ID not found in token claims")
				return
			}

			// Convert float64 to int
			userIDInt := int(userID)

			// Set header if needed (converting to string only for header)
			r.Header.Set("X-User-ID", strconv.Itoa(userIDInt))

			// Store claims and userID (as int) in context
			ctx := r.Context()
			ctx = context.WithValue(ctx, ContextKey("claims"), claims)
			ctx = context.WithValue(ctx, ContextKey("userID"), userIDInt)
			r = r.WithContext(ctx)

			utils.Log(fmt.Sprintf("User ID: %d", userIDInt), utils.ColorGreen)

			//check if the user must change their password
			mustChangePassword, ok := claims["must_change_password"].(bool)
			if ok && mustChangePassword {
				payload := map[string]interface{}{"message": "Password change required"}
				w.Header().Set("X-Password-Change-Required", "true")
				responses.SendResponse(w, responses.StatusError, "Password change required", payload, http.StatusForbidden)
				return
			}

		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func JwtMiddlewareNoReset(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("Attempting to authenticate request")
		token := r.Header.Get("Authorization")
		log.Printf("Token: %s", token)
		//strip the "Bearer " prefix
		token = strings.TrimPrefix(token, "Bearer ")
		if token == "" {
			responses.SendResponse(w, responses.StatusError, "Unauthorized", nil, http.StatusUnauthorized)
			log.Println("No token provided")
			return
		}

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JwtSecret), nil
		})
		if err != nil {
			responses.SendResponse(w, responses.StatusError, "Unauthorized", nil, http.StatusUnauthorized)
			log.Printf("Failed to parse token: %v", err)
			return
		}

		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			userID, ok := claims["sub"].(float64)
			if !ok {
				responses.SendResponse(w, responses.StatusError, "Unauthorized", nil, http.StatusUnauthorized)
				log.Println("User ID not found in token claims")
				return
			}

			// Convert float64 to int
			userIDInt := int(userID)

			// Set header if needed (converting to string only for header)
			r.Header.Set("X-User-ID", strconv.Itoa(userIDInt))

			// Store claims and userID (as int) in context
			ctx := r.Context()
			ctx = context.WithValue(ctx, ContextKey("claims"), claims)
			ctx = context.WithValue(ctx, ContextKey("userID"), userIDInt)
			r = r.WithContext(ctx)

			utils.Log(fmt.Sprintf("User ID: %d", userIDInt), utils.ColorGreen)

		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if claims, ok := r.Context().Value(ContextKey("claims")).(jwt.MapClaims); ok {
			isAdmin, ok := claims["admin"].(bool)
			if !ok || !isAdmin {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	})
}
