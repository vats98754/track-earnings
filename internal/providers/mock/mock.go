package mock

import (
	"math/rand"
	"time"

	"github.com/vats98754/track-earnings/internal/providers"
)

type Client struct{}

func New() *Client { return &Client{} }

func (c *Client) GetQuote(ticker string) (providers.PriceQuote, error) {
	rand.Seed(time.Now().UnixNano())
	price := 100 + rand.Float64()*50
	return providers.PriceQuote{
		Price:  price,
		Open:   price * 0.99,
		High:   price * 1.02,
		Low:    price * 0.98,
		Volume: 1_000_000 + rand.Float64()*500_000,
		EV:     price*10_000_000 + 5_000_000_000, // fake EV
	}, nil
}

func (c *Client) GetFundamentals(ticker string) (providers.Fundamentals, error) {
	return providers.Fundamentals{
		EPS:         5.0,
		BookValuePS: 20.0,
		SalesPS:     40.0,
		EBITDA:      50_000_000_000,
		GrowthEPS:   0.12,
		ROE:         0.25,
		ROIC:        0.18,
		FCF:         15_000_000_000,
		MarketCap:   2_000_000_000_000,
		Debt:        120_000_000_000,
		Equity:      400_000_000_000,
	}, nil
}
