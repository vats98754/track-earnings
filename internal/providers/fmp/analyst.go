package fmp

import (
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vats98754/track-earnings/internal/repo"
)

// GetAnalystEstimates returns minimal snapshot for a ticker.
func (c *Client) GetAnalystEstimates(ticker string) (repo.AnalystRow, error) {
	// Endpoint: /api/v3/analyst-estimates/{ticker}
	var raw []struct{
		Date string `json:"date"`
		EstimatedEPSAvg float64 `json:"estimatedEPSAvg"`
		NumberAnalystEvaluating int `json:"numberAnalystEvaluating"`
	}
	if err := c.get("/api/v3/analyst-estimates/"+ticker, url.Values{"limit": []string{"1"}}, &raw); err != nil {
		return repo.AnalystRow{}, err
	}
	row := repo.AnalystRow{Ticker: ticker, Provider: "fmp"}
	if len(raw) > 0 {
		row.AsOf, _ = time.Parse("2006-01-02", raw[0].Date)
		row.NextYearEPS = raw[0].EstimatedEPSAvg
	}
	// Price target (different endpoint)
	var pt []struct{ TargetMedian float64 `json:"targetMedian"`; Updated string `json:"updated"` }
	if err := c.get("/api/v3/price-target/"+ticker, nil, &pt); err == nil && len(pt) > 0 {
		row.TargetPrice = pt[0].TargetMedian
	}
	// Ratings summary (optional)
	var rt []struct{ Buy int `json:"buy"`; Hold int `json:"hold"`; Sell int `json:"sell"` }
	if err := c.get("/api/v3/rating/"+ticker, url.Values{"limit": []string{"1"}}, &rt); err == nil && len(rt) > 0 {
		row.RatingBuy = rt[0].Buy
		row.RatingHold = rt[0].Hold
		row.RatingSell = rt[0].Sell
	}
	if row.AsOf.IsZero() { row.AsOf = time.Now() }
	log.Debug().Str("ticker", ticker).Msg("analyst estimates fetched")
	return row, nil
}
