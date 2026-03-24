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
	var bookres BookResponse
	err = result.Scan(&bookres.Id, &bookres.Title, &bookres.Author, &bookres.Created_On_UTC)
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
	context.JSON(http.StatusOK, bookres)
}

func Book_Insert(context *gin.Context) {
	var bookreq BookRequest
	// check if body valid
	if err := context.BindJSON(&bookreq); err != nil {
		log.Printf("BindJSON book failed. Err:\n%s\n", err)
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Malformed request body.",
		})
		return
	}
	// check that title has valid value
	if len(bookreq.Title) < 1 {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Required field title not provided",
		})
		return
	}
	// insert values by db query
	bookres, err := Insert_Book(bookreq.Title, bookreq.Author)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error. Please try again later.",
		})
		return
	}

	// send back 201 & newly inserted book
	context.JSON(http.StatusCreated, bookres)
}

// TODO
func Book_Update(context *gin.Context) {
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
