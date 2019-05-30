package main

import (
	"fmt"
	"log"
)

func main() {

	db := &Database{}
	err := db.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
s := &Searcher{}
	result, err := s.Search("http://github.com/kskitek/toodler")
	if err != nil {
		log.Fatal(err.Error())
	}

	todos := make([]Todo, len(result))
	for i, v := range result {
		todos[i] = Todo{
			Hash:       v.Hash,
			Author:     v.Author,
			Date:       v.Date,
			FileName:   v.FileName,
			LineNumber: v.LineNumber,
			Line:       v.Line,
			Priority:   2,
		}
	}

	err = db.Save("p1", todos)
	if err != nil {
		log.Fatal(err.Error())
	}

	all, err := db.GetAll()
	for _, v := range all {
		fmt.Println(v)
	}
}
