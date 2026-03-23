package domain

// Location represents a geographic point from geocode.json
type Location struct {
	CountryCode string  `json:"country_code"`
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Population  int     `json:"population"`
	State       string  `json:"state,omitempty"`
	County      string  `json:"county,omitempty"`
	
	// Country name populated via countries.csv match
	Country     string  `json:"country"`
}
