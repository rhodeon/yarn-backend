package routes

import (
	controller "github.com/Mutay1/chat-backend/cmd/api/controllers"
	"github.com/Mutay1/chat-backend/cmd/api/internal"
	"github.com/gin-gonic/gin"
)

//UserRoutes function
func UserRoutes(app internal.Application, incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controller.SignUp(app))
	incomingRoutes.POST("/users/login", controller.Login(app))
	incomingRoutes.POST("/users/refresh-token", controller.RefreshToken(app))
}
