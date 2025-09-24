package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"majesticcoding.com/api/models"
)

const geocodeURL = "https://maps.googleapis.com/maps/api/geocode/json"

type googleResp struct {
	Results []gResult `json:"results"`
	Status  string    `json:"status"`
}

type gResult struct {
	FormattedAddress  string           `json:"formatted_address"`
	AddressComponents []gAddrComponent `json:"address_components"`
	Geometry          struct {
		Location struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
	} `json:"geometry"`
}

type gAddrComponent struct {
	LongName string   `json:"long_name"`
	Types    []string `json:"types"`
}

func Geocode(ctx context.Context, query string) (*models.GeocodeResult, error) {
	apiKey := strings.TrimSpace(os.Getenv("GCP_API_KEY"))
	if apiKey == "" {
		return nil, fmt.Errorf("missing GCP_API_KEY env var")
	}

	u, _ := url.Parse(geocodeURL)
	q := u.Query()
	q.Set("address", query)
	q.Set("key", apiKey)
	u.RawQuery = q.Encode()

	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	req.Header.Set("User-Agent", "MajesticCoding-Geocode/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("geocoding request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding API returned status %d", resp.StatusCode)
	}

	var gResp googleResp
	if err := json.NewDecoder(resp.Body).Decode(&gResp); err != nil {
		return nil, err
	}
	if gResp.Status != "OK" || len(gResp.Results) == 0 {
		return nil, fmt.Errorf("geocoding failed: status=%s", gResp.Status)
	}

	r := gResp.Results[0]
	return &models.GeocodeResult{
		Formatted:  r.FormattedAddress,
		Location:   models.GeoPoint{Lat: r.Geometry.Location.Lat, Lng: r.Geometry.Location.Lng},
		Components: parseComponents(r.AddressComponents),
		Source:     "google-geocoding",
	}, nil
}

func parseComponents(ac []gAddrComponent) models.AddressComponents {
	var comp models.AddressComponents
	for _, c := range ac {
		switch {
		case contains(c.Types, "locality"):
			comp.City = c.LongName
		case contains(c.Types, "administrative_area_level_1"):
			comp.State = c.LongName
		case contains(c.Types, "country"):
			comp.Country = c.LongName
		case contains(c.Types, "postal_code"):
			comp.Postcode = c.LongName
		case contains(c.Types, "route"):
			comp.Route = c.LongName
		case contains(c.Types, "street_number"):
			comp.StreetNumber = c.LongName
		}
	}
	return comp
}

func contains(list []string, val string) bool {
	for _, x := range list {
		if x == val {
			return true
		}
	}
	return false
}
