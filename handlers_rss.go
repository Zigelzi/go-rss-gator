package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Zigelzi/go-rss-gator/internal/database"
	"github.com/Zigelzi/go-rss-gator/internal/rss"
	"github.com/google/uuid"
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

func handleAddFeed(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("RSS feed name is required argument (1)")
	}
	if len(cmd.Args) < 2 {
		return errors.New("RSS feed URL is required argument (2)")
	}
	if len(cmd.Args) > 2 {
		return fmt.Errorf("got over 2 arguments (%d): %v", len(cmd.Args), cmd.Args)
	}
	fmt.Printf("Adding RSS feed with name [%s] and URL [%s]\n", cmd.Args[0], cmd.Args[1])

	// Add validation that URL starts with http(s) and ends in .xml to ensure correct format.

	currentUserName := s.currentConfig.CurrentUserName
	user, err := s.db.GetUser(context.Background(), currentUserName)
	if err != nil {
		return fmt.Errorf("unable to get current user: %w", err)
	}

	feed, err := s.db.CreateFeed(context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      cmd.Args[0],
			Url:       cmd.Args[1],
			UserID:    user.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("unable to create new feed: %w", err)
	}
	fmt.Printf("Successfully added RSS feed with name [%s] and URL [%s] to user [%s]", feed.Name, feed.Url, user.Name)
	return nil
}

func handleListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found")
		return nil
	}

	fmt.Println("All RSS feeds")
	fmt.Println("[Username] RSS feed title - RSS feed URL")
	fmt.Println(strings.Repeat("-", 10))
	for _, feed := range feeds {
		printFeed(feed)
	}
	fmt.Println(strings.Repeat("-", 10))
	return nil
}

func printFeed(feed database.GetFeedsRow) {
	fmt.Printf("[%s] %s - %s\n", feed.UserName.String, feed.FeedName, feed.FeedUrl)
}
