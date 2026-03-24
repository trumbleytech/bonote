package main

import (
	"time"
)

type BookRequest struct {
	Id     uint   `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author,omitempty"`
}

type BookResponse struct {
	Id              uint       `json:"id"`
	Title           string     `json:"title"`
	Author          string     `json:"author,omitempty"`
	Created_On_UTC  *time.Time `json:"createdOnUTC,omitempty"`
	Modified_On_UTC *time.Time `json:"modifiedOnUTC,omitempty"`
}

type ChapterRequest struct {
	Id              uint       `json:"id"`
	Book_id         uint       `json:"bookId"`
	Number          *uint16    `json:"number,omitempty"`
	Name            string     `json:"name,omitempty"`
	Created_On_UTC  *time.Time `json:"createdOnUTC,omitempty"`
	Modified_On_UTC *time.Time `json:"modifiedOnUTC,omitempty"`
}

type ChapterResponse struct {
	Id              uint       `json:"id"`
	Book_id         uint       `json:"bookId"`
	Number          *uint16    `json:"number,omitempty"`
	Name            string     `json:"name,omitempty"`
	Created_On_UTC  *time.Time `json:"createdOnUTC,omitempty"`
	Modified_On_UTC *time.Time `json:"modifiedOnUTC,omitempty"`
}

type NoteRequest struct {
	Id         uint   `json:"id"`
	Chapter_id uint   `json:"chapterId"`
	Name       string `json:"name"`
	Content    string `json:"content"`
}

type NoteResponse struct {
	Id              uint       `json:"id"`
	Chapter_id      uint       `json:"chapterId"`
	Name            string     `json:"name"`
	Content         string     `json:"content"`
	Created_On_UTC  *time.Time `json:"createdOnUTC,omitempty"`
	Modified_On_UTC *time.Time `json:"modifiedOnUTC,omitempty"`
}
