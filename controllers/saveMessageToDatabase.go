package controllers

import (
	"context"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"instant-messenger-backend/database"
	"instant-messenger-backend/models"
)

func SaveMessageToDatabase(messageJSON []byte) (primitive.ObjectID, error) {
	// Decode JSON message
	var message models.Message
	err := json.Unmarshal(messageJSON, &message)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("error unmarshalling message: %w", err)
	}

	// Generate a new ObjectID
	message.ID = primitive.NewObjectID()

	// Insert the message into the database
	collection := database.GetCollection(database.DB, "messages")
	_, err = collection.InsertOne(context.TODO(), message)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("error inserting message into database: %w", err)
	}

	fmt.Printf("Message with ID %v saved to database.\n", message.ID)

	//Fungction publisher rabbitmq untuk mengirim pesan ke RabbitMQ setelah berhasil disimpan ke database

	// err = rabbitmq.PublishMessageToRabbitMQ(message)
	// if err!= nil {
	//     fmt.Printf("Error publishing message to RabbitMQ: %v\n", err)
	//     return primitive.NilObjectID, fmt.Errorf("error publishing message to RabbitMQ: %w", err)
	// }
	return message.ID, nil // No need to return the message here
}
