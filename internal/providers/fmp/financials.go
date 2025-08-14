package fmp

import (
	"net/url"
	"time"
	"strconv"

	"github.com/vats98754/track-earnings/internal/repo"
)

// GetIncomeStatements fetches historical income statements (annual or quarter)
func (c *Client) GetIncomeStatements(ticker, period string, limit int) ([]repo.IncomeRow, error) {
	q := url.Values{}
	if period != "" { q.Set("period", period) }
	if limit > 0 { q.Set("limit", strconv.Itoa(limit)) }
	var raw []struct{
		Date string `json:"date"`
		Revenue float64 `json:"revenue"`
		CostOfRevenue float64 `json:"costOfRevenue"`
		GrossProfit float64 `json:"grossProfit"`
		OperatingIncome float64 `json:"operatingIncome"`
		NetIncome float64 `json:"netIncome"`
		EPS float64 `json:"eps"`
		EBITDA float64 `json:"ebitda"`
	}
	if err := c.get("/api/v3/income-statement/"+ticker, q, &raw); err != nil { return nil, err }
	rows := make([]repo.IncomeRow, 0, len(raw))
	for _, r := range raw {
		pe, _ := time.Parse("2006-01-02", r.Date)
		rows = append(rows, repo.IncomeRow{Ticker: ticker, PeriodEnd: pe, Period: periodOr(period, "annual"), Revenue: r.Revenue, CostOfRevenue: r.CostOfRevenue, GrossProfit: r.GrossProfit, OperatingIncome: r.OperatingIncome, NetIncome: r.NetIncome, EPS: r.EPS, EBITDA: r.EBITDA, Provider: "fmp"})
	}
	return rows, nil
}

// GetBalanceSheets fetches historical balance sheets (annual or quarter)
func (c *Client) GetBalanceSheets(ticker, period string, limit int) ([]repo.BalanceRow, error) {
	q := url.Values{}
	if period != "" { q.Set("period", period) }
	if limit > 0 { q.Set("limit", strconv.Itoa(limit)) }
	var raw []struct{
		Date string `json:"date"`
		TotalAssets float64 `json:"totalAssets"`
		TotalLiabilities float64 `json:"totalLiabilities"`
		TotalDebt float64 `json:"totalDebt"`
		TotalEquity float64 `json:"totalStockholdersEquity"`
		CashAndEquivalents float64 `json:"cashAndCashEquivalents"`
		CurrentAssets float64 `json:"totalCurrentAssets"`
		CurrentLiabilities float64 `json:"totalCurrentLiabilities"`
	}
	if err := c.get("/api/v3/balance-sheet-statement/"+ticker, q, &raw); err != nil { return nil, err }
	rows := make([]repo.BalanceRow, 0, len(raw))
	for _, r := range raw {
		pe, _ := time.Parse("2006-01-02", r.Date)
		rows = append(rows, repo.BalanceRow{Ticker: ticker, PeriodEnd: pe, Period: periodOr(period, "annual"), TotalAssets: r.TotalAssets, TotalLiabilities: r.TotalLiabilities, TotalDebt: r.TotalDebt, TotalEquity: r.TotalEquity, CashAndEquivalents: r.CashAndEquivalents, CurrentAssets: r.CurrentAssets, CurrentLiabilities: r.CurrentLiabilities, Provider: "fmp"})
	}
	return rows, nil
}

// GetCashFlows fetches historical cash flow statements (annual or quarter)
func (c *Client) GetCashFlows(ticker, period string, limit int) ([]repo.CashflowRow, error) {
	q := url.Values{}
	if period != "" { q.Set("period", period) }
	if limit > 0 { q.Set("limit", strconv.Itoa(limit)) }
	var raw []struct{
		Date string `json:"date"`
		OperatingCF float64 `json:"netCashProvidedByOperatingActivities"`
		InvestingCF float64 `json:"netCashUsedForInvestingActivites"`
		FinancingCF float64 `json:"netCashUsedProvidedByFinancingActivities"`
		FreeCashFlow float64 `json:"freeCashFlow"`
		DividendsPaid float64 `json:"dividendsPaid"`
	}
	if err := c.get("/api/v3/cash-flow-statement/"+ticker, q, &raw); err != nil { return nil, err }
	rows := make([]repo.CashflowRow, 0, len(raw))
	for _, r := range raw {
		pe, _ := time.Parse("2006-01-02", r.Date)
		rows = append(rows, repo.CashflowRow{Ticker: ticker, PeriodEnd: pe, Period: periodOr(period, "annual"), OperatingCF: r.OperatingCF, InvestingCF: r.InvestingCF, FinancingCF: r.FinancingCF, FreeCashFlow: r.FreeCashFlow, DividendsPaid: r.DividendsPaid, Provider: "fmp"})
	}
	return rows, nil
}

// GetCompanyProfile fetches company metadata
func (c *Client) GetCompanyProfile(ticker string) (repo.CompanyProfile, error) {
	var raw []struct{
		CompanyName string `json:"companyName"`
		Sector string `json:"sector"`
		Industry string `json:"industry"`
		Exchange string `json:"exchange"`
		Country string `json:"country"`
		MktCap float64 `json:"mktCap"`
		Shares float64 `json:"sharesOutstanding"`
	}
	if err := c.get("/api/v3/profile/"+ticker, nil, &raw); err != nil { return repo.CompanyProfile{Ticker: ticker}, err }
	p := repo.CompanyProfile{Ticker: ticker}
	if len(raw) > 0 {
		p.Name = raw[0].CompanyName
		p.Sector = raw[0].Sector
		p.Industry = raw[0].Industry
		p.Exchange = raw[0].Exchange
		p.Country = raw[0].Country
		p.MarketCap = raw[0].MktCap
		p.SharesOutstanding = raw[0].Shares
	}
	return p, nil
}

func periodOr(p, d string) string { if p == "" { return d }; return p }
