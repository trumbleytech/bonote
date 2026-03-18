package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// validate env vars & init global DB
	if err := Validate_db_config(); err != nil {
		log.Fatal(err)
	}
	err := Open_DB_Connection()
	if err != nil {
		log.Fatal(err)
	}

	// init gin
	router := gin.Default()

	// define basic get method for book
	router.GET("/book/:id", Book_Get)
	router.GET("/chapter/:id", Chapter_Get)
	router.GET("/note/:id", Note_Get)

	// start the server
	fmt.Printf("Starting server on port %s...\n", PORT)
	if err := router.Run(":" + PORT); err != nil {
		fmt.Printf("Server failed to start. Error:\n%s", err)
	}
}
