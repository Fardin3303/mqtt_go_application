package mqtt

import (
	"encoding/json"
	"log"
	"time"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/go-pg/pg/v10"
	"github.com/mqtt_go_application/pkg/models"
)

func InitMQTT(db *pg.DB) error {
	opts := mqtt.NewClientOptions().AddBroker("tcp://mqtt.eclipse.org:1883")
	opts.SetClientID("go-mqtt-example")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}
	defer client.Disconnect(250)

	topic := "charger/1/connector/1/session/1"
	if token := client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		var mqttMsg models.MQTTMessage
		err := json.Unmarshal(msg.Payload(), &mqttMsg)
		if err != nil {
			log.Printf("Error unmarshalling MQTT message: %s", err.Error())
			return
		}

		mqttMsg.Timestamp = time.Now()

		err = db.Insert(&mqttMsg)
		if err != nil {
			log.Printf("Error inserting message to database: %s", err.Error())
			return
		}

		log.Printf("Received message on topic %s: %s\n", msg.Topic(), msg.Payload())
	}); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to MQTT topic %s: %w", topic, token.Error())
	}
	defer client.Unsubscribe(topic)

	return nil
}
