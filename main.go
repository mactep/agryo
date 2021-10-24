package main

import (
	"log"

	"github.com/mactep/agryo/server"
	"github.com/mactep/agryo/util"
)

func main() {
	server, err := server.NewServer(
		util.GetEnvAndRaise("ACCOUNT_ID"),
		util.GetEnvAndRaise("PRIVATE_KEY"),
		util.GetEnvAndRaise("DB_USERNAME"),
		util.GetEnvAndRaise("DB_PASSWORD"),
		util.GetEnvAndRaise("DB_HOST"),
		util.GetEnvAndRaise("DB_PORT"),
		util.GetEnvAndRaise("DB_NAME"),
	)

	if err != nil {
		log.Fatalf("Failed to start the API: %v", err)
	}

	server.Run()
}
