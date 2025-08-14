package repo

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repo struct { DB *sqlx.DB }

func New(db *sqlx.DB) *Repo { return &Repo{DB: db} }

func (r *Repo) InsertPrice(ctx context.Context, p PriceRow) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO prices (ticker, time, open, high, low, close, volume, provider)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (ticker, time) DO UPDATE SET open=EXCLUDED.open, high=EXCLUDED.high, low=EXCLUDED.low, close=EXCLUDED.close, volume=EXCLUDED.volume, provider=EXCLUDED.provider
	`, p.Ticker, p.Time, p.Open, p.High, p.Low, p.Close, p.Volume, p.Provider)
	return err
}

func (r *Repo) GetRecentPrices(ctx context.Context, ticker string, limit int) ([]PriceRow, error) {
	rows := []PriceRow{}
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT ticker, time, open, high, low, close, volume, provider
		FROM prices WHERE ticker=$1 ORDER BY time DESC LIMIT $2
	`, ticker, limit)
	return rows, err
}

type PriceRow struct {
	Ticker   string    `db:"ticker"`
	Time     time.Time `db:"time"`
	Open     float64   `db:"open"`
	High     float64   `db:"high"`
	Low      float64   `db:"low"`
	Close    float64   `db:"close"`
	Volume   float64   `db:"volume"`
	Provider string    `db:"provider"`
}
