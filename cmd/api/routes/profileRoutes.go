package routes

import (
	controller "github.com/Mutay1/chat-backend/cmd/api/controllers"
	"github.com/gin-gonic/gin"
)

//ProfileRoutes Function
func ProfileRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/profile", controller.UpdateProfile())
	incomingRoutes.GET("/users/profile", controller.GetProfile())
}
