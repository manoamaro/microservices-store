package main

import (
	"log"
	"manoamaro.github.com/auth_service/internal"
)

func main() {

	app := internal.NewApplication()
	app.RunMigrations()
	app.RegisterControllers()

	err := make(chan error)
	app.Run(err)
	log.Fatal(<-err)
}
