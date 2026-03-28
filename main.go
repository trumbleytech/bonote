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

	// BOOK
	router.GET("/book/:id", GetBook)
	router.POST("/book", SaveBook)
	router.PUT("/book/:id", UpdateBook)

	// CHAPTERS BY BOOK
	router.GET("/book/:id/chapters", GetChaptersByBookID)

	// CHAPTER
	router.GET("/chapter/:id", GetChapter)
	router.POST("/chapter", SaveChapter)
	router.PUT("chapter/:id", UpdateChapter)

	// NOTES BY CHAPTER
	router.GET("/chapter/:id/notes", GetNotesByChapterID)

	// NOTE
	router.GET("/note/:id", GetNote)
	router.POST("/note", SaveNote)
	router.PUT("/note/:id", UpdateNote)

	// start the server
	fmt.Printf("Starting server on port %s...\n", PORT)
	if err := router.Run(":" + PORT); err != nil {
		fmt.Printf("Server failed to start. Error:\n%s", err)
	}
}
