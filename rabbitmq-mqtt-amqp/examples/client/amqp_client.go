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

	// Declare exchange (must match the host)
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

	// Set QoS to prefetch only one message at a time
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Fatalf("Failed to set QoS: %v", err)
	}

	// Register consumer
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack (set to false for manual ack)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	fmt.Println("AMQP Client started. Waiting for messages...")
	fmt.Println("Press CTRL+C to exit")

	// Process messages
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var msg Message
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				d.Nack(false, false) // Reject and don't requeue
				continue
			}

			fmt.Printf("Received: ID=%d, Content=%s, Timestamp=%s\n",
				msg.ID, msg.Content, msg.Timestamp.Format(time.RFC3339))

			// Acknowledge message
			if err := d.Ack(false); err != nil {
				log.Printf("Failed to acknowledge message: %v", err)
			}
		}
	}()

	<-forever
}

