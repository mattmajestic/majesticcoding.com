package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"majesticcoding.com/api/models"
)

func MetricsHandler(c *gin.Context) {
	cpuPercent, _ := cpu.Percent(0, false)
	vmStat, _ := mem.VirtualMemory()
	swapStat, _ := mem.SwapMemory()

	memFreePercent := 100 - vmStat.UsedPercent
	uptimeSeconds, _ := host.Uptime()
	uptimeHours := float64(uptimeSeconds) / 3600

	metrics := models.Metrics{
		CPUPercent:      cpuPercent,
		MemTotal:        vmStat.Total,
		MemUsed:         vmStat.Used,
		MemUsedPercent:  vmStat.UsedPercent,
		MemFree:         vmStat.Free,
		MemFreePercent:  memFreePercent,
		SwapUsedPercent: swapStat.UsedPercent,
		UptimeHours:     uptimeHours,
	}

	c.JSON(http.StatusOK, metrics)
}
