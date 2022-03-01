package controllers

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	helper "github.com/Mutay1/chat-backend/helpers"
	"github.com/Mutay1/chat-backend/models"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cloudName string = os.Getenv("CLOUDINARY_CLOUD_NAME")
var apiKey string = os.Getenv("CLOUDINARY_API_KEY")
var apiSecret string = os.Getenv("CLOUDINARY_API_SECRET")

func uploadAvatar(file multipart.File, ctx context.Context, uid string, fileTags string) (*uploader.UploadResult, error) {
	cld, _ := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)

	result, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: uid,
		// Split the tags by comma
		Tags: strings.Split(",", fileTags),
	})
	return result, err
}

func deleteAvatar(ctx context.Context, uid string) (*uploader.DestroyResult, error) {
	cld, _ := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)

	result, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: uid,
	})
	fmt.Println("DELETED")
	return result, err
}

func updateAvatar(file multipart.File, ctx context.Context, uid string, fileTags string) (*uploader.UploadResult, error) {
	cld, _ := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	_, err := deleteAvatar(ctx, uid)
	if err != nil {
		return nil, err
	}
	result, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: uid,
		// Split the tags by comma
		Tags: strings.Split(",", fileTags),
	})
	return result, err
}

//UpdateProfile is used to update status, about and avatar
func UpdateProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		var foundUser models.User
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		err := userCollection.FindOne(ctx, bson.M{"email": c.GetString("email")}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid User"})
			return
		}
		fileTags := c.PostForm("tags")
		status := c.PostForm("status")
		city := c.PostForm("city")
		about := c.PostForm("about")
		file, _, err := c.Request.FormFile("selectedFile")
		var result *uploader.UploadResult
		var updateObj primitive.D
		var secureURL string
		if err == nil {
			if foundUser.AvatarURL != "" {
				result, err = updateAvatar(file, c, c.GetString("uid"), fileTags)
			} else {
				result, err = uploadAvatar(file, c, c.GetString("uid"), fileTags)
			}
			if err != nil {
				c.String(http.StatusConflict, "Upload to cloudinary failed")
			}
			secureURL = result.SecureURL
		}
		updateObj = append(updateObj, bson.E{"avatarURL", secureURL})
		updateObj = append(updateObj, bson.E{"about", about})
		updateObj = append(updateObj, bson.E{"status", status})
		updateObj = append(updateObj, bson.E{"city", city})

		UpdatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updatedAt", UpdatedAt})

		upsert := true
		filter := bson.M{"email": foundUser.Email}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		_, err = userCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		id, err := primitive.ObjectIDFromHex(foundUser.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		var friend models.Friend
		err = userCollection.FindOne(ctx, bson.M{"email": c.GetString("email")}).Decode(&friend)
		defer cancel()
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid User"})
			return
		}
		data, err := helper.ToDoc(friend)
		_, err = friendshipCollection.UpdateMany(
			ctx,
			bson.M{"requester._id": id},
			bson.D{{"$set",
				bson.D{
					{"requester", data},
				},
			}},
		)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		_, err = friendshipCollection.UpdateMany(
			ctx,
			bson.M{"recipient._id": id},
			bson.D{{"$set",
				bson.D{
					{"recipient", data},
				},
			}},
		)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		defer cancel()

		c.JSON(http.StatusCreated, gin.H{
			"message": "Successfully uploaded the file",
		})
	}
}

//GetProfile returns user Profile
func GetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		var foundUser models.User
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		err := userCollection.FindOne(ctx, bson.M{"email": c.GetString("email")}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid User"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"city":      foundUser.City,
			"about":     foundUser.About,
			"status":    foundUser.Status,
			"firstName": foundUser.FirstName,
			"lastName":  foundUser.LastName,
			"avatar":    foundUser.AvatarURL,
		})
	}
}

func GetUploadedFiles(c *gin.Context) {}

func UpdateFile(c *gin.Context) {}

func DeleteFile(c *gin.Context) {}
