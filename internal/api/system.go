package api

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

func run(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)
	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = &out
	err := c.Run()
	return strings.TrimSpace(out.String()), err
}

func GetSystemInfo(c *gin.Context) {
	uptime, _ := run("uptime", "-p")
	load, _ := run("uptime")
	name, _ := run("hostname")

	c.JSON(200, gin.H{
		"hostname": name,
		"uptime":   uptime,
		"load":     load,
	})
}
