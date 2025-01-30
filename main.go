package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"os"

	"github.com/DeanWard/erugo/config"
	"github.com/DeanWard/erugo/db"
	"github.com/DeanWard/erugo/handlers"
	"github.com/DeanWard/erugo/middleware"
	"github.com/DeanWard/erugo/setup"
	"github.com/DeanWard/erugo/utils"
	"github.com/gorilla/mux"
	_ "modernc.org/sqlite"
)

//go:embed frontend/dist/*
var staticFiles embed.FS

//go:embed migrations/*
var migrations embed.FS

func main() {

	// Load configuration
	configErr := config.LoadConfig("config.json")
	if configErr != nil {
		log.Printf("Error loading config: %v. Creating default config.", configErr)
		defaultConfig := utils.GetDefaultConfig()
		os.WriteFile("config.json", []byte(defaultConfig), 0644)
		config.LoadConfig("config.json")
	}

	//load migrations from embed.FS
	migrations, err := fs.Sub(migrations, "migrations")
	if err != nil {
		log.Fatalf("Failed to create embedded filesystem: %v", err)
	}

	// Connect to the database and run migrations
	database := db.Connect()
	db.Migrate(database, migrations)

	// Embed static files for the frontend
	embeddedFS, embeddedErr := fs.Sub(staticFiles, "frontend/dist")
	if embeddedErr != nil {
		log.Fatalf("Failed to create embedded filesystem: %v", embeddedErr)
	}

	runSetupIfNeeded(database)

	// Bring up the server
	bringUpServer(database, embeddedFS)
}

func runSetupIfNeeded(database *sql.DB) {
	//are there any users in the database?
	users := db.GetUsers(database)
	if len(users) == 0 {
		//run the setup
		log.Println("No users found, running setup")
		setup.RunSetup(database)
	} else {
		log.Println("Users found, skipping setup")
	}
}

func bringUpServer(database *sql.DB, embeddedFS fs.FS) {
	// Initialize the gorilla/mux router
	router := mux.NewRouter()
	router.StrictSlash(false)

	// Define API routes

	//POST /api/auth/login - login
	router.Handle("/api/auth/login/", handlers.LoginHandler(database)).Methods("POST")
	//POST /api/auth/refresh - refresh a token
	router.Handle("/api/auth/refresh/", handlers.RefreshTokenHandler(database)).Methods("POST")
	//POST /api/auth/logout - logout
	router.Handle("/api/auth/logout/", handlers.LogoutHandler(database)).Methods("POST")
	//POST /api/shares - create a new share
	router.Handle("/api/shares/", middleware.JwtMiddleware(handlers.CreateShareHandler(database))).Methods("POST")
	//GET /api/shares/{longId} - get a share by its longId
	router.Handle("/api/shares/{longId}/", handlers.GetShareHandler(database)).Methods("GET")
	//POST /api/shares/{longId}/download - download a share by its longId
	router.Handle("/api/shares/{longId}/download/", handlers.DownloadShareHandler(database)).Methods("GET")
	//GET /api/health - health check
	router.Handle("/api/health/", handlers.HealthCheckHandler()).Methods("GET")

	// Serve the frontend (static files)
	router.PathPrefix("/").Handler(handlers.ServeFrontendHandler(embeddedFS))

	// Apply CORS middleware
	handlerWithCORS := middleware.CorsMiddleware(router)

	// Start the server
	log.Println("Starting server on", fmt.Sprintf(":%d", config.AppConfig.BindPort))
	listenErr := http.ListenAndServe(fmt.Sprintf(":%d", config.AppConfig.BindPort), handlerWithCORS)
	if listenErr != nil {
		log.Fatalf("Error starting server: %v", listenErr)
	}
}
