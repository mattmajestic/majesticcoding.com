package handlers

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func DeployIACHandler(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider required"})
		return
	}

	dir := fmt.Sprintf("./iac/%s", provider)
	cmd := exec.Command("terraform", "init")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": string(out)})
		return
	}

	cmd = exec.Command("terraform", "apply", "-auto-approve")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": string(out)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Deployed to %s successfully!", provider)})
}
