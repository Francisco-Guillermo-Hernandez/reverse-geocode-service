package search

import (
	"github.com/Francisco-Guillermo-Hernandez/reverse-geocode-service/internal/domain"
)

// LocationSearcher describes an interface for geographically finding the closest location.
type LocationSearcher interface {
	Search(lat, lng float64) domain.Location
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

// Search returns the nearest location to the given latitude and longitude.
func (g *kdTreeGeocoder) Search(lat, lng float64) domain.Location {
	var best domain.Location
	bestDist := -1.0
	if g.root != nil {
		g.root.Nearest(lat, lng, &best, &bestDist)
	}
	return best
}
