package main

import (
	"database/sql"
	"fmt"
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

func UpdateBook(context *gin.Context) {
	// use URL param as id source
	id := context.Params.ByName("id")

	bookID, err := strconv.Atoi(id)
	if (err != nil) || ((bookID < 0) || (bookID > 2147483647)) {
		// handle bad data
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid bookId",
		})
		return
	}

	// check body for update values
	var book Book
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

	book, err = UpdateBookDB(bookID, book.Title, book.Author)

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

func DeleteBook(context *gin.Context) {
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

	err = DeleteBookDB(bookID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error. Please try again later.",
		})
		return
	}

	context.Status(http.StatusNoContent)
}
