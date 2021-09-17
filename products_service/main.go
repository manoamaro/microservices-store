package main

import (
	"os"

	"manoamaro.github.com/products_service/internal"
)

func main() {
	mongoDBClient := internal.ConnectMongoDB(os.Getenv("MONGO_URL"))
	defer mongoDBClient.Disconnect(nil)

	internal.StartMQ(os.Getenv("AMQP_URL"))
}
