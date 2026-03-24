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

func Validate_db_config() error {
	if len(DB_HOST) < 1 ||
		len(DB_NAME) < 1 ||
		len(DB_PORT) < 1 ||
		len(DB_USER) < 1 ||
		len(DB_PASS) < 1 {
		return errors.New("Invalid DB Env Vars")
	}
	return nil
}

func Open_DB_Connection() error {
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

func Query_Book_By_Id(book_id int) *sql.Row {
	result := DB.QueryRow("SELECT id,title,author,created_on_utc FROM BOOK WHERE ID = $1", book_id)
	return result
}

func Query_Chapter_By_Id(chapter_id int) *sql.Row {
	result := DB.QueryRow("SELECT id,book_id,number,name,created_on_utc,modified_on_utc FROM chapter WHERE ID = $1", chapter_id)
	return result
}

func Query_Note_By_Id(note_id int) *sql.Row {
	result := DB.QueryRow("SELECT id,chapter_id,name,content,created_on_utc,modified_on_utc FROM note WHERE ID = $1", note_id)
	return result
}

func Insert_Book(title, author string) (BookResponse, error) {
	var query string = "INSERT INTO book (title,author) VALUES ($1,$2) RETURNING id,title,author,created_on_utc"
	var bookRes BookResponse
	err := DB.QueryRow(query, title, author).Scan(&bookRes.Id, &bookRes.Title, &bookRes.Author, &bookRes.Created_On_UTC)
	if err != nil {
		log.Printf("Book insert failed. Err:\n%s\n", err)
		return BookResponse{}, err
	}
	log.Printf("INSERT BOOK: title: %s || author: %s\n", title, author)
	return bookRes, nil
}

func Insert_Chapter(book_id uint, number *uint16, name string) error {
	_, err := DB.Exec("INSERT INTO chapter (book_id,number,name) VALUES ($1,$2,$3)", book_id, number, name)
	if err != nil {
		log.Printf("Chapter insert failed. Err:\n%s\n", err)
		return err
	}
	return nil
}
