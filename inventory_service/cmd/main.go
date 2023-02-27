package main

import (
	"log"

	"github.com/manoamaro/microservices-store/inventory_service/internal"
)

func main() {
	app := internal.NewApplication()
	app.RunMigrations()
	app.RegisterControllers()
	err := make(chan error)
	app.Run(err)
	log.Fatal(<-err)
}
