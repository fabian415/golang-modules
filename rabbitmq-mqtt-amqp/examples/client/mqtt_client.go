package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	broker   = "tcp://localhost:1883"
	clientID = "mqtt_client_subscriber"
	topic    = "mqtt/topic/messages"
	username = "admin"
	password = "admin123"
)

type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var message Message
	if err := json.Unmarshal(msg.Payload(), &message); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	fmt.Printf("Received from topic '%s': ID=%d, Content=%s, Timestamp=%s\n",
		msg.Topic(), message.ID, message.Content, message.Timestamp.Format(time.RFC3339))
}

func main() {
	fmt.Printf("Connecting to MQTT broker at %s...\n", broker)

	// MQTT connection options
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetConnectTimeout(10 * time.Second) // Set connection timeout
	opts.SetPingTimeout(5 * time.Second)     // Set ping timeout
	opts.SetKeepAlive(30 * time.Second)      // Set keep alive interval
	opts.SetDefaultPublishHandler(messageHandler)
	
	// Add connection status callbacks for debugging
	opts.OnConnect = func(client mqtt.Client) {
		fmt.Println("✓ Successfully connected to MQTT broker")
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Printf("✗ Connection lost: %v", err)
	}

	// Create MQTT client
	client := mqtt.NewClient(opts)

	// Connect to broker with timeout
	fmt.Println("Attempting to connect...")
	token := client.Connect()
	
	// Wait for connection with timeout
	connected := token.WaitTimeout(15 * time.Second)
	if !connected {
		log.Fatalf("Failed to connect to MQTT broker: connection timeout after 15 seconds")
	}
	
	if token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}
	
	defer client.Disconnect(250)

	// Subscribe to topic
	if token := client.Subscribe(topic, 1, messageHandler); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic: %v", token.Error())
	}

	fmt.Printf("MQTT Client subscribed to topic '%s'. Waiting for messages...\n", topic)
	fmt.Println("Press CTRL+C to exit")

	// Wait for interrupt signal to gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("\nShutting down...")

	// Unsubscribe
	if token := client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		log.Printf("Failed to unsubscribe: %v", token.Error())
	}
}

