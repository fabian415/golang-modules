package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	amqp "github.com/rabbitmq/amqp091-go"
)

// App struct
type App struct {
	ctx context.Context

	// AMQP connections per role
	amqpStates map[string]*amqpState
	amqpMutex  sync.Mutex

	// MQTT connections per role
	mqttClients map[string]mqtt.Client
	mqttMutex   sync.Mutex

	// Message queue for frontend
	messageQueue    []MessageItem
	messageQueueMux sync.RWMutex
	maxQueueSize    int
}

// Message represents a message structure
type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// MessageItem represents a received message item
type MessageItem struct {
	Protocol  string    `json:"protocol"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

const (
	roleHost   = "host"
	roleClient = "client"
)

type amqpState struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

// ConnectionConfig represents connection configuration
type ConnectionConfig struct {
	Protocol   string `json:"protocol"`   // "amqp" or "mqtt"
	Host       string `json:"host"`       // broker host
	Port       string `json:"port"`       // broker port
	Username   string `json:"username"`   // username
	Password   string `json:"password"`   // password
	Exchange   string `json:"exchange"`   // AMQP exchange (for AMQP only)
	Queue      string `json:"queue"`      // AMQP queue (for AMQP only)
	RoutingKey string `json:"routingKey"` // AMQP routing key (for AMQP only)
	Role       string `json:"role"`       // "host" or "client"
}

// PublishConfig represents publish configuration
type PublishConfig struct {
	Protocol string `json:"protocol"` // "amqp" or "mqtt"
	Topic    string `json:"topic"`    // topic/routing key
	Message  string `json:"message"`  // message content
}

// SubscribeConfig represents subscribe configuration
type SubscribeConfig struct {
	Protocol string `json:"protocol"` // "amqp" or "mqtt"
	Topic    string `json:"topic"`    // topic/routing key
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		messageQueue: make([]MessageItem, 0),
		maxQueueSize: 100,
		amqpStates:   make(map[string]*amqpState),
		mqttClients:  make(map[string]mqtt.Client),
	}
}

func normalizeRole(role string) string {
	if role == roleClient {
		return roleClient
	}
	return roleHost
}

func (a *App) ensureAMQPState(role string) *amqpState {
	if a.amqpStates == nil {
		a.amqpStates = make(map[string]*amqpState)
	}
	state, ok := a.amqpStates[role]
	if !ok {
		state = &amqpState{}
		a.amqpStates[role] = state
	}
	return state
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Connect establishes connection to AMQP or MQTT broker using host role by default.
func (a *App) Connect(config ConnectionConfig) error {
	config.Role = normalizeRole(config.Role)
	switch config.Protocol {
	case "amqp":
		return a.connectAMQP(config)
	case "mqtt":
		return a.connectMQTT(config)
	default:
		return fmt.Errorf("unsupported protocol: %s", config.Protocol)
	}
}

// ConnectHost establishes a host-side connection explicitly.
func (a *App) ConnectHost(config ConnectionConfig) error {
	config.Role = roleHost
	return a.Connect(config)
}

// ConnectClient establishes a client-side connection explicitly.
func (a *App) ConnectClient(config ConnectionConfig) error {
	config.Role = roleClient
	return a.Connect(config)
}

// Disconnect closes the host connection for the given protocol.
func (a *App) Disconnect(protocol string) error {
	return a.disconnect(protocol, roleHost)
}

// DisconnectHost closes the host connection for the given protocol.
func (a *App) DisconnectHost(protocol string) error {
	return a.disconnect(protocol, roleHost)
}

// DisconnectClient closes the client connection for the given protocol.
func (a *App) DisconnectClient(protocol string) error {
	return a.disconnect(protocol, roleClient)
}

func (a *App) disconnect(protocol, role string) error {
	switch protocol {
	case "amqp":
		return a.disconnectAMQP(role)
	case "mqtt":
		return a.disconnectMQTT(role)
	default:
		return fmt.Errorf("unsupported protocol: %s", protocol)
	}
}

// Publish sends a message
func (a *App) Publish(config PublishConfig) error {
	switch config.Protocol {
	case "amqp":
		return a.publishAMQP(config)
	case "mqtt":
		return a.publishMQTT(config)
	default:
		return fmt.Errorf("unsupported protocol: %s", config.Protocol)
	}
}

// Subscribe subscribes to a topic
func (a *App) Subscribe(config SubscribeConfig) error {
	switch config.Protocol {
	case "amqp":
		return a.subscribeAMQP(config)
	case "mqtt":
		return a.subscribeMQTT(config)
	default:
		return fmt.Errorf("unsupported protocol: %s", config.Protocol)
	}
}

// Unsubscribe unsubscribes from a topic
func (a *App) Unsubscribe(protocol, topic string) error {
	switch protocol {
	case "amqp":
		// AMQP doesn't need explicit unsubscribe
		return nil
	case "mqtt":
		return a.unsubscribeMQTT(topic)
	default:
		return fmt.Errorf("unsupported protocol: %s", protocol)
	}
}

// GetMessages retrieves and clears all messages from the queue
func (a *App) GetMessages() []MessageItem {
	a.messageQueueMux.Lock()
	defer a.messageQueueMux.Unlock()

	messages := make([]MessageItem, len(a.messageQueue))
	copy(messages, a.messageQueue)
	a.messageQueue = a.messageQueue[:0] // Clear the queue

	return messages
}

// addMessage adds a message to the queue
func (a *App) addMessage(protocol, content string) {
	a.messageQueueMux.Lock()
	defer a.messageQueueMux.Unlock()

	item := MessageItem{
		Protocol:  protocol,
		Content:   content,
		Timestamp: time.Now(),
	}

	a.messageQueue = append(a.messageQueue, item)

	// Keep queue size limited
	if len(a.messageQueue) > a.maxQueueSize {
		a.messageQueue = a.messageQueue[len(a.messageQueue)-a.maxQueueSize:]
	}
}

// AMQP functions
func (a *App) connectAMQP(config ConnectionConfig) error {
	role := normalizeRole(config.Role)
	a.amqpMutex.Lock()
	defer a.amqpMutex.Unlock()

	state := a.ensureAMQPState(role)
	if state.channel != nil {
		state.channel.Close()
		state.channel = nil
	}
	if state.conn != nil {
		state.conn.Close()
		state.conn = nil
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", config.Username, config.Password, config.Host, config.Port)

	conn, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %v", err)
	}

	// Declare exchange
	if config.Exchange != "" {
		err = ch.ExchangeDeclare(
			config.Exchange,
			"topic",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			ch.Close()
			conn.Close()
			return fmt.Errorf("failed to declare exchange: %v", err)
		}
	}

	// Declare queue
	if config.Queue != "" {
		_, err = ch.QueueDeclare(
			config.Queue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			ch.Close()
			conn.Close()
			return fmt.Errorf("failed to declare queue: %v", err)
		}

		// Bind queue to exchange
		if config.Exchange != "" && config.RoutingKey != "" {
			err = ch.QueueBind(
				config.Queue,
				config.RoutingKey,
				config.Exchange,
				false,
				nil,
			)
			if err != nil {
				ch.Close()
				conn.Close()
				return fmt.Errorf("failed to bind queue: %v", err)
			}
		}
	}

	state.conn = conn
	state.channel = ch
	if config.Exchange != "" {
		state.exchange = config.Exchange
	} else {
		state.exchange = "amqp_exchange" // Default
	}
	return nil
}

func (a *App) disconnectAMQP(role string) error {
	role = normalizeRole(role)
	a.amqpMutex.Lock()
	defer a.amqpMutex.Unlock()

	state := a.ensureAMQPState(role)
	if state.channel != nil {
		state.channel.Close()
		state.channel = nil
	}
	if state.conn != nil {
		state.conn.Close()
		state.conn = nil
	}
	state.exchange = ""
	return nil
}

func (a *App) publishAMQP(config PublishConfig) error {
	a.amqpMutex.Lock()
	defer a.amqpMutex.Unlock()

	state := a.ensureAMQPState(roleHost)
	if state.channel == nil {
		return fmt.Errorf("AMQP host not connected")
	}

	// Parse message or use as-is
	var body []byte
	var err error

	// Try to parse as JSON, if fails use as plain text
	var msg Message
	if err := json.Unmarshal([]byte(config.Message), &msg); err != nil {
		// Not JSON, use as plain text
		body = []byte(config.Message)
	} else {
		// Valid JSON, use it
		body, err = json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %v", err)
		}
	}

	// Use stored exchange and topic as routing key
	exchange := state.exchange
	if exchange == "" {
		exchange = "amqp_exchange" // Fallback
	}
	routingKey := config.Topic

	err = state.channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}

func (a *App) subscribeAMQP(config SubscribeConfig) error {
	a.amqpMutex.Lock()
	defer a.amqpMutex.Unlock()

	state := a.ensureAMQPState(roleClient)
	if state.channel == nil {
		return fmt.Errorf("AMQP client not connected")
	}
	ch := state.channel

	// Use topic as queue name, or create a queue
	queueName := config.Topic
	if queueName == "" {
		queueName = "amqp_queue"
	}

	// Declare queue
	_, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	// Bind to exchange if exists
	exchange := "amqp_exchange"
	if err := ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err == nil {
		ch.QueueBind(
			queueName,
			config.Topic,
			exchange,
			false,
			nil,
		)
	}

	// Set QoS
	err = ch.Qos(1, 0, false)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %v", err)
	}

	// Consume messages
	msgs, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %v", err)
	}

	// Start goroutine to handle messages
	go func(channel *amqp.Channel) {
		for d := range msgs {
			fmt.Println(string(d.Body))
			a.addMessage("amqp", string(d.Body))

			// Acknowledge message
			if channel != nil {
				d.Ack(false)
			}
		}
	}(ch)

	return nil
}

// MQTT functions
func (a *App) connectMQTT(config ConnectionConfig) error {
	role := normalizeRole(config.Role)
	a.mqttMutex.Lock()
	defer a.mqttMutex.Unlock()

	broker := fmt.Sprintf("tcp://%s:%s", config.Host, config.Port)
	clientID := fmt.Sprintf("mqtt_client_%d", time.Now().Unix())

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetConnectTimeout(10 * time.Second)
	opts.SetPingTimeout(5 * time.Second)
	opts.SetKeepAlive(30 * time.Second)

	opts.OnConnect = func(client mqtt.Client) {
		log.Println("MQTT connected successfully")
	}

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Printf("MQTT connection lost: %v", err)
	}

	client := mqtt.NewClient(opts)
	token := client.Connect()

	connected := token.WaitTimeout(15 * time.Second)
	if !connected {
		return fmt.Errorf("MQTT connection timeout")
	}

	if token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %v", token.Error())
	}

	if existing, ok := a.mqttClients[role]; ok && existing != nil && existing.IsConnected() {
		existing.Disconnect(250)
	}
	a.mqttClients[role] = client
	return nil
}

func (a *App) disconnectMQTT(role string) error {
	role = normalizeRole(role)
	a.mqttMutex.Lock()
	defer a.mqttMutex.Unlock()

	if client, ok := a.mqttClients[role]; ok && client != nil {
		if client.IsConnected() {
			client.Disconnect(250)
		}
		delete(a.mqttClients, role)
	}
	return nil
}

func (a *App) publishMQTT(config PublishConfig) error {
	a.mqttMutex.Lock()
	defer a.mqttMutex.Unlock()

	client, ok := a.mqttClients[roleHost]
	if !ok || client == nil || !client.IsConnected() {
		return fmt.Errorf("MQTT host not connected")
	}

	// Parse message or use as-is
	var body []byte

	// Try to parse as JSON, if fails use as plain text
	var msg Message
	if err := json.Unmarshal([]byte(config.Message), &msg); err != nil {
		// Not JSON, use as plain text
		body = []byte(config.Message)
	} else {
		// Valid JSON, use it
		body, err = json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %v", err)
		}
	}

	token := client.Publish(config.Topic, 1, false, body)
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("failed to publish message: %v", token.Error())
	}

	return nil
}

func (a *App) subscribeMQTT(config SubscribeConfig) error {
	a.mqttMutex.Lock()
	defer a.mqttMutex.Unlock()

	client, ok := a.mqttClients[roleClient]
	if !ok || client == nil || !client.IsConnected() {
		return fmt.Errorf("MQTT client not connected")
	}

	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		a.addMessage("mqtt", string(msg.Payload()))
	}

	token := client.Subscribe(config.Topic, 1, messageHandler)
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic: %v", token.Error())
	}

	return nil
}

func (a *App) unsubscribeMQTT(topic string) error {
	a.mqttMutex.Lock()
	defer a.mqttMutex.Unlock()

	client, ok := a.mqttClients[roleClient]
	if !ok || client == nil || !client.IsConnected() {
		return fmt.Errorf("MQTT client not connected")
	}

	token := client.Unsubscribe(topic)
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("failed to unsubscribe from topic: %v", token.Error())
	}

	return nil
}
