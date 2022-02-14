package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Mutay1/chat-backend/database"
	"github.com/Mutay1/chat-backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var requestCollection *mongo.Collection = database.OpenCollection(database.Client, "request")

type body struct {
	Username string `json:"username"`
}

//SendRequest generates a Friend Request
func SendRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.Request
		var requester models.RequestUser
		var recipient models.RequestUser
		body := body{}
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

		if recipient.Username == requester.Username {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You can't send Friend Request to yourself."})
			return
		}
		request.ID = primitive.NewObjectID()
		request.Recipient = recipient
		request.Requester = requester
		request.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		request.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		resultInsertionNumber, insertErr := requestCollection.InsertOne(ctx, request)
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
		cursor, err := requestCollection.Find(ctx, bson.D{{"requester._id", id}})
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
		cursor, err := requestCollection.Find(ctx, bson.D{{"recipient._id", id}})
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
