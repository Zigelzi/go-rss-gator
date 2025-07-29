package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Zigelzi/go-rss-gator/internal/rss"
)

func scrapeFeed(s *state) {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No feeds exist")
			return
		}
		log.Fatalf("unable to get next feed to fetch: %v", err)
		return
	}

	feedContent, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		log.Fatalf("unable to fetch feed from: %v", err)
		return
	}

	s.db.MarkFeedFetched(context.Background(), nextFeed.ID)

	printFeedContent(feedContent)
	log.Printf("Scraped feed %s with %d posts", nextFeed.Name, len(feedContent.Channel.Items))
}
