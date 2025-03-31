package db

import (

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
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
	db, err := gorm.Open(postgres.Open(conf.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &Db{
		DB: db,
	}
}
