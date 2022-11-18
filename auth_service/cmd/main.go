package main

import (
	"github.com/manoamaro/microservices-store/auth_service/internal"
	"log"
)

func main() {

	app := internal.NewApplication()
	app.RunMigrations()
	app.RegisterControllers()

	err := make(chan error)
	app.Run(err)
	log.Fatal(<-err)
}
