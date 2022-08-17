package db

import (
	"database/sql"

	"github.com/agflow/tools/log"
)

// Client is a wrapper of a sql.DB client
type Client struct {
	DB *sql.DB
}

// Service is an interface od db.Service
type Service interface {
	Select(interface{}, string, ...interface{}) error
	Close() error
	Exec(string, ...interface{}) error
}

// Select selects from `client` using the `query` and `args` and set the result on `dest`
func (c *Client) Select(dest interface{}, query string, args ...interface{}) error {
	return Select(c.DB, dest, query, args...)
}

// Exec executes from `client` using the `query` and `args`
func (c *Client) Exec(query string, args ...interface{}) error {
	_, err := c.DB.Exec(query, args...)
	return err
}

// New return a new db.Client
func New(url string) (*Client, error) {
	db, err := sql.Open("postgres", url)
	return &Client{DB: db}, err
}

// MustNew return a new db.Client without an error
func MustNew(url string) *Client {
	dbSvc, err := New(url)
	if err != nil {
		log.Fatal(err)
	}
	return dbSvc
}

// Close closes db connection from the client
func (c *Client) Close() error {
	return c.DB.Close()
}
