package ingest

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vats98754/track-earnings/internal/providers"
)

type MacroJob struct {
	Series   []string
	Provider providers.MacroProvider
	Sink     func(ctx context.Context, series string, at time.Time, value float64) error
}

func (j MacroJob) Run(ctx context.Context) error {
	for _, s := range j.Series {
		pts, err := j.Provider.GetSeries(s)
		if err != nil { log.Error().Str("series", s).Err(err).Msg("macro fetch"); continue }
		for _, p := range pts {
			if j.Sink != nil {
				_ = j.Sink(ctx, p.Series, time.Unix(p.Time, 0), p.Value)
			}
		}
	}
	return nil
}
