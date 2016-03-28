package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	go client()
	go server()

	var a string
	fmt.Scanln(&a)
}

func client() {
	conn, ch, q := getQueue()
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		q.Name, // queue
		"", // Consumer
		true, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil) // arg amqp.Table

	failOnError(err, "Fialed to register a consumer")

	for msg := range msgs {
		log.Printf("Received message with message: %s", msg.Body)
	}
}

func server() {
	conn, ch, q := getQueue()
	defer conn.Close()
	defer ch.Close()

	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body: []byte("Hello RabbitMQ"),
	}

	ch.Publishin("", q.Name, false, false, msg)
}

func getQueue() (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial("amqp:// guest@192.168.0.27:5672")
	failOnError(err, "Failes to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"Sensor", // Name of the Queue
		false, // Durable
		false, // autoDelete
		false, // esclusive
		false, // noWait
		nil)// args amqp.Table
	failOnError(err, "Failes to declare a queue")

	return conn, ch, &q
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}