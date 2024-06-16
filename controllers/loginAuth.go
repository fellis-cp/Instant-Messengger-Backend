package controllers

import (
	"errors"
	"instant-messenger-backend/database"
	"instant-messenger-backend/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var usersCollections *mongo.Collection = database.GetCollection(database.DB, "users")

func LoginAuth(c *gin.Context) {

	var loginData struct {
		NIM      string `json:"nim"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid data",
		})
		return
	}

	if len(loginData.NIM) != 11 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Nim Harus 11 karakter",
		})
		return
	}

	var user models.User
	filter := bson.M{"nim": loginData.NIM}

	err := usersCollections.FindOne(c, filter).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "User tidak ditemukan",
			})
			return
		}

		log.Printf("Error finding user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Internal server error",
		})
	}

	if user.Password != loginData.Password {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Password salah",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  user,
	})
}
