package ingest

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/vats98754/track-earnings/internal/repo"
)

// TranscriptsProvider returns transcript metadata rows (mapped to filings table with Form=TRANSCRIPT for simplicity)
type TranscriptsProvider interface { GetTranscripts(ticker string, limit int) ([]repo.FilingRow, error) }

type TranscriptsJob struct {
	Tickers []string
	Provider TranscriptsProvider
	Repo *repo.Repo
}

func (j TranscriptsJob) Run(ctx context.Context) error {
	if j.Repo == nil || j.Provider == nil { return nil }
	for _, t := range j.Tickers {
		rows, err := j.Provider.GetTranscripts(t, 10)
		if err != nil { log.Error().Str("ticker", t).Err(err).Msg("transcripts fetch"); continue }
		for _, r := range rows {
			if err := j.Repo.UpsertFiling(ctx, r); err != nil {
				log.Error().Str("ticker", t).Err(err).Msg("transcript upsert")
			}
		}
	}
	return nil
}
