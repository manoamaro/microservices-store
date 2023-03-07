package main

import (
	"github.com/manoamaro/microservices-store/order_service/internal"
	"log"
)

func main() {
	app := internal.NewApplication()
	err := make(chan error)
	app.Run(err)
	log.Fatal(<-err)
}
