package server

import (
	"net/http"

	"github.com/poprih/ur-monitor/internal/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	r := gin.Default()

	r.POST("/line/webhook", func(c *gin.Context) {
		controllers.HandleLineWebhook(c.Writer, c.Request)
	})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	r.Run(":8080")
}
