package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vats98754/track-earnings/internal/processing/technicals"
	"github.com/vats98754/track-earnings/internal/repo"
)

type TechnicalsHandler struct{ Repo *repo.Repo }

func (h TechnicalsHandler) Routes(r chi.Router) {
	r.Get("/technicals/{ticker}", h.get)
}

func (h TechnicalsHandler) get(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "ticker")
	if h.Repo == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "no storage"})
		return
	}
	rows, err := h.Repo.GetRecentPrices(r.Context(), t, 300)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	closes := make([]float64, len(rows))
	for i := range rows {
		closes[len(rows)-1-i] = rows[i].Close
	}
	macd, signal := technicals.MACD(closes)
	mid, up, low := technicals.Bollinger(closes, 20, 2)
	resp := map[string]any{
		"sma20":   technicals.SMA(closes, 20),
		"sma50":   technicals.SMA(closes, 50),
		"ema20":   technicals.EMA(closes, 20),
		"rsi14":   technicals.RSI(closes, 14),
		"macd":    macd,
		"signal":  signal,
		"bb_mid":  mid,
		"bb_upper": up,
		"bb_lower": low,
	}
	_ = json.NewEncoder(w).Encode(resp)
}
