package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

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

/*
SELECT FUNCS
*/

// to be used by login function
func GetUserByUsernameDB(username string) (User, error) {
	var user User
	query := "SELECT id,hashed_password FROM users WHERE username = $1"
	err := DB.QueryRow(query, username).Scan(&user.Id, &user.Password)
	if err != nil {
		log.Printf("GetUserByUsername failed. Err: %s\n", err)
		return User{}, err
	}
	return user, nil

}

// to be used to pull relevant user data
func GetUserByIDDB(userID int) (User, error) {
	var user User
	query := "SELECT id,email,username,created_on_utc,modified_on_utc FROM users WHERE id = $1"
	err := DB.QueryRow(query, userID).Scan(&user.Id, &user.Email, &user.Username, &user.CreatedOnUTC, &user.ModifiedOnUTC)
	if err != nil {
		log.Printf("GetUserByID failed. Err: %s\n", err)
		return User{}, err
	}
	return user, nil
}

func GetBookByIDDB(bookID int) (Book, error) {
	var book Book
	query := "SELECT id,title,author,created_on_utc,modified_on_utc FROM BOOK WHERE ID = $1"
	err := DB.QueryRow(query, bookID).Scan(&book.Id, &book.Title, &book.Author, &book.CreatedOnUTC, &book.ModifiedOnUTC)
	if err != nil {
		log.Printf("GetBookByID failed. Err: %s\n", err)
		return Book{}, err

	}
	return book, nil
}

func GetChapterByIDDB(chapterID int) (Chapter, error) {
	var chapter Chapter
	query := "SELECT id,book_id,number,name,created_on_utc,modified_on_utc FROM chapter WHERE ID = $1"
	err := DB.QueryRow(query, chapterID).Scan(&chapter.Id, &chapter.BookID, &chapter.Number, &chapter.Name, &chapter.CreatedOnUTC, &chapter.ModifiedOnUTC)
	if err != nil {
		log.Printf("GetChapterByID failed. Err: %s\n", err)
		return Chapter{}, err
	}
	return chapter, nil
}

func GetNoteByIDDB(noteID int) (Note, error) {
	var note Note
	query := "SELECT id,chapter_id,name,content,created_on_utc,modified_on_utc FROM note WHERE ID = $1"
	err := DB.QueryRow(query, noteID).Scan(&note.Id, &note.ChapterID, &note.Name, &note.Content, &note.CreatedOnUTC, &note.ModifiedOnUTC)
	if err != nil {
		log.Printf("GetNoteByID failed. Err: %s\n", err)
		return Note{}, err

	}
	return note, nil
}

func GetBooksByUserIDDB(userID int) ([]Book, error) {
	query := "SELECT id,title,author,created_on_utc,modified_on_utc FROM book WHERE user_id = $1 ORDER BY id DESC"
	rows, err := DB.Query(query, userID)
	if err != nil {
		log.Printf("GetBooksByUserIDDB query failed. Err: %s\n", err)
		return nil, err
	}
	defer rows.Close()

	books := []Book{}

	for rows.Next() {
		book := Book{}
		if err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.CreatedOnUTC, &book.ModifiedOnUTC); err != nil {
			log.Printf("rows.Scan failed. Err: %s\n", err)
			continue
		}
		books = append(books, book)
	}
	return books, nil
}

func GetChaptersByBookIDDB(bookID int) ([]Chapter, error) {
	query := "SELECT id,book_id,number,name,created_on_utc,modified_on_utc FROM chapter WHERE book_id = $1 ORDER BY number DESC"
	rows, err := DB.Query(query, bookID)
	if err != nil {
		log.Printf("GetChaptersByBookIDDB query failed. Err: %s\n", err)
		return nil, err
	}
	defer rows.Close()

	chapters := []Chapter{}

	for rows.Next() {
		chapter := Chapter{}
		if err := rows.Scan(&chapter.Id, &chapter.BookID, &chapter.Number, &chapter.Name, &chapter.CreatedOnUTC, &chapter.ModifiedOnUTC); err != nil {
			log.Printf("rows.Scan failed. Err: %s\n", err)
			return nil, err
		}
		chapters = append(chapters, chapter)
	}
	return chapters, nil
}

