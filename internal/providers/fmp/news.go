package fmp

import (
	"fmt"
	"net/url"

	"github.com/vats98754/track-earnings/internal/providers"
)

type newsItem struct {
	PublishedDate string  `json:"publishedDate"`
	Title         string  `json:"title"`
	URL           string  `json:"url"`
	Site          string  `json:"site"`
	Sentiment     *float64 `json:"sentimentScore"`
}

func (c *Client) GetNews(ticker string, limit int) ([]providers.NewsItem, error) {
	if limit <= 0 { limit = 20 }
	q := url.Values{"tickers": []string{ticker}, "limit": []string{itoa(limit)}}
	var raw []newsItem
	if err := c.get("/api/v3/stock_news", q, &raw); err != nil { return nil, err }
	out := make([]providers.NewsItem, 0, len(raw))
	for _, n := range raw {
		var s float64
		if n.Sentiment != nil { s = *n.Sentiment }
		out = append(out, providers.NewsItem{
			Time:      0,
			Title:     n.Title,
			URL:       n.URL,
			Sentiment: s,
		})
	}
	return out, nil
}

func itoa(i int) string { return fmt.Sprintf("%d", i) }
