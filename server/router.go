package server

import (
  "github.com/gin-gonic/gin"
  "tournament/controllers"
)

func Router() *gin.Engine {
  router := gin.Default()

	v1 := router.Group("v1")
	{
    playerGroup := v1.Group("players")
    {
      player := new(controllers.PlayerController)
      playerGroup.POST("/fund", player.Fund)
      playerGroup.POST("/take", player.Take)
      playerGroup.GET("/:id/balance", player.Balance)
    }

    tournamentGroup := v1.Group("tournaments")
    {
      tournament := new(controllers.TournamentController)
      tournamentGroup.POST("/announce", tournament.Announce)
      tournamentGroup.POST("/result", tournament.Result)
      tournamentGroup.POST("/join", tournament.Join)
    }

    systemGroup := v1.Group("system")
    {
      system := new(controllers.SystemController)
      systemGroup.DELETE("/reset", system.Reset)
    }
	}

	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(404)
	})

	return router
}
