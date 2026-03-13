package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Book_Get(context *gin.Context) {
	// pull id from params
	id := context.Params.ByName("id")

	// convert string to id to validate input
	book_id, err := strconv.Atoi(id)
	if (err != nil) || ((book_id < 0) || (book_id > 2147483647)) {
		// handle bad data
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid book Id",
		})
		return
	}

	// query db for book by book id
	result := Query_Book_By_Id(book_id)
	var book Book
	err = result.Scan(&book.Id, &book.Title, &book.Author)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("No book found with Id %d", book_id),
			})
			return
		} else {
			log.Printf("SQL Error:\n%s\n", err)
			context.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error. Please try again later.",
			})
			return
		}
	}
	context.JSON(http.StatusOK, book)
}

func Chapter_Get(context *gin.Context) {
	// pull id from params
	id := context.Params.ByName("id")

	// convert string to id to validate input
	c_id, err := strconv.Atoi(id)
	if err != nil {
		// handle bad data
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Chapter Id",
		})
		return
	}
	fmt.Printf("chapter method hit. id = %d\n", c_id)
	context.JSON(http.StatusOK, chapter)
}
