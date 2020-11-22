package calcapp

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func NewDB(sqlitePath string) (*sql.DB, error) {
	var isEmpty bool
	if _, err := os.Stat(sqlitePath); os.IsNotExist(err) {
		isEmpty = true
	} else if err != nil {
		return nil, fmt.Errorf("open sqlite path error: %w", err)
	}

	db, err := sql.Open("sqlite3", sqlitePath)
	if err != nil {
		return nil, fmt.Errorf("opening database error: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database error: %w", err)
	}
	log.Println("Database connection established")

	if !isEmpty {
		return db, nil
	}

	err = CreateSchema(db)
	if err != nil {
		return nil, fmt.Errorf("creating schema error: %w", err)
	}

	return db, nil
}

func CreateSchema(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS history (
		id INTEGER,
		event_time TEXT NOT NULL,
		expression TEXT NOT NULL,
		result TEXT NOT NULL,
		PRIMARY KEY("id")
);`)
	if err == nil {
		log.Println(`New table "history" successfully created`)
	}
	return err
}
