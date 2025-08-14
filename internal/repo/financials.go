package repo

import (
	"context"
	"time"
)

type CompanyProfile struct {
	Ticker            string
	Name              string
	Sector            string
	Industry          string
	Exchange          string
	Country           string
	MarketCap         float64
	SharesOutstanding float64
}

func (r *Repo) UpsertCompanyProfile(ctx context.Context, p CompanyProfile) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO company_profile (ticker, name, sector, industry, exchange, country, market_cap, shares_outstanding, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,now())
		ON CONFLICT (ticker) DO UPDATE SET name=EXCLUDED.name, sector=EXCLUDED.sector, industry=EXCLUDED.industry, exchange=EXCLUDED.exchange, country=EXCLUDED.country, market_cap=EXCLUDED.market_cap, shares_outstanding=EXCLUDED.shares_outstanding, updated_at=now()
	`, p.Ticker, p.Name, p.Sector, p.Industry, p.Exchange, p.Country, p.MarketCap, p.SharesOutstanding)
	return err
}

type IncomeRow struct {
	Ticker string
	PeriodEnd time.Time
	Period string
	Revenue float64
	CostOfRevenue float64
	GrossProfit float64
	OperatingIncome float64
	NetIncome float64
	EPS float64
	EBITDA float64
	Provider string
}

func (r *Repo) UpsertIncome(ctx context.Context, row IncomeRow) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO income_statements (ticker, period_end, period, revenue, cost_of_revenue, gross_profit, operating_income, net_income, eps, ebitda, provider)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (ticker, period_end, period) DO UPDATE SET revenue=EXCLUDED.revenue, cost_of_revenue=EXCLUDED.cost_of_revenue, gross_profit=EXCLUDED.gross_profit, operating_income=EXCLUDED.operating_income, net_income=EXCLUDED.net_income, eps=EXCLUDED.eps, ebitda=EXCLUDED.ebitda, provider=EXCLUDED.provider
	`, row.Ticker, row.PeriodEnd, row.Period, row.Revenue, row.CostOfRevenue, row.GrossProfit, row.OperatingIncome, row.NetIncome, row.EPS, row.EBITDA, row.Provider)
	return err
}

type BalanceRow struct {
	Ticker string
	PeriodEnd time.Time
	Period string
	TotalAssets float64
	TotalLiabilities float64
	TotalDebt float64
	TotalEquity float64
	CashAndEquivalents float64
	CurrentAssets float64
	CurrentLiabilities float64
	Provider string
}

func (r *Repo) UpsertBalance(ctx context.Context, row BalanceRow) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO balance_sheets (ticker, period_end, period, total_assets, total_liabilities, total_debt, total_equity, cash_and_equivalents, current_assets, current_liabilities, provider)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (ticker, period_end, period) DO UPDATE SET total_assets=EXCLUDED.total_assets, total_liabilities=EXCLUDED.total_liabilities, total_debt=EXCLUDED.total_debt, total_equity=EXCLUDED.total_equity, cash_and_equivalents=EXCLUDED.cash_and_equivalents, current_assets=EXCLUDED.current_assets, current_liabilities=EXCLUDED.current_liabilities, provider=EXCLUDED.provider
	`, row.Ticker, row.PeriodEnd, row.Period, row.TotalAssets, row.TotalLiabilities, row.TotalDebt, row.TotalEquity, row.CashAndEquivalents, row.CurrentAssets, row.CurrentLiabilities, row.Provider)
	return err
}

type CashflowRow struct {
	Ticker string
	PeriodEnd time.Time
	Period string
	OperatingCF float64
	InvestingCF float64
	FinancingCF float64
	FreeCashFlow float64
	DividendsPaid float64
	Provider string
}

func (r *Repo) UpsertCashflow(ctx context.Context, row CashflowRow) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO cash_flows (ticker, period_end, period, operating_cash_flow, investing_cash_flow, financing_cash_flow, free_cash_flow, dividends_paid, provider)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (ticker, period_end, period) DO UPDATE SET operating_cash_flow=EXCLUDED.operating_cash_flow, investing_cash_flow=EXCLUDED.investing_cash_flow, financing_cash_flow=EXCLUDED.financing_cash_flow, free_cash_flow=EXCLUDED.free_cash_flow, dividends_paid=EXCLUDED.dividends_paid, provider=EXCLUDED.provider
	`, row.Ticker, row.PeriodEnd, row.Period, row.OperatingCF, row.InvestingCF, row.FinancingCF, row.FreeCashFlow, row.DividendsPaid, row.Provider)
	return err
}

