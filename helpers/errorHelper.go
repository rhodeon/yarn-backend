package helper

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// HandleInternalServerError logs the error and sends
// a generic 500 error response to the client.
func HandleInternalServerError(ctx *gin.Context, err error) {
	log.Printf("internal server error: %s", err.Error())
	ctx.AbortWithStatusJSON(
		http.StatusInternalServerError,
		gin.H{"error": "internal server error"},
	)
}
