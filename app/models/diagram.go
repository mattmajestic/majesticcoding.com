package models

type Node struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`     // e.g. "compute", "db"
	Provider string                 `json:"provider"` // "aws", "gcp", "azure"
	Region   string                 `json:"region"`
	Config   map[string]interface{} `json:"config"` // service-specific fields
	UI       NodeUI                 `json:"ui"`
	PriceKey string                 `json:"price_key"`
}

type NodeUI struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Label string `json:"label"`
}

type Diagram struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Nodes  []Node `json:"nodes"`
}
