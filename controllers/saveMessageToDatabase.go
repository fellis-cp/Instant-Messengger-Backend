package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
    
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"instant-messenger-backend/database"
	"instant-messenger-backend/models"
)

func SaveMessageToDatabase(messageJSON []byte) (bool, error) {
	// Decode JSON message
	var message models.Message
	err := json.Unmarshal(messageJSON, &message)
	if err != nil {
		return false, fmt.Errorf("error unmarshalling message: %w", err)
	}

	// Generate a new ObjectID
	message.ID = primitive.NewObjectID()

	// Insert the message into the database
	collection := database.GetCollection(database.DB, "messages")
	_, err = collection.InsertOne(context.TODO(), message)
	if err != nil {
		return false, fmt.Errorf("error inserting message into database: %w", err)
	}

	fmt.Printf("Message with ID %v saved to database.\n", message.ID)

	// Update chat document with chatId
	updateResult, err := chatCollections.UpdateOne(context.TODO(), bson.M{"chatId": message.ChatID}, bson.M{
		"$set": bson.M{
			"lastMessage":     message.Content,
			"lastMessageTime": message.SentAt,
		},
	})

	if err != nil {
		log.Printf("Error updating chat: %v\n", err)
		return false, err // Propagate the error
	}

	// Check if any document was updated
	if updateResult.MatchedCount == 0 {
		log.Println("No chat document updated with chatId: ", message.ChatID)
	} else {
		log.Println("Chat updated successfully")
	}

	return true, nil // Indicate success and any encountered error
}
