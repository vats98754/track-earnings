package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Services struct {
	Valuation  ValuationHandler
	Universes  UniverseHandler
	Technicals TechnicalsHandler
	News       NewsHandler
}

func NewRouter(svcs Services) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Route("/api", func(api chi.Router) {
		svcs.Valuation.Routes(api)
		svcs.Universes.Routes(api)
		svcs.Technicals.Routes(api)
		svcs.News.Routes(api)
	})
	return r
}
