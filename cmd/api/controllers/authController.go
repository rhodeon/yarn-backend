package controllers

import (
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

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var userCollection = database.OpenCollection(database.Client, "user")
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
		msg = fmt.Sprintf("invalid user credentials")
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
		accessToken, refreshToken, _ := helper.GenerateTokens(app.Config.JwtSecret, user.UserID)
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
			"token":          accessToken,
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

//Login authenticates a single user.
func Login(app internal.Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User

		if err := ctx.BindJSON(&user); err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
		}

		// retrieve user from repository
		foundUser, err := app.Repositories.Users.GetByEmail(*user.Email)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrRecordNotFound):
				ctx.AbortWithStatusJSON(
					http.StatusUnprocessableEntity,
					gin.H{"error": "invalid user credentials"},
				)

			default:
				helper.HandleInternalServerError(ctx, err)
			}

			return
		}

		// check password
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if passwordIsValid != true {
			ctx.AbortWithStatusJSON(
				http.StatusUnprocessableEntity,
				gin.H{"error": msg},
			)
			return
		}

		// generate and update user tokens
		accessToken, refreshToken, _ := helper.GenerateTokens(app.Config.JwtSecret, foundUser.UserID)

		if err = app.Repositories.Users.UpdateRefreshToken(foundUser.UserID, refreshToken); err != nil {
			helper.HandleInternalServerError(ctx, err)
			return
		}

		h, _ := time.ParseDuration("24h")
		ctx.JSON(http.StatusOK, gin.H{
			"token":          accessToken,
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

// RefreshToken refreshes both the access and refresh tokens of a user.
func RefreshToken(app internal.Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User
		if err := ctx.BindJSON(&user); err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
			return
		}

		// retrieve user associated with refresh token
		foundUser, err := app.Repositories.Users.GetByRefreshToken(*user.RefreshToken)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrRecordNotFound):
				ctx.AbortWithStatusJSON(
					http.StatusUnprocessableEntity,
					gin.H{"error": "invalid or expired refresh token"},
				)

			default:
				helper.HandleInternalServerError(ctx, err)
			}

			return
		}

		// generate and update user tokens
		accessToken, refreshToken, _ := helper.GenerateTokens(app.Config.JwtSecret, foundUser.UserID)

		if err = app.Repositories.Users.UpdateRefreshToken(foundUser.UserID, refreshToken); err != nil {
			helper.HandleInternalServerError(ctx, err)
			return
		}

		h, _ := time.ParseDuration("24h")
		ctx.JSON(http.StatusOK, gin.H{
			"token":          accessToken,
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
