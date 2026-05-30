// internal/parser/worker.go
package parser

import (
	"log"
	"native-rss/db"
	"sync"
	"time"
)

// StartWorker begins the background polling process
func StartWorker(interval time.Duration) {
	log.Printf("Starting background RSS worker (Interval: %v)", interval)
	
	// Run immediately on startup
	fetchAlFeeds()

	// Then run on the interval ticker
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			fetchAlFeeds()
		}
	}()
}

func fetchAlFeeds() {
	feeds, err := db.GetAllFeeds()
	if err != nil {
		log.Printf("Worker error fetching feeds from DB: %v", err)
		return
	}

	// Use a WaitGroup to fetch feeds concurrently
	// This makes fetching 100 feeds take the same time as fetching 1
	var wg sync.WaitGroup
	for _, feed := range feeds {
		wg.Add(1)
		go func(f db.Feed) {
			defer wg.Done()
			FetchAndSaveFeed(f.ID, f.URL)
		}(feed)
	}
	
	wg.Wait()
	log.Println("Background sync complete.")
}