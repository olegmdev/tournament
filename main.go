package main

import (
	"flag"
	"fmt"
	"os"

	"tournament/config"
	"tournament/server"
	"tournament/db"
)

func main() {
	environment := flag.String("env", "development", "application environment")

	flag.Usage = func() {
		fmt.Println("Usage: server -env={mode}")
		os.Exit(1)
	}

	flag.Parse()

	config.Init(*environment)
	db.Init()
	server.Init()
}
