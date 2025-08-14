package ingest

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vats98754/track-earnings/internal/repo"
)

// BreadthJob computes simple market breadth metrics from latest prices in DB.
// For each run, it loads last close for each ticker and computes advancing/declining counts and percent above 50SMA.

type BreadthJob struct {
	Repo *repo.Repo
	Universe []string
}

func (j BreadthJob) Run(ctx context.Context) error {
	if j.Repo == nil || len(j.Universe) == 0 { return nil }
	adv, dec := 0, 0
	above50 := 0
	for _, t := range j.Universe {
		rows, err := j.Repo.GetRecentPrices(ctx, t, 60)
		if err != nil || len(rows) < 2 { continue }
		last := rows[0].Close
		prev := rows[1].Close
		if last > prev { adv++ } else if last < prev { dec++ }
		// 50SMA
		if len(rows) >= 50 {
			sum := 0.0
			for i := 0; i < 50; i++ { sum += rows[i].Close }
			ma := sum / 50.0
			if last >= ma { above50++ }
		}
	}
	log.Info().Int("adv", adv).Int("dec", dec).Int("above50", above50).Msg("breadth snapshot")
	// Could persist into a table if needed (omitted for brevity)
	_ = time.Now()
	return nil
}
