package api

import (
	"net/http"

	"system-services-ui/internal/systemd"

	"github.com/gin-gonic/gin"
)

func ListServices(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"services": systemd.List(),
	})
}

func RestartService(c *gin.Context) {
	name := c.Param("name")

	if err := systemd.Restart(name); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "restarted"})
}

func ServiceStatus(c *gin.Context) {
	name := c.Param("name")

	enabled, running, err := systemd.Status(name)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"service": name,
		"enabled": enabled,
		"running": running,
	})
}

func AllServicesStatus(c *gin.Context) {
	result := make(map[string]gin.H)

	for _, name := range systemd.List() {
		enabled, running, _ := systemd.Status(name)
		result[name] = gin.H{
			"enabled": enabled,
			"running": running,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"services": result,
	})
}
