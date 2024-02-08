package database

import (
	"errors"
	"fmt"
	"github.com/mqtt_go_application/pkg/models"
	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/pg/v10"
)

func ConnectDB() (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		Addr:     "host.docker.internal:5432",
		User:     "postgres_user",
		Password: "postgres_password",
		Database: "enersense",
	})
	if db == nil {
		return nil, errors.New("failed to connect to database")
	}
	fmt.Println("Connected to database")
	return db, nil
}

func CreateTables(db *pg.DB) error {
	models := []interface{}{(*models.MQTTMessage)(nil)}
	fmt.Println("Creating tables...")

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}