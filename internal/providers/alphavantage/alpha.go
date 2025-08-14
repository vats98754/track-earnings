package alphavantage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/vats98754/track-earnings/internal/providers"
)

type Client struct {
	Key string
	HTTP *http.Client
}

func New(key string) *Client {
	return &Client{Key: key, HTTP: &http.Client{Timeout: 15 * time.Second}}
}

func (c *Client) endpoint(function string, params map[string]string) string {
	q := url.Values{}
	q.Set("function", function)
	q.Set("apikey", c.Key)
	for k, v := range params { q.Set(k, v) }
	return "https://www.alphavantage.co/query?" + q.Encode()
}

// Minimal demo: GLOBAL_QUOTE
func (c *Client) GetQuote(ticker string) (providers.PriceQuote, error) {
	url := c.endpoint("GLOBAL_QUOTE", map[string]string{"symbol": ticker})
	resp, err := c.HTTP.Get(url)
	if err != nil { return providers.PriceQuote{}, err }
	defer resp.Body.Close()
	var raw map[string]map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil { return providers.PriceQuote{}, err }
	q := raw["Global Quote"]
	if q == nil { return providers.PriceQuote{}, fmt.Errorf("no quote") }
	parse := func(k string) float64 { var f float64; _, _ = fmt.Sscanf(q[k], "%f", &f); return f }
	return providers.PriceQuote{
		Price:  parse("05. price"),
		Open:   parse("02. open"),
		High:   parse("03. high"),
		Low:    parse("04. low"),
		Volume: parse("06. volume"),
	}, nil
}

// Placeholders for extended endpoints (BALANCE_SHEET, INCOME_STATEMENT, CASH_FLOW)
func (c *Client) GetFundamentals(ticker string) (providers.Fundamentals, error) {
	log.Warn().Str("provider", "alphavantage").Msg("GetFundamentals not implemented: returning zero values")
	return providers.Fundamentals{}, nil
}
