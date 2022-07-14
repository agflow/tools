package sql

import (
	"database/sql"
	"log"
)

type DBService struct {
	DB *sql.DB
}

type API interface {
	Select(interface{}, string, ...interface{}) error
}

func (db *DBService) Select(dest interface{}, query string, args ...interface{}) error {
	return Select(db.DB, dest, query, args...)
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
	return &DBService{DB: db}, nil
}

func MustNew(url string) *DBService {
	dbSvc, err := New(url)
	if err != nil {
		log.Fatal(err)
	}
	return dbSvc
}
