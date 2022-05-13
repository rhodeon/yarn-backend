package controllers

import (
	"net/http"

	"github.com/Mutay1/chat-backend/cmd/api/internal"
	"github.com/Mutay1/chat-backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetFriends(app internal.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.GetString("uid"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UserID"})
			return
		}
		requestsLoaded, err := app.Repositories.Friendships.GetFriends(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var friendsLoaded []models.Friend
		ID := c.GetString(("uid"))
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
		c.JSON(http.StatusOK, friendsLoaded)
	}
}
