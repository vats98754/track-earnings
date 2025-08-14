package ingest

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vats98754/track-earnings/internal/providers"
	"github.com/vats98754/track-earnings/internal/repo"
)

type NewsJob struct {
	Tickers  []string
	Provider providers.NewsProvider
	Repo     *repo.Repo
	Source   string
}

func (j NewsJob) Run(ctx context.Context) error {
	for _, t := range j.Tickers {
		items, err := j.Provider.GetNews(t, 50)
		if err != nil {
			log.Error().Str("ticker", t).Err(err).Msg("news fetch")
			continue
		}
		for _, it := range items {
			// Normalize sentiment: if FMP source, assume 0..1 -> map to -1..1, else clamp [-1,1]
			s := it.Sentiment
			if j.Source == "fmp" {
				s = 2*s - 1
			}
			if s > 1 {
				s = 1
			} else if s < -1 {
				s = -1
			}
			row := repo.NewsRow{Source: j.Source, Ticker: t, Time: time.Unix(it.Time, 0), Title: it.Title, URL: it.URL, Sentiment: s}
			if j.Repo != nil {
				if err := j.Repo.UpsertNews(ctx, row); err != nil {
					log.Error().Str("ticker", t).Err(err).Msg("db upsert news")
				}
			}
		}
	}
	return nil
}
