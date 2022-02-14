package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Request is the model that governs all notes objects retrived or inserted into the DB
type Request struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	Requester RequestUser        `json:"requester" bson:"requester" validate:"required"`
	Recipient RequestUser        `json:"recipient" bson:"recipient" validate:"required"`
	Accepted  bool               `json:"accepted" bson:"accepted"`
}

//RequestUser is model that connects friend request to users
type RequestUser struct {
	ID        primitive.ObjectID `bson:"_id"`
	FirstName *string            `json:"firstName" validate:"required,min=2,max=100" bson:"firstName"`
	LastName  *string            `json:"lastName" validate:"required,min=2,max=100" bson:"lastName"`
	Username  *string            `json:"username" validate:"required"`
	AvatarURL string             `json:"avatarURL" bson:"avatarURL"`
	Status    string             `json:"status" bson:"status"`
	About     string             `json:"about" bson:"about"`
	City      string             `json:"city" bson:"city"`
}
