package ingest

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vats98754/track-earnings/internal/repo"
)

// FinancialsProvider is implemented by providers that can fetch historical financial statements.
type FinancialsProvider interface {
	GetIncomeStatements(ticker, period string, limit int) ([]repo.IncomeRow, error)
	GetBalanceSheets(ticker, period string, limit int) ([]repo.BalanceRow, error)
	GetCashFlows(ticker, period string, limit int) ([]repo.CashflowRow, error)
	GetCompanyProfile(ticker string) (repo.CompanyProfile, error)
}

// FinancialsHistoryJob backfills quarterly and annual financials with simple rate limiting.
// Rate limiting: sleep between tickers and between API calls.

type FinancialsHistoryJob struct {
	Tickers  []string
	Provider FinancialsProvider
	Repo     *repo.Repo
	PerTickerSleep time.Duration
	PerCallSleep time.Duration
	Limit int // per endpoint, per ticker
}

func (j FinancialsHistoryJob) Run(ctx context.Context) error {
	if j.Repo == nil || j.Provider == nil { return nil }
	if j.PerTickerSleep == 0 { j.PerTickerSleep = 2 * time.Second }
	if j.PerCallSleep == 0 { j.PerCallSleep = 750 * time.Millisecond }
	if j.Limit == 0 { j.Limit = 40 }
	for _, t := range j.Tickers {
		// profile
		if p, err := j.Provider.GetCompanyProfile(t); err == nil {
			_ = j.Repo.UpsertCompanyProfile(ctx, p)
		}
		// annual
		if rows, err := j.Provider.GetIncomeStatements(t, "annual", j.Limit); err == nil {
			for _, r := range rows { _ = j.Repo.UpsertIncome(ctx, r) }
		}
		time.Sleep(j.PerCallSleep)
		if rows, err := j.Provider.GetBalanceSheets(t, "annual", j.Limit); err == nil {
			for _, r := range rows { _ = j.Repo.UpsertBalance(ctx, r) }
		}
		time.Sleep(j.PerCallSleep)
		if rows, err := j.Provider.GetCashFlows(t, "annual", j.Limit); err == nil {
			for _, r := range rows { _ = j.Repo.UpsertCashflow(ctx, r) }
		}
		// quarterly
		time.Sleep(j.PerCallSleep)
		if rows, err := j.Provider.GetIncomeStatements(t, "quarter", j.Limit); err == nil {
			for _, r := range rows { _ = j.Repo.UpsertIncome(ctx, r) }
		}
		time.Sleep(j.PerCallSleep)
		if rows, err := j.Provider.GetBalanceSheets(t, "quarter", j.Limit); err == nil {
			for _, r := range rows { _ = j.Repo.UpsertBalance(ctx, r) }
		}
		time.Sleep(j.PerCallSleep)
		if rows, err := j.Provider.GetCashFlows(t, "quarter", j.Limit); err == nil {
			for _, r := range rows { _ = j.Repo.UpsertCashflow(ctx, r) }
		}

		log.Info().Str("ticker", t).Msg("financials backfilled")
		time.Sleep(j.PerTickerSleep)
	}
	return nil
}
