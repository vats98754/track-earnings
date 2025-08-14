package ingest

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vats98754/track-earnings/internal/repo"
)

// InsiderOwnershipProvider provides insiders and institutional ownership snapshots.
type InsiderOwnershipProvider interface {
	GetInsiders(ticker string, limit int) ([]repo.InsiderRow, error)
	GetInstitutionalOwnership(ticker string) (repo.InstitutionalOwnershipRow, error)
}

type InsiderOwnershipJob struct {
	Tickers []string
	Provider InsiderOwnershipProvider
	Repo *repo.Repo
}

func (j InsiderOwnershipJob) Run(ctx context.Context) error {
	if j.Repo == nil || j.Provider == nil { return nil }
	for _, t := range j.Tickers {
		ins, err := j.Provider.GetInsiders(t, 100)
		if err == nil {
			for _, r := range ins { if err := j.Repo.UpsertInsider(ctx, r); err != nil { log.Error().Str("ticker", t).Err(err).Msg("insider upsert") } }
		} else { log.Error().Str("ticker", t).Err(err).Msg("insiders fetch") }
		own, err := j.Provider.GetInstitutionalOwnership(t)
		if err == nil && !own.AsOf.IsZero() {
			if err := j.Repo.UpsertInstitutionalOwnership(ctx, own); err != nil { log.Error().Str("ticker", t).Err(err).Msg("ownership upsert") }
		}
		// sleep a little to be kind to API
		time.Sleep(300 * time.Millisecond)
	}
	return nil
}
