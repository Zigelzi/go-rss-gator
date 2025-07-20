package main

import (
	"context"
	"database/sql"
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

	feedURL := cmd.Args[1]
	feed, err := s.db.CreateFeed(context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      cmd.Args[0],
			Url:       feedURL,
			UserID:    user.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("unable to create new feed: %w", err)
	}
	fmt.Printf("Successfully added RSS feed with name [%s] and URL [%s] to user [%s]", feed.Name, feed.Url, user.Name)

	_, err = s.db.FollowFeed(context.Background(), database.FollowFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to follow feed [%s]: %w", feedURL, err)
	}
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

func handleFollowFeed(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("feed URL is required argument")
	}
	if len(cmd.Args) > 1 {
		return fmt.Errorf("got over 1 argument (%d): %v", len(cmd.Args), cmd.Args)
	}
	feedURL := cmd.Args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Feed with url [%s] doesn't exist. Add it by using 'addfeed [url]' command first", feedURL)
			return nil
		}
		return fmt.Errorf("unable to get feed with URL [%s]: %w", feedURL, err)
	}
	user, err := s.db.GetUser(context.Background(), s.currentConfig.CurrentUserName)
	if err != nil {
		return fmt.Errorf("unable to get user [%s]: %w", s.currentConfig.CurrentUserName, err)
	}

	followedFeed, err := s.db.FollowFeed(context.Background(), database.FollowFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to follow feed [%s]: %w", feedURL, err)
	}

	fmt.Printf("%s started to follow feed: %s - %s\n",
		followedFeed.UserName,
		followedFeed.FeedName,
		feedURL)
	return nil
}

func handleListFollowedFeeds(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.currentConfig.CurrentUserName)
	if err != nil {
		return fmt.Errorf("unable to get user: %w", err)
	}
	feeds, err := s.db.GetUserFeedFollows(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("unable to get feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("You are not following any feeds. Follow on by using 'addfeed [url]' command")
		return nil
	}
	fmt.Println("The RSS feeds you're following:")
	for _, feed := range feeds {
		fmt.Println(feed.FeedName)
	}
	return nil
}
