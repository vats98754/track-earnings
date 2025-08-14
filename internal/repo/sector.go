package repo

import (
	"context"
)

type SectorMedian struct {
	AsOf   string  // YYYY-MM-DD
	Sector string
	Metric string
	Value  float64
}

// UpsertSectorMedian stores or updates a sector median metric.
func (r *Repo) UpsertSectorMedian(ctx context.Context, sm SectorMedian) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO sector_medians (as_of, sector, metric, value)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (as_of, sector, metric)
		DO UPDATE SET value=EXCLUDED.value
	`, sm.AsOf, sm.Sector, sm.Metric, sm.Value)
	return err
}

// GetLatestSectorMedians returns a map sector->metric->value for latest as_of.
func (r *Repo) GetLatestSectorMedians(ctx context.Context) (map[string]map[string]float64, error) {
	rows, err := r.DB.QueryContext(ctx, `
		WITH latest AS (
			SELECT MAX(as_of) AS as_of FROM sector_medians
		)
		SELECT m.sector, m.metric, m.value
		FROM sector_medians m
		JOIN latest l ON l.as_of = m.as_of
	`)
	if err != nil { return nil, err }
	defer rows.Close()
	out := map[string]map[string]float64{}
	for rows.Next() {
		var s, m string
		var v float64
		if err := rows.Scan(&s, &m, &v); err != nil { return nil, err }
		if _, ok := out[s]; !ok { out[s] = map[string]float64{} }
		out[s][m] = v
	}
	return out, rows.Err()
}
