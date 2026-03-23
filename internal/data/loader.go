package data

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Francisco-Guillermo-Hernandez/reverse-geocode-service/internal/domain"
)

// LoadCountries parses the countries.csv file and returns a map of CountryCode to CountryName.
func LoadCountries(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open countries file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	countries := make(map[string]string)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading csv record: %w", err)
		}

		// Assuming structure: code,name
		if len(record) >= 2 {
			countries[record[0]] = record[1]
		}
	}

	return countries, nil
}

// LoadLocations parses the geocode.json file into a slice of Locations.
func LoadLocations(filePath string) ([]domain.Location, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open locations file: %w", err)
	}
	defer file.Close()

	var locations []domain.Location
	if err := json.NewDecoder(file).Decode(&locations); err != nil {
		return nil, fmt.Errorf("failed to decode json locations: %w", err)
	}

	return locations, nil
}

// PrepareData loads both sources and attaches the country name to each location.
func PrepareData(locationsPath, countriesPath string) ([]domain.Location, error) {
	countries, err := LoadCountries(countriesPath)
	if err != nil {
		return nil, fmt.Errorf("could not load countries: %w", err)
	}

	locations, err := LoadLocations(locationsPath)
	if err != nil {
		return nil, fmt.Errorf("could not load locations: %w", err)
	}

	// Populate the Country field for each location
	for i := range locations {
		if name, ok := countries[locations[i].CountryCode]; ok {
			locations[i].Country = name
		}
	}

	return locations, nil
}
