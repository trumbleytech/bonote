package main

import (
	"time"
)

type Book struct {
	Id              uint   `json:"id"`
	Title           string `json:"title"`
	Author          string `json:"author,omitempty"`
	Created_On_UTC  *time.Time
	Modified_On_UTC *time.Time
}

type Chapter struct {
	Id              uint    `json:"id"`
	Book_id         uint    `json:"book_id"`
	Number          *uint16 `json:"number"`
	Name            string  `json:"name"`
	Created_On_UTC  *time.Time
	Modified_On_UTC *time.Time
}

type Note struct {
	Id              uint   `json:"id"`
	Chapter_id      uint   `json:"chapter_id"`
	Name            string `json:"name"`
	Content         string `json:"content"`
	Created_On_UTC  *time.Time
	Modified_On_UTC *time.Time
}
