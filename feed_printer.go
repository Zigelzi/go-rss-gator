package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/Zigelzi/go-rss-gator/internal/rss"
)

func printFeedContent(rssFeed *rss.RSSFeed) {
	fmt.Println(rssFeed.Channel.Title)
	fmt.Println(rssFeed.Channel.Description)
	fmt.Println()
	for _, item := range rssFeed.Channel.Items[:5] {
		fmt.Println(item.Title)
		fmt.Println(item.Description)
		fmt.Printf("Read more: %s\n", item.Link)
		pubTimestamp, err := parseTimestamp(item.PublishDate)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Published %.0f h ago\n",
			time.Since(pubTimestamp).Hours())
		fmt.Println(strings.Repeat("-", 30))
		fmt.Println()
	}
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
