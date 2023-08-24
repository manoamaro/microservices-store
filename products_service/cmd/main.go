package main

import (
	"log"
	"log/slog"

	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"github.com/manoamaro/microservices-store/products_service/internal"
)

func main() {
	helpers.SetLogger()
	if err := helpers.LoadEnv(); err != nil {
		slog.Error("Error loading .env file: %s", err.Error())
	}

	application := internal.NewApplication()
	err := make(chan error)
	application.Run(err)
	log.Fatal(<-err)
}
