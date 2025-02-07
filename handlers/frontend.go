package handlers

import (
	"database/sql"
	"fmt"

	"github.com/DeanWard/erugo/db"

	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

// ServeFrontendHandler returns an http.Handler wrapping ServeFrontend with embeddedFS.
func ServeFrontendHandler(embeddedFS fs.FS, database *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeFrontend(w, r, embeddedFS, database)
	})
}

func ServeFrontend(w http.ResponseWriter, r *http.Request, embeddedFS fs.FS, database *sql.DB) {
	log.Printf("Serving frontend for %s", r.URL.Path)
	if wantsAsset(r.URL.Path) {
		http.FileServer(http.FS(embeddedFS)).ServeHTTP(w, r)
		return
	}

	cssVariables, err := getCssVariables(database)
	if err != nil {
		http.Error(w, "failed to get css variables", http.StatusInternalServerError)
		return
	}

	// Serve index.html for all other paths
	indexFile, err := embeddedFS.Open("index.html")
	if err != nil {
		http.Error(w, "index.html not found", http.StatusInternalServerError)
		// Log a list of all files in the embeddedFS
		files, err := fs.ReadDir(embeddedFS, ".")
		if err != nil {
			http.Error(w, "failed to read embeddedFS", http.StatusInternalServerError)
			return
		}
		log.Printf("Embedded files: %v", files)
		return
	}
	defer indexFile.Close()

	// Read the entire file into memory
	content, err := io.ReadAll(indexFile)
	if err != nil {
		http.Error(w, "failed to read index.html", http.StatusInternalServerError)
		return
	}

	// Convert content to string for manipulation
	htmlContent := string(content)

	// Inject CSS variables just before the closing head tag
	if idx := strings.Index(htmlContent, "</head>"); idx != -1 {
		htmlContent = fmt.Sprintf("%s%s%s",
			htmlContent[:idx], // Everything before </head>
			cssVariables,
			htmlContent[idx:], // From </head> to end
		) // Closing the fmt.Sprintf function call
	} else {
		// If no <head> tag is found, inject at the beginning of the file
		htmlContent = cssVariables + htmlContent
	}

	// Set content type and serve the modified file
	w.Header().Set("Content-Type", "text/html")
	if _, err := io.Copy(w, strings.NewReader(htmlContent)); err != nil {
		http.Error(w, "failed to serve index.html", http.StatusInternalServerError)
	}
}

func wantsAsset(path string) bool {
	return strings.HasPrefix(path, "/assets/") || strings.Contains(path, ".")
}

func getCssVariables(database *sql.DB) (string, error) {
	settings, err := db.SettingsByGroup(database, "ui.css")
	if err != nil {
		return "", err
	}

	// Create map for easier access
	settingsMap := make(map[string]string)
	for _, setting := range settings {
		settingsMap[setting.Id] = setting.Value
	}

	cssVariables := fmt.Sprintf(`
	<style>
	:root {
			--primary-color: %s;
			--secondary-color: %s;
			--accent-color: %s;
			--accent-color-light: %s;
	}
	</style>`,
		settingsMap["css_primary_color"],
		settingsMap["css_secondary_color"],
		settingsMap["css_accent_color"],
		settingsMap["css_accent_color_light"])

	return cssVariables, nil
}
