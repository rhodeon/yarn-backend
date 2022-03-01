package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Friendship is the model that governs all notes objects retrived or inserted into the DB
type Friendship struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	Requester Friend             `json:"requester" bson:"requester" validate:"required"`
	Recipient Friend             `json:"recipient" bson:"recipient" validate:"required"`
	Accepted  bool               `json:"accepted" bson:"accepted"`
}

//Friend is model that connects friend request to users
type Friend struct {
	ID        primitive.ObjectID `bson:"_id"`
	FirstName *string            `json:"firstName" validate:"required,min=2,max=100" bson:"firstName"`
	LastName  *string            `json:"lastName" validate:"required,min=2,max=100" bson:"lastName"`
	Username  *string            `json:"username" validate:"required"`
	AvatarURL string             `json:"avatarURL" bson:"avatarURL"`
	Status    string             `json:"status" bson:"status"`
	About     string             `json:"about" bson:"about"`
	City      string             `json:"city" bson:"city"`
	Messages  []Message          `json:"messages" bson:"messages"`
	Archived  bool               `json:"archived,omitempty" bson:"archived"`
	Favorite  bool               `json:"favorite,omitempty" bson:"favorite"`
	Blocked   bool               `json:"blocked,omitempty" bson:"blocked"`
}

// Message is return msg
type Message struct {
	Sender      string    `json:"sender,omitempty" bson:"sender"`
	RecipientID string    `json:"recipientID" bson:"recipientID"`
	Content     string    `json:"content,omitempty" bson:"content"`
	Delivered   bool      `json:"delivered" bson:"delivered"`
	Read        bool      `json:"read" bson:"read"`
	CreatedAt   time.Time `json:"createdAt,omitempty" bson:"createdAt"`
	MessageType string    `json:"messageType"`
	UpdateType  string    `json:"updateType"`
}
