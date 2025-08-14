package ingest

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vats98754/track-earnings/internal/providers"
	"github.com/vats98754/track-earnings/internal/repo"
)

type FundamentalsJob struct {
	Tickers      []string
	Provider     providers.FundamentalsProvider
	Repo         *repo.Repo
	ProviderName string
}

func (j FundamentalsJob) Run(ctx context.Context) error {
	for _, t := range j.Tickers {
		f, err := j.Provider.GetFundamentals(t)
		if err != nil {
			log.Error().Str("ticker", t).Err(err).Msg("fundamentals fetch")
			continue
		}
		// Persist minimal fields to support ratios quickly (TTM snapshot)
		if j.Repo != nil {
			_ = j.Repo.UpsertCompanyProfile(ctx, repo.CompanyProfile{Ticker: t, MarketCap: f.MarketCap})
			now := time.Now()
			_ = j.Repo.UpsertIncome(ctx, repo.IncomeRow{Ticker: t, PeriodEnd: now, Period: "ttm", EPS: f.EPS, EBITDA: f.EBITDA, Provider: j.ProviderName})
			_ = j.Repo.UpsertBalance(ctx, repo.BalanceRow{Ticker: t, PeriodEnd: now, Period: "ttm", TotalEquity: f.Equity, TotalDebt: f.Debt, Provider: j.ProviderName})
			_ = j.Repo.UpsertCashflow(ctx, repo.CashflowRow{Ticker: t, PeriodEnd: now, Period: "ttm", FreeCashFlow: f.FCF, Provider: j.ProviderName})
			// Also write to fundamentals_metrics for sector medians job
			_ = j.Repo.UpsertFundamentals(ctx, repo.FundamentalsRow{Ticker: t, AsOf: now, EPS: f.EPS, BookValuePS: f.BookValuePS, SalesPS: f.SalesPS, EBITDA: f.EBITDA, GrowthEPS: f.GrowthEPS, ROE: f.ROE, ROIC: f.ROIC, FCF: f.FCF, MarketCap: f.MarketCap, Debt: f.Debt, Equity: f.Equity, Provider: j.ProviderName})
		}
	}
	return nil
}
