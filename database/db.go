package database

import (
	"database/sql"
	"log"

	"github.com/ahmedalialphasquad123/calculationService/config"
	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	connStr := config.Config("DB_URL")
	log.Printf("db URL: %v", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
