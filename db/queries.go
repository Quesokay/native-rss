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

// Article represents a row in the articles table
type Article struct {
	ID          int
	FeedID      int
	Title       string
	URL         string
	Description string
	Content		string
	PublishedAt time.Time
	IsRead      bool
	EnclosureURL string
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

func SaveArticle(feedID int, title, url, description, content string, publishedAt time.Time, enclosureURL string) error {
	query := `
		INSERT OR IGNORE INTO articles (feed_id, title, url, description, content, published_at, enclosure_url)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := DB.Exec(query, feedID, title, url, description, content, publishedAt, enclosureURL)
	return err
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
		SELECT id, title, url, description, content, published_at, is_read 
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
		if err := rows.Scan(&a.ID, &a.Title, &a.URL, &a.Description, &a.Content, &a.PublishedAt, &a.IsRead); err != nil {
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

// GetArticleByID fetches a single article for the reading view
func GetArticleByID(id int) (Article, error) {
	var a Article
	query := `
		SELECT id, feed_id, title, url, description, COALESCE(content, ''), published_at, is_read, COALESCE(enclosure_url, '')
		FROM articles 
		WHERE id = ?
	`
	err := DB.QueryRow(query, id).Scan(
		&a.ID, &a.FeedID, &a.Title, &a.URL, &a.Description, &a.Content, &a.PublishedAt, &a.IsRead, &a.EnclosureURL,
	)
	return a, err
}

// DeleteFeed removes a feed and (due to CASCADE) all its associated articles
func DeleteFeed(feedID int) error {
	_, err := DB.Exec("DELETE FROM feeds WHERE id = ?", feedID)
	return err
}
