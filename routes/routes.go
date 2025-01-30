package routes

import (
	"database/sql"
	"io/fs"

	"github.com/DeanWard/erugo/handlers"
	"github.com/DeanWard/erugo/middleware"
	"github.com/gorilla/mux"
)

// RegisterRoutes sets up all routes for the application
func RegisterRoutes(router *mux.Router, database *sql.DB, embeddedFS fs.FS) {
	// Auth routes
	registerAuthRoutes(router, database)

	// Share routes
	registerShareRoutes(router, database)

	// User routes
	registerUserRoutes(router, database)

	// Health check
	registerHealthRoutes(router)

	// Frontend routes
	registerFrontendRoutes(router, embeddedFS)
}

func registerAuthRoutes(router *mux.Router, database *sql.DB) {
	router.Handle("/api/auth/login",
		handlers.LoginHandler(database),
	).Methods("POST")

	router.Handle("/api/auth/refresh",
		handlers.RefreshTokenHandler(database),
	).Methods("POST")

	router.Handle("/api/auth/logout",
		middleware.JwtMiddleware(
			handlers.LogoutHandler(database),
		),
	).Methods("POST")
}

func registerShareRoutes(router *mux.Router, database *sql.DB) {
	router.Handle("/api/shares",
		middleware.JwtMiddleware(
			handlers.CreateShareHandler(database),
		),
	).Methods("POST")

	router.Handle("/api/shares/{longId}",
		handlers.GetShareHandler(database),
	).Methods("GET")

	router.Handle("/api/shares/{longId}/download",
		handlers.DownloadShareHandler(database),
	).Methods("GET")
}

func registerUserRoutes(router *mux.Router, database *sql.DB) {
	router.Handle("/api/users",
		middleware.JwtMiddleware(
			middleware.AdminMiddleware(
				handlers.GetUsersHandler(database),
			),
		),
	).Methods("GET")
}

func registerHealthRoutes(router *mux.Router) {
	router.Handle("/api/health",
		handlers.HealthCheckHandler(),
	).Methods("GET")
}

func registerFrontendRoutes(router *mux.Router, embeddedFS fs.FS) {
	router.PathPrefix("/").Handler(handlers.ServeFrontendHandler(embeddedFS))
}
