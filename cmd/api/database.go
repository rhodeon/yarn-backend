package main

import (
	"context"
	"github.com/Mutay1/chat-backend/cmd/api/internal"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// connectClient returns the database for the server
func openDb(config internal.Config) (*mongo.Database, error) {
	client, err := connectClient(config)
	if err != nil {
		return nil, err
	}

	db := client.Database(config.Db.Name)
	return db, nil
}

// connectClient returns a connected database client
func connectClient(config internal.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Db.Uri))
	if err != nil {
		return nil, err
	}

	// verify database connection
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
