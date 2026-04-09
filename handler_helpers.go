package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleInternalServerError(context *gin.Context) {
	context.JSON(http.StatusInternalServerError, gin.H{
		"message": "Internal Server Error. Please try again later.",
	})
}
