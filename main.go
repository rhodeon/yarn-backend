package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Mutay1/chat-backend/controllers"
	middleware "github.com/Mutay1/chat-backend/middlewares"
	routes "github.com/Mutay1/chat-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	go controllers.Manager.Start()
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://192.168.43.236:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))
	routes.UserRoutes(router)
	routes.WsRoutes(router)

	router.Use(middleware.Authentication())
	routes.ProfileRoutes(router)
	routes.RequestRoutes(router)
	routes.FriendRoutes(router)

	// API-2
	router.GET("/api-1", func(c *gin.Context) {
		fmt.Println(c.Get("email"))
		c.JSON(200, gin.H{"success": "Access granted for api-1"})

	})

	// API-1
	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})
	router.Run(":" + port)
}