func GetNotesByChapterIDDB(chapterID int) ([]Note, error) {
	query := "SELECT id,name,content,chapter_id,created_on_utc,modified_on_utc FROM note WHERE chapter_id = $1"
	rows, err := DB.Query(query, chapterID)
	if err != nil {
		log.Printf("GetNotesByChapterID query failed. Err: %s\n", err)
		return nil, err
	}

	defer rows.Close()

	notes := []Note{}

	for rows.Next() {
		note := Note{}
		if err := rows.Scan(&note.Id, &note.Name, &note.Content, &note.ChapterID, &note.CreatedOnUTC, &note.ModifiedOnUTC); err != nil {
			log.Printf("rows.Scan failed. Err: %s\n", err)
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

/*
INSERT FUNCS
*/

func SaveUserDB(username, email, hashedPassword string) (User, error) {
	query := "INSERT INTO users (username, email, hashed_password) VALUES ($1, $2, $3) RETURNING id,username,email,created_on_utc,modified_on_utc"
	var user User
	err := DB.QueryRow(query, username, email, hashedPassword).Scan(&user.Id, &user.Username, &user.Email, &user.CreatedOnUTC, &user.ModifiedOnUTC)
	if err != nil {
		log.Printf("SaveUser failed. Err: %\n", err)
		return User{}, err
	}
	return user, nil
}

func SaveBookDB(title, author string) (Book, error) {
	var query string = "INSERT INTO book (title,author) VALUES ($1,$2) RETURNING id,title,author,created_on_utc"
	var book Book
	err := DB.QueryRow(query, title, author).Scan(&book.Id, &book.Title, &book.Author, &book.CreatedOnUTC)
	if err != nil {
		log.Printf("SaveBook failed. Err:%s\n", err)
		return Book{}, err
	}
	return book, nil
}

func SaveChapterDB(bookID uint, number *uint16, name string) (Chapter, error) {
	var chapter Chapter
	query := "INSERT INTO chapter (book_id,number,name) VALUES ($1,$2,$3) RETURNING id,book_id,number,name,created_on_utc,modified_on_utc"
	err := DB.QueryRow(query, bookID, number, name).Scan(&chapter.Id, &chapter.BookID, &chapter.Number, &chapter.Name, &chapter.CreatedOnUTC, &chapter.ModifiedOnUTC)
	if err != nil {
		log.Printf("SaveChapter failed. Err:%s\n", err)
		return Chapter{}, err
	}
	return chapter, nil
}

func SaveNoteDB(name, content string, chapterID uint) (Note, error) {
	var note Note
	query := "INSERT INTO note (name,content,chapter_id) VALUES ($1,$2,$3) RETURNING id,name,content,chapter_id,created_on_utc,modified_on_utc"
	err := DB.QueryRow(query, name, content, chapterID).Scan(&note.Id, &note.Name, &note.Content, &note.ChapterID, &note.CreatedOnUTC, &note.ModifiedOnUTC)
	if err != nil {
		log.Printf("SaveNote failed. Err:%s\n", err)
		return Note{}, err
	}
	return note, nil
}

/*
UPDATE FUNCS
*/

func UpdateBookDB(id int, title, author string) (Book, error) {
	var query string = "UPDATE book SET title = $1, author = $2 WHERE id = $3 RETURNING id,title,author,created_on_utc,modified_on_utc"
	var book Book
	err := DB.QueryRow(query, title, author, id).Scan(&book.Id, &book.Title, &book.Author, &book.CreatedOnUTC, &book.ModifiedOnUTC)
	if err != nil {
		log.Printf("UpdateBook failed. Err:%s\n", err)
		return Book{}, err
	}
	return book, nil
}

func UpdateChapterDB(id int, number *uint16, name string) (Chapter, error) {
	var chapter Chapter
	query := "UPDATE chapter SET number = $1, name = $2 WHERE id = $3 RETURNING id,book_id,number,name,created_on_utc,modified_on_utc"
	err := DB.QueryRow(query, number, name, id).Scan(&chapter.Id, &chapter.BookID, &chapter.Number, &chapter.Name, &chapter.CreatedOnUTC, &chapter.ModifiedOnUTC)
	if err != nil {
		log.Printf("UpdateChapter failed. Err:%s\n", err)
		return Chapter{}, err
	}
	return chapter, nil
}

func UpdateNoteDB(id int, name, content string) (Note, error) {
	var note Note
	query := "UPDATE note SET name = $1, content = $2 WHERE id = $3 RETURNING id,name,content,chapter_id,created_on_utc,modified_on_utc"
	err := DB.QueryRow(query, name, content, id).Scan(&note.Id, &note.Name, &note.Content, &note.ChapterID, &note.CreatedOnUTC, &note.ModifiedOnUTC)
	if err != nil {
		log.Printf("SaveNote failed. Err:%s\n", err)
		return Note{}, err
	}
	return note, nil
}

/*
DELETE FUNCS
*/

func DeleteUserDB(id int) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := DB.Exec(query, id)
	return err
}
func DeleteBookDB(id int) error {
	query := "DELETE FROM book WHERE id = $1"
	_, err := DB.Exec(query, id)
	return err
}

func DeleteChapterDB(id int) error {
	query := "DELETE FROM chapter WHERE id = $1"
	_, err := DB.Exec(query, id)
	return err
}

func DeleteNoteDB(id int) error {
	query := "DELETE FROM note WHERE id = $1"
	_, err := DB.Exec(query, id)
	return err
}

/*
SESSION FUNCS
*/

func GetUserIDBySessionTokenHashDB(tokenHash string) (int, error) {
	var userID int
	query := "SELECT user_id FROM sessions WHERE token = $1"
	err := DB.QueryRow(query, tokenHash).Scan(&userID)
	if err != nil {
		log.Printf("GetUserIDBySessionTokenDB Failed. Err: %s\n", err)
		return 0, err
	}
	return userID, nil
}

func CreateNewSession(userID int, token_hash string, expiresAt time.Time) error {
	query := "INSERT INTO sessions (user_id, token_hash, expires_at) VALUES ($1,$2,$3)"
	_, err := DB.Exec(query, userID, token_hash, expiresAt)
	return err
}
