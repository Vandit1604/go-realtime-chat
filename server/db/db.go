package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

// returns a new connection to the database
func NewDatabase() (*Database, error) {
	db, err := sql.Open("postgres", "postgresql://root:password@localhost:5432/go-chat?sslmode=disable")
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDb() *sql.DB {
	return d.db
}
