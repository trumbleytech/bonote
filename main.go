package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// validate env vars & init global DB
	if err := ValidateDBConfig(); err != nil {
		log.Fatal(err)
	}
	err := OpenDBConnection()
	if err != nil {
		log.Fatal(err)
	}

	// init gin
	router := gin.Default()
	router.POST("/login", LoginUser)
	router.POST("/logout", LogoutUser)
	user := router.Group("/user")
	{
		user.GET("/:id", GetUser)
		user.POST("", SaveUser)
		user.DELETE("/:id", DeleteUser)
	}
	book := router.Group("/book")
	{
		book.GET("/:id", RequireAuth, GetBook)
		book.POST("", RequireAuth, SaveBook)
		book.PUT("/:id", RequireAuth, UpdateBook)
		book.DELETE("/:id", RequireAuth, DeleteBook)
		// pull chapters by book id
		book.GET("/:id/chapters", RequireAuth, GetChaptersByBookID)
	}

	chapter := router.Group("/chapter")
	{
		chapter.GET("/:id", RequireAuth, GetChapter)
		chapter.POST("", RequireAuth, SaveChapter)
		chapter.PUT("/:id", RequireAuth, UpdateChapter)
		chapter.DELETE("/:id", RequireAuth, DeleteChapter)
		// get notes by chapter id
		chapter.GET("/:id/notes", RequireAuth, GetNotesByChapterID)
	}

	note := router.Group("/note")
	{
		note.GET("/:id", RequireAuth, GetNote)
		note.POST("", RequireAuth, SaveNote)
		note.PUT("/:id", RequireAuth, UpdateNote)
		note.DELETE("/:id", RequireAuth, DeleteNote)
	}

	// start the server
	log.Printf("Starting server on port %s...\n", PORT)
	if err := router.Run(":" + PORT); err != nil {
		log.Printf("Server failed to start. Error:\n%s", err)
	}
}
