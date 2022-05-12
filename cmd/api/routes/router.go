package routes

import (
	"fmt"
	"github.com/Mutay1/chat-backend/cmd/api/internal"
	middleware "github.com/Mutay1/chat-backend/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func Router(app internal.Application) *gin.Engine {
	gin.EnableJsonDecoderDisallowUnknownFields()
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))
	UserRoutes(router)
	WsRoutes(router)

	router.Use(middleware.Authentication())
	ProfileRoutes(router)
	RequestRoutes(router)
	FriendRoutes(router)

	// API-1
	router.GET("/api-1", func(c *gin.Context) {
		fmt.Println(c.Get("email"))
		c.JSON(200, gin.H{"success": "Access granted for api-1"})

	})

	// API-2
	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	return router
}
