// cmd/server/main.go
package main

import (
	"log"
	"native-rss/db" // Make sure this matches your module name in go.mod
)

func main() {
	log.Println("Starting Native RSS server...")

	// Initialize the SQLite database
	// This will create a file named "rss.db" in your root directory
	err := db.InitDB("rss.db")
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	
	// We will add the HTTP server and background worker here in the next phases
	log.Println("Setup complete. Server is ready (shutting down for now).")
}