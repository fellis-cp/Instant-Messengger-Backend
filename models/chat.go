package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	ID           primitive.ObjectID   `bson:"_id" json:"_id"`       					// Unique chat ID
	Participants []primitive.ObjectID `bson:"participants" json:"participants"`        	// Array of user IDs participating in the chat
	ChatType     string               `bson:"chatType" json:"chatType"`            		// Type of chat (e.g., "private", "group")
	Name         string               `bson:"name" json:"name"`      					// Chat name (optional for group chats)
	CreatedAt    time.Time            `bson:"createdAt" json:"createdAt"` 				// Timestamp of chat creation
}