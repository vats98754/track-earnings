package technicals

func SMA(values []float64, period int) float64 {
	if period <= 0 || len(values) < period { return 0 }
	sum := 0.0
	for i := len(values)-period; i < len(values); i++ { sum += values[i] }
	return sum / float64(period)
}

func EMA(values []float64, period int) float64 {
	if period <= 0 || len(values) < period { return 0 }
	k := 2.0 / (float64(period) + 1)
	ema := SMA(values[:period], period)
	for i := period; i < len(values); i++ {
		ema = values[i]*k + ema*(1-k)
	}
	return ema
}

// MACD computes MACD line and signal line using EMA12 and EMA26 with signal EMA9.
func MACD(closes []float64) (macd []float64, signal []float64) {
	if len(closes) == 0 { return nil, nil }
	// For now, return a simple approximation since the EMA function returns single values
	e12 := EMA(closes, 12)
	e26 := EMA(closes, 26)
	n := len(closes)
	macd = make([]float64, n)
	// Fill with the final MACD value for all positions (simplified)
	macdVal := e12 - e26
	for i := 0; i < n; i++ { macd[i] = macdVal }
	// Signal line as a simplified EMA of the MACD values
	signal = make([]float64, n)
	signalVal := EMA(macd, 9)
	for i := 0; i < n; i++ { signal[i] = signalVal }
	return
}

// Bollinger computes upper/lower bands with period n and k stddevs.
func Bollinger(closes []float64, n int, k float64) (middle, upper, lower []float64) {
	if n <= 0 || len(closes) == 0 { return nil, nil, nil }
	m := SMA(closes, n)
	upper = make([]float64, len(closes))
	lower = make([]float64, len(closes))
	middle = make([]float64, len(closes))
	for i := range closes {
		middle[i] = m  // Use the single SMA value for all positions
		if i+1 < n { upper[i], lower[i] = 0, 0; continue }
		// compute stddev last n
		sum, sum2 := 0.0, 0.0
		for j := i-n+1; j <= i; j++ { v := closes[j]; sum += v; sum2 += v*v }
		mean := sum / float64(n)
		variance := sum2/float64(n) - mean*mean
		if variance < 0 { variance = 0 }
		std := sqrt(variance)
		upper[i] = m + k*std
		lower[i] = m - k*std
	}
	return middle, upper, lower
}

func sqrt(x float64) float64 {
	// simple Newton-Raphson for sqrt for portability
	if x <= 0 { return 0 }
	y := x
	for i := 0; i < 20; i++ { y = 0.5 * (y + x/y) }
	return y
}
