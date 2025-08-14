package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vats98754/track-earnings/internal/services"
	"github.com/vats98754/track-earnings/internal/valuation"
)

type ValuationHandler struct {
	Service services.ValuationService
}

func (h ValuationHandler) Routes(r chi.Router) {
	r.Get("/valuation/{ticker}", h.getValuation)
}

func (h ValuationHandler) getValuation(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	comp := valuation.Comparator{SectorPE: 20, SectorPB: 3, SectorEVEBITDA: 12}
	report, err := h.Service.Build(r.Context(), ticker, comp)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(struct {
		At     time.Time       `json:"at"`
		Ticker string          `json:"ticker"`
		Report services.Report `json:"report"`
	}{At: time.Now(), Ticker: ticker, Report: report})
}
