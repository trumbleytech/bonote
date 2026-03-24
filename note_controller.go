package main

import (
	"database/sql"

	"fmt"

	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Note_Get(context *gin.Context) {
	// pull id from params
	id := context.Params.ByName("id")

	// convert string to id to validate input
	note_id, err := strconv.Atoi(id)
	if err != nil {
		// handle bad data
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Note Id",
		})
		return
	}
	// query db for note by note id
	result := Query_Note_By_Id(note_id)
	var note Note
	err = result.Scan(&note.Id, &note.Chapter_id, &note.Name, &note.Content, &note.Created_On_UTC, &note.Modified_On_UTC)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("No note found with Id %d", note_id),
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
	context.JSON(http.StatusOK, note)
}
