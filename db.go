package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func Validate_db_config() error {
	if len(DB_HOST) < 1 ||
		len(DB_NAME) < 1 ||
		len(DB_PORT) < 1 ||
		len(DB_USER) < 1 ||
		len(DB_PASS) < 1 {
		return errors.New("Invalid DB Env Vars")
	}
	return nil
}

func Open_DB_Connection() (*sql.DB, error) {
	conn_str := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME, DB_SSL)
	db, err := sql.Open("postgres", conn_str)
	if err != nil {
		// logging location for easier debug
		log.Printf("Error making DB Connection. Err:\n%s\n", err)
		return nil, err
	}
	log.Printf("DB Connected Successfully.\n")
	return db, nil
}
