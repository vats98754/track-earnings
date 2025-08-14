package repo

import (
	"context"
	"time"
)

type NewsRow struct {
	Source    string
	Ticker    string
	Time      time.Time
	Title     string
	URL       string
	Sentiment float64
}

func (r *Repo) UpsertNews(ctx context.Context, n NewsRow) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO news_items (source, ticker, time, title, url, sentiment)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (source, url) DO UPDATE SET ticker=EXCLUDED.ticker, time=EXCLUDED.time, title=EXCLUDED.title, sentiment=EXCLUDED.sentiment
	`, n.Source, n.Ticker, n.Time, n.Title, n.URL, n.Sentiment)
	return err
}

func (r *Repo) GetNewsByTicker(ctx context.Context, ticker string, limit int) ([]NewsRow, error) {
	if limit <= 0 { limit = 20 }
	rows := []NewsRow{}
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT source, ticker, time, title, url, sentiment
		FROM news_items WHERE ticker=$1 ORDER BY time DESC LIMIT $2
	`, ticker, limit)
	return rows, err
}
