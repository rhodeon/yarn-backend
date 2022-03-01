package routes

import (
	controller "github.com/Mutay1/chat-backend/controllers"

	"github.com/gin-gonic/gin"
)

//UserRoutes function
func WsRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/ws", controller.WsHandler)
	incomingRoutes.GET("/pong", controller.Pong())
}
