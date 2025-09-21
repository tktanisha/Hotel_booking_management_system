package db

import (
	"database/sql"
)

//go:generate mockgen -source=db.go -destination=../mocks/mock_db.go -package=mocks
type DB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Close() error
}
