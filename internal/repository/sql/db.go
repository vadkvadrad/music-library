package sql

import (
	"database/sql"
	"fmt"
)

type DB struct {
	*sql.DB
}

func NewDB(db *sql.DB) *DB {
	return &DB{db}
}

// Transaction interface для работы с транзакциями
type Transaction interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// Проверяем, что *sql.Tx реализует Transaction
var _ Transaction = (*sql.Tx)(nil)
// Проверяем, что *DB реализует Transaction
var _ Transaction = (*DB)(nil)