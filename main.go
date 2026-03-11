package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	// init gin
	router := gin.Default()

	// define basic get method for book
	router.GET("/book/:id", func(context *gin.Context) {
		// pull id from params
		id := context.Params.ByName("id")

		// convert string to id to validate input
		c_id, err := strconv.Atoi(id)
		if err != nil {
			// handle bad data
			context.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid book id",
				"params":  context.Params,
			})
			return
		}
		fmt.Printf("method hit. id = %d", c_id)
		context.JSON(http.StatusOK, gin.H{
			"message": "sending back something",
			"params":  context.Params,
		})
	})

	// start the server
	fmt.Printf("Starting server on port %s...\n", PORT)
	if err := router.Run(":" + PORT); err != nil {
		fmt.Printf("Server failed to start. Error:\n%s", err)
	}
}
