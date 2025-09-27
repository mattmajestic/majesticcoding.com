package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"majesticcoding.com/api/models"
	"majesticcoding.com/api/services"
)

func MetricsHandler(c *gin.Context) {
	cacheKey := "metrics:system:2m"

	// Try to get from Redis cache first (2 minutes TTL = 120 seconds)
	var metrics models.Metrics
	err := services.RedisGetJSON(cacheKey, &metrics)
	if err == nil {
		log.Printf("✅ System metrics cache HIT")
		c.JSON(http.StatusOK, metrics)
		return
	}
	log.Printf("🔍 System metrics cache MISS, collecting fresh data")

	cpuPercent, _ := cpu.Percent(0, false)
	vmStat, _ := mem.VirtualMemory()
	swapStat, _ := mem.SwapMemory()

	memFreePercent := 100 - vmStat.UsedPercent
	uptimeSeconds, _ := host.Uptime()
	uptimeHours := float64(uptimeSeconds) / 3600

	metrics = models.Metrics{
		CPUPercent:      cpuPercent,
		MemTotal:        vmStat.Total,
		MemUsed:         vmStat.Used,
		MemUsedPercent:  vmStat.UsedPercent,
		MemFree:         vmStat.Free,
		MemFreePercent:  memFreePercent,
		SwapUsedPercent: swapStat.UsedPercent,
		UptimeHours:     uptimeHours,
	}

	// Cache the metrics for 2 minutes (120 seconds)
	if err := services.RedisSetJSON(cacheKey, metrics, 120); err != nil {
		log.Printf("⚠️ Failed to cache system metrics: %v", err)
	} else {
		log.Printf("💾 Cached system metrics for 2 minutes")
	}

	c.JSON(http.StatusOK, metrics)
}
