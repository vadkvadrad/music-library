package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Db struct {
	*sql.DB
}

type DbConfig struct {
	Dsn string
}

func NewDbConfig(dsn string) *DbConfig {
	return &DbConfig{
		Dsn: dsn,
	}
}

func NewDb(conf *DbConfig) *Db {
	db, err := sql.Open("postgres", conf.Dsn)
	if err != nil {
		panic(err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		panic(err)
	}

	return &Db{
		DB: db,
	}
}
