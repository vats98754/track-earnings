package technicals

// Simple RSI placeholder; replace with robust implementation later
func RSI(prices []float64, period int) float64 {
	if len(prices) < period+1 || period <= 0 { return 0 }
	gain, loss := 0.0, 0.0
	for i := 1; i <= period; i++ {
		delta := prices[len(prices)-i] - prices[len(prices)-i-1]
		if delta > 0 { gain += delta } else { loss -= delta }
	}
	if loss == 0 { return 100 }
	rs := gain / loss
	return 100 - (100 / (1 + rs))
}
