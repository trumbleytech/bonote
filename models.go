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
	Id              uint
	Book_id         uint
	Number          *uint16
	Name            string
	Created_On_UTC  *time.Time
	Modified_On_UTC *time.Time
}

type Note struct {
	Id              uint
	Chapter_id      uint
	Name            string
	Content         string
	Created_On_UTC  *time.Time
	Modified_On_UTC *time.Time
}
