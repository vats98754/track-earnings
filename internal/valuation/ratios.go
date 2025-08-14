package valuation

import "math"

type Ratios struct {
	PE        float64
	PB        float64
	PS        float64
	EVEBITDA  float64
	PEG       float64
	ROE       float64
	ROIC      float64
	FCFYield  float64
	DebtEquity float64
}

func safeDiv(a, b float64) float64 {
	if math.IsNaN(a) || math.IsNaN(b) || b == 0 {
		return math.NaN()
	}
	return a / b
}
