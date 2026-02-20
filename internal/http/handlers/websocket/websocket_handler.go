package websocket

import (
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/gin-gonic/gin"
)

func WebsocketHandler(c *gin.Context) {
	userID, exists := c.Get("id")
	if !exists {
		c.AbortWithStatus(401)
		return
	}

	if app.MelodyInstance == nil {
		c.JSON(500, gin.H{"error": "Websocket service not available"})
		return
	}
	if c.Request.Method != "GET" {
		c.JSON(405, gin.H{"error": "Method not allowed"})
		return
	}
	err := app.MelodyInstance.HandleRequestWithKeys(c.Writer, c.Request, map[string]interface{}{"userID": userID})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to handle websocket request"})
		return
	}
}
