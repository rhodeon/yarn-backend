package controllers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Mutay1/chat-backend/database"
	"github.com/Mutay1/chat-backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var friendshipCollection *mongo.Collection = database.OpenCollection(database.Client, "friendships")

type Body struct {
	Username string `json:"username"`
}

type RequestBody struct {
	ID string
}

//SendRequest generates a Friend Request
func SendRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.Friendship
		var requester models.Friend
		var recipient models.Friend
		var messages []models.Message
		body := Body{}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		err := userCollection.FindOne(ctx, bson.M{"email": c.GetString("email")}).Decode(&requester)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Username"})
			return
		}
		err = userCollection.FindOne(ctx, bson.M{"username": body.Username}).Decode(&recipient)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Username"})
			return
		}
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
		count, err := friendshipCollection.CountDocuments(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request already sent"})
			return
		}
		if recipient.ID == requester.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You can't send friend request to yourself."})
			return
		}

		requester.Messages = messages
		recipient.Messages = messages
		request.ID = primitive.NewObjectID()
		request.Recipient = recipient
		request.Requester = requester
		request.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		request.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		resultInsertionNumber, insertErr := friendshipCollection.InsertOne(ctx, request)

		if insertErr != nil {
			msg := fmt.Sprintf("Request item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

//GetSentRequest retrieves all requests sent by signed in user
func GetSentRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.GetString("uid"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UserID"})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		filter := bson.D{
			{"$and",
				bson.A{
					bson.D{{"accepted", bson.D{{"$eq", false}}}},
					bson.D{{"requester._id", id}},
				},
			},
		}
		cursor, err := friendshipCollection.Find(ctx, filter)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Username"})
			return
		}
		var requestsLoaded []bson.M
		if err = cursor.All(ctx, &requestsLoaded); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		}
		defer cancel()
		c.JSON(http.StatusOK, requestsLoaded)
	}
}

//GetReceivedRequest retrieves all requests received by signed in user
func GetReceivedRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.GetString("uid"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UserID"})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		filter := bson.D{
			{"$and",
				bson.A{
					bson.D{{"accepted", bson.D{{"$eq", false}}}},
					bson.D{{"recipient._id", id}},
				},
			},
		}
		cursor, err := friendshipCollection.Find(ctx, filter)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Username"})
			return
		}
		var requestsLoaded []bson.M
		if err = cursor.All(ctx, &requestsLoaded); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		}
		defer cancel()
		c.JSON(http.StatusOK, requestsLoaded)
	}
}

func AcceptRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		body := RequestBody{}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(body)
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			// Handle error
		}
		fmt.Println((jsonData))
		id, err := primitive.ObjectIDFromHex(body.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		_, err = friendshipCollection.UpdateOne(
			ctx,
			bson.M{"_id": id},
			bson.D{{"$set",
				bson.D{
					{"accepted", true},
				},
			}},
		)
		defer cancel()
		c.JSON(http.StatusOK, gin.H{
			"message": "Request successfully accepted",
		})
	}
}

func DeleteRequest() gin.HandlerFunc {
	return func(c *gin.Context) {

		body := RequestBody{}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id, err := primitive.ObjectIDFromHex(string(body.ID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		_, err = friendshipCollection.DeleteOne(
			ctx,
			bson.M{"_id": id},
		)
		defer cancel()
	}
}
