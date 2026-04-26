package search

import (
	"net/url"
	"strings"

	"github.com/Francisco-Guillermo-Hernandez/reverse-geocode-service/internal/domain"
)

// LocationSearcher describes an interface for geographically finding the closest location.
type LocationSearcher interface {
	Search(lat, lng float64) domain.Location
	SearchCityByName(cityName string, countryCode string) domain.Location
}

type kdTreeGeocoder struct {
	root *KDNode
}

// NewKDTreeGeocoder creates a LocationSearcher powered by a 2D KD-Tree.
func NewKDTreeGeocoder(locations []domain.Location) LocationSearcher {
	// Create a copy of the slice as BuildKDTree mutates it (sorting)
	locs := make([]domain.Location, len(locations))
	copy(locs, locations)

	root := BuildKDTree(locs, 0)
	return &kdTreeGeocoder{root: root}
}

var accentReplacer = strings.NewReplacer(
	"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u",
	"Á", "A", "É", "E", "Í", "I", "Ó", "O", "Ú", "U",
	"ü", "u", "Ü", "U", "ñ", "n", "Ñ", "N",
)

func normalizeString(s string) string {
	return accentReplacer.Replace(s)
}

func levenshtein(s, t string) int {
	if s == "" {
		return len(t)
	}
	if t == "" {
		return len(s)
	}
	v0 := make([]int, len(t)+1)
	v1 := make([]int, len(t)+1)
	for i := 0; i < len(v0); i++ {
		v0[i] = i
	}
	for i := 0; i < len(s); i++ {
		v1[0] = i + 1
		for j := 0; j < len(t); j++ {
			cost := 1
			if s[i] == t[j] {
				cost = 0
			}
			v1[j+1] = min(v1[j]+1, min(v0[j+1]+1, v0[j]+cost))
		}
		for j := 0; j < len(v0); j++ {
			v0[j] = v1[j]
		}
	}
	return v1[len(t)]
}

// Search returns the nearest location to the given latitude and longitude.
func (g *kdTreeGeocoder) Search(lat, lng float64) domain.Location {
	var best domain.Location
	bestDist := -1.0
	if g.root != nil {
		g.root.Nearest(lat, lng, &best, &bestDist)
	}
	return best
}

func (g *kdTreeGeocoder) SearchCityByName(cityName string, countryCode string) domain.Location {
	decodedCityName, err := url.QueryUnescape(cityName)
	if err != nil {
		decodedCityName = cityName
	}
	decodedCityName = strings.ToLower(normalizeString(decodedCityName))
	countryCode = strings.ToUpper(countryCode)

	var found domain.Location
	var bestDist = 999
	var searchNode func(n *KDNode) bool
	searchNode = func(n *KDNode) bool {
		if n == nil {
			return false
		}
		
		if strings.ToUpper(n.Location.CountryCode) == countryCode {
			normalizedCity := strings.ToLower(normalizeString(n.Location.City))
			if strings.Contains(normalizedCity, decodedCityName) {
				found = n.Location
				return true
			}
			// Fallback: check fuzzy match distance (useful for slight typos like "Caracha" instead of "Carcha")
			dist := levenshtein(normalizedCity, decodedCityName)
			if dist <= 2 && dist < bestDist {
				found = n.Location
				bestDist = dist
			}
		}
		
		if searchNode(n.Left) {
			return true
		}
		return searchNode(n.Right)
	}

	searchNode(g.root)
	return found
}
