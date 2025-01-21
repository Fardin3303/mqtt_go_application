package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/go-pg/pg/v10"
	"github.com/mqtt_go_application/pkg/models"
)

// SessionIDGenerator represents a generator for unique session IDs
type SessionIDGenerator struct {
	counter int
	mu      sync.Mutex // for thread safety
}

// Define a function to attempt reconnection
func reconnect(client mqtt.Client) {
	for {
		// Attempt to connect to the MQTT broker
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Printf("Failed to reconnect to MQTT broker: %v", token.Error())
			// Wait before attempting to reconnect
			time.Sleep(5 * time.Second)
		} else {
			log.Println("Successfully reconnected to the MQTT broker")
			break
		}
	}
}

func InitMQTT(db *pg.DB) error {
	opts := mqtt.NewClientOptions().AddBroker("host.docker.internal:1883")
	opts.SetClientID("go-mqtt-example")

	client := mqtt.NewClient(opts)

	// Define a connection lost handler to handle connection losses
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("Connection lost: %v", err)
		// Attempt to reconnect here if necessary
		reconnect(client)
	})

	// Attempt to connect to the MQTT broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	// Subscribe to the MQTT topic and handle incoming messages
	topic := "charger/1/connector/1/session/1"
	log.Printf("Subscribing to MQTT topic: %s\n", topic)
	if token := client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		// Pass db as a parameter to handleMessage function
		log.Printf("Received message on topic %s: %s\n", msg.Topic(), msg.Payload())
		handleMessage(msg, db)
	}); token.Wait() && token.Error() != nil {
		log.Printf("Error subscribing to MQTT topic %s: %v\n", topic, token.Error())
		return fmt.Errorf("failed to subscribe to MQTT topic %s: %w", topic, token.Error())
	} else {
		log.Printf("Subscribed to MQTT topic: %s\n", topic)
	}

	// Start a goroutine to publish new sessions every 1 minute
	generator := NewSessionIDGenerator()
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			// Check if the client is connected before publishing messages
			if !client.IsConnected() {
				log.Println("MQTT client is not connected")
				continue
			}

			// Create a new session message
			sessionMsg := models.MQTTMessage{
				SessionID:            generator.GenerateSessionID(),
				EnergyDeliveredInKWh: 30,
				DurationInSeconds:    45,
				SessionCostInCents:   70,
				Timestamp:            time.Now(),
			}

			// Convert the message to JSON
			msgJSON, err := json.Marshal(sessionMsg)
			if err != nil {
				log.Printf("Error marshalling session message: %s", err.Error())
				continue
			}

			// Publish the message to the MQTT broker
			token := client.Publish(topic, 0, false, msgJSON)
			token.Wait()
			if token.Error() != nil {
				log.Printf("Error publishing session message: %s", token.Error())
				continue
			}

			log.Printf("Published new session message: %s", msgJSON)
		}
	}()

	return nil
}

// Function to handle incoming MQTT messages
func handleMessage(msg mqtt.Message, db *pg.DB) {
	log.Printf("Handling MQTT message: %s\n", msg.Payload())
	var mqttMsg models.MQTTMessage
	err := json.Unmarshal(msg.Payload(), &mqttMsg)
	if err != nil {
		log.Printf("Error unmarshalling MQTT message: %s", err.Error())
		return
	}

	mqttMsg.Timestamp = time.Now()

	// Insert message into the database
	fmt.Printf("Inserting message into database: %+v\n", mqttMsg)
	_, err = db.Model(&mqttMsg).Insert()
	if err != nil {
		log.Printf("Error inserting message to database: %s", err.Error())
		return
	}
	fmt.Println("Message inserted into database")

	log.Printf("Received message on topic %s: %s\n", msg.Topic(), msg.Payload())
}

// NewSessionIDGenerator creates a new session ID generator
func NewSessionIDGenerator() *SessionIDGenerator {
	return &SessionIDGenerator{
		counter: 1, // start from 1
	}
}

// GenerateSessionID generates a unique session ID
func (gen *SessionIDGenerator) GenerateSessionID() int {
	gen.mu.Lock()
	defer gen.mu.Unlock()
	id := gen.counter
	gen.counter++
	return id
}
