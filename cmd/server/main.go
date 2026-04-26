package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Francisco-Guillermo-Hernandez/reverse-geocode-service/internal/api"
	"github.com/Francisco-Guillermo-Hernandez/reverse-geocode-service/internal/data"
	"github.com/Francisco-Guillermo-Hernandez/reverse-geocode-service/internal/search"
)

func main() {
	locationsPath := "data/geocode.json"
	countriesPath := "data/countries.csv"

	log.Println("Loading data sources...")
	locations, err := data.PrepareData(locationsPath, countriesPath)
	if err != nil {
		log.Fatalf("Error loading data: %v", err)
	}

	log.Printf("Loaded %d locations. Building KD-Tree...", len(locations))
	searcher := search.NewKDTreeGeocoder(locations)
	log.Println("KD-Tree built successfully.")

	handler := api.NewGeocodeHandler(searcher)

	mux := http.NewServeMux()


	mux.HandleFunc("/geocode", handler.HandleGeocode)
	mux.HandleFunc("/search-city", handler.HandleSearchCity)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("Starting reverse geocode service on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
