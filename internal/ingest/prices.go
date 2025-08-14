package ingest

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vats98754/track-earnings/internal/providers"
	"github.com/vats98754/track-earnings/internal/repo"
	"github.com/vats98754/track-earnings/internal/storage"
)

type PricesJob struct {
	Tickers      []string
	Provider     providers.PriceProvider
	ProviderName string
	Repo         *repo.Repo
	Redis        *storage.Redis
}

func (j PricesJob) Run(ctx context.Context) error {
	for _, t := range j.Tickers {
		q, err := j.Provider.GetQuote(t)
		if err != nil { log.Error().Str("ticker", t).Err(err).Msg("quote error"); continue }
		row := repo.PriceRow{
			Ticker:   t,
			Time:     time.Now(),
			Open:     q.Open,
			High:     q.High,
			Low:      q.Low,
			Close:    q.Price,
			Volume:   q.Volume,
			Provider: j.ProviderName,
		}
		if j.Repo != nil {
			if err := j.Repo.InsertPrice(ctx, row); err != nil {
				log.Error().Str("ticker", t).Err(err).Msg("db insert price")
			}
		}
		if j.Redis != nil {
			key := fmt.Sprintf("%s:price", t)
			if err := j.Redis.C.Set(ctx, key, q.Price, 0).Err(); err != nil {
				log.Error().Str("key", key).Err(err).Msg("redis set price")
			}
		}
	}
	return nil
}
