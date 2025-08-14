package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vats98754/track-earnings/internal/repo"
)

type NewsService interface {
	GetNewsByTicker(ticker string, limit int) ([]repo.NewsRow, error)
}

type NewsHandler struct { Repo *repo.Repo }

func (h NewsHandler) Routes(r chi.Router) {
	r.Get("/news/{ticker}", func(w http.ResponseWriter, r *http.Request) {
		t := chi.URLParam(r, "ticker")
		limit := 50
		if h.Repo == nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{"error":"repo not configured"})
			return
		}
		items, err := h.Repo.GetNewsByTicker(r.Context(), t, limit)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		_ = json.NewEncoder(w).Encode(items)
	})
}
