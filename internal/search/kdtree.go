package search

import (
	"sort"

	"github.com/Francisco-Guillermo-Hernandez/reverse-geocode-service/internal/domain"
)

// KDNode represents a node in the KD-Tree for 2D spatial search.
type KDNode struct {
	Location domain.Location
	Left     *KDNode
	Right    *KDNode
	Axis     int // 0 for Latitude, 1 for Longitude
}

// BuildKDTree recursively builds a balanced KD-Tree from a slice of locations.
// It mutates the order of the locations slice.
func BuildKDTree(locations []domain.Location, depth int) *KDNode {
	if len(locations) == 0 {
		return nil
	}

	axis := depth % 2

	// Sort locations based on the current axis to find the median
	if axis == 0 {
		sort.SliceStable(locations, func(i, j int) bool {
			return locations[i].Latitude < locations[j].Latitude
		})
	} else {
		sort.SliceStable(locations, func(i, j int) bool {
			return locations[i].Longitude < locations[j].Longitude
		})
	}

	median := len(locations) / 2

	return &KDNode{
		Location: locations[median],
		Left:     BuildKDTree(locations[:median], depth+1),
		Right:    BuildKDTree(locations[median+1:], depth+1),
		Axis:     axis,
	}
}

// distanceSq computes the squared Euclidean distance between two coordinate pairs.
// Standard python reverse_geocode behavior uses Euclidean space on unprojected lat/lng.
func distanceSq(lat1, lng1, lat2, lng2 float64) float64 {
	dLat := lat1 - lat2
	dLng := lng1 - lng2
	return dLat*dLat + dLng*dLng
}

// Nearest recursively finds the nearest location in the KD-Tree.
func (n *KDNode) Nearest(lat, lng float64, best *domain.Location, bestDist *float64) {
	if n == nil {
		return
	}

	dist := distanceSq(lat, lng, n.Location.Latitude, n.Location.Longitude)
	if *bestDist < 0 || dist < *bestDist {
		*bestDist = dist
		*best = n.Location
	}

	var primary, secondary *KDNode
	var targetCoord, nodeCoord float64

	if n.Axis == 0 {
		targetCoord = lat
		nodeCoord = n.Location.Latitude
	} else {
		targetCoord = lng
		nodeCoord = n.Location.Longitude
	}

	if targetCoord < nodeCoord {
		primary = n.Left
		secondary = n.Right
	} else {
		primary = n.Right
		secondary = n.Left
	}

	if primary != nil {
		primary.Nearest(lat, lng, best, bestDist)
	}

	// Check if we need to explore the secondary side of the splitting plane
	diff := targetCoord - nodeCoord
	if diff*diff < *bestDist && secondary != nil {
		secondary.Nearest(lat, lng, best, bestDist)
	}
}
