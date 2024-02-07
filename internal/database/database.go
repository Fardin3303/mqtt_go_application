package database

import (
	"errors"
	"github.com/go-pg/pg/v10"
)

func ConnectDB() (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		Addr:     ":5432",
		User:     "postgres",
		Password: "password",
		Database: "mqtt_db",
	})
	if db == nil {
		return nil, errors.New("failed to connect to database")
	}
	return db, nil
}
