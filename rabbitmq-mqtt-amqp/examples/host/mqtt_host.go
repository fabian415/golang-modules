package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	broker   = "tcp://localhost:1883"
	clientID = "mqtt_host_client"
	topic    = "mqtt/topic/messages"
	username = "admin"
	password = "admin123"
)

type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	fmt.Println("Starting MQTT Host...")
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

	fmt.Println("MQTT Host connected. Sending messages...")

	// Send messages
	for i := 1; i <= 10; i++ {
		msg := Message{
			ID:        i,
			Content:   fmt.Sprintf("MQTT Message #%d", i),
			Timestamp: time.Now(),
		}

		body, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Failed to marshal message: %v", err)
			continue
		}

		token := client.Publish(topic, 1, false, body)
		token.Wait()

		if token.Error() != nil {
			log.Printf("Failed to publish message: %v", token.Error())
			continue
		}

		fmt.Printf("Sent to topic '%s': %s\n", topic, string(body))
		time.Sleep(1 * time.Second)
	}

	fmt.Println("All messages sent successfully!")
}

