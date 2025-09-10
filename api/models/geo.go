package models

type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type AddressComponents struct {
	StreetNumber string `json:"street_number,omitempty"`
	Route        string `json:"route,omitempty"`
	City         string `json:"city,omitempty"`
	County       string `json:"county,omitempty"`
	State        string `json:"state,omitempty"`
	Postcode     string `json:"postcode,omitempty"`
	Country      string `json:"country,omitempty"`
}

type GeocodeResult struct {
	Formatted  string            `json:"formatted"`
	Location   GeoPoint          `json:"location"`
	Components AddressComponents `json:"components"`
	Source     string            `json:"source"`
}
