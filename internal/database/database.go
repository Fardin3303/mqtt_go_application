package database

import (
	"github.com/go-pg/pg/v10"
)

func ConnectDB() *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     ":5432",
		User:     "postgres",
		Password: "password",
		Database: "mqtt_db",
	})
	return db
}
