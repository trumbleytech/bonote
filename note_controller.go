package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetNote(context *gin.Context) {
	// pull id from params
	id := context.Params.ByName("id")

	// convert string to id to validate input
	noteID, err := strconv.Atoi(id)
	if err != nil {
		// handle bad data
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Note Id",
		})
		return
	}
	// query db for note by note id
	note, err := GetNoteByIDDB(noteID)

	// handle db func error
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("Note with %d not found.", noteID),
			})
			return
		} else {
			log.Printf("SQL Error:\n%s\n", err)
			HandleInternalServerError(context)
			return
		}
	}
	context.JSON(http.StatusOK, note)
}

func UpdateNote(context *gin.Context) {
	// pull id from params
	id := context.Params.ByName("id")

	// convert string to id to validate input
	noteID, err := strconv.Atoi(id)
	if err != nil {
		// handle bad data
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Note Id",
		})
		return
	}
	var note Note
	// check valid body
	if err := context.ShouldBindJSON(&note); err != nil {
		log.Printf("ShouldBindJSON note failed. Err: %s\n", err)
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Malformed request body.",
		})
		return
	}

	// verify note name is still valid
	if len(note.Name) < 1 {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Note name is a required field.",
		})
		return
	}

	note, err = UpdateNoteDB(noteID, note.Name, note.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("Note with %d not found.", noteID),
			})
			return
		} else {
			log.Printf("SQL Error:\n%s\n", err)
			HandleInternalServerError(context)
			return
		}
	}

	context.JSON(http.StatusOK, note)
}

func SaveNote(context *gin.Context) {
	var note Note
	// check valid body
	if err := context.ShouldBindJSON(&note); err != nil {
		log.Printf("ShouldBindJSON note failed. Err: %s\n", err)
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Malformed request body.",
		})
		return
	}

	// make sure chapterID is valid
	if note.ChapterID < 1 {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Note chapterId is a required field.",
		})
		return
	}

	if len(note.Name) < 1 {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Note name is a required field.",
		})
		return
	}

	note, err := SaveNoteDB(note.Name, note.Content, note.ChapterID)
	if err != nil {
		HandleInternalServerError(context)
		return
	}

	context.JSON(http.StatusCreated, note)

}

func GetNotesByChapterID(context *gin.Context) {
	id := context.Params.ByName("id")

	chapterID, err := strconv.Atoi(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Chapter ID.",
		})
		return
	}
	notes, err := GetNotesByChapterIDDB(chapterID)
	if err != nil {
		log.Printf("GetNotesByChapterID failed. Err: %s\n", err)
		HandleInternalServerError(context)
		return
	}
	context.JSON(http.StatusOK, notes)
}

func DeleteNote(context *gin.Context) {
	// pull id from params
	id := context.Params.ByName("id")

	// convert string to id to validate input
	noteID, err := strconv.Atoi(id)
	if err != nil {
		// handle bad data
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Note Id",
		})
		return
	}

	err = DeleteNoteDB(noteID)
	if err != nil {
		HandleInternalServerError(context)
		return
	}

	context.Status(http.StatusNoContent)
}
