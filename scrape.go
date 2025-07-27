package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Zigelzi/go-rss-gator/internal/database"
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

	lastFetchedAt := sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}

	s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:            nextFeed.ID,
		UpdatedAt:     time.Now().UTC(),
		LastFetchedAt: lastFetchedAt,
	})

	printFeedContent(feedContent)
}
