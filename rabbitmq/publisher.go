package rabbitmq

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"

	"instant-messenger-backend/models"
)

func PublishToDatabase(c *gin.Context) {
	// Get the request body and convert it to message
	var message models.Message
	if err := c.BindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert the message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Error marshalling message to JSON: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Publish message to the exchange
	ch, err := rabbitMQ.Channel()
	if err != nil {
		fmt.Printf("Error getting RabbitMQ channel: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer ch.Close()

	err = ch.Publish(
		"DatabaseExchange", // Exchange name
		"database_key",     // Routing key (matches queue binding)
		true,               // Mandatory (don't fail if no queue bound)
		false,              // Immediate (don't wait for ack)
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageJSON,
		},
	)
	if err != nil {
		fmt.Printf("Error publishing message to RabbitMQ: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// fmt.Printf("Message '%s' published successfully!\n", string(messageJSON))
	// No return statement needed, as the function now doesn't return anything

	c.JSON(http.StatusOK, gin.H{
        "error":    false,
        "messages": "Message send successfully, try to saving to database",
		"status": "pending",
    })
}

func PublishToClient(messageJSON []byte) {
	// Publish message to the exchange
	ch, err := rabbitMQ.Channel()
	if err != nil {
		fmt.Printf("Error getting RabbitMQ channel: %v\n", err)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer ch.Close()

	err = ch.Publish(
		"ClientExchange", // Exchange name
		"client_key",     // Routing key (matches queue binding)
		true,             // Mandatory (don't fail if no queue bound)
		false,            // Immediate (don't wait for ack)
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageJSON,
		},
	)
	if err != nil {
		fmt.Printf("Error publishing message to RabbitMQ: %v\n", err)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// fmt.Printf("Message '%s' published successfully!\n", string(messageJSON))
	// No return statement needed, as the function now doesn't return anything
}
