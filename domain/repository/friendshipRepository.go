package repository

import (
	"github.com/Mutay1/chat-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FriendshipRepository interface {
	GetFriends(id primitive.ObjectID) ([]models.Friendship, error)
	CountFriendRequest(requester, recipient models.Friend) (int64, error)
	CreateRequest(request models.Friendship) (*mongo.InsertOneResult, error)
	GetRequestSent(id primitive.ObjectID) ([]bson.M, error)
	GetRequestReceived(id primitive.ObjectID) ([]bson.M, error)
	AcceptRequest(id primitive.ObjectID) error
	DeleteRequest(id primitive.ObjectID) error
}
