package server

import (
	"github.com/Mengkzhaoyun/gostream/src/version"
	"github.com/gin-gonic/gin"
)

// Health endpoint returns a 500 if the server state is unhealthy.
func Health(c *gin.Context) {
	c.String(200, "")
}

// Version endpoint returns the server version and build information.
func Version(c *gin.Context) {
	c.JSON(200, gin.H{
		"source":  "https://github.com/Mengkzhaoyun/gostream",
		"version": version.Version.String(),
	})
}
