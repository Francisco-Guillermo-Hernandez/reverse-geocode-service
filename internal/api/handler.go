package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Francisco-Guillermo-Hernandez/reverse-geocode-service/internal/search"
)

// GeocodeHandler exposes the reverse geocoding functionality over HTTP.
type GeocodeHandler struct {
	searcher search.LocationSearcher
}

// NewGeocodeHandler creates a new GeocodeHandler.
func NewGeocodeHandler(searcher search.LocationSearcher) *GeocodeHandler {
	return &GeocodeHandler{searcher: searcher}
}

// HandleGeocode handles the GET request to reverse geocode a coordinate.
// Expected query parameters: lat, lng
func (h *GeocodeHandler) HandleGeocode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")

	if latStr == "" || lngStr == "" {
		http.Error(w, "Missing lat or lng parameter", http.StatusBadRequest)
		return
	}

	lat, errLat := strconv.ParseFloat(latStr, 64)
	lng, errLng := strconv.ParseFloat(lngStr, 64)

	if errLat != nil || errLng != nil {
		http.Error(w, "Invalid lat or lng parameter", http.StatusBadRequest)
		return
	}

	location := h.searcher.Search(lat, lng)

	// If no location was found, we return a 404
	if location.City == "" && location.CountryCode == "" {
		http.Error(w, "No location found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(location); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
