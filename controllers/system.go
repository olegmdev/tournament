package controllers

import (
	"tournament/models"

	"github.com/gin-gonic/gin"
)

type SystemController struct{}

// Truncates database tables (had not be implemented bc of risks)
func (ctrl SystemController) Reset(c *gin.Context) {
	if err := models.Reset(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(200)
}
