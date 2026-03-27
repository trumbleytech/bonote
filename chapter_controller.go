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

func GetChapter(context *gin.Context) {
	// pull id from params
	id := context.Params.ByName("id")

	// convert string to id to validate input
	chapterID, err := strconv.Atoi(id)
	if err != nil {
		// handle bad data
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Chapter ID.",
		})
		return
	}
	// query db for chapter by chapter id
	chapter, err := GetChapterByID(chapterID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("Chapter with ID %d not found.", chapterID),
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

func GetChaptersByBookID(context *gin.Context) {
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

	chapters, err := GetChaptersByBookIDDB(bookID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("No chapters found with bookId %d", bookID),
			})
			return
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error. Please try again later.",
			})
			return
		}
	}

	context.JSON(http.StatusOK, chapters)

}

func SaveChapter(context *gin.Context) {
	var chapter Chapter
	// check if body valid
	if err := context.ShouldBindJSON(&chapter); err != nil {
		log.Printf("ShouldBindJSON chapter failed. Err:\n%s\n", err)
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Malformed request body.",
		})
		return
	}

	if err := ValidateChapterInsert(chapter.BookID, chapter.Number, chapter.Name); err != nil {
		log.Printf("Chapter validation failed. Err:\n%s\n", err)
		context.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// insert values by db query
	chapter, err := SaveChapterDB(chapter.BookID, chapter.Number, chapter.Name)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error. Please try again later.",
		})
		return
	}

	// send back 204
	context.JSON(http.StatusCreated, chapter)
}

// helper func to check if chapter number is valid
func ValidChapterNumber(chapterNumber uint16) bool {
	if chapterNumber > 0 {
		return true
	}
	return false
}

// helper func to check valid book id
func ValidBookID(bookID uint) bool {
	if bookID > 0 {
		return true
	}
	return false
}

// consolidated logic to validate required fields for chapter insert
func ValidateChapterInsert(bookID uint, number *uint16, name string) error {
	// require book id rel
	if !ValidBookID(bookID) {
		return errors.New("Invalid bookId.")
	}

	// name or chapter must be supplied
	if len(name) < 1 {
		// check chapter number if no name supplied
		if number != nil {
			if !ValidChapterNumber(*number) {
				return errors.New("Invalid Chapter number.")
			}
		} else {
			return errors.New("Either Chapter Name or Chapter Number is required.")
		}
	}

	return nil
}
