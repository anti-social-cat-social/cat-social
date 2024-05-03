package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

type DB interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

func InitDb() *sqlx.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	params := os.Getenv("DB_PARAMS")

	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s", host, port, username, password, dbname, params)
	db, err := sqlx.Connect("postgres", connection)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("Database connected")
	}
	return db
}
