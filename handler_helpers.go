package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleForbidden(context *gin.Context) {
	context.JSON(http.StatusForbidden, gin.H{
		"message": "Action forbidden.",
	})
}
func HandleInternalServerError(context *gin.Context) {
	context.JSON(http.StatusInternalServerError, gin.H{
		"message": "Internal Server Error. Please try again later.",
	})
}

func UserCanModifyResource(usermin UserMin, resourceUserID uint) bool {
	return usermin.Role == "admin" || usermin.Id == int(resourceUserID)
}

func HandleNotLoggedIn(context *gin.Context) {
	context.JSON(http.StatusUnauthorized, gin.H{
		"message": "You are not logged in.",
	})
}
