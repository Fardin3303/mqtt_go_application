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
	db := database.ConnectDB()
	defer db.Close()

	// Subscribe to MQTT topic and handle messages
	mqtt.InitMQTT(db)

	// Setup and start REST API server
	api.StartServer(db)

	// Graceful termination
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Exiting...")
}
