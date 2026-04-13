package main

import (
	"errors"
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

func UserCanAccessResource(usermin UserMin, resourceUserID int) bool {
	if usermin.Role == "admin" || usermin.Id == resourceUserID {
		return true
	}
	return false
}

func HandleNotLoggedIn(context *gin.Context) {
	context.JSON(http.StatusUnauthorized, gin.H{
		"message": "You are not logged in.",
	})
}

func GetUserMinFromContext(context *gin.Context) (UserMin, error) {
	value, exists := context.Get("usermin")
	if !exists {
		return UserMin{}, errors.New("usermin does not exist")
	}
	return value.(UserMin), nil
}
