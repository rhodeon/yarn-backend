package main

import (
	"github.com/Mutay1/chat-backend/cmd/api/internal"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// load environment variables from dotenv file
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}

	// setup server configurations
	config := internal.Config{}
	config.Parse()
	if err := config.Validate(); err != nil {
		log.Fatalln(err)
	}

	// set Gin to release mode on production
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// start server
	if err := serveApp(config); err != nil {
		log.Fatalln(err)
	}
}
