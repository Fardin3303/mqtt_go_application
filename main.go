package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/go-pg/pg/v10"
	"github.com/mqtt_go_application/internal/database"
	"github.com/mqtt_go_application/internal/mqtt"
	"github.com/mqtt_go_application/api"
)

func main() {
	// Connect to Postgres database
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err.Error())
	}
	defer db.Close()

	// Subscribe to MQTT topic and handle messages
	err = mqtt.InitMQTT(db)
	if err != nil {
		log.Fatalf("Error initializing MQTT: %s", err.Error())
	}

	// Setup and start REST API server
	err = api.StartServer(db)
	if err != nil {
		log.Fatalf("Error starting API server: %s", err.Error())
	}

	// Graceful termination
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Exiting...")
}
