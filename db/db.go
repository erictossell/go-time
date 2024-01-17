package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./timetracker.db")
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
	CREATE TABLE IF NOT EXISTS time_entries (
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
    );`
	_, err := db.Exec(createTablesSQL)
	if err != nil {
		return fmt.Errorf("error creating tables: %w", err)
	}
	return nil
}

func EditTimeEntry(db *sql.DB, id int, name, description string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	statement, err := tx.Prepare("UPDATE time_entries SET name = ?, description = ? WHERE id = ?")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error preparing update statement: %w", err)
	}
	defer statement.Close()

	_, err = statement.Exec(name, description, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing update statement: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func ListTimeEntries(db *sql.DB) ([]TimeEntry, error) {
	var entries []TimeEntry
	rows, err := db.Query("SELECT id, name, description, start_time, end_time FROM time_entries")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry TimeEntry
		if err := rows.Scan(&entry.ID, &entry.Name, &entry.Description, &entry.StartTime, &entry.EndTime); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func SaveTimeEntry(db *sql.DB, name, description string, start, end time.Time) error {
	// Input validation (example: check for empty name, invalid dates, etc.)
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if end.Before(start) {
		return fmt.Errorf("end time cannot be before start time")
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	// Prepare statement within the transaction
	statement, err := tx.Prepare("INSERT INTO time_entries (name, description, start_time, end_time) VALUES (?, ?, ?, ?)")
	if err != nil {
		tx.Rollback() // Rollback in case of an error
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer statement.Close()

	// Execute the statement
	_, err = statement.Exec(name, description, start, end)
	if err != nil {
		tx.Rollback() // Rollback in case of an error
		return fmt.Errorf("error executing statement: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
