// internal/parser/parser.go
package parser

import (
	"log"
	"net/http"
	"time"
	"native-rss/db"

	"github.com/go-shiori/go-readability"
	"github.com/mmcdole/gofeed"
)

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
		pubDate := time.Now()
		if item.PublishedParsed != nil {
			pubDate = *item.PublishedParsed
		}

		// --- UPGRADED READABILITY LOGIC ---
		var fullContent string
		
		client := &http.Client{Timeout: 15 * time.Second}
		req, err := http.NewRequest("GET", item.Link, nil)
		
		if err == nil {
			// Disguise our app as a standard Chrome browser to bypass basic bot-blockers
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
			
			resp, err := client.Do(req)
			if err == nil {
				article, err := readability.FromReader(resp.Body, req.URL)
				if err != nil {
					log.Printf("Readability parsing failed for %s: %v", item.Link, err)
					fullContent = item.Description
				} else {
					fullContent = article.Content
				}
				resp.Body.Close()
			} else {
				log.Printf("HTTP fetch failed for %s: %v", item.Link, err)
				fullContent = item.Description
			}
		} else {
			fullContent = item.Description
		}

		// Save everything to the database
		err = db.SaveArticle(feedID, item.Title, item.Link, item.Description, fullContent, pubDate)
		if err == nil {
			newArticles++
		} else {
			// ADD THIS LINE so you know exactly why it fails in the future!
			log.Printf("Database error saving '%s': %v", item.Title, err)
		}
	}
	
	log.Printf("Finished %s - Found %d new items", feedURL, newArticles)
}