package handlers

import (
	"io"
	"log"
	"net/http"
	"strings"

	"io/fs"
)

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

	// Set content type and serve the file
	w.Header().Set("Content-Type", "text/html")
	if _, err := io.Copy(w, indexFile); err != nil {
		http.Error(w, "failed to serve index.html", http.StatusInternalServerError)
	}
}

func wantsAsset(path string) bool {
	return strings.HasPrefix(path, "/assets/") || strings.Contains(path, ".")
}
