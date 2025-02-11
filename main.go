package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"

	"os"

	"github.com/DeanWard/erugo/config"
	"github.com/DeanWard/erugo/db"
	"github.com/DeanWard/erugo/middleware"
	"github.com/DeanWard/erugo/routes"
	"github.com/DeanWard/erugo/utils"
	"github.com/gorilla/mux"
	_ "modernc.org/sqlite"
)

//go:embed frontend/dist/*
var staticFiles embed.FS

//go:embed migrations/*
var migrations embed.FS

//go:embed swagger.html
var docsHTML string

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

	createLogoFile(embeddedFS)

	// Bring up the server
	bringUpServer(database, embeddedFS)
}

func createLogoFile(embeddedFS fs.FS) {
	//check if the logo file and private directory exist
	if _, err := os.Stat(config.AppConfig.PrivateDataPath + "/logo.png"); os.IsNotExist(err) {
		//create the private directory if it doesn't exist
		os.MkdirAll(config.AppConfig.PrivateDataPath, 0755)

		//the logo file is in the embeddedFS at frontend/dist/erugo.png
		logoFile, err := embeddedFS.Open("erugo.png")
		if err != nil {
			log.Fatalf("Failed to open logo file: %v", err)
		}
		defer logoFile.Close()

		//create the logo file
		logoFile, err = os.Create(config.AppConfig.PrivateDataPath + "/logo.png")
		if err != nil {
			log.Fatalf("Failed to create logo file: %v", err)
		}
		defer logoFile.Close()
		//copy the logo file to the private directory
		srcFile, err := embeddedFS.Open("erugo.png")
		if err != nil {
			log.Fatalf("Failed to open source logo file: %v", err)
		}
		defer srcFile.Close()

		destFile, err := os.Create(config.AppConfig.PrivateDataPath + "/logo.png")
		if err != nil {
			log.Fatalf("Failed to create destination logo file: %v", err)
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			log.Fatalf("Failed to copy logo file: %v", err)
		}
		log.Println("Logo file created")
	} else {
		log.Println("Logo file already exists")
	}
}

func bringUpServer(database *sql.DB, embeddedFS fs.FS) {
	// Initialize the gorilla/mux router
	router := mux.NewRouter()
	router.StrictSlash(false)

	router.Handle("/api-docs", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(docsHTML))
	}))

	// Register all routes
	routes.RegisterRoutes(router, database, embeddedFS)

	// Apply CORS middleware
	handlerWithCORS := middleware.CorsMiddleware(router)

	// Start the server
	log.Println("Starting server on", fmt.Sprintf(":%d", config.AppConfig.BindPort))
	listenErr := http.ListenAndServe(fmt.Sprintf(":%d", config.AppConfig.BindPort), handlerWithCORS)
	if listenErr != nil {
		log.Fatalf("Error starting server: %v", listenErr)
	}
}
