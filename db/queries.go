// db/queries.go
package db

import (
	"time"
)

// db/queries.go

// Feed represents a row in the feeds table, now with an UnreadCount
type Feed struct {
	ID          int
	URL         string
	UnreadCount int
}

// GetAllFeeds retrieves feeds and counts their unread articles
func GetAllFeeds() ([]Feed, error) {
	// This query joins the articles table and counts only where is_read = 0
	query := `
		SELECT f.id, f.url, COUNT(a.id) as unread_count
		FROM feeds f
		LEFT JOIN articles a ON f.id = a.feed_id AND a.is_read = 0
		GROUP BY f.id
	`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []Feed
	for rows.Next() {
		var f Feed
		if err := rows.Scan(&f.ID, &f.URL, &f.UnreadCount); err != nil {
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

// Article represents a row in the articles table
type Article struct {
	ID          int
	FeedID      int
	Title       string
	URL         string
	Description string
	PublishedAt time.Time
	IsRead      bool
}

// AddFeed inserts a new feed URL. 
// We temporarily use the URL as the title until the background parser updates it.
func AddFeed(url string) (int, error) {
	res, err := DB.Exec("INSERT OR IGNORE INTO feeds (title, url) VALUES (?, ?)", url, url)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

// GetArticlesByFeed fetches the latest 50 articles for a specific feed
func GetArticlesByFeed(feedID int) ([]Article, error) {
	rows, err := DB.Query(`
		SELECT id, title, url, description, published_at, is_read 
		FROM articles 
		WHERE feed_id = ? 
		ORDER BY published_at DESC LIMIT 50
	`, feedID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var a Article
		if err := rows.Scan(&a.ID, &a.Title, &a.URL, &a.Description, &a.PublishedAt, &a.IsRead); err != nil {
			continue
		}
		articles = append(articles, a)
	}
	return articles, nil
}

// MarkArticleAsRead sets the is_read flag to true for a specific article
func MarkArticleAsRead(articleID int) error {
	_, err := DB.Exec("UPDATE articles SET is_read = 1 WHERE id = ?", articleID)
	return err
}