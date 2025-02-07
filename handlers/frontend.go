package handlers

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

// CSS variables to inject
const cssVariables = `
<style>
:root {
  --primary-color: rgb(238, 193, 84);
  --secondary-color: rgb(34, 34, 34);
  --accent-color: rgb(84, 129, 238);
  --accent-color-light: rgb(238, 238, 238);
}
</style>
`

// ServeFrontendHandler returns an http.Handler wrapping ServeFrontend with embeddedFS.
func ServeFrontendHandler(embeddedFS fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeFrontend(w, r, embeddedFS)
	})
}

func ServeFrontend(w http.ResponseWriter, r *http.Request, embeddedFS fs.FS) {
	log.Printf("Serving frontend for %s", r.URL.Path)
	if wantsAsset(r.URL.Path) {
		http.FileServer(http.FS(embeddedFS)).ServeHTTP(w, r)
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
