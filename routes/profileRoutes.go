package routes

import (
	"github.com/Mutay1/chat-backend/controllers"
	controller "github.com/Mutay1/chat-backend/controllers"
	"github.com/gin-gonic/gin"
)

//ProfileRoutes Function
func ProfileRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/profile", controller.UpdateProfile())
	incomingRoutes.GET("/users/profile", controllers.GetProfile())
}
