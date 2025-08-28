package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type servicesResp struct {
	Services []struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	} `json:"services"`
}

type money struct {
	CurrencyCode string `json:"currencyCode"`
	Units        string `json:"units"`
	Nanos        int64  `json:"nanos"`
}

type tieredRate struct {
	UnitPrice        money   `json:"unitPrice"`
	StartUsageAmount float64 `json:"startUsageAmount"`
}

type pricingExpression struct {
	TieredRates []tieredRate `json:"tieredRates"`
}

type pricingInfo struct {
	PricingExpression pricingExpression `json:"pricingExpression"`
}

type sku struct {
	SkuId          string        `json:"skuId"`
	Description    string        `json:"description"`
	ServiceRegions []string      `json:"serviceRegions"`
	PricingInfo    []pricingInfo `json:"pricingInfo"`
}

type skusPage struct {
	Skus          []sku  `json:"skus"`
	NextPageToken string `json:"nextPageToken"`
}

type CloudRunAverages struct {
	Region              string   `json:"region"`
	Currency            string   `json:"currency"`
	AvgVCPUSecondUSD    float64  `json:"avgVCPUSecondUSD"`
	AvgGiBSecondUSD     float64  `json:"avgGiBSecondUSD"`
	AvgReqPerMillionUSD float64  `json:"avgReqPerMillionUSD"`
	SampledCount        int      `json:"sampledCount"`
	SampledSKUs         []string `json:"sampledSkus"` // skuId:description
}

var httpClient = &http.Client{Timeout: 20 * time.Second}

func FetchCloudRunAverages(apiKey, currency, region string) (*CloudRunAverages, error) {
	if apiKey == "" {
		return nil, errors.New("GCP_API_KEY empty")
	}
	if currency == "" {
		currency = "USD"
	}
	if region == "" {
		region = "us-central1"
	}

	// 1) Find the Cloud Run service
	svURL := fmt.Sprintf("https://cloudbilling.googleapis.com/v1/services?pageSize=5000&key=%s",
		url.QueryEscape(apiKey))
	svRes, err := httpClient.Get(svURL)
	if err != nil {
		return nil, err
	}
	defer svRes.Body.Close()

	var s servicesResp
	if err := json.NewDecoder(svRes.Body).Decode(&s); err != nil {
		return nil, err
	}

	serviceName := ""
	for _, v := range s.Services {
		if v.DisplayName == "Cloud Run" {
			serviceName = v.Name
			break
		}
	}
	if serviceName == "" {
		return nil, errors.New("Cloud Run service not found")
	}

	// 2) List SKUs (paginate) and filter
	var (
		sumCPU, nCPU float64
		sumMem, nMem float64
		sumReq, nReq float64
		sampled      []string
		pageToken    string
	)

	for {
		skuURL := fmt.Sprintf("https://cloudbilling.googleapis.com/v1/%s/skus?currencyCode=%s&key=%s",
			serviceName, url.QueryEscape(currency), url.QueryEscape(apiKey))
		if pageToken != "" {
			skuURL += "&pageToken=" + url.QueryEscape(pageToken)
		}

		resp, err := httpClient.Get(skuURL)
		if err != nil {
			return nil, err
		}
		var pg skusPage
		if err := json.NewDecoder(resp.Body).Decode(&pg); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		for _, k := range pg.Skus {
			desc := k.Description
			dl := strings.ToLower(desc)

			// Must be Cloud Run; skip commitments/free
			if !strings.Contains(desc, "Cloud Run") {
				continue
			}
			if strings.Contains(dl, "commitment") || strings.Contains(dl, "free") {
				continue
			}

			// Region match: serviceRegions contains region or "global"
			if !regionMatch(k.ServiceRegions, region) {
				continue
			}

			// First tier unit price
			p, ok := firstTierUSD(k.PricingInfo)
			if !ok || p <= 0 {
				continue
			}

			switch {
			case (strings.Contains(dl, "cpu") || strings.Contains(dl, "vcpu")) && strings.Contains(dl, "second"):
				sumCPU += p
				nCPU++
				sampled = append(sampled, k.SkuId+": "+desc)
			case (strings.Contains(dl, "memory") || strings.Contains(dl, "gib")) && strings.Contains(dl, "second"):
				sumMem += p
				nMem++
				sampled = append(sampled, k.SkuId+": "+desc)
			case strings.Contains(dl, "request"):
				sumReq += p
				nReq++
				sampled = append(sampled, k.SkuId+": "+desc)
			default:
				// ignore other Cloud Run SKUs (e.g., networking)
			}
		}

		if pg.NextPageToken == "" {
			break
		}
		pageToken = pg.NextPageToken
	}

	return &CloudRunAverages{
		Region:              region,
		Currency:            currency,
		AvgVCPUSecondUSD:    safeAvg(sumCPU, nCPU),
		AvgGiBSecondUSD:     safeAvg(sumMem, nMem),
		AvgReqPerMillionUSD: safeAvg(sumReq, nReq),
		SampledCount:        len(sampled),
		SampledSKUs:         sampled,
	}, nil
}

func regionMatch(regs []string, want string) bool {
	for _, r := range regs {
		if r == "global" || r == want {
			return true
		}
	}
	return false
}

func firstTierUSD(pis []pricingInfo) (float64, bool) {
	if len(pis) == 0 || len(pis[0].PricingExpression.TieredRates) == 0 {
		return 0, false
	}
	u := pis[0].PricingExpression.TieredRates[0].UnitPrice
	return usdFromMoney(u), true
}

func usdFromMoney(m money) float64 {
	sign := 1.0
	s := m.Units
	if strings.HasPrefix(s, "-") {
		sign = -1
		s = s[1:]
	}
	var units int64
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			continue
		}
		units = units*10 + int64(ch-'0')
	}
	return sign * (float64(units) + float64(m.Nanos)/1e9)
}

func safeAvg(sum float64, n float64) float64 {
	if n == 0 {
		return 0
	}
	// keep a sensible precision for unit prices
	return float64(int((sum/n)*1e9+0.5)) / 1e9
}
