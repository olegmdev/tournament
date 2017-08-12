package controllers

import (
	"tournament/forms"
	"tournament/models"

	"github.com/gin-gonic/gin"
)

type PlayerController struct{}

// Fund adds points to a player balance
// If player not exists, it will be created with that points amount
func (ctrl PlayerController) Fund(c *gin.Context) {
	var body forms.Player
	var model models.Player

	if c.BindJSON(&body) == nil {
		if _, err := model.Fund(body); err != nil {
			c.JSON(422, gin.H{"error": err.Error()})
			return
		}

		c.Status(200)
	}
}

// Take simply takes specified amount of points from a player.
// In case player doesn't have proper amount of points, it'll be rejected
func (ctrl PlayerController) Take(c *gin.Context) {
	var body forms.Player
	var model models.Player

	if c.BindJSON(&body) == nil {
		if _, err := model.Take(body); err != nil {
			c.JSON(422, gin.H{"error": err.Error()})
			return
		}

		c.Status(200)
	}
}

// Balance retrieve current player balance
func (ctrl PlayerController) Balance(c *gin.Context) {
	var player models.Player

	result, err := player.Get(c.Param("id"))
	if err != nil {
		c.JSON(422, gin.H{"error": err.Error()})
		return
	}

	player = result.(models.Player)
	c.JSON(200, gin.H{"id": player.ID, "balance": player.Balance})
}
