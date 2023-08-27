package handlers

import "github.com/gin-gonic/gin"

func CreateSegment(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
