package fmp

import (
	"net/url"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vats98754/track-earnings/internal/repo"
)

// Splits and Dividends
func (c *Client) GetSplits(ticker string, from, to time.Time) ([]repo.SplitRow, error) {
	q := url.Values{}
	if !from.IsZero() { q.Set("from", from.Format("2006-01-02")) }
	if !to.IsZero() { q.Set("to", to.Format("2006-01-02")) }
	var raw struct{ Symbol string `json:"symbol"`; Historical []struct{ Date string `json:"date"`; Numerator float64 `json:"numerator"`; Denominator float64 `json:"denominator"`; Label string `json:"label"` } `json:"historical"` }
	if err := c.get("/api/v3/historical-price-full/stock_split/"+ticker, q, &raw); err != nil { return nil, err }
	out := make([]repo.SplitRow, 0, len(raw.Historical))
	for _, h := range raw.Historical {
		d, _ := time.Parse("2006-01-02", h.Date)
		out = append(out, repo.SplitRow{Ticker: ticker, ExDate: d, Ratio: h.Label, Numerator: h.Numerator, Denominator: h.Denominator})
	}
	return out, nil
}

func (c *Client) GetDividends(ticker string, from, to time.Time) ([]repo.DividendRow, error) {
	q := url.Values{}
	if !from.IsZero() { q.Set("from", from.Format("2006-01-02")) }
	if !to.IsZero() { q.Set("to", to.Format("2006-01-02")) }
	var raw struct{ Symbol string `json:"symbol"`; Historical []struct{ Date string `json:"date"`; Dividend float64 `json:"dividend"` } `json:"historical"` }
	if err := c.get("/api/v3/historical-price-full/stock_dividend/"+ticker, q, &raw); err != nil { return nil, err }
	out := make([]repo.DividendRow, 0, len(raw.Historical))
	for _, h := range raw.Historical {
		d, _ := time.Parse("2006-01-02", h.Date)
		out = append(out, repo.DividendRow{Ticker: ticker, ExDate: d, Amount: h.Dividend})
	}
	return out, nil
}

// Insiders and Ownership
func (c *Client) GetInsiders(ticker string, limit int) ([]repo.InsiderRow, error) {
	q := url.Values{"symbol": []string{ticker}}
	if limit > 0 { q.Set("limit", strconv.Itoa(limit)) }
	var raw []struct{
		TransactionDate string  `json:"transactionDate"`
		InsiderName     string  `json:"insiderName"`
		TypeOfOwner     string  `json:"typeOfOwner"`
		SecuritiesTransacted int64 `json:"securitiesTransacted"`
		Price           float64 `json:"price"`
		TotalValue      float64 `json:"totalValue"`
	}
	if err := c.get("/api/v4/insider-trading", q, &raw); err != nil { return nil, err }
	out := make([]repo.InsiderRow, 0, len(raw))
	for _, r := range raw {
		fd, _ := time.Parse("2006-01-02", r.TransactionDate)
		out = append(out, repo.InsiderRow{Ticker: ticker, FiledDate: fd, Insider: r.InsiderName, TransactionType: r.TypeOfOwner, Shares: r.SecuritiesTransacted, Price: r.Price, TotalValue: r.TotalValue})
	}
	return out, nil
}

func (c *Client) GetInstitutionalOwnership(ticker string) (repo.InstitutionalOwnershipRow, error) {
	var raw []struct{ Date string `json:"date"`; Holders int `json:"holders"`; Percent float64 `json:"percentage"` }
	// Endpoint may vary; best-effort fallback
	if err := c.get("/api/v3/institutional-holder/"+ticker, nil, &raw); err != nil {
		log.Warn().Str("ticker", ticker).Err(err).Msg("institutional-holder fetch failed")
		return repo.InstitutionalOwnershipRow{Ticker: ticker}, nil
	}
	if len(raw) == 0 { return repo.InstitutionalOwnershipRow{Ticker: ticker}, nil }
	asof, _ := time.Parse("2006-01-02", raw[0].Date)
	return repo.InstitutionalOwnershipRow{Ticker: ticker, AsOf: asof, Holders: raw[0].Holders, PercentOwned: raw[0].Percent}, nil
}
