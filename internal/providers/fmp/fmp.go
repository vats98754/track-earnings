package fmp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vats98754/track-earnings/internal/providers"
)

type Client struct {
	Key  string
	HTTP *http.Client
}

func New(key string) *Client {
	return &Client{Key: key, HTTP: &http.Client{Timeout: 15 * time.Second}}
}

func (c *Client) get(path string, q url.Values, v any) error {
	if q == nil { q = url.Values{} }
	q.Set("apikey", c.Key)
	u := url.URL{Scheme: "https", Host: "financialmodelingprep.com", Path: path, RawQuery: q.Encode()}
	resp, err := c.HTTP.Get(u.String())
	if err != nil { return err }
	defer resp.Body.Close()
	if resp.StatusCode >= 400 { return fmt.Errorf("fmp http %d", resp.StatusCode) }
	return json.NewDecoder(resp.Body).Decode(v)
}

// GetFundamentals pulls a minimal set of metrics using FMP endpoints.
func (c *Client) GetFundamentals(ticker string) (providers.Fundamentals, error) {
	var out providers.Fundamentals

	// 1) key-metrics-ttm
	var km []struct{
		EPSTTM float64 `json:"epsTTM"`
		BookValuePerShareTTM float64 `json:"bookValuePerShareTTM"`
		RevenuePerShareTTM float64 `json:"revenuePerShareTTM"`
		EBITDAttm float64 `json:"ebitdaTTM"`
		ROETTM float64 `json:"roeTTM"`
		ROICTTM float64 `json:"roicTTM"`
		FreeCashFlowTTM float64 `json:"freeCashFlowTTM"`
		DebtToEquityTTM float64 `json:"debtToEquityTTM"`
	}
	if err := c.get("/api/v3/key-metrics-ttm/"+ticker, nil, &km); err != nil {
		log.Warn().Err(err).Str("ticker", ticker).Msg("fmp key-metrics-ttm failed")
	} else if len(km) > 0 {
		out.EPS = km[0].EPSTTM
		out.BookValuePS = km[0].BookValuePerShareTTM
		out.SalesPS = km[0].RevenuePerShareTTM
		out.EBITDA = km[0].EBITDAttm
		out.ROE = km[0].ROETTM
		out.ROIC = km[0].ROICTTM
		out.FCF = km[0].FreeCashFlowTTM
	}

	// 2) quote for market cap
	var quote []struct{ MarketCap float64 `json:"marketCap"` }
	if err := c.get("/api/v3/quote/"+ticker, nil, &quote); err != nil {
		log.Warn().Err(err).Str("ticker", ticker).Msg("fmp quote failed")
	} else if len(quote) > 0 {
		out.MarketCap = quote[0].MarketCap
	}

	// 3) enterprise-values for EV and net debt (optional)
	var evResp struct{ EnterpriseValues []struct{ EnterpriseValue float64 `json:"enterpriseValue"`; AddTotalDebt float64 `json:"addTotalDebt"`; MinusCashAndCashEquivalents float64 `json:"minusCashAndCashEquivalents"` } `json:"enterpriseValues"` }
	if err := c.get("/api/v3/enterprise-values/"+ticker, url.Values{"limit": []string{"1"}}, &evResp); err == nil && len(evResp.EnterpriseValues) > 0 {
		// We don't store EV in fundamentals; price provider may include EV optionally.
		// Map debt and equity approximately if not present
		if out.Debt == 0 { out.Debt = evResp.EnterpriseValues[0].AddTotalDebt }
	}

	// 4) balance sheet for total equity (optional)
	var bs []struct{ TotalStockholdersEquity float64 `json:"totalStockholdersEquity"` }
	if err := c.get("/api/v3/balance-sheet-statement/"+ticker, url.Values{"limit": []string{"1"}}, &bs); err == nil && len(bs) > 0 {
		out.Equity = bs[0].TotalStockholdersEquity
	}

	return out, nil
}
