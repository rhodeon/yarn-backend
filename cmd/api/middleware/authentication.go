package middleware

import (
	"errors"
	"fmt"
	"github.com/Mutay1/chat-backend/cmd/api/internal"
	"github.com/Mutay1/chat-backend/domain/repository"
	"net/http"

	helper "github.com/Mutay1/chat-backend/helpers"

	"github.com/gin-gonic/gin"
)

// Authentication validates the provided JWT and authenticates users.
func Authentication(app internal.Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Vary", "Authorization")

		// check for the existence of the Authorization header
		clientToken := ctx.Request.Header.Get("Authorization")
		if clientToken == "" {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": fmt.Sprintf("no Authorization header provided")},
			)
			return
		}

		// validate JWT if it exists
		claims, err := helper.ValidateToken(app.Config.JwtSecret, clientToken)
		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": err.Error()},
			)
			return
		}

		// retrieve associated user
		user, err := app.Repositories.Users.GetById(claims.Subject)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrRecordNotFound):
				ctx.AbortWithStatusJSON(
					http.StatusUnprocessableEntity,
					gin.H{"error": "user not found"},
				)

			default:
				helper.HandleInternalServerError(ctx, err)
			}

			return
		}

		// set user ID and email in context for further use
		ctx.Set("uid", user.UserID)
		ctx.Set("email", *user.Email)
		ctx.Next()
	}
}
