package routes

import (
  "github.com/gin-gonic/gin"
  "instant-messenger-backend/rabbitmq"
)

func SetupRouter(r *gin.Engine) {
  r.POST("/messages", rabbitmq.PublishMessage)
}