// internal/parser/parser.go
package parser

import (
	"log"
	"native-rss/db"
	"time"

	"github.com/mmcdole/gofeed"
)

// FetchAndSaveFeed downloads a single RSS feed and saves its articles
func FetchAndSaveFeed(feedID int, feedURL string) {
	log.Printf("Fetching feed: %s", feedURL)
	
	fp := gofeed.NewParser()
	parsedFeed, err := fp.ParseURL(feedURL)
	if err != nil {
		log.Printf("Failed to parse feed %s: %v", feedURL, err)
		return
	}

	newArticles := 0
	for _, item := range parsedFeed.Items {
		// Handle missing dates gracefully
		pubDate := time.Now()
		if item.PublishedParsed != nil {
			pubDate = *item.PublishedParsed
		}

		// Insert into the database (duplicates are ignored automatically)
		err := db.SaveArticle(feedID, item.Title, item.Link, item.Description, pubDate)
		if err != nil {
			log.Printf("Failed to save article %s: %v", item.Title, err)
		} else {
			newArticles++
		}
	}
	
	log.Printf("Finished %s - Found %d items", feedURL, newArticles)
}