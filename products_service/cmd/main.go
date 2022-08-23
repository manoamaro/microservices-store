package main

import (
	"log"
	"manoamaro.github.com/products_service/internal"
	"net/http"
	"time"
)

func main() {

	application := internal.NewApplication()

	r := application.SetupRoutes()

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
