package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"majesticcoding.com/api/models"
)

func TestMetricsHandler(t *testing.T) {
	router := SetupTestRouter()
	router.GET("/api/metrics", MetricsHandler)

	t.Run("should return system metrics", func(t *testing.T) {
		w := PerformRequest(router, "GET", "/api/metrics", nil)
		
		AssertJSONResponse(t, w, http.StatusOK)
		
		var metrics models.Metrics
		err := json.Unmarshal(w.Body.Bytes(), &metrics)
		assert.NoError(t, err)
		
		// Verify metrics structure and reasonable values
		assert.NotNil(t, metrics.CPUPercent)
		assert.True(t, len(metrics.CPUPercent) > 0, "CPU percent should have at least one value")
		
		// Memory values should be positive
		assert.Greater(t, metrics.MemTotal, uint64(0), "Total memory should be positive")
		assert.GreaterOrEqual(t, metrics.MemUsed, uint64(0), "Used memory should be non-negative")
		assert.GreaterOrEqual(t, metrics.MemFree, uint64(0), "Free memory should be non-negative")
		
		// Percentages should be between 0 and 100
		assert.GreaterOrEqual(t, metrics.MemUsedPercent, 0.0, "Memory used percent should be >= 0")
		assert.LessOrEqual(t, metrics.MemUsedPercent, 100.0, "Memory used percent should be <= 100")
		assert.GreaterOrEqual(t, metrics.MemFreePercent, 0.0, "Memory free percent should be >= 0")
		assert.LessOrEqual(t, metrics.MemFreePercent, 100.0, "Memory free percent should be <= 100")
		assert.GreaterOrEqual(t, metrics.SwapUsedPercent, 0.0, "Swap used percent should be >= 0")
		assert.LessOrEqual(t, metrics.SwapUsedPercent, 100.0, "Swap used percent should be <= 100")
		
		// Uptime should be positive
		assert.GreaterOrEqual(t, metrics.UptimeHours, 0.0, "Uptime should be non-negative")
		
		// Memory used + free should approximately equal total (allowing for some variance)
		totalApprox := metrics.MemUsed + metrics.MemFree
		variance := float64(metrics.MemTotal) * 0.1 // 10% variance allowed
		assert.InDelta(t, metrics.MemTotal, totalApprox, variance, "Memory used + free should approximately equal total")
	})

	t.Run("should handle multiple requests consistently", func(t *testing.T) {
		// Make multiple requests to ensure consistency
		for i := 0; i < 3; i++ {
			w := PerformRequest(router, "GET", "/api/metrics", nil)
			AssertJSONResponse(t, w, http.StatusOK)
			
			var metrics models.Metrics
			err := json.Unmarshal(w.Body.Bytes(), &metrics)
			assert.NoError(t, err)
			assert.Greater(t, metrics.MemTotal, uint64(0))
		}
	})
}

func TestMetricsHandlerResponseFormat(t *testing.T) {
	router := SetupTestRouter()
	router.GET("/api/metrics", MetricsHandler)

	w := PerformRequest(router, "GET", "/api/metrics", nil)
	AssertJSONResponse(t, w, http.StatusOK)

	// Verify all expected fields are present
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	expectedFields := []string{
		"cpu_percent", "mem_total", "mem_used", "mem_used_percent",
		"mem_free", "mem_free_percent", "swap_used_percent", "uptime_hours",
	}

	for _, field := range expectedFields {
		assert.Contains(t, response, field, "Response should contain field: %s", field)
	}
}