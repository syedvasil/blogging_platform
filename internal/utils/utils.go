package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleError(ctx *gin.Context, err error) {
	if errors.Is(err, errors.New("not allowed")) {
		ctx.JSON(http.StatusForbidden, gin.H{"err": err.Error()})
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
