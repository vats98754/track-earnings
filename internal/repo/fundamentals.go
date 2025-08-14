package repo

import (
	"context"
	"time"
)

type FundamentalsRow struct {
	Ticker     string
	AsOf       time.Time
	EPS        float64
	BookValuePS float64
	SalesPS    float64
	EBITDA     float64
	GrowthEPS  float64
	ROE        float64
	ROIC       float64
	FCF        float64
	MarketCap  float64
	Debt       float64
	Equity     float64
	Provider   string
}

func (r *Repo) UpsertFundamentals(ctx context.Context, f FundamentalsRow) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO fundamentals_metrics (ticker, asof, eps, book_value_ps, sales_ps, ebitda, growth_eps, roe, roic, fcf, market_cap, debt, equity, provider)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		ON CONFLICT (ticker, asof) DO UPDATE SET eps=EXCLUDED.eps, book_value_ps=EXCLUDED.book_value_ps, sales_ps=EXCLUDED.sales_ps, ebitda=EXCLUDED.ebitda, growth_eps=EXCLUDED.growth_eps, roe=EXCLUDED.roe, roic=EXCLUDED.roic, fcf=EXCLUDED.fcf, market_cap=EXCLUDED.market_cap, debt=EXCLUDED.debt, equity=EXCLUDED.equity, provider=EXCLUDED.provider
	`, f.Ticker, f.AsOf, f.EPS, f.BookValuePS, f.SalesPS, f.EBITDA, f.GrowthEPS, f.ROE, f.ROIC, f.FCF, f.MarketCap, f.Debt, f.Equity, f.Provider)
	return err
}
