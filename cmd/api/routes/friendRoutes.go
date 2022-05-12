package routes

import (
	controller "github.com/Mutay1/chat-backend/cmd/api/controllers"
	"github.com/gin-gonic/gin"
)

//FriendRoutes Function
func FriendRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/friends", controller.GetFriends())
}
