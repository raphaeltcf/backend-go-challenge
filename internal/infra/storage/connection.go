package storage

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func NewConnection(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	queryTable := `CREATE TABLE IF NOT EXISTS orders (
	id           TEXT PRIMARY KEY,
    user_id      TEXT NOT NULL,
    amount       REAL NOT NULL,
    status       TEXT NOT NULL,
    error        TEXT,
    created_at   DATETIME NOT NULL,
    processed_at DATETIME NOT NULL
);`

	_, err = db.Exec(queryTable)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	fmt.Println("Database connection established and table created successfully.")
	return db, nil
}
