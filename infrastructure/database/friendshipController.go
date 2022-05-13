package database

import (
	"context"
	"errors"
	"time"

	"github.com/Mutay1/chat-backend/domain/repository"
	"github.com/Mutay1/chat-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FriendshipController struct {
	Db *mongo.Database
}

func NewFriendshipController(db *mongo.Database) repository.FriendshipRepository {
	return FriendshipController{
		Db: db,
	}
}

const friendshipCollection = "friendships"

// get user friends
func (f FriendshipController) GetFriends(id primitive.ObjectID) ([]models.Friendship, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"accepted", bson.D{{"$eq", true}}}},
				bson.M{
					"$or": []interface{}{
						bson.M{"requester._id": id},
						bson.M{"recipient._id": id},
					},
				},
			},
		},
	}

	cursor, err := f.Db.Collection(friendshipCollection).Find(ctx, filter)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	var requestsLoaded []models.Friendship
	if err = cursor.All(ctx, &requestsLoaded); err != nil {
		return nil, errors.New("an error occurred")
	}

	return requestsLoaded, nil
}

// count number of friend requests
// filter is hard coded here as only one instance of the query is found
func (f FriendshipController) CountFriendRequest(requester, recipient models.Friend) (int64, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.D{
		{"$or",
			bson.A{
				bson.M{
					"$and": []interface{}{
						bson.M{"requester.username": recipient.Username},
						bson.M{"recipient.username": requester.Username},
					},
				},
				bson.M{
					"$and": []interface{}{
						bson.M{"requester.username": requester.Username},
						bson.M{"recipient.username": recipient.Username},
					},
				},
			},
		},
	}

	return f.Db.Collection(friendshipCollection).CountDocuments(ctx, filter)
}

// create new friendship request
func (f FriendshipController) CreateRequest(request models.Friendship) (*mongo.InsertOneResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	return f.Db.Collection(friendshipCollection).InsertOne(ctx, request)
}

// get request sent by user
func (f FriendshipController) GetRequestSent(id primitive.ObjectID) ([]bson.M, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"accepted", bson.D{{"$eq", false}}}},
				bson.D{{"requester._id", id}},
			},
		},
	}

	cursor, err := f.Db.Collection(friendshipCollection).Find(ctx, filter)
	if err != nil {
		return nil, errors.New("invalid Username")
	}

	var requestsLoaded []bson.M
	if err = cursor.All(ctx, &requestsLoaded); err != nil {
		return nil, errors.New("an error occurred")
	}

	return requestsLoaded, nil
}

// get requests received by user
func (f FriendshipController) GetRequestReceived(id primitive.ObjectID) ([]bson.M, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"accepted", bson.D{{"$eq", false}}}},
				bson.D{{"recipient._id", id}},
			},
		},
	}
	cursor, err := f.Db.Collection(friendshipCollection).Find(ctx, filter)
	if err != nil {
		return nil, errors.New("invalid Username")
	}

	var requestsLoaded []bson.M
	if err = cursor.All(ctx, &requestsLoaded); err != nil {
		return nil, errors.New("an error occurred")
	}

	return requestsLoaded, nil
}

// update friendship to signify acceptance
func (f FriendshipController) AcceptRequest(id primitive.ObjectID) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	_, err := f.Db.Collection(friendshipCollection).UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{{"$set",
			bson.D{
				{"accepted", true},
			},
		}},
	)
	return err
}

// delete friendship instance
func (f FriendshipController) DeleteRequest(id primitive.ObjectID) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	_, err := f.Db.Collection(friendshipCollection).DeleteOne(
		ctx,
		bson.M{"_id": id},
	)
	return err
}
