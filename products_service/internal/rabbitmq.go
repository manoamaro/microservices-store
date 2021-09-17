package internal

import (
	"log"

	"github.com/streadway/amqp"
)

var connection *amqp.Connection

func ConnectRabbitMQ(url string) *amqp.Connection {
	conn, err := amqp.Dial(url)
	FailOnError(err, "Failed to connect to RabbitMQ")
	return conn
}

func StartMQ(url string) {
	forever := make(chan bool)

	conn := ConnectRabbitMQ(url)
	ch, err := conn.Channel()
	FailOnError(err, "")

	mainQ, err := ch.QueueDeclare("Products", true, false, false, false, nil)
	FailOnError(err, "")

	productMsgs, err := ch.Consume(mainQ.Name, "", true, false, false, false, amqp.Table{})
	FailOnError(err, "")

	go func() {
		for msg := range productMsgs {
			handleMessage(ch, &mainQ, &msg)
		}
	}()

	<-forever
}

func handleMessage(ch *amqp.Channel, q *amqp.Queue, msg *amqp.Delivery) {
	log.Printf("message received: %s", msg.Body)
}
