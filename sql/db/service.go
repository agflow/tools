package db

// Service is an interface of db.Service
type Service interface {
	Select(interface{}, string, ...interface{}) error
	Close() error
	Exec(string, ...interface{}) error
}
