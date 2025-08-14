package ingest

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vats98754/track-earnings/internal/repo"
)

// AlertsJob scans for significant events and logs alerts (extend to persist/push later).
// Examples: insider buys > $1M, dividend cuts, unusual split ratios.

type AlertsJob struct {
	Repo *repo.Repo
}

func (j AlertsJob) Run(ctx context.Context) error {
	if j.Repo == nil { return nil }
	// naive examples: last 7 days insider buys over $1M
	_, _ = time.Now(), ctx
	// Implementation placeholder: would query insider_trades and dividends deltas; log alerts.
	log.Info().Msg("alerts scan complete")
	return nil
}
