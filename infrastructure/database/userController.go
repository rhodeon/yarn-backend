package database

import (
	"context"
	"errors"
	"github.com/Mutay1/chat-backend/domain/repository"
	"github.com/Mutay1/chat-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserController struct {
	Db *mongo.Database
}

const collectionUsers = "user"

// Create registers a new user, returning an error if a duplicate username or email is found.
// repository.ErrDuplicateDetails is returned if at least the username or the email already exists in the database.
func (u UserController) Create(user models.User) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// check if any pre-existing user with the same username or email exists
	count, err := u.Db.Collection(collectionUsers).CountDocuments(ctx, bson.M{
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

	if _, err := u.Db.Collection(collectionUsers).InsertOne(ctx, user); err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetByEmail retrieves an existing user via their email.
// repository.ErrRecordNotFound is returned if no qualifying user is found.
func (u UserController) GetByEmail(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// empty struct to populate with fetched user data
	foundUser := models.User{}

	err := u.Db.Collection(collectionUsers).FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&foundUser)

	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return models.User{}, repository.ErrRecordNotFound

		default:
			return models.User{}, err
		}
	}

	return foundUser, nil
}
