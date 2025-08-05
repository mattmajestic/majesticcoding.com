package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// MetricsHandler godoc
// @Summary Metrics Status
// @Description Returns the current metrics status
// @Tags Metrics
// @Success 200 {object} map[string]interface{}
// @Router /metrics [get]
func MetricsHandler(c *gin.Context) {
	cpuPercent, _ := cpu.Percent(0, false)
	vmStat, _ := mem.VirtualMemory()
	swapStat, _ := mem.SwapMemory()

	memFreePercent := 100 - vmStat.UsedPercent

	uptimeSeconds, _ := host.Uptime()
	uptimeHours := float64(uptimeSeconds) / 3600

	c.JSON(http.StatusOK, gin.H{
		"cpu_percent":       cpuPercent,
		"mem_total":         vmStat.Total,
		"mem_used":          vmStat.Used,
		"mem_used_percent":  vmStat.UsedPercent,
		"mem_free":          vmStat.Free,
		"mem_free_percent":  memFreePercent,
		"swap_used_percent": swapStat.UsedPercent,
		"uptime_hours":      uptimeHours,
	})
}
