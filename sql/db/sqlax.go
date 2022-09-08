package db

import (
	"github.com/jmoiron/sqlx"

	"github.com/agflow/tools/log"
)

// SQLXClient is a wrapper of a sqlx.DB client
type SQLXClient struct {
	DB *sqlx.DB
}

// Select selects from `client` using the `query` and `args` and set the result on `dest`
func (c *SQLXClient) Select(dest interface{}, query string, args ...interface{}) error {
	return c.DB.Select(dest, query, args...)
}

// Exec executes from `client` using the `query` and `args`
func (c *SQLXClient) Exec(query string, args ...interface{}) error {
	_, err := c.DB.Exec(query, args...)
	return err
}

// NewSQLXClient return a new db.Client
func NewSQLXClient(url string) (*SQLXClient, error) {
	db, err := sqlx.Connect("postgres", url)
	return &SQLXClient{DB: db}, err
}

// MustNewSQLXClient return a new db.Client without an error
func MustNewSQLXClient(url string) *SQLXClient {
	dbSvc, err := NewSQLXClient(url)
	if err != nil {
		log.Fatal(err)
	}
	return dbSvc
}

// Close closes db connection from the client
func (c *SQLXClient) Close() error {
	return c.DB.Close()
}
