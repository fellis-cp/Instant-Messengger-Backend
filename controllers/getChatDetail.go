package controllers

import (
	"context"
	"instant-messenger-backend/database"
	"instant-messenger-backend/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var messageCollections *mongo.Collection = database.GetCollection(database.DB, "messages")

func GetChatDetail(c *gin.Context) {
	// Get chatId from URL parameter
	chatId := c.Query("chatId")

	if chatId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Missing required parameter: chatId",
		})
		return
	}

	log.Printf("Extracted chatId: %s\n", chatId)

	// Get All message with chatId matches And Sorting Asc Based SentAt field
	var messageList []models.Message
	curMessage, err := messageCollections.Find(context.TODO(), bson.M{"chatId": chatId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	defer curMessage.Close(context.TODO())

	// Check for empty results
	if !curMessage.Next(context.TODO()) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "No chat found with chatId: " + chatId,
		})
		return
	}

	// Use curMessage.All to decode all documents at once
	if err := curMessage.All(context.TODO(), &messageList); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":    false,
		"messages": messageList,
	})
}
