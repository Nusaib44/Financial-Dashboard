package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func SendInternalError(c *gin.Context) {
	SendError(c, http.StatusInternalServerError, "Internal Server Error")
}
