package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         		primitive.ObjectID 		`bson:"_id" json:"_id"`            				// Unique user ID
	Email      		string             		`bson:"email" json:"email"`                    	// User's email address
	Password   		string             		`bson:"password" json:"password"`               // User's password (hashed and secured)
	ChatList		[]primitive.ObjectID 	`bson:"chatList" json:"chatList"`				// Array of chat IDs
	Name       		string             		`bson:"name" json:"name"`                     	// User's full name
	ProfilePicture 	string             		`bson:"profilePicture" json:"profilePicture"` 	// URL or path to the user's profile picture
	NIM				string             		`bson:"nim" json:"nim"`                       	// User's NIM
	Faculty    		string             		`bson:"faculty" json:"faculty"`        			// User's Faculty
	Department 		string            		`bson:"department" json:"department"`     		// User's Department
	Subjects   		[]string           		`bson:"subject" json:"subject"`        			// Array of user's subjects
	CreatedAt  		time.Time          		`bson:"createdAt" json:"createdAt"`      		// Timestamp of user creation
}