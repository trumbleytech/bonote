package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// global DB obj
var DB *sql.DB

func ValidateDBConfig() error {
	if len(DB_HOST) < 1 ||
		len(DB_NAME) < 1 ||
		len(DB_PORT) < 1 ||
		len(DB_USER) < 1 ||
		len(DB_PASS) < 1 {
		return errors.New("Invalid DB Env Vars")
	}
	return nil
}

func OpenDBConnection() error {
	var err error
	conn_str := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME, DB_SSL)
	DB, err = sql.Open("postgres", conn_str)
	if err != nil {
		// logging location for easier debug
		log.Printf("Error making DB Connection. Err:\n%s\n", err)
		return err
	}
	if err = DB.Ping(); err != nil {
		return err
	} else {
		log.Printf("DB Connected Successfully.\n")
	}
	return nil
}

func GetBookByID(bookID int) (Book, error) {
	var book Book
	query := "SELECT id,title,author,created_on_utc FROM BOOK WHERE ID = $1"
	err := DB.QueryRow(query, bookID).Scan(&book.Id, &book.Title, &book.Author, &book.Created_On_UTC, &book.Modified_On_UTC)
	if err != nil {
		log.Printf("GetBookByID failed. Err: %s\n", err)
		return Book{}, GenerateResourceNotFoundError("Book", bookID)

	}
	return book, nil
}

func GetChapterByID(chapterID int) (Chapter, error) {
	var chapter Chapter
	query := "SELECT id,book_id,number,name,created_on_utc,modified_on_utc FROM chapter WHERE ID = $1"
	err := DB.QueryRow(query, chapterID).Scan(&chapter.Id, &chapter.Book_id, &chapter.Number, chapter.Name, chapter.Created_On_UTC, chapter.Modified_On_UTC)
	if err != nil {
		log.Printf("GetChapterByID failed. Err: %s\n", err)
		return Chapter{}, err
	}
	return chapter, nil
}

func GetNoteByID(noteID int) (Note, error) {
	var note Note
	query := "SELECT id,chapter_id,name,content,created_on_utc,modified_on_utc FROM note WHERE ID = $1"
	err := DB.QueryRow(query, noteID).Scan(&note.Id, &note.Chapter_id, &note.Name, &note.Content, &note.Created_On_UTC, &note.Modified_On_UTC)
	if err != nil {
		log.Printf("GetNoteByID failed. Err: %s\n", err)
		return Note{}, GenerateResourceNotFoundError("Note", noteID)

	}
	return note, nil
}

func SaveBookDB(title, author string) (Book, error) {
	var query string = "INSERT INTO book (title,author) VALUES ($1,$2) RETURNING id,title,author,created_on_utc"
	var book Book
	err := DB.QueryRow(query, title, author).Scan(&book.Id, &book.Title, &book.Author, &book.Created_On_UTC)
	if err != nil {
		log.Printf("SaveBook failed. Err:%s\n", err)
		return Book{}, err
	}
	return book, nil
}

func SaveChapterDB(bookID uint, number *uint16, name string) (Chapter, error) {
	var chapter Chapter
	query := "INSERT INTO chapter (book_id,number,name) VALUES ($1,$2,$3) RETURNING id,book_id,number,name,created_on_utc,modified_on_utc"
	err := DB.QueryRow(query, bookID, number, name).Scan(&chapter.Id, &chapter.Book_id, &chapter.Number, &chapter.Name, &chapter.Created_On_UTC, &chapter.Modified_On_UTC)
	if err != nil {
		log.Printf("SaveChapter failed. Err:%s\n", err)
		return Chapter{}, err
	}
	return chapter, nil
}

func GenerateResourceNotFoundError(resourceName string, resourceID int) error {
	return fmt.Errorf("%s with ID %d not found.", resourceName, resourceID)
}
