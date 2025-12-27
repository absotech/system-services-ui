package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/system", GetSystemInfo)
	rg.GET("/services", ListServices)
	rg.GET("/services/:name/status", ServiceStatus)
	rg.GET("/services/status", AllServicesStatus)
	rg.POST("/services/:name/restart", RestartService)
}
