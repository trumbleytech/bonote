package main

import (
	"database/sql"

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
