package kree

import (
	"log"
	"os"

	"database/sql"
)

var (
	DB *sql.DB

	MYSQL_CONNECTION_URI = os.Getenv("MYSQL_CONNECTION_URI")
)

func InitDB() {
	db, err := sql.Open("mysql", MYSQL_CONNECTION_URI)
	if err != nil {
		log.Fatal(err)
	}

	DB = db
}

func GetDB() *sql.DB {
	return DB
}
