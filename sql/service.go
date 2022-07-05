package sql

import (
	"database/sql"
	"log"
)

type DBService struct {
	db *sql.DB
	API
}

type API interface {
	Select(interface{}, string, ...interface{}) error
}

func (db *DBService) Select(dest interface{}, query string, args ...interface{}) error {
	return Select(db.db, dest, query, args...)
}

func New(url string) (*DBService, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	// check that the connection is valid
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DBService{db: db}, nil
}

func MustNew(url string) *DBService {
	dbSvc, err := New(url)
	if err != nil {
		log.Fatal(err)
	}
	return dbSvc
}
