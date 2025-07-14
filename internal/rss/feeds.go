package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	fmt.Printf("Fetching RSS feed from: %s\n", feedURL)
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
