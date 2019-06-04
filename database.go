package main

import "github.com/kskitek/toodler/database"

type Database interface {
	Connect() error
	Save(url string, todos []database.Todo) error
	GetAll() ([]database.Todo, error)
}
