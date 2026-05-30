// cmd/server/main.go
package main

import (
	"log"
	"net/http"
	"native-rss/db"
	"native-rss/internal/handler"
	"native-rss/internal/parser"
	"time"
)

func main() {
	log.Println("Starting Native RSS server...")

	// 1. Initialize SQLite
	err := db.InitDB("rss.db")
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// 2. Start Background Worker
	parser.StartWorker(15 * time.Minute)

	// 3. Set Up Route Handlers
	mux := http.NewServeMux()
	
	// Serve compiled static Tailwind styles
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("ui/assets"))))
	
	// Main application view
	mux.HandleFunc("/", handler.HandleIndex)

	// 4. Start HTTP Server
	serverAddr := ":8080"
	log.Printf("Web UI available at http://localhost%s", serverAddr)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}