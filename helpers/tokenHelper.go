package helper

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Mutay1/chat-backend/database"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

// GenerateTokens generates both the detailed token and refresh token
func GenerateTokens(jwtSecret string, userId string) (signedAccessToken string, signedRefreshToken string, err error) {
	// generate access token with a lifetime of 24 hours
	claims := jwt.StandardClaims{
		Subject:   userId,
		ExpiresAt: time.Now().Local().Add(24 * time.Hour).Unix(),
	}

	// generate refresh token with a lifetime of 1 week
	refreshClaims := jwt.StandardClaims{
		ExpiresAt: time.Now().Local().Add(7 * 24 * time.Hour).Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtSecret))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(jwtSecret))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

// ValidateToken validates the provided JWT.
// An error is returned if the token is invalid or expired.
func ValidateToken(jwtSecret string, signedToken string) (*jwt.StandardClaims, error) {
	// attempt to parse token
	token, err := jwt.ParseWithClaims(
		signedToken,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		},
	)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	// extract claims from token
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}

//UpdateAllTokens renews the user tokens when they login
func UpdateAllTokens(signedToken string, signedRefreshToken string, userID string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refreshToken", signedRefreshToken})

	UpdatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updatedAt", UpdatedAt})

	upsert := true
	filter := bson.M{"userID": userID}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	fmt.Println(userID)

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)
	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}

	return
}
