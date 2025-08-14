package ingest

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vats98754/track-earnings/internal/repo"
)

// CorporateActionsProvider covers splits and dividends.
type CorporateActionsProvider interface {
	GetSplits(ticker string, from, to time.Time) ([]repo.SplitRow, error)
	GetDividends(ticker string, from, to time.Time) ([]repo.DividendRow, error)
}

type CorporateActionsJob struct {
	Tickers []string
	Provider CorporateActionsProvider
	Repo *repo.Repo
}

func (j CorporateActionsJob) Run(ctx context.Context) error {
	if j.Repo == nil || j.Provider == nil { return nil }
	to := time.Now()
	from := to.AddDate(-10,0,0)
	for _, t := range j.Tickers {
		splits, err := j.Provider.GetSplits(t, from, to)
		if err == nil {
			for _, s := range splits { if err := j.Repo.UpsertSplit(ctx, s); err != nil { log.Error().Str("ticker", t).Err(err).Msg("split upsert") } }
		} else { log.Error().Str("ticker", t).Err(err).Msg("splits fetch") }
		divs, err := j.Provider.GetDividends(t, from, to)
		if err == nil {
			for _, d := range divs { if err := j.Repo.UpsertDividend(ctx, d); err != nil { log.Error().Str("ticker", t).Err(err).Msg("dividend upsert") } }
		} else { log.Error().Str("ticker", t).Err(err).Msg("dividends fetch") }
	}
	return nil
}
