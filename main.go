package main

import (
	"instant-messenger-backend/database"
	"instant-messenger-backend/rabbitmq"
	"instant-messenger-backend/routes"
	"instant-messenger-backend/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	database.ConnectDB()
	go websocket.InitGorillaWebsocket()
	rabbitmq.ConnectToRabbitMQ()
	routes.SetupRouter(r)
	r.Run(":5000")
}
