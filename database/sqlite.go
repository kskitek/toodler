package database

import (
	"database/sql"
	"io/ioutil"

	_ "github.com/mattn/go-sqlite3"
)

func NewSqlite() *Sqlite {
	return &Sqlite{path: "/tmp/toodler.db"}
}

type Sqlite struct {
	path           string
	migrationsPath string
	db             *sql.DB
}

// TODO switch all functions to use ctx

func (s *Sqlite) Connect() error {
	db, err := sql.Open("sqlite3", s.path)
	if err != nil {
		return err
	}
	s.db = db
	return s.createSchema()
}

func (s *Sqlite) Save(url string, todos []Todo) error {
	return nil
}

func (s *Sqlite) GetAll() ([]Todo, error) {
	return make([]Todo, 0), nil
}

func (s *Sqlite) createSchema() error {
	buf, err := ioutil.ReadFile(s.migrationsPath)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(string(buf))
	if err != nil {
		return err
	}

	return nil
}
