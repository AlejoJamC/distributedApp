// Queue Utils
package qutils

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func GetChannel(url string)( *amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to establish connection to message broker")
	ch, err := conn.Channel()
	failOnError(err, "Failed to get channel for connection")

	return conn, ch
}

func GetQueue(name string, ch *amqp.Channel) *amqp.Queue {
	q, err := ch.QueueDeclare(
		name, // name
		false, // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil) // args amqp.Table

	failOnError(err, "Failed to declare Queue")

	return &q
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}