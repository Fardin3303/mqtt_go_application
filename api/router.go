package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/mqtt_go_application/pkg/models"
)

func StartServer(db *pg.DB) {
	r := gin.Default()

	r.GET("/messages", func(c *gin.Context) {
		var messages []models.MQTTMessage
		err := db.Model(&messages).Select()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, messages)
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
