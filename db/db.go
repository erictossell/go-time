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
		err := db.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	tableCreators := []func(*sql.DB) error{
		createEntriesTable,
		createTimersTable,
		createTagsTable,
		createEntryTagsTable,
		createTimerTagsTable,
	}

	for _, createFunc := range tableCreators {
		if err := createFunc(db); err != nil {
			return err // Stops at the first error
		}
	}
	return nil
}

func createEntriesTable(db *sql.DB) error {
	sql := `
    CREATE TABLE IF NOT EXISTS entries (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        description TEXT,
        start_time DATETIME NOT NULL,
        end_time DATETIME NOT NULL
    );`
	_, err := db.Exec(sql)
	return err
}

func createTimersTable(db *sql.DB) error {
	sql := `
    CREATE TABLE IF NOT EXISTS timers (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        is_running BOOLEAN NOT NULL,
        name TEXT,
        start_time DATETIME
    );`
	_, err := db.Exec(sql)
	return err
}

func createTagsTable(db *sql.DB) error {
	sql := `
    CREATE TABLE IF NOT EXISTS tags (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL UNIQUE
    );`
	_, err := db.Exec(sql)
	return err
}

func createEntryTagsTable(db *sql.DB) error {
	sql := `
    CREATE TABLE IF NOT EXISTS entry_tags (
        entry_id INTEGER NOT NULL,
        tag_id INTEGER NOT NULL,
        FOREIGN KEY (entry_id) REFERENCES entries(id),
        FOREIGN KEY (tag_id) REFERENCES tags(id),
        PRIMARY KEY (entry_id, tag_id)
    );`
	_, err := db.Exec(sql)
	return err
}

func createTimerTagsTable(db *sql.DB) error {
	sql := `
    CREATE TABLE IF NOT EXISTS timer_tags (
        timer_id INTEGER NOT NULL,
        tag_id INTEGER NOT NULL,
        PRIMARY KEY (timer_id, tag_id),
        FOREIGN KEY (timer_id) REFERENCES timers(id) ON DELETE CASCADE,
        FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
    );`
	_, err := db.Exec(sql)
	return err
}
