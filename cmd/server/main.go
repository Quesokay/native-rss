// cmd/server/main.go
package main

import (
	"log"
	"native-rss/db"
	"native-rss/internal/parser"
	"time"
)

func main() {
	log.Println("Starting Native RSS server...")

	// 1. Initialize DB
	err := db.InitDB("rss.db")
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// [Optional for testing] Let's inject a test feed manually just so the parser has something to do!
	db.DB.Exec("INSERT OR IGNORE INTO feeds (title, url) VALUES ('TechCrunch', 'https://techcrunch.com/feed/')")

	// 2. Start Background Worker (runs every 15 minutes)
	// Because this runs in a goroutine, it won't block the rest of our app
	parser.StartWorker(15 * time.Minute)

	// 3. Keep the application running (temporary hack until we add the HTTP server)
	// We are just telling the main thread to sleep forever so it doesn't exit immediately
	select {} 
}