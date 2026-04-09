package main

import (
	"time"
)

type Book struct {
	Id            uint       `json:"id"`
	Title         string     `json:"title"`
	Author        string     `json:"author,omitempty"`
	UserID        int        `json:"userId"`
	CreatedOnUTC  *time.Time `json:"createdOnUTC,omitempty"`
	ModifiedOnUTC *time.Time `json:"modifiedOnUTC,omitempty"`
}

type Chapter struct {
	Id            uint       `json:"id"`
	BookID        uint       `json:"bookId"`
	Number        *uint16    `json:"number,omitempty"`
	Name          string     `json:"name,omitempty"`
	UserID        int        `json:"userId"`
	CreatedOnUTC  *time.Time `json:"createdOnUTC,omitempty"`
	ModifiedOnUTC *time.Time `json:"modifiedOnUTC,omitempty"`
}

type Note struct {
	Id            uint       `json:"id"`
	ChapterID     uint       `json:"chapterId"`
	Name          string     `json:"name"`
	Content       string     `json:"content"`
	UserID        int        `json:"userId"`
	CreatedOnUTC  *time.Time `json:"createdOnUTC,omitempty"`
	ModifiedOnUTC *time.Time `json:"modifiedOnUTC,omitempty"`
}

type User struct {
	Id            uint       `json:"id"`
	Email         string     `json:"email"`
	Username      string     `json:"username"`
	Role          string     `json:"role,omitempty"`
	Password      string     `json:"password,omitempty"`
	CreatedOnUTC  *time.Time `json:"createdOnUTC,omitempty"`
	ModifiedOnUTC *time.Time `json:"modifiedOnUTC,omitempty"`
}

type UserMin struct {
	Id   int
	Role string
}
