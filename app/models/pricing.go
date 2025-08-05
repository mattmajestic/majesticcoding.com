package models

type Pricing struct {
	Key       string  `json:"key"`        // e.g., "aws:compute:t3.micro:us-east-1"
	Service   string  `json:"service"`    // e.g., "compute"
	Provider  string  `json:"provider"`   // "aws", "gcp", "azure"
	Region    string  `json:"region"`     // "us-east-1"
	UnitPrice float64 `json:"unit_price"` // e.g., 0.023 per hour
	Currency  string  `json:"currency"`   // e.g., "USD"
}
