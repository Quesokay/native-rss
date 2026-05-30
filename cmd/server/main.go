// cmd/server/main.go
package main

import (
	"log"
	"net/http"
	"os"
	"native-rss/db"
	"native-rss/internal/handler"
	"native-rss/internal/parser"
	"time"
)

func main() {
	log.Println("Starting Native RSS server...")

	// Determine database path based on environment
	dbPath := "rss.db"
	if os.Getenv("ENV") == "production" {
		dbPath = "/app/data/rss.db"
	}

	err := db.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// 2. Start Background Worker
	parser.StartWorker(15 * time.Minute)

	// 3. Set Up Route Handlers
	mux := http.NewServeMux()
	
	// Serve compiled static Tailwind styles
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("ui/assets"))))
	
	// Page Routes
	mux.HandleFunc("GET /", handler.HandleIndex)

	// HTMX API Routes
	mux.HandleFunc("GET /feeds", handler.HandleGetFeedList)
	mux.HandleFunc("POST /feeds", handler.HandleAddFeed)
	mux.HandleFunc("GET /feeds/{id}/articles", handler.HandleGetArticles)
	mux.HandleFunc("POST /articles/{id}/read", handler.HandleMarkRead)

	// 4. Start HTTP Server
	serverAddr := ":8080"
	log.Printf("Web UI available at http://localhost%s", serverAddr)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}