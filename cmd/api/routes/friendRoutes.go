package routes

import (
	controller "github.com/Mutay1/chat-backend/cmd/api/controllers"
	"github.com/Mutay1/chat-backend/cmd/api/internal"
	"github.com/gin-gonic/gin"
)

//FriendRoutes Function
func FriendRoutes(app internal.Application, incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/friends", controller.GetFriends(app))
}