type DividendRow struct { Ticker string; ExDate time.Time; Amount float64 }

func (r *Repo) UpsertDividend(ctx context.Context, row DividendRow) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO dividends (ticker, ex_date, amount) VALUES ($1,$2,$3)
		ON CONFLICT (ticker, ex_date) DO UPDATE SET amount=EXCLUDED.amount
	`, row.Ticker, row.ExDate, row.Amount)
	return err
}

type AnalystRow struct {
	Ticker string
	AsOf time.Time
	NextYearEPS float64
	CurrentYearEPS float64
	RatingBuy int
	RatingHold int
	RatingSell int
	TargetPrice float64
	Provider string
}

func (r *Repo) UpsertAnalyst(ctx context.Context, row AnalystRow) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO analyst_estimates (ticker, as_of, next_year_eps, current_year_eps, rating_buy, rating_hold, rating_sell, target_price, provider)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (ticker, as_of) DO UPDATE SET next_year_eps=EXCLUDED.next_year_eps, current_year_eps=EXCLUDED.current_year_eps, rating_buy=EXCLUDED.rating_buy, rating_hold=EXCLUDED.rating_hold, rating_sell=EXCLUDED.rating_sell, target_price=EXCLUDED.target_price, provider=EXCLUDED.provider
	`, row.Ticker, row.AsOf, row.NextYearEPS, row.CurrentYearEPS, row.RatingBuy, row.RatingHold, row.RatingSell, row.TargetPrice, row.Provider)
	return err
}

type InsiderRow struct {
	Ticker string
	FiledDate time.Time
	Insider string
	TransactionType string
	Shares int64
	Price float64
	TotalValue float64
}

func (r *Repo) UpsertInsider(ctx context.Context, row InsiderRow) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO insider_trades (ticker, filed_date, insider, transaction_type, shares, price, total_value)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (ticker, filed_date, insider, transaction_type) DO UPDATE SET shares=EXCLUDED.shares, price=EXCLUDED.price, total_value=EXCLUDED.total_value
	`, row.Ticker, row.FiledDate, row.Insider, row.TransactionType, row.Shares, row.Price, row.TotalValue)
	return err
}

type FilingRow struct {
	Ticker string
	FiledDate time.Time
	Form string
	Title string
	URL string
}

func (r *Repo) UpsertFiling(ctx context.Context, row FilingRow) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO filings (ticker, filed_date, form, title, url)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (filed_date, form, url) DO NOTHING
	`, row.Ticker, row.FiledDate, row.Form, row.Title, row.URL)
	return err
}

type SplitRow struct {
	Ticker string
	ExDate time.Time
	Ratio string
	Numerator float64
	Denominator float64
}

func (r *Repo) UpsertSplit(ctx context.Context, row SplitRow) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO splits (ticker, ex_date, ratio, numerator, denominator)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (ticker, ex_date) DO UPDATE SET ratio=EXCLUDED.ratio, numerator=EXCLUDED.numerator, denominator=EXCLUDED.denominator
	`, row.Ticker, row.ExDate, row.Ratio, row.Numerator, row.Denominator)
	return err
}

type InstitutionalOwnershipRow struct {
	Ticker string
	AsOf time.Time
	Holders int
	PercentOwned float64
}

func (r *Repo) UpsertInstitutionalOwnership(ctx context.Context, row InstitutionalOwnershipRow) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO institutional_ownership (ticker, as_of, holders, percent_owned)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (ticker, as_of) DO UPDATE SET holders=EXCLUDED.holders, percent_owned=EXCLUDED.percent_owned
	`, row.Ticker, row.AsOf, row.Holders, row.PercentOwned)
	return err
}

func (r *Repo) GetCompanyProfile(ctx context.Context, ticker string) (CompanyProfile, error) {
	var p CompanyProfile
	err := r.DB.GetContext(ctx, &p, `SELECT ticker, name, sector, industry, exchange, country, market_cap, shares_outstanding FROM company_profile WHERE ticker=$1`, ticker)
	return p, err
}
