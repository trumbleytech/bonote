package main

import "time"

type Book struct {
	Id     uint
	Title  string
	Author string
}

type Chapter struct {
	Id      uint
	Number  uint16
	Name    string
	Book_id uint
}

type Note struct {
	Id                uint
	Name              string
	Content           string
	Last_Modified_UTC time.Time
	Created_On_UTC    time.Time
	Chapter_id        uint
}
