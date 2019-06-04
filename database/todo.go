package database

import "time"

type Todo struct {
	Hash       string
	Author     string
	Date       time.Time
	FileName   string
	LineNumber int
	Line       string
	Priority   int
}
