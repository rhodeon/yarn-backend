package routes

import (
	controller "github.com/Mutay1/chat-backend/controllers"
	"github.com/gin-gonic/gin"
)

//RequestRoutes Function
func RequestRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/request", controller.SendRequest())
	incomingRoutes.GET("/users/request/sent", controller.GetSentRequest())
	incomingRoutes.GET("/users/request/received", controller.GetReceivedRequest())
}
