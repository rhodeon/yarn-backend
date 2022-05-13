package routes

import (
	controller "github.com/Mutay1/chat-backend/cmd/api/controllers"
	"github.com/Mutay1/chat-backend/cmd/api/internal"
	"github.com/gin-gonic/gin"
)

//RequestRoutes Function
func RequestRoutes(app internal.Application, incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/request", controller.SendRequest(app))
	incomingRoutes.GET("/users/request/sent", controller.GetSentRequest(app))
	incomingRoutes.GET("/users/request/received", controller.GetReceivedRequest(app))
	incomingRoutes.POST("/users/request/accept", controller.AcceptRequest(app))
	incomingRoutes.POST("/users/request/delete", controller.DeleteRequest(app))
}
