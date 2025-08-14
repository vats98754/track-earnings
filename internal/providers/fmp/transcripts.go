package fmp

import (
	"net/url"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vats98754/track-earnings/internal/repo"
)

// GetTranscripts fetches latest earnings call transcripts metadata.
func (c *Client) GetTranscripts(ticker string, limit int) ([]repo.FilingRow, error) {
	q := url.Values{}
	if ticker != "" { q.Set("symbol", ticker) }
	if limit > 0 { q.Set("limit", strconv.Itoa(limit)) }
	var raw []struct{
		Symbol string `json:"symbol"`
		Date string `json:"date"`
		Title string `json:"title"`
		Url string `json:"url"`
	}
	if err := c.get("/api/v3/earning_call_transcript", q, &raw); err != nil {
		return nil, err
	}
	out := make([]repo.FilingRow, 0, len(raw))
	for _, r := range raw {
		pe, _ := time.Parse("2006-01-02", r.Date)
		out = append(out, repo.FilingRow{Ticker: r.Symbol, FiledDate: pe, Form: "TRANSCRIPT", Title: r.Title, URL: r.Url})
	}
	log.Debug().Str("ticker", ticker).Int("items", len(out)).Msg("transcripts fetched")
	return out, nil
}
