package main

import (
	"log"

	"github.com/mactep/agryo/server"
	"github.com/mactep/agryo/util"
)

func main() {
	server, err := server.NewServer(
		util.GetEnvOrPanic("ACCOUNT_ID"),
		util.GetEnvOrPanic("PRIVATE_KEY"),
		util.GetEnvOrPanic("DB_USERNAME"),
		util.GetEnvOrPanic("DB_PASSWORD"),
		util.GetEnvOrPanic("DB_HOST"),
		util.GetEnvOrPanic("DB_PORT"),
		util.GetEnvOrPanic("DB_NAME"),
	)
	defer server.Close()

	if err != nil {
		log.Fatalf("Failed to start the API: %v", err)
	}

	server.Run()
}
