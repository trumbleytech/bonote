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

func Chapter_Get(context *gin.Context) {
	// pull id from params
	id := context.Params.ByName("id")

	// convert string to id to validate input
	chapter_id, err := strconv.Atoi(id)
	if err != nil {
		// handle bad data
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Chapter Id",
		})
		return
	}
	// query db for chapter by chapter id
	result := Query_Chapter_By_Id(chapter_id)
	var chapter Chapter
	err = result.Scan(&chapter.Id, &chapter.Book_id, &chapter.Number, &chapter.Name, &chapter.Created_On_UTC, &chapter.Modified_On_UTC)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("No chapter found with Id %d", chapter_id),
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
	context.JSON(http.StatusOK, chapter)

}

func Chapter_Insert(context *gin.Context) {
	var chapter Chapter
	// check if body valid
	if err := context.BindJSON(&chapter); err != nil {
		log.Printf("BindJSON book failed. Err:\n%s\n", err)
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "malformed request body",
		})
		return
	}

	if err := Validate_Chapter_Insert(chapter.Book_id, chapter.Number, chapter.Name); err != nil {
		log.Printf("Chapter validation failed. Err:\n%s\n", err)
		context.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// insert values by db query
	err := Insert_Chapter(chapter.Book_id, chapter.Number, chapter.Name)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error. Please try again later.",
		})
		return
	}

	// send back 204
	context.Status(http.StatusNoContent)
}

// helper func to check if chapter number is valid
func Valid_chapter_number(chapter_number uint16) bool {
	if chapter_number > 0 {
		return true
	}
	return false
}

// helper func to check valid book id
func Valid_book_id(book_id uint) bool {
	if book_id > 0 {
		return true
	}
	return false
}

// consolidated logic to validate required fields for chapter insert
func Validate_Chapter_Insert(book_id uint, number *uint16, name string) error {
	// require book id rel
	if !Valid_book_id(book_id) {
		return errors.New("Invalid book_id.")
	}

	// name or chapter must be supplied
	if len(name) < 1 {
		// check chapter number if no name supplied
		if number != nil {
			if !Valid_chapter_number(*number) {
				return errors.New("Invalid Chapter number.")
			}
		} else {
			return errors.New("Either Chapter Name or Chapter Number is required.")
		}
	}

	return nil
}
