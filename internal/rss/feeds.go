package rss

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Zigelzi/go-rss-gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	fmt.Printf("Fetching RSS feed from: %s\n\n", feedURL)
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}
	req.Header.Set("User-Agent", "rss-gator")

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request RSS feed from [%s]: %w", feedURL, err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read the response: %w", err)
	}

	var rssFeed RSSFeed
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return nil, fmt.Errorf("unable to parse the XML from response: %w", err)
	}
	cleanFeedContent(&rssFeed)
	return &rssFeed, nil
}

func cleanFeedContent(rssFeed *RSSFeed) *RSSFeed {
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i, item := range rssFeed.Channel.Items {

		rssFeed.Channel.Items[i].Title = html.UnescapeString(item.Title)
		rssFeed.Channel.Items[i].Description = html.UnescapeString(item.Description)
	}
	return rssFeed
}

// Save the posts in a RSS feed to the database.
func SaveFeedPosts(db *database.Queries, rssFeed *RSSFeed, feedId uuid.UUID) error {
	newPostCount := 0
	// Check if there are any posts in the RSS feed.
	// Log a feedback for developer if there aren't any posts.
	// Loop through all the posts and try to save them.
	// Continue to next post if saving fails. Optionally store which posts weren't saved.
	for _, post := range rssFeed.Channel.Items {
		desc := sql.NullString{
			String: post.Description,
			Valid:  true,
		}
		publishedAt, err := parseTimestamp(post.PublishDate)
		if err != nil {
			return err
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       post.Title,
			Description: desc,
			Url:         post.Link,
			PublishedAt: publishedAt,
			FeedID:      feedId,
		})
		if err != nil {
			var pqErr *pq.Error
			// Skip errors about duplicate values (Postgres error code 23505)
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				continue
			}
			log.Printf("unable to save post [%s]: %v", post.Title, err)
			continue
		}
		newPostCount++
	}
	log.Printf("Saved %d new posts (total %d)", newPostCount, len(rssFeed.Channel.Items))
	return nil
}

func parseTimestamp(timestampStr string) (time.Time, error) {
	formats := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC822,
		time.RFC822Z,
	}

	for _, format := range formats {
		if timestamp, err := time.Parse(format, timestampStr); err == nil {
			return timestamp, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse timestamp: %s", timestampStr)
}
