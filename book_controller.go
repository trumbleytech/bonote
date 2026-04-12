package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetBook(context *gin.Context) {
	// id parsing first to fail fast

	// pull id from params
	id := context.Params.ByName("id")

	// convert string to id to validate input
	bookID, err := GetBookIDIntFromURL(id)
	if err != nil {
		// handle bad data
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid bookId",
		})
		return
	}
	value, exists := context.Get("usermin")
	if !exists {
		HandleInternalServerError(context)
		return
	}
	usermin := value.(UserMin)
	// query db for book by book id
	book, err := GetBookByIDDB(bookID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("Book with ID %d not found.", bookID),
			})
			return
		} else {
			go log.Printf("SQL Error:\n%s\n", err)
			HandleInternalServerError(context)
			return
		}
	}

	if !UserCanModifyResource(usermin, book.UserID) {
		HandleForbidden(context)
	} else {
		context.JSON(http.StatusOK, book)
	}
}

func UpdateBook(context *gin.Context) {
	// use URL param as id source
	id := context.Params.ByName("id")

	bookID, err := GetBookIDIntFromURL(id)
	if err != nil {
		// handle bad data
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid bookId",
		})
		return
	}
	value, exists := context.Get("usermin")
	if !exists {
		HandleInternalServerError(context)
		return
	}
	usermin := value.(UserMin)
	// query db for book by book id
	book, err := GetBookByIDDB(bookID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("Book with ID %d not found.", bookID),
			})
			return
		} else {
			go log.Printf("SQL Error:\n%s\n", err)
			HandleInternalServerError(context)
			return
		}
	}

	if !UserCanModifyResource(usermin, book.UserID) {
		HandleForbidden(context)
		return
	}
	// check body for update values
	var reqbook Book
	if err := context.ShouldBindJSON(&reqbook); err != nil {
		go log.Printf("ShouldBindJSON book failed. Err:\n%s\n", err)
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Malformed request body.",
		})
		return
	}

	// check that title has valid value
	if len(reqbook.Title) < 1 {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Required field title not provided",
		})
		return
	}

	book, err = UpdateBookDB(bookID, reqbook.Title, reqbook.Author)

	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("Book with ID %d not found.", bookID),
			})
			return
		} else {
			go log.Printf("SQL Error:\n%s\n", err)
			HandleInternalServerError(context)
			return
		}
	}

	context.JSON(http.StatusOK, book)
}

func SaveBook(context *gin.Context) {
	value, exists := context.Get("usermin")
	if !exists {
		HandleInternalServerError(context)
		return
	}
	usermin := value.(UserMin)

	var book Book
	// check if body valid
	if err := context.ShouldBindJSON(&book); err != nil {
		go log.Printf("ShouldBindJSON book failed. Err:\n%s\n", err)
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
	book, err := SaveBookDB(book.Title, book.Author, usermin.Id)
	if err != nil {
		HandleInternalServerError(context)
		return
	}

	// send back 201 & newly inserted book
	context.JSON(http.StatusCreated, book)
}

func DeleteBook(context *gin.Context) {
	// pull id from params
	id := context.Params.ByName("id")

	// convert string to id to validate input
	bookID, err := GetBookIDIntFromURL(id)
	if err != nil {
		// handle bad data
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid bookId",
		})
		return
	}

	value, exists := context.Get("usermin")
	if !exists {
		HandleInternalServerError(context)
		return
	}
	usermin := value.(UserMin)

	book, err := GetBookByIDDB(bookID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("Book with ID %d not found.", bookID),
			})
			return
		}
		HandleInternalServerError(context)
		return
	}

	// handle auth
	if !UserCanModifyResource(usermin, book.UserID) {
		HandleForbidden(context)
		return
	}

	err = DeleteBookDB(bookID)
	if err != nil {
		HandleInternalServerError(context)
		return
	}

	context.Status(http.StatusNoContent)
}

/*
	HELPER FUNCS
*/

func GetBookIDIntFromURL(urlParam string) (int, error) {
	id, err := strconv.Atoi(urlParam)
	if err != nil {
		return 0, err
	}
	if (id < 0) || (id > 2147483647) {
		return 0, errors.New("Invalid bookId")
	}
	return id, nil
}
