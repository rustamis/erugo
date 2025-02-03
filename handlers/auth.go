package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/DeanWard/erugo/auth"
	"github.com/DeanWard/erugo/config"
	"github.com/DeanWard/erugo/db"
	"github.com/golang-jwt/jwt/v5"
)

// TokenType distinguishes between access and refresh tokens
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// AuthRequest represents the login credentials
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TokenResponse represents the successful auth response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"` // seconds until expiration
}

// Claims extends jwt.MapClaims to add type safety
type Claims struct {
	UserID    int       `json:"sub"`
	TokenType TokenType `json:"type"`
	IsAdmin   bool      `json:"admin"`
	jwt.RegisteredClaims
}

// createToken generates a new JWT token with appropriate claims
func createToken(userID int, isAdmin bool, tokenType TokenType) (*jwt.Token, error) {
	var expirationTime time.Time
	switch tokenType {
	case AccessToken:
		expirationTime = time.Now().Add(15 * time.Minute)
	case RefreshToken:
		expirationTime = time.Now().Add(24 * time.Hour)
	default:
		return nil, fmt.Errorf("invalid token type: %s", tokenType)
	}

	claims := Claims{
		UserID:    userID,
		TokenType: tokenType,
		IsAdmin:   isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims), nil
}

// setRefreshTokenCookie adds a secure HTTP-only cookie with the refresh token
func setRefreshTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400, // 24 hours in seconds
	})
}

// LoginHandler authenticates users and issues tokens
func LoginHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendResponse(w, StatusError, "Invalid request format", nil, http.StatusBadRequest)
			return
		}

		// Validate input
		if req.Username == "" || req.Password == "" {
			sendResponse(w, StatusError, "Username and password are required", nil, http.StatusBadRequest)
			return
		}

		// Get user and verify password
		user, err := db.UserByName(database, req.Username)
		if err != nil || !auth.CheckPassword(user.PasswordHash, req.Password) {
			// Use a generic error message to avoid user enumeration
			sendResponse(w, StatusError, "Invalid credentials", nil, http.StatusUnauthorized)
			return
		}

		// Generate access token
		accessToken, err := createToken(user.ID, user.Admin, AccessToken)
		if err != nil {
			sendResponse(w, StatusError, "Authentication failed", nil, http.StatusInternalServerError)
			return
		}

		accessTokenString, err := accessToken.SignedString([]byte(config.AppConfig.JwtSecret))
		if err != nil {
			sendResponse(w, StatusError, "Authentication failed", nil, http.StatusInternalServerError)
			return
		}

		// Generate refresh token
		refreshToken, err := createToken(user.ID, user.Admin, RefreshToken)
		if err != nil {
			sendResponse(w, StatusError, "Authentication failed", nil, http.StatusInternalServerError)
			return
		}

		refreshTokenString, err := refreshToken.SignedString([]byte(config.AppConfig.JwtSecret))
		if err != nil {
			sendResponse(w, StatusError, "Authentication failed", nil, http.StatusInternalServerError)
			return
		}

		// Set refresh token in HTTP-only cookie
		setRefreshTokenCookie(w, refreshTokenString)

		// Return access token in response body
		response := TokenResponse{
			AccessToken: accessTokenString,
			TokenType:   "Bearer",
			ExpiresIn:   15 * 60, // 15 minutes in seconds
		}

		sendResponse(w, StatusSuccess, "Login successful", response, http.StatusOK)
	}
}

// RefreshTokenHandler issues new access tokens using refresh tokens
func RefreshTokenHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get refresh token from cookie
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			sendResponse(w, StatusError, "No refresh token provided", nil, http.StatusUnauthorized)
			return
		}

		// Parse and validate the refresh token
		token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.AppConfig.JwtSecret), nil
		})

		if err != nil {
			sendResponse(w, StatusError, "Invalid refresh token", nil, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid || claims.TokenType != RefreshToken {
			sendResponse(w, StatusError, "Invalid refresh token", nil, http.StatusUnauthorized)
			return
		}

		// Get current user information
		user, err := db.UserByID(database, claims.UserID)
		if err != nil {
			if errors.Is(err, db.ErrUserNotFound) {
				sendResponse(w, StatusError, "User not found", nil, http.StatusUnauthorized)
				return
			}
			sendResponse(w, StatusError, "Authentication failed", nil, http.StatusInternalServerError)
			return
		}

		// Generate new access token
		newAccessToken, err := createToken(user.ID, user.Admin, AccessToken)
		if err != nil {
			sendResponse(w, StatusError, "Failed to create access token", nil, http.StatusInternalServerError)
			return
		}

		accessTokenString, err := newAccessToken.SignedString([]byte(config.AppConfig.JwtSecret))
		if err != nil {
			sendResponse(w, StatusError, "Failed to create access token", nil, http.StatusInternalServerError)
			return
		}

		response := TokenResponse{
			AccessToken: accessTokenString,
			TokenType:   "Bearer",
			ExpiresIn:   15 * 60,
		}

		sendResponse(w, StatusSuccess, "Token refreshed successfully", response, http.StatusOK)
	}
}

// LogoutHandler invalidates the refresh token
func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Clear the refresh token cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   -1,
		})

		sendResponse(w, StatusSuccess, "Logout successful", nil, http.StatusOK)
	}
}
