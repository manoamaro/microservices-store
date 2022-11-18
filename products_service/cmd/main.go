package main

import (
	"fmt"
	"github.com/manoamaro/microservices-store/commons/pkg/helpers"
	"github.com/manoamaro/microservices-store/products_service/internal"
	"log"
	"net/http"
	"time"
)

func main() {

	application := internal.NewApplication()

	r := application.SetupRoutes()

	port := helpers.GetEnv("PORT", "8080")

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
