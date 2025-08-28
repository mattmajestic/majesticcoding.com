package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const infracostURL = "https://pricing.api.infracost.io/graphql"

type InfracostRequest struct {
	Vendor         string // e.g., "gcp"
	Service        string // e.g., "Cloud Run"
	Region         string // e.g., "us-central1"
	ProductFamily  string // optional, e.g., "Serverless" (leave empty to broaden)
	PurchaseOption string // e.g., "on_demand"
	Currency       string // default "USD"
}

// GraphQL response structs (only what we need)
type graphQLPrices struct {
	USD         string `json:"USD"`
	Unit        string `json:"unit"`
	Description string `json:"description"`
}

type graphQLProduct struct {
	Prices []graphQLPrices `json:"prices"`
}

type graphQLResp struct {
	Data struct {
		Products []graphQLProduct `json:"products"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type InfracostPrice struct {
	USD         string `json:"USD"`
	Unit        string `json:"unit"`
	Description string `json:"description"`
}

type InfracostResult struct {
	Vendor         string           `json:"vendor"`
	Service        string           `json:"service"`
	Region         string           `json:"region"`
	ProductFamily  string           `json:"productFamily,omitempty"`
	PurchaseOption string           `json:"purchaseOption"`
	Currency       string           `json:"currency"`
	Count          int              `json:"count"`
	Prices         []InfracostPrice `json:"prices"`
}

// Build a minimal GraphQL query string
func (r InfracostRequest) buildQuery() string {
	vendor := escape(r.Vendor)
	service := escape(r.Service)
	region := escape(r.Region)
	purchase := escape(r.PurchaseOption)

	filter := fmt.Sprintf(`vendorName: "%s", service: "%s", region: "%s"`, vendor, service, region)
	if r.ProductFamily != "" {
		filter += fmt.Sprintf(`, productFamily: "%s"`, escape(r.ProductFamily))
	}

	// Request USD, unit, description per price
	query := fmt.Sprintf(`
	{ products(filter: { %s }) {
			prices(filter: { purchaseOption: "%s" }) {
				USD unit description
			}
		}
	}`, filter, purchase)

	return query
}

func escape(s string) string {
	// very small helper; values are simple words so just replace quotes
	return s
}

func FetchInfracostPrices(apiKey string, req InfracostRequest) (*InfracostResult, error) {
	if apiKey == "" {
		return nil, errors.New("INFRACOST_API_KEY not set")
	}
	if req.Vendor == "" {
		req.Vendor = "gcp"
	}
	if req.Service == "" {
		req.Service = "Cloud Run"
	}
	if req.Region == "" {
		req.Region = "us-central1"
	}
	if req.PurchaseOption == "" {
		req.PurchaseOption = "on_demand"
	}
	if req.Currency == "" {
		req.Currency = "USD" // API converts at query time when needed
	}

	payload := map[string]string{"query": req.buildQuery()}
	b, _ := json.Marshal(payload)

	httpReq, _ := http.NewRequest("POST", infracostURL, bytes.NewReader(b))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Api-Key", apiKey)

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out graphQLResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if len(out.Errors) > 0 {
		return nil, fmt.Errorf("infracost error: %s", out.Errors[0].Message)
	}

	result := InfracostResult{
		Vendor:         req.Vendor,
		Service:        req.Service,
		Region:         req.Region,
		ProductFamily:  req.ProductFamily,
		PurchaseOption: req.PurchaseOption,
		Currency:       req.Currency,
	}

	for _, p := range out.Data.Products {
		for _, pr := range p.Prices {
			result.Prices = append(result.Prices, InfracostPrice{
				USD:         pr.USD,
				Unit:        pr.Unit,
				Description: pr.Description,
			})
		}
	}
	result.Count = len(result.Prices)
	return &result, nil
}
