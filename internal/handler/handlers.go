// internal/handler/handlers.go
package handler

import (
	"net/http"
	"strconv"
	"native-rss/db"
	"native-rss/internal/parser"
	"native-rss/ui"

	"github.com/a-h/templ"
)

// HandleIndex serves the main dashboard
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	feeds, _ := db.GetAllFeeds()
	templ.Handler(ui.Layout("Native RSS Reader", ui.AppShell(feeds))).ServeHTTP(w, r)
}

// HandleAddFeed processes the form submission from HTMX
func HandleAddFeed(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url != "" {
		feedID, err := db.AddFeed(url)
		if err == nil && feedID > 0 {
			// Kick off a background job to fetch this new feed immediately
			go parser.FetchAndSaveFeed(feedID, url)
		}
	}

	// Fetch the updated list of feeds
	feeds, _ := db.GetAllFeeds()
	
	// HTMX expects just the HTML snippet that changed, so we ONLY return the FeedList
	ui.FeedList(feeds).Render(r.Context(), w)
}

// HandleGetArticles fetches articles when a feed is clicked in the sidebar
func HandleGetArticles(w http.ResponseWriter, r *http.Request) {
	// Go 1.22 gives us PathValue to easily grab the {id} from the URL
	idStr := r.PathValue("id")
	feedID, _ := strconv.Atoi(idStr)

	articles, _ := db.GetArticlesByFeed(feedID)

	// Return ONLY the article list to swap into the right pane
	ui.ArticleList(articles).Render(r.Context(), w)
}


// HandleMarkRead updates the DB and returns a "Read" badge to swap into the UI
func HandleMarkRead(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	articleID, _ := strconv.Atoi(idStr)

	db.MarkArticleAsRead(articleID)

	// Instead of an empty response, return our new static badge!
	ui.ReadBadge().Render(r.Context(), w)
}