package repo

import (
	"context"
	"time"
)

func (r *Repo) UpsertMacroPoint(ctx context.Context, series string, at time.Time, value float64, source string) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO macro_points (series, time, value, source)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (series, time) DO UPDATE SET value=EXCLUDED.value, source=EXCLUDED.source
	`, series, at, value, source)
	return err
}
