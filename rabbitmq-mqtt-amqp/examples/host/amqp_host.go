package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	rabbitMQURL = "amqp://admin:admin123@localhost:5672/"
	queueName   = "amqp_queue"
	exchange    = "amqp_exchange"
	routingKey  = "amqp.routing.key"
)

type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare exchange
	err = ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %v", err)
	}

	// Declare queue
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		queueName,  // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind a queue: %v", err)
	}

	fmt.Println("AMQP Host started. Sending messages...")

	// Send messages
	for i := 1; i <= 10; i++ {
		msg := Message{
			ID:        i,
			Content:   fmt.Sprintf("AMQP Message #%d", i),
			Timestamp: time.Now(),
		}

		body, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Failed to marshal message: %v", err)
			continue
		}

		err = ch.Publish(
			exchange,   // exchange
			routingKey, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Persistent, // Make message persistent
				Body:         body,
				Timestamp:    time.Now(),
			},
		)
		if err != nil {
			log.Printf("Failed to publish a message: %v", err)
			continue
		}

		fmt.Printf("Sent: %s\n", string(body))
		time.Sleep(1 * time.Second)
	}

	fmt.Println("All messages sent successfully!")
}

