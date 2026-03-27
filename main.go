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
	router.POST("/book", SaveBook)
	router.GET("/book/:id", GetBook)

	// CHAPTER
	router.GET("/book/:id/chapters", GetChaptersByBookID)
	router.GET("/chapter/:id", GetChapter)
	router.POST("/chapter", SaveChapter)

	// NOTE
	router.GET("/chapter/:id/notes", GetNotesByChapterID)
	router.GET("/note/:id", GetNote)
	router.POST("/note", SaveNote)

	// start the server
	fmt.Printf("Starting server on port %s...\n", PORT)
	if err := router.Run(":" + PORT); err != nil {
		fmt.Printf("Server failed to start. Error:\n%s", err)
	}
}
