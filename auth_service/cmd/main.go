package main

import (
	"github.com/manoamaro/microservices-store/auth_service/internal"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"golang.org/x/exp/slog"
)

func main() {

	helpers.SetLogger()

	if err := helpers.LoadEnv(); err != nil {
		slog.Error("Error loading .env file: %s", err.Error())
	}

	app := internal.NewApplication()
	app.RunMigrations()

	err := make(chan error)
	app.Run(err)
	slog.Error("Error running application: %s", <-err)
}
