package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Francisco-Guillermo-Hernandez/reverse-geocode-service/internal/api"
	"github.com/Francisco-Guillermo-Hernandez/reverse-geocode-service/internal/domain"
)

type mockSearcher struct{}

func (m *mockSearcher) Search(lat, lng float64) domain.Location {
	if lat == -37.81 && lng == 144.96 {
		return domain.Location{
			City:        "Melbourne",
			CountryCode: "AU",
			Country:     "Australia",
			Latitude:    -37.81,
			Longitude:   144.96,
		}
	}
	// return empty location to simulate not found
	return domain.Location{}
}

func (m *mockSearcher) SearchCityByName(cityName string, countryCode string) domain.Location {
	if cityName == "Melbourne" && countryCode == "AU" {
		return domain.Location{
			City:        "Melbourne",
			CountryCode: "AU",
			Country:     "Australia",
			Latitude:    -37.81,
			Longitude:   144.96,
		}
	}
	return domain.Location{}
}

func TestGeocodeHandler_HandleGeocode(t *testing.T) {
	searcher := &mockSearcher{}
	handler := api.NewGeocodeHandler(searcher)

	tests := []struct {
		name         string
		method       string
		url          string
		expectedCode int
		expectedCity string
	}{
		{
			name:         "Valid Query",
			method:       http.MethodGet,
			url:          "/api/v1/geocode?lat=-37.81&lng=144.96",
			expectedCode: http.StatusOK,
			expectedCity: "Melbourne",
		},
		{
			name:         "Missing Lat",
			method:       http.MethodGet,
			url:          "/api/v1/geocode?lng=144.96",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid Lat",
			method:       http.MethodGet,
			url:          "/api/v1/geocode?lat=abc&lng=144.96",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Post Method Not Allowed",
			method:       http.MethodPost,
			url:          "/api/v1/geocode?lat=-37.81&lng=144.96",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "Not Found",
			method:       http.MethodGet,
			url:          "/api/v1/geocode?lat=0&lng=0",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			rec := httptest.NewRecorder()
			handler.HandleGeocode(rec, req)

			if rec.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d", tt.expectedCode, rec.Code)
			}

			if tt.expectedCode == http.StatusOK {
				var loc domain.Location
				if err := json.NewDecoder(rec.Body).Decode(&loc); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if loc.City != tt.expectedCity {
					t.Errorf("Expected city %s, got %s", tt.expectedCity, loc.City)
				}
			}
		})
	}
}
