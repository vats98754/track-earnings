package services

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/vats98754/track-earnings/internal/repo"
	"github.com/vats98754/track-earnings/internal/providers"
	"github.com/vats98754/track-earnings/internal/valuation"
)

type ValuationService struct {
	Prices providers.PriceProvider
	Fund   providers.FundamentalsProvider
	Repo   *repo.Repo // optional, to fetch sector medians
}

type Report struct {
	Ratios   valuation.Ratios `json:"ratios"`
	Signal   valuation.Signal `json:"signal"`
	Inputs   valuation.Inputs `json:"inputs"`
}

func (s ValuationService) Build(ctx context.Context, ticker string, comp valuation.Comparator) (Report, error) {
	q, err := s.Prices.GetQuote(ticker)
	if err != nil { return Report{}, err }
	f, err := s.Fund.GetFundamentals(ticker)
	if err != nil { return Report{}, err }
	in := valuation.Inputs{
		Price:       q.Price,
		EPS:         f.EPS,
		BookValuePS: f.BookValuePS,
		SalesPS:     f.SalesPS,
		EV:          q.EV,
		EBITDA:      f.EBITDA,
		GrowthEPS:   f.GrowthEPS,
		ROE:         f.ROE,
		ROIC:        f.ROIC,
		FCF:         f.FCF,
		MarketCap:   f.MarketCap,
		Debt:        f.Debt,
		Equity:      f.Equity,
	}
	r := valuation.ComputeRatios(in)

	// Try to replace comparator thresholds with sector medians if available
	if s.Repo != nil {
		meds, err := s.Repo.GetLatestSectorMedians(ctx)
		if err == nil {
			// Fetch company sector
			var sector string
			row := s.Repo.DB.QueryRowContext(ctx, `SELECT sector FROM company_profile WHERE ticker=$1`, ticker)
			_ = row.Scan(&sector)
			if sector != "" {
				if m, ok := meds[sector]; ok {
					if v, ok := m["pe"]; ok { comp.SectorPE = v }
					if v, ok := m["pb"]; ok { comp.SectorPB = v }
					if v, ok := m["evebitda"]; ok { comp.SectorEVEBITDA = v }
				}
			}
		}
	}

	signal := comp.Classify(r)
	log.Info().Str("ticker", ticker).Interface("ratios", r).Str("signal", string(signal)).Msg("valuation report")
	return Report{Ratios: r, Signal: signal, Inputs: in}, nil
}
