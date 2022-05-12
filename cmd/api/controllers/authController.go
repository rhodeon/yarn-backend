package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/Mutay1/chat-backend/cmd/api/internal"
	"github.com/Mutay1/chat-backend/domain/repository"
	"log"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/Mutay1/chat-backend/database"

	helper "github.com/Mutay1/chat-backend/helpers"
	"github.com/Mutay1/chat-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

//HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

//VerifyPassword checks the input password while verifying it with the passward in the DB.
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or passowrd is incorrect")
		check = false
	}

	return check, msg
}

//SignUp creates a user account
func SignUp(app internal.Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// validate request
		var user models.User
		if err := ctx.BindJSON(&user); err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
			return
		}

		if err := validate.Struct(user); err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
			return
		}

		// set user details
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, user.UserID)
		user.Token = &token
		user.RefreshToken = &refreshToken
		user.Status = "Hello There! Connect with me on Yarn!"
		password := HashPassword(*user.Password)
		user.Password = &password

		// create new user in repository
		newUser, err := app.Repositories.Users.Create(user)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrDuplicateDetails):
				ctx.AbortWithStatusJSON(
					http.StatusUnprocessableEntity,
					gin.H{"error": "the email or username already exists"},
				)

			default:
				helper.HandleInternalServerError(ctx, err)
			}

			return
		}

		h, _ := time.ParseDuration("24h")
		ctx.JSON(http.StatusOK, gin.H{
			"token":          token,
			"refreshToken":   refreshToken,
			"expirationTime": h.Milliseconds(),
			"userID":         newUser.UserID,
			"profile": gin.H{
				"city":      newUser.City,
				"about":     newUser.About,
				"status":    newUser.Status,
				"firstName": newUser.FirstName,
				"lastName":  newUser.LastName,
				"avatar":    newUser.AvatarURL,
			},
		})
	}
}

//Login is the api used to get a single user
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or passowrd is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, foundUser.UserID)

		helper.UpdateAllTokens(token, refreshToken, foundUser.UserID)
		h, _ := time.ParseDuration("24h")
		c.JSON(http.StatusOK, gin.H{
			"token":          token,
			"refreshToken":   refreshToken,
			"expirationTime": h.Milliseconds(),
			"userID":         foundUser.UserID,
			"profile": gin.H{
				"city":      foundUser.City,
				"about":     foundUser.About,
				"status":    foundUser.Status,
				"firstName": foundUser.FirstName,
				"lastName":  foundUser.LastName,
				"avatar":    foundUser.AvatarURL,
			},
		})
	}
}

//RefreshToken api is used to refresh user token
func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		var foundUser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"refreshToken": *user.RefreshToken}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Token"})
			return
		}
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, foundUser.UserID)
		helper.UpdateAllTokens(token, refreshToken, foundUser.UserID)
		h, _ := time.ParseDuration("24h")
		c.JSON(http.StatusOK, gin.H{
			"token":          token,
			"refreshToken":   refreshToken,
			"expirationTime": h.Milliseconds(),
			"userID":         foundUser.UserID,
			"profile": gin.H{
				"city":      foundUser.City,
				"about":     foundUser.About,
				"status":    foundUser.Status,
				"firstName": foundUser.FirstName,
				"lastName":  foundUser.LastName,
				"avatar":    foundUser.AvatarURL,
			},
		})
	}
}
