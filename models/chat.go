package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	ID           		primitive.ObjectID   `bson:"_id" json:"_id"`       						// Unique chat ID
	ChatID				string				 `bson:"chatId" json:"chatId"`
	Participants 		[]string 			 `bson:"participants" json:"participants"`        	// Array of user IDs participating in the chat
	ChatType     		string               `bson:"chatType" json:"chatType"`            		// Type of chat (e.g., "private", "group")
	Name         		string               `bson:"name" json:"name"`      					// Chat name (optional for group chats)
	LastMessage	 		string				 `bson:"lastMessage" json:"lastMessage"`
	LastMessageTime    	time.Time            `bson:"lastMessageTime" json:"lastMessageTime"` 	// Timestamp of chat creation
}