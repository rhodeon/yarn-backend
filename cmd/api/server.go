package main

import (
	"fmt"
	"github.com/Mutay1/chat-backend/cmd/api/internal"
	"github.com/Mutay1/chat-backend/cmd/api/routes"
	"github.com/Mutay1/chat-backend/controllers"
	"log"
	"net/http"
	"time"
)

// serveApp launches the server and handles its shutdown
func serveApp(config internal.Config) error {
	// launch WebSocket server manager
	go controllers.Manager.Start()

	app := internal.Application{
		Config: config,
	}

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Port),
		Handler:      routes.Router(app),
		IdleTimeout:  1 * time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// start server
	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	log.Println("stopped server")
	return nil
}
