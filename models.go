package main

import "time"

type Book struct {
	Id                uint
	Title             string
	Author            string
	Created_On_UTC    time.Time
	Last_Modified_UTC time.Time
}

type Chapter struct {
	Id                uint
	Book_id           uint
	Number            uint16
	Name              string
	Created_On_UTC    time.Time
	Last_Modified_UTC time.Time
}

type Note struct {
	Id                uint
	Chapter_id        uint
	Name              string
	Content           string
	Created_On_UTC    time.Time
	Last_Modified_UTC time.Time
}
