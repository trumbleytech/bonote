package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RequireAuth(context *gin.Context) {
	// pull usermin via session cookie
	sid, err := context.Cookie("sid")
	if err != nil {
		HandleNotLoggedIn(context)
		context.Abort()
		return
	}
	usermin, err := GetUserMinBySessionToken(sid)
	if err != nil {
		if err == sql.ErrNoRows {
			HandleNotLoggedIn(context)
			context.Abort()
			return
		}
		HandleInternalServerError(context)
		context.Abort()
		return
	}
	context.Set("usermin", usermin)
}
