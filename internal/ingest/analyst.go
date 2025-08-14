package ingest

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/vats98754/track-earnings/internal/repo"
)

type AnalystProvider interface { GetAnalystEstimates(ticker string) (repo.AnalystRow, error) }

type AnalystJob struct {
	Tickers []string
	Provider AnalystProvider
	Repo *repo.Repo
}

func (j AnalystJob) Run(ctx context.Context) error {
	for _, t := range j.Tickers {
		row, err := j.Provider.GetAnalystEstimates(t)
		if err != nil { log.Error().Str("ticker", t).Err(err).Msg("analyst fetch"); continue }
		if j.Repo != nil { if err := j.Repo.UpsertAnalyst(ctx, row); err != nil { log.Error().Str("ticker", t).Err(err).Msg("db analyst upsert") } }
	}
	return nil
}
