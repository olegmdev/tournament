package controllers

import (
	"tournament/forms"
	"tournament/models"

	"github.com/gin-gonic/gin"
)

type TournamentController struct{}

// Announce new tournament
func (ctrl TournamentController) Announce(c *gin.Context) {
	var body forms.Tournament
	var model models.Tournament

	if c.BindJSON(&body) == nil {
		if _, err := model.Announce(body); err != nil {
			c.JSON(422, gin.H{"error": err.Error()})
			return
		}

		c.Status(200)
	}
}

// Finish tournament by ID and reward winners with relevant prizes
func (ctrl TournamentController) Result(c *gin.Context) {
	var body forms.TournamentResult
	var model models.Tournament

	if c.BindJSON(&body) == nil {
		if _, err := model.Result(body); err != nil {
			c.JSON(422, gin.H{"error": err.Error()})
			return
		}

		c.Status(200)
	}
}

// Join player (with backers) to a tournament.
func (ctrl TournamentController) Join(c *gin.Context) {
	var body forms.TournamentJoin
	var model models.Tournament

	if c.BindJSON(&body) == nil {
		if _, err := model.Join(body); err != nil {
			c.JSON(422, gin.H{"error": err.Error()})
			return
		}

		c.Status(200)
	}
}
