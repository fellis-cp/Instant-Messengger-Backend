package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
    ID          primitive.ObjectID  `bson:"_id" json:"_id"`
    ChatID      string              `bson:"chatID" json:"chatID"`
    SenderID    string              `bson:"senderID" json:"senderID"`
    Content     string              `bson:"content" json:"content"`
    SentAt      time.Time           `bson:"sentAt" json:"sentAt"`
    Attachments []string            `bson:"attachments" json:"attachments"`
}
