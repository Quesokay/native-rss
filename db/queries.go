// db/queries.go
package db

import (
	"log"
	"time"
)

// Feed represents a row in the feeds table
type Feed struct {
	ID  int
	URL string
}

// GetAllFeeds retrieves all subscribed RSS feeds from the database
func GetAllFeeds() ([]Feed, error) {
	rows, err := DB.Query("SELECT id, url FROM feeds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []Feed
	for rows.Next() {
		var f Feed
		if err := rows.Scan(&f.ID, &f.URL); err != nil {
			log.Printf("Error scanning feed row: %v", err)
			continue
		}
		feeds = append(feeds, f)
	}
	return feeds, nil
}

// SaveArticle inserts a new article. If the URL already exists, it ignores it.
func SaveArticle(feedID int, title, url, description string, publishedAt time.Time) error {
	query := `
		INSERT OR IGNORE INTO articles (feed_id, title, url, description, published_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := DB.Exec(query, feedID, title, url, description, publishedAt)
	return err
}