package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//creating test versions of models before db is setup

var (
	chapter = Chapter{
		1,
		0,
		"Chapter Name",
		1,
	}

	note = Note{
		1,
		"Name of note",
		"Some content for a note",
		time.Now().UTC(),
		time.Now().UTC(),
		1,
	}
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
	router.GET("/note/:id", func(context *gin.Context) {
		// pull id from params
		id := context.Params.ByName("id")

		// convert string to id to validate input
		c_id, err := strconv.Atoi(id)
		if err != nil {
			// handle bad data
			context.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid note id",
				"params":  context.Params,
			})
			return
		}
		fmt.Printf("note method hit. id = %d\n", c_id)
		context.JSON(http.StatusOK, note)
	})

	// start the server
	fmt.Printf("Starting server on port %s...\n", PORT)
	if err := router.Run(":" + PORT); err != nil {
		fmt.Printf("Server failed to start. Error:\n%s", err)
	}
}
