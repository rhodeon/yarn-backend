package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Mutay1/chat-backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetFriends() gin.HandlerFunc {
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
		cursor, err := friendshipCollection.Find(ctx, filter)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var requestsLoaded []models.Friendship
		var friendsLoaded []models.Friend
		ID := c.GetString(("uid"))
		if err = cursor.All(ctx, &requestsLoaded); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		}
		for _, user := range requestsLoaded {
			if user.Requester.ID.Hex() != ID {

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				friendsLoaded = append(friendsLoaded, user.Requester)
			} else {
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				friendsLoaded = append(friendsLoaded, user.Recipient)
			}
		}
		defer cancel()
		c.JSON(http.StatusOK, friendsLoaded)
	}
}
