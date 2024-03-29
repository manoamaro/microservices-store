package main

import (
	"log"
	"log/slog"

	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"github.com/manoamaro/microservices-store/order_service/internal"
)

func main() {
	helpers.SetLogger()
	if err := helpers.LoadEnv(); err != nil {
		slog.Error("Error loading .env file: %s", err.Error())
	}

	app := internal.NewApplication()
	err := make(chan error)
	app.Run(err)
	log.Fatal(<-err)
}
