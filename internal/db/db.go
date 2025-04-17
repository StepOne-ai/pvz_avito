package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	createTables()
	return nil
}

func createTables() {
	userTable := `
    CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        role TEXT NOT NULL
    );`

	pvzTable := `
    CREATE TABLE IF NOT EXISTS pvzs (
        id TEXT PRIMARY KEY,
        registration_date DATETIME NOT NULL,
        city TEXT NOT NULL
    );`

	receptionTable := `
    CREATE TABLE IF NOT EXISTS receptions (
        id TEXT PRIMARY KEY,
        date_time DATETIME NOT NULL,
        pvz_id TEXT NOT NULL,
        status TEXT NOT NULL,
        FOREIGN KEY (pvz_id) REFERENCES pvzs(id)
    );`

	productTable := `
    CREATE TABLE IF NOT EXISTS products (
        id TEXT PRIMARY KEY,
        date_time DATETIME NOT NULL,
        type TEXT NOT NULL,
        reception_id TEXT NOT NULL,
        FOREIGN KEY (reception_id) REFERENCES receptions(id)
    );`

	tables := []string{userTable, pvzTable, receptionTable, productTable}
	for _, table := range tables {
		_, err := DB.Exec(table)
		if err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}
	}
}
