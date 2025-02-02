package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/DeanWard/erugo/auth"
	"github.com/DeanWard/erugo/config"
	"github.com/DeanWard/erugo/db"
	"github.com/DeanWard/erugo/responses/json_response"
	"github.com/golang-jwt/jwt/v5"
)

func LoginHandler(database *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleLogin(database, w, r)
	})
}

func handleLogin(database *sql.DB, w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//try to find the user in the database, then check if the password is correct using the auth.CheckPassword function
	user := db.UserByName(database, req.Username)
	if user == nil {
		json_response.New(json_response.ErrorStatus, "Invalid username or password", nil, http.StatusUnauthorized).Send(w)
		log.Printf("Failed to find user %s in the database", req.Username)
		return
	}

	if !auth.CheckPassword(user.Password, req.Password) {
		json_response.New(json_response.ErrorStatus, "Invalid username or password", nil, http.StatusUnauthorized).Send(w)
		log.Printf("Failed to verify password for user %s", req.Username)
		return
	}

	//create a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"admin": user.Admin,
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"admin": user.Admin,
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.JwtSecret))

	if err != nil {
		json_response.New(json_response.ErrorStatus, "Failed to create token", nil, http.StatusInternalServerError).Send(w)
		log.Printf("Failed to create token for user %s", req.Username)
		return
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(config.AppConfig.JwtSecret))
	if err != nil {
		json_response.New(json_response.ErrorStatus, "Failed to create refresh token", nil, http.StatusInternalServerError).Send(w)
		log.Printf("Failed to create refresh token for user %s", req.Username)
		return
	}

	//add an http only cookie to the response containing the refresh token
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshTokenString,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	json_response.New(json_response.SuccessStatus, "Login successful", map[string]string{"token": tokenString}, http.StatusOK).Send(w)

}

func RefreshTokenHandler(database *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleRefreshToken(database, w, r)
	})
}

func handleRefreshToken(database *sql.DB, w http.ResponseWriter, r *http.Request) {
	//read the refresh token from the cookie
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		json_response.New(json_response.ErrorStatus, "No refresh token provided", nil, http.StatusUnauthorized).Send(w)
		return
	}

	//parse the refresh token
	parsedToken, err := jwt.Parse(refreshToken.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JwtSecret), nil
	})

	if err != nil {
		json_response.New(json_response.ErrorStatus, "Failed to parse refresh token", nil, http.StatusUnauthorized).Send(w)
		log.Printf("Failed to parse refresh token: %v", err)
		return
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		user := db.UserByID(database, int(claims["sub"].(float64)))
		if user == nil {
			json_response.New(json_response.ErrorStatus, "User not found", nil, http.StatusUnauthorized).Send(w)
			return
		}

		//create a new jwt token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":   user.ID,
			"exp":   time.Now().Add(time.Minute * 15).Unix(),
			"admin": user.Admin,
		})

		tokenString, err := token.SignedString([]byte(config.AppConfig.JwtSecret))
		if err != nil {
			json_response.New(json_response.ErrorStatus, "Failed to create token", nil, http.StatusInternalServerError).Send(w)
			log.Printf("Failed to create token for user %s", user.Username)
			return
		}

		json_response.New(json_response.SuccessStatus, "Token refreshed successfully", map[string]string{"token": tokenString}, http.StatusOK).Send(w)

	}

}

func LogoutHandler(database *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleLogout(w)
	})
}

func handleLogout(w http.ResponseWriter) {
	//delete the refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "refresh_token",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	})

	json_response.New(json_response.SuccessStatus, "Logout successful", nil, http.StatusOK).Send(w)
}
