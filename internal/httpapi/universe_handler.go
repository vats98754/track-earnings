package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vats98754/track-earnings/internal/universe"
)

type UniverseHandler struct { Reg *universe.Registry }

func (h UniverseHandler) Routes(r chi.Router) {
	r.Get("/universes", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(h.Reg.List())
	})
	r.Get("/universe/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		if u, ok := h.Reg.Get(name); ok {
			_ = json.NewEncoder(w).Encode(u)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
	})
}
