package valuation

import (
	"math"
)

type Inputs struct {
	Price       float64
	EPS         float64
	BookValuePS float64
	SalesPS     float64
	EV          float64
	EBITDA      float64
	GrowthEPS   float64 // expected EPS growth (y/y) as decimal
	ROE         float64
	ROIC        float64
	FCF         float64
	MarketCap   float64
	Debt        float64
	Equity      float64
}

func ComputeRatios(in Inputs) Ratios {
	var r Ratios
	r.PE = safeDiv(in.Price, in.EPS)
	r.PB = safeDiv(in.Price, in.BookValuePS)
	r.PS = safeDiv(in.Price, in.SalesPS)
	r.EVEBITDA = safeDiv(in.EV, in.EBITDA)
	r.PEG = safeDiv(r.PE, in.GrowthEPS)
	r.ROE = in.ROE
	r.ROIC = in.ROIC
	r.FCFYield = safeDiv(in.FCF, in.MarketCap)
	r.DebtEquity = safeDiv(in.Debt, in.Equity)
	return r
}

type Signal string

const (
	SignalUndervalued Signal = "undervalued"
	SignalFair        Signal = "fair"
	SignalOvervalued  Signal = "overvalued"
)

type Comparator struct {
	SectorPE    float64
	SectorPB    float64
	SectorEVEBITDA float64
}

func (c Comparator) Classify(r Ratios) Signal {
	// Simple rule-based; extend with z-scores or Bayesian model later
	score := 0
	if !math.IsNaN(r.PE) && c.SectorPE > 0 {
		if r.PE < 0.8*c.SectorPE { score-- }
		if r.PE > 1.2*c.SectorPE { score++ }
	}
	if !math.IsNaN(r.PB) && c.SectorPB > 0 {
		if r.PB < 0.8*c.SectorPB { score-- }
		if r.PB > 1.2*c.SectorPB { score++ }
	}
	if !math.IsNaN(r.EVEBITDA) && c.SectorEVEBITDA > 0 {
		if r.EVEBITDA < 0.8*c.SectorEVEBITDA { score-- }
		if r.EVEBITDA > 1.2*c.SectorEVEBITDA { score++ }
	}
	if r.FCFYield > 0.06 { score-- } // generous cash generator
	if r.DebtEquity > 2.0 { score++ }

	if score <= -2 { return SignalUndervalued }
	if score >= 2 { return SignalOvervalued }
	return SignalFair
}
