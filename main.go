package main

import (
	"log"

	config "github.com/brianwu291/repo-changes-analyzer/config"
	di "github.com/brianwu291/repo-changes-analyzer/internal/di"
)

func main() {
	config := config.NewConfig()

	di := di.NewDI(config)

	if err := di.HTTPServer.Start(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
