package main

import (
	"context"
	"log"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestConnectionRabbitMQ(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		t.Fatal(err.Error())
	}

	defer conn.Close()

	log.Print("connected to rabbitmq server")

	chanel, err := conn.Channel()

	if err != nil {
		t.Fatal(err.Error())
	}

	defer chanel.Close()

	queue, err := chanel.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		t.Fatal(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "test 12345"

	err = chanel.PublishWithContext(ctx,
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)

	if err != nil {
		t.Fatal(err.Error())
	}

	log.Printf("send %v to rabbitmq server", body)
}
