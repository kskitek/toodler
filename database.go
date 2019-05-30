package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"go.etcd.io/bbolt"
)

type Todo struct {
	Hash       string
	Author     string
	Date       time.Time
	FileName   string
	LineNumber int
	Line       string
	Priority   int
}

type Database struct {
	db *bbolt.DB
}

func (d *Database) Connect() error {
	// dbPath := "/var/blames.db"
	dbPath := "/tmp/blames.db"
	db, err := bbolt.Open(dbPath, 0666, &bbolt.Options{Timeout: time.Second * 2})
	if err != nil {
		log.Println("Cannot connect to db: " + err.Error())
	}
	d.db = db
	log.Println("Connected to DB: " + dbPath)

	return err
}

func (d *Database) Save(url string, todos []Todo) error {
	todosInFiles := make(map[string][]Todo)
	for _, t := range todos {
		l, ok := todosInFiles[t.FileName]
		if !ok {
			l = make([]Todo, 1)
		}
		l = append(l, t)
		todosInFiles[t.FileName] = l
	}

	return d.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(url))
		if err != nil {
			return err
		}
		// gob.Encode()
		for file, tl := range todosInFiles {
			fb, err := b.CreateBucketIfNotExists([]byte(file))
			if err != nil {
				return err
			}
			for _, t := range tl {
				key := fmt.Sprintf("%s:%d", t.Hash, t.LineNumber)
				var buff bytes.Buffer
				enc := gob.NewEncoder(&buff)
				err := enc.Encode(t)
				if err != nil {
					return err
				}
				fb.Put([]byte(key), buff.Bytes())
			}

		}

		return nil
	})
}

func (d *Database) GetAll() ([]Todo, error) {
	todos := make([]Todo, 0)

	err := d.db.View(func(tx *bbolt.Tx) error {
		tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
			fmt.Println(string(name))
			b.ForEach(func(name []byte, v []byte) error {
				fmt.Println(string(name))
				b.Bucket(name).ForEach(func(name []byte, v []byte) error {
					fmt.Println("# " + string(name))
					buff := bytes.NewBuffer(v)
					dec := gob.NewDecoder(buff)
					t := Todo{}
					err := dec.Decode(&t)
					if err != nil {
						return err
					}
					todos = append(todos, t)
					return nil
				})

				return nil
			})

			return nil
		})
		return nil
	})

	if err != nil {
		return make([]Todo, 0), err
	}
	return todos, nil
}
