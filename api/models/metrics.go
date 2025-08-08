package models

type Metrics struct {
	CPUPercent      []float64 `json:"cpu_percent"`
	MemTotal        uint64    `json:"mem_total"`
	MemUsed         uint64    `json:"mem_used"`
	MemUsedPercent  float64   `json:"mem_used_percent"`
	MemFree         uint64    `json:"mem_free"`
	MemFreePercent  float64   `json:"mem_free_percent"`
	SwapUsedPercent float64   `json:"swap_used_percent"`
	UptimeHours     float64   `json:"uptime_hours"`
}
