package db

import (
	"database/sql"
	"log"
)

// Client is a wrapper of a sql.DB client
type Client struct {
	DB *sql.DB
}

// Service is an interface od db.Service
type Service interface {
	Select(interface{}, string, ...interface{}) error
}

// Select selects from `client`` using the `query` and `args`
func (c *Client) Select(dest interface{}, query string, args ...interface{}) error {
	return Select(c.DB, dest, query, args...)
}

// New return a new db.Client
func New(url string) (*Client, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	// check that the connection is valid
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Client{DB: db}, nil
}

// MustNew return a new db.Client without an error
func MustNew(url string) *Client {
	dbSvc, err := New(url)
	if err != nil {
		log.Fatal(err)
	}
	return dbSvc
}
