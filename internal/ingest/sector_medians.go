package ingest

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vats98754/track-earnings/internal/repo"
	"github.com/vats98754/track-earnings/internal/valuation"
)

// SectorMediansJob computes per-sector medians for valuation metrics from latest available data.
// Requires company_profile.sector and fundamentals_metrics or computed ratios per ticker.

type SectorMediansJob struct {
	Repo *repo.Repo
}

func (j SectorMediansJob) Run(ctx context.Context) error {
	if j.Repo == nil { return nil }
	asOf := time.Now().Format("2006-01-02")
	// Query latest fundamentals_metrics and join with company_profile for sector.
	rows, err := j.Repo.DB.QueryContext(ctx, `
		WITH latest_fm AS (
			SELECT DISTINCT ON (ticker) ticker, asof, eps, book_value_ps, sales_ps, ebitda, growth_eps, roe, roic, fcf, market_cap, debt, equity
			FROM fundamentals_metrics
			ORDER BY ticker, asof DESC
		)
		SELECT p.sector,
		       fm.market_cap, fm.eps, fm.book_value_ps, fm.sales_ps, fm.ebitda, fm.growth_eps, fm.roe, fm.roic, fm.fcf, fm.debt, fm.equity
		FROM latest_fm fm
		JOIN company_profile p ON p.ticker = fm.ticker
		WHERE p.sector IS NOT NULL AND p.sector <> ''
	`)
	if err != nil { return err }
	defer rows.Close()

	// Aggregate per sector
	type acc struct{ pe, pb, ps, eve, peg, fcfy, de []float64 }
	agg := map[string]*acc{}
	for rows.Next() {
		var sector string
		var mc, eps, bvps, sps, ebitda, gr, roe, roic, fcf, debt, equity float64
		if err := rows.Scan(&sector, &mc, &eps, &bvps, &sps, &ebitda, &gr, &roe, &roic, &fcf, &debt, &equity); err != nil { return err }
		in := valuation.Inputs{Price: 0, EPS: eps, BookValuePS: bvps, SalesPS: sps, EV: mc+debt-fcf, EBITDA: ebitda, GrowthEPS: gr, ROE: roe, ROIC: roic, FCF: fcf, MarketCap: mc, Debt: debt, Equity: equity}
		r := valuation.ComputeRatios(in)
		a := agg[sector]
		if a == nil { a = &acc{}; agg[sector] = a }
		push := func(slice *[]float64, v float64) { if !isNaN(v) { *slice = append(*slice, v) } }
		push(&a.pe, r.PE); push(&a.pb, r.PB); push(&a.ps, r.PS); push(&a.eve, r.EVEBITDA); push(&a.peg, r.PEG); push(&a.fcfy, r.FCFYield); push(&a.de, r.DebtEquity)
	}
	if err := rows.Err(); err != nil { return err }

	median := func(vals []float64) float64 { return p50(vals) }
	for sector, a := range agg {
		pairs := []struct{m string; v float64}{
			{"pe", median(a.pe)}, {"pb", median(a.pb)}, {"ps", median(a.ps)}, {"evebitda", median(a.eve)}, {"peg", median(a.peg)}, {"fcfyield", median(a.fcfy)}, {"debt_equity", median(a.de)},
		}
		for _, pr := range pairs {
			if err := j.Repo.UpsertSectorMedian(ctx, repo.SectorMedian{AsOf: asOf, Sector: sector, Metric: pr.m, Value: pr.v}); err != nil {
				log.Error().Err(err).Str("sector", sector).Str("metric", pr.m).Msg("upsert sector median")
			}
		}
	}
	log.Info().Msg("sector medians computed")
	return nil
}

func isNaN(v float64) bool { return v != v }

func p50(vals []float64) float64 {
	if len(vals) == 0 { return 0 }
	// simple nth-element via sort
	s := make([]float64, len(vals))
	copy(s, vals)
	// insertion sort is fine for small n
	for i := 1; i < len(s); i++ {
		j := i
		for j > 0 && s[j-1] > s[j] {
			s[j-1], s[j] = s[j], s[j-1]
			j--
		}
	}
	mid := len(s)/2
	if len(s)%2 == 1 { return s[mid] }
	return 0.5*(s[mid-1]+s[mid])
}
