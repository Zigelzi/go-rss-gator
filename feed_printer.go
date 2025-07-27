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
		pubTimestamp, err := time.Parse(time.RFC1123, item.PublishDate)
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
