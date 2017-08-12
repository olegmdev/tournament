package server

import "tournament/config"

func Init() {
  options := config.GetConfig()

  router := Router()
  router.Run(options.GetString("server.port"))
}
