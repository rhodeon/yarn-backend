package database

import (
	"context"
	"github.com/Mutay1/chat-backend/domain/repository"
	"github.com/Mutay1/chat-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserController struct {
	Db *mongo.Database
}

const collectionUser = "user"

// Create registers a new user, returning an error if a duplicate username or email is found.
func (u UserController) Create(user models.User) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// check if any pre-existing user with the same username or email exists
	count, err := u.Db.Collection(collectionUser).CountDocuments(ctx, bson.M{
		"$or": bson.A{
			bson.M{"username": user.Username},
			bson.M{"email": user.Email},
		},
	})

	if err != nil {
		return models.User{}, err
	}

	if count > 0 {
		return models.User{}, repository.ErrDuplicateDetails
	}

	if _, err := u.Db.Collection(collectionUser).InsertOne(ctx, user); err != nil {
		return models.User{}, err
	}

	return user, nil
}
