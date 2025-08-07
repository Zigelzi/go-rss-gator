package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Zigelzi/go-rss-gator/internal/database"
	"github.com/google/uuid"
)

func handleAggregate(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("feed URL is required argument (1), usage: 'agg TIME_BETWEEN_REQUESTS'")
	}
	if len(cmd.Args) > 1 {
		return fmt.Errorf("got over 1 argument (%d): %v,  usage: 'agg TIME_BETWEEN_REQUESTS'", len(cmd.Args), cmd.Args)
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("unable to parse time between requests argument: %w", err)
	}
	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeed(s)
	}
}

func handleAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("RSS feed name is required argument (1), usage: 'addfeed NAME URL'")
	}
	if len(cmd.Args) < 2 {
		return errors.New("RSS feed URL is required argument (2), usage: 'addfeed NAME URL'")
	}
	if len(cmd.Args) > 2 {
		return fmt.Errorf("got over 2 arguments (%d): %v, usage: 'addfeed NAME URL'", len(cmd.Args), cmd.Args)
	}

	// Add validation that URL starts with http(s) and ends in .xml to ensure correct format.
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
		printFeedDetails(feed)
	}
	fmt.Println(strings.Repeat("-", 10))
	return nil
}

func handleFollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("feed URL is required argument (1), usage: 'follow URL'")
	}
	if len(cmd.Args) > 1 {
		return fmt.Errorf("got over 1 argument (%d): %v,  usage: 'follow URL'", len(cmd.Args), cmd.Args)
	}
	feedURL := cmd.Args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("feed with url [%s] doesn't exist. Add it by using 'addfeed URL' command first", feedURL)
		}
		return fmt.Errorf("unable to get feed with URL [%s]: %w", feedURL, err)
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

	log.Printf("%s started to follow feed: %s - %s\n",
		followedFeed.UserName,
		followedFeed.FeedName,
		feedURL)
	return nil
}

func handleListFollowedFeeds(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetUserFeedFollows(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("unable to get feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("You are not following any feeds. Follow on by using 'addfeed URL' command")
		return nil
	}
	fmt.Println("The RSS feeds you're following:")
	for _, feed := range feeds {
		fmt.Println(feed.FeedName)
	}
	return nil
}

func handleUnfollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("feed URL is required argument (1), usage: 'unfollow URL'")
	}
	if len(cmd.Args) > 1 {
		return fmt.Errorf("got over 1 argument (%d): %v,  usage: 'unfollow URL'", len(cmd.Args), cmd.Args)
	}
	feedURL := cmd.Args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no feed with URL %s exists", feedURL)
		}
		return fmt.Errorf("unable to get feed: %w", err)
	}
	err = s.db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to unfollow feed with URL [%s]: %w", feedURL, err)
	}
	fmt.Printf("Successfully unfollowed feed %s - %s", feed.Name, feed.Url)
	return nil
}

func handleBrowsePosts(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("got over 1 argument (%d): %v,  usage: 'browse [LIMIT]'", len(cmd.Args), cmd.Args)
	}

	postLimit := 2
	var err error
	if len(cmd.Args) == 1 {
		postLimit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("unable to parse the limit argument: %w", err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(postLimit),
	})
	if err != nil {
		return fmt.Errorf("unable to get posts for user [%s]: %w", user.Name, err)
	}

	for _, post := range posts {
		hoursSincePublished := time.Since(post.PublishedAt).Hours()
		fmt.Println(post.Title)
		fmt.Println(post.Description.String)
		fmt.Printf("Published: %v (%.f h ago)\n", post.PublishedAt.Format(time.DateTime), hoursSincePublished)
		fmt.Printf("Read more: %s\n\n", post.Url)
	}

	return nil
}
