package main

import (
	"fmt"
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

	book := router.Group("/book")
	{
		book.GET("/:id", GetBook)
		book.POST("", SaveBook)
		book.PUT("/:id", UpdateBook)
		book.DELETE("/:id", DeleteBook)
		// pull chapters by book id
		book.GET("/:id/chapters", GetChaptersByBookID)
	}

	chapter := router.Group("/chapter")
	{
		chapter.GET("/:id", GetChapter)
		chapter.POST("", SaveChapter)
		chapter.PUT("/:id", UpdateChapter)
		chapter.DELETE("/:id", DeleteChapter)
		// get notes by chapter id
		chapter.GET("/:id/notes", GetNotesByChapterID)
	}

	note := router.Group("/note")
	{
		note.GET("/:id", GetNote)
		note.POST("", SaveNote)
		note.PUT("/:id", UpdateNote)
		note.DELETE("/:id", DeleteNote)
	}

	// start the server
	fmt.Printf("Starting server on port %s...\n", PORT)
	if err := router.Run(":" + PORT); err != nil {
		log.Printf("Server failed to start. Error:\n%s", err)
	}
}
