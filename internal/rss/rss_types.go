package rss

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PublishDate string `xml:"pubDate"`
}
