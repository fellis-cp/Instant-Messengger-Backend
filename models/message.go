package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	ChatID      string             `bson:"chatId" json:"chatId"`
	SenderID    string             `bson:"senderId" json:"senderId"`
	Content     string             `bson:"content" json:"content"`
	SentAt      time.Time          `bson:"sentAt" json:"sentAt"`
	Attachments []string           `bson:"attachments" json:"attachments"`
}
