package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetBook(context *gin.Context) {
	// pull id from params
	id := context.Params.ByName("id")

	// convert string to id to validate input
	bookID, err := strconv.Atoi(id)
	if (err != nil) || ((bookID < 0) || (bookID > 2147483647)) {
		// handle bad data
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid bookId",
		})
		return
	}

	// query db for book by book id
	book, err := GetBookByID(bookID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("Book with ID %d not found.", bookID),
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

func SaveBook(context *gin.Context) {
	var book Book
	// check if body valid
	if err := context.ShouldBindJSON(&book); err != nil {
		log.Printf("ShouldBindJSON book failed. Err:\n%s\n", err)
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Malformed request body.",
		})
		return
	}
	// check that title has valid value
	if len(book.Title) < 1 {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Required field title not provided",
		})
		return
	}

	// insert values by db query
	book, err := SaveBookDB(book.Title, book.Author)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error. Please try again later.",
		})
		return
	}

	// send back 201 & newly inserted book
	context.JSON(http.StatusCreated, book)
}

// TODO
func UpdateBook(context *gin.Context) {
	// read body as byte array
	barray, err := io.ReadAll(context.Request.Body)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error. Please try again later.",
		})
		return
	}

	// create str str map to unmarshall body into
	reqbodymap := make(map[string]string)

	// unmarshall byte array into map
	json.Unmarshal(barray, &reqbodymap)
	log.Printf("Book update map: %v+", reqbodymap)

}
