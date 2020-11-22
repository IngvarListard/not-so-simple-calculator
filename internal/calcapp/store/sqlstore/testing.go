package sqlstore

import (
	"database/sql"
	"fmt"
	server "github.com/IngvarListard/not-so-simple-calculator/internal/calcapp"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
	"testing"
)

func TestDB(t *testing.T) (db *sql.DB, teardown func(...string)) {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	if err := server.CreateSchema(db); err != nil {
		log.Fatal(err)
	}
	return db, func(tables ...string) {
		if len(tables) > 0 {
			_, _ = db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", ")))
		}

		_ = db.Close()
	}
}
