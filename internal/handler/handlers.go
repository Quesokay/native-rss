// internal/handler/handlers.go
package handler

import (
	"net/http"
	"native-rss/db"
	"native-rss/ui"

	"github.com/a-h/templ"
)

// HandleIndex serves the main initial dashboard page
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	feeds, err := db.GetAllFeeds()
	if err != nil {
		http.Error(w, "Failed to load feeds", http.StatusInternalServerError)
		return
	}

	// Compose the layout with the AppShell component
	component := ui.Layout("Native RSS Reader", ui.AppShell(feeds))
	
	// Turn the template into an active HTTP response
	templ.Handler(component).ServeHTTP(w, r)
}