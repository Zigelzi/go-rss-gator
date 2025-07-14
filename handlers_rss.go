package main

import (
	"context"
	"fmt"

	"github.com/Zigelzi/go-rss-gator/internal/rss"
)

func handleAggregate(s *state, cmd command) error {

	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("unable to fetch feed from: %w", err)
	}
	fmt.Println(feed.Channel.Title)
	fmt.Println(feed.Channel.Description)
	for _, item := range feed.Channel.Items {
		fmt.Println(item.Title)
		fmt.Println(item.Description)
		fmt.Println(item.PublishDate)
		fmt.Println()
	}
	return nil
}
