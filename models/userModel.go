package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//User is the model that governs all notes objects retrived or inserted into the DB
type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	FirstName    *string            `json:"firstName" validate:"required,min=2,max=100" bson:"firstName"`
	LastName     *string            `json:"lastName" validate:"required,min=2,max=100" bson:"lastName"`
	Password     *string            `json:"Password" validate:"required,min=6"`
	Email        *string            `json:"email" validate:"email,required"`
	Username     *string            `json:"username" validate:"required"`
	Token        *string            `json:"token"`
	RefreshToken *string            `json:"refreshToken" bson:"refreshToken"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
	UserID       string             `json:"userID" bson:"userID"`
	AvatarURL    string             `json:"avatarURL" bson:"avatarURL"`
	Status       string             `json:"status" bson:"status"`
	About        string             `json:"about" bson:"about"`
	City         string             `json:"city" bson:"city"`
}
