package handlers

import (
	"bytes"
	"database/sql"
	"fmt"

	"golang.org/x/net/html"

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

	userCount, err := db.UserCount(database)
	if err != nil {
		http.Error(w, "failed to get user count", http.StatusInternalServerError)
		return
	}
	isSetupNeeded := userCount == 0

	if wantsAsset(r.URL.Path) {
		http.FileServer(http.FS(embeddedFS)).ServeHTTP(w, r)
		return
	}

	cssVariables, err := getCssVariables(database)
	if err != nil {
		http.Error(w, "failed to get css variables", http.StatusInternalServerError)
		return
	}

	applicationName, err := getApplicationName(database)
	if err != nil {
		http.Error(w, "failed to get application name", http.StatusInternalServerError)
		return
	}

	htmlDataAtributeSettings, err := db.SettingsByIds(database, []string{"logo_width", "application_name", "css_primary_color", "css_secondary_color", "css_accent_color", "css_accent_color_light"})
	if err != nil {
		log.Printf("failed to get settings: %v", err)
		http.Error(w, "failed to get settings", http.StatusInternalServerError)
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

	if idx := strings.Index(htmlContent, "<title>"); idx != -1 {
		endIdx := strings.Index(htmlContent, "</title>")
		if endIdx != -1 {
			htmlContent = fmt.Sprintf("%s<title>%s</title>%s",
				htmlContent[:idx],      // Everything before <title>
				applicationName,        // New title
				htmlContent[endIdx+8:], // Everything after </title>
			)
		}
	}

	htmlDataAttributeSettings := []Setting{}
	for _, setting := range htmlDataAtributeSettings {
		htmlDataAttributeSettings = append(htmlDataAttributeSettings, Setting{Id: setting.Id, Value: setting.Value})
	}

	if isSetupNeeded {
		htmlDataAttributeSettings = append(htmlDataAttributeSettings, Setting{Id: "setup_needed", Value: "true"})
	}
	htmlContent = injectSettings(htmlContent, htmlDataAttributeSettings)

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
	<style id="erugo-css-variables">
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

func getApplicationName(database *sql.DB) (string, error) {
	setting, err := db.SettingById(database, "application_name")
	if err != nil {
		return "", err
	}
	return setting.Value, nil
}

func injectSettings(htmlContent string, htmlDataAttributeSettings []Setting) string {
	// Find elements with data-settings="" and process each one
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return htmlContent
	}

	var processNode func(*html.Node)
	processNode = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// Check if node has data-settings attribute
			for i, attr := range n.Attr {
				if attr.Key == "data-settings" && attr.Val == "" {
					// Remove the data-settings attribute
					n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)

					// Add individual setting attributes
					for _, setting := range htmlDataAttributeSettings {
						n.Attr = append(n.Attr, html.Attribute{
							Key: fmt.Sprintf("data-setting-%s", setting.Id),
							Val: setting.Value,
						})
					}
					break
				}
			}
		}

		// Process child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processNode(c)
		}
	}

	processNode(doc)

	// Convert back to string
	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return htmlContent
	}

	return buf.String()
}

// Setting represents a single data attribute setting
type Setting struct {
	Id    string
	Value string
}
