package tests

import (
	"testing"

	"github.com/Francisco-Guillermo-Hernandez/reverse-geocode-service/internal/domain"
	"github.com/Francisco-Guillermo-Hernandez/reverse-geocode-service/internal/search"
)

func TestKDTreeGeocoder_Search(t *testing.T) {
	// Sample dataset mimicking a small portion of geocode.json
	dataset := []domain.Location{
		{City: "Melbourne", Latitude: -37.81, Longitude: 144.96, CountryCode: "AU"},
		{City: "Sydney", Latitude: -33.86, Longitude: 151.21, CountryCode: "AU"},
		{City: "New York", Latitude: 40.71, Longitude: -74.00, CountryCode: "US"},
		{City: "London", Latitude: 51.51, Longitude: -0.13, CountryCode: "GB"},
		{City: "Tokyo", Latitude: 35.69, Longitude: 139.69, CountryCode: "JP"},
	}

	geocoder := search.NewKDTreeGeocoder(dataset)

	tests := []struct {
		name         string
		lat          float64
		lng          float64
		expectedCity string
	}{
		{
			name:         "Exact Match Melbourne",
			lat:          -37.81,
			lng:          144.96,
			expectedCity: "Melbourne",
		},
		{
			name:         "Close to Melbourne",
			lat:          -38.3401,
			lng:          144.7365, // Geelong area, should resolve to Melbourne in our tiny dataset
			expectedCity: "Melbourne",
		},
		{
			name:         "Exact Match New York",
			lat:          40.71,
			lng:          -74.00,
			expectedCity: "New York",
		},
		{
			name:         "Close to London",
			lat:          51.50,
			lng:          -0.10,
			expectedCity: "London",
		},
		{
			name:         "Middle of Pacific Ocean",
			lat:          0.0,
			lng:          -140.0,
			expectedCity: "New York", // Given our small dataset, NY might be closest or Tokyo depending on wrapping (Euclidean won't wrap though)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := geocoder.Search(tt.lat, tt.lng)
			if result.City != tt.expectedCity {
				t.Errorf("expected %s, got %s", tt.expectedCity, result.City)
			}
		})
	}
}
