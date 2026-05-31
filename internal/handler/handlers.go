// internal/handler/handlers.go
package handler

import (
	"net/http"
	"strconv"
	"native-rss/db"
	"native-rss/internal/parser"
	"native-rss/ui"

	"github.com/a-h/templ"
	"io"
	"log"
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

// HandleGetFeedList returns only the sidebar feed list for auto-refreshing
func HandleGetFeedList(w http.ResponseWriter, r *http.Request) {
	feeds, _ := db.GetAllFeeds()
	ui.FeedList(feeds).Render(r.Context(), w)
}

// HandleGetSingleArticle fetches and renders the full reading view
func HandleGetSingleArticle(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	articleID, _ := strconv.Atoi(idStr)

	article, err := db.GetArticleByID(articleID)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	// Automatically mark it as read when opened!
	db.MarkArticleAsRead(articleID)

	ui.FullArticle(article).Render(r.Context(), w)
}


func HandleDeleteFeed(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	feedID, _ := strconv.Atoi(idStr)

	db.DeleteFeed(feedID)

	// THESE TWO LINES ARE CRITICAL
	// If they are missing, your sidebar deletes itself from the page!
	feeds, _ := db.GetAllFeeds()
	ui.FeedList(feeds).Render(r.Context(), w)
}

// HandleImportOPML accepts an uploaded XML file and bulk-adds the feeds
func HandleImportOPML(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, fileHeader, err := r.FormFile("opml_file")
	if err != nil {
		log.Printf("Upload Error: %v", err)
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, _ := io.ReadAll(file)
	log.Printf("Received file: %s (%d bytes)", fileHeader.Filename, len(data))

	// Extract the URLs
	urls := parser.ExtractFeeds(data)
	log.Printf("OPML Parser found %d URLs inside the file", len(urls))

	// Loop through and save them, logging ANY database errors
	successCount := 0
	for _, u := range urls {
		// 1. Insert the feed into the database
		res, err := db.DB.Exec("INSERT OR IGNORE INTO feeds (title, url) VALUES (?, ?)", u, u)
		if err != nil {
			log.Printf("Database rejected feed '%s': %v", u, err)
			continue
		}

		// Check if it was actually a new insert (not ignored)
		rows, _ := res.RowsAffected()
		if rows > 0 {
			successCount++
			
			// 2. GET THE INSERTED FEED ID
			var feedID int
			err = db.DB.QueryRow("SELECT id FROM feeds WHERE url = ?", u).Scan(&feedID)
			
			if err == nil {
				// 3. BLAST OFF A BACKGROUND WORKER INSTANTLY!
				// This fetches all articles in parallel without making the browser wait
				go func(id int, urlStr string) {
					log.Printf("🚀 Instant sync triggered for new feed [%d]: %s", id, urlStr)
					// Call your existing parser function here!
					parser.FetchAndSaveFeed(id, urlStr)
				}(feedID, u)
			}
		}
	}
	log.Printf("Successfully imported %d new feeds to the database", successCount)

	// Fetch updated list and render
	feeds, _ := db.GetAllFeeds()
	ui.FeedList(feeds).Render(r.Context(), w)
}