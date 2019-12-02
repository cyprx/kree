package kree

import (
	"log"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DB *sql.DB
)

func InitDB(uri string) {
	db, err := sql.Open("mysql", uri)
	if err != nil {
		log.Fatal(err)
	}

	DB = db
}

func GetDB() *sql.DB {
	return DB
}
