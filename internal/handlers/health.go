package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"filters": gin.H{},
		"status":  "ok",
		"message": "Northwind MCP server is running!",
	})
}
