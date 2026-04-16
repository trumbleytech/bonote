package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUser(context *gin.Context) {
	usermin, err := GetUserMinFromContext(context)
	if err != nil {
		HandleInternalServerError(context)
	}
	id := context.Params.ByName("id")

	userID, err := strconv.Atoi(id)
	if err != nil {
		// handle bad data
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid User ID.",
		})
		return
	}

	User, err := GetUserByIDDB(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("User with ID %d not found.", userID),
			})
			return
		} else {
			log.Printf("SQL Error:\n%s\n", err)
			HandleInternalServerError(context)
			return
		}
	}

	if !UserCanAccessResource(usermin, int(User.Id)) {
		HandleForbidden(context)
		return
	}
	context.JSON(http.StatusOK, User)
}

func RegisterUser(context *gin.Context) {
	var user User
	// check that req binds to model
	if err := context.ShouldBindJSON(&user); err != nil {
		log.Printf("ShouldBindJSON user failed. Err:\n%s\n", err)
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Malformed request body.",
		})
		return
	}

	if len(user.Username) < 1 || len(user.Email) < 1 || len(user.Password) < 1 {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Username, Email, and Password are required",
		})
		return
	}

	password_bytes := []byte(user.Password)
	if len(password_bytes) > 72 {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Password too long.",
		})
		return
	}
	hashed_password, err := bcrypt.GenerateFromPassword(password_bytes, 10)
	if err != nil {
		log.Printf("Error hashing user password. Err: %s\n", err)
		HandleInternalServerError(context)
		return
	}
	user, err = SaveUserDB(user.Username, user.Email, string(hashed_password))
	if err != nil {
		log.Printf("User Save DB failed. Err: %s\n", err)
		HandleInternalServerError(context)
		return
	}
	context.JSON(http.StatusCreated, user)

}

func DeleteUser(context *gin.Context) {
	usermin, err := GetUserMinFromContext(context)
	if err != nil {
		HandleInternalServerError(context)
		return
	}
	err = DeleteUserDB(usermin.Id)
	if err != nil {
		HandleInternalServerError(context)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully.",
	})

}

func LoginUser(context *gin.Context) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var body LoginRequest
	// validate user data
	if err := context.ShouldBindJSON(&body); err != nil {
		log.Printf("ShouldBindJSON login body failed. Err:\n%s\n", err)
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Malformed request body.",
		})
		return
	}

	// handle missing fields
	if !ValidLoginRequestBody(body.Username, body.Password) {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Username and Password are required to login.",
		})
		return
	}

	user, err := GetUserByUsernameDB(body.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid username or password.",
			})

		} else {
			HandleInternalServerError(context)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid username or password.",
		})
		return
	}

	// username and password valid, generate session
	rawToken := GenerateRawToken()
	hashToken := HashToken(rawToken)
	// validate token for 48 hours
	expiresAt := GenerateSessionExpiresAtTime()

	// save db vals and handle error
	if err = CreateNewSession(user.Id, hashToken, expiresAt); err != nil {
		HandleInternalServerError(context)
		return
	}

	// send session token back to client
	context.SetCookie(
		"sid",
		rawToken,
		172800,
		"/",
		"localhost",
		false,
		true,
	)
	context.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func LogoutUser(context *gin.Context) {
	// check if user has sid cookie
	sid, err := context.Cookie("sid")
	// if sid doesn't exist, user is not logged in
	if err != nil {
		ClearSessionID(context)
		// not logged in, end state is reached
		context.Status(http.StatusNoContent)
		return
	}

	tokenHash := HashToken(sid)
	err = DeleteSessionByTokenHashDB(tokenHash)
	if err != nil {
		go log.Printf("Error deleting session by token hash. Err: %s\n", err)
	}
	ClearSessionID(context)
	context.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully.",
	})
}

/*
HELPER FUNCS
*/

func ValidLoginRequestBody(username, password string) bool {
	return len(username) > 0 && len(password) > 0
}

func ClearSessionID(context *gin.Context) {
	context.SetCookie(
		"sid",
		"",
		-1,
		"/",
		"localhost",
		false,
		true,
	)
}

func GenerateSessionExpiresAtTime() time.Time {
	return time.Now().UTC().Add(48 * time.Hour)
}

func GetUserMinBySessionToken(sessionToken string) (UserMin, error) {
	tokenHash := HashToken(sessionToken)
	usermin, err := GetUserMinBySessionTokenHashDB(tokenHash)
	if err != nil {
		return UserMin{}, err
	}
	return usermin, err

}

func GenerateRawToken() string {
	return rand.Text()
}

func HashToken(token string) string {
	bytes := sha256.Sum256([]byte(token))
	return hex.EncodeToString(bytes[:])
}

func GetBooksByUserID(context *gin.Context) {
	usermin, err := GetUserMinFromContext(context)
	if err != nil {
		HandleInternalServerError(context)
	}

	books, err := GetBooksByUserIDDB(usermin.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("No books found with userId %d", usermin.Id),
			})
			return
		}
		HandleInternalServerError(context)
		return
	}
	context.JSON(http.StatusOK, books)
}
