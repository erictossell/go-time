package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbFile string) (*sql.DB, error) {

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := createTables(db); err != nil {
		db.Close() // Close the database connection on error
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	createTablesSQL := `
	CREATE TABLE IF NOT EXISTS entries (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        description TEXT,
        start_time DATETIME NOT NULL,
        end_time DATETIME NOT NULL
    );
    CREATE TABLE IF NOT EXISTS timer_state (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        is_running BOOLEAN NOT NULL,
        task_name TEXT,
        start_time DATETIME
    );
    CREATE TABLE IF NOT EXISTS tags (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	createTablesSQL +=
	name TEXT NOT NULL
    );
    CREATE TABLE IF NOT EXISTS entry_tags (
	    entry_id INTEGER NOT NULL,
	    tag_id INTEGER NOT NULL,
	    FOREIGN KEY (entry_id) REFERENCES entries(id),
	    FOREIGN KEY (tag_id) REFERENCES tags(id)
	    PRIMARY KEY (entry_id, tag_id)
    );
	`
	_, err := db.Exec(createTablesSQL)
	if err != nil {
		return fmt.Errorf("error creating tables: %w", err)
	}
	return nil
}
