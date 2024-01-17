package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func EditTimeEntry(db *sql.DB, id int, name, description string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	statement, err := tx.Prepare("UPDATE time_entries SET name = ?, description = ? WHERE id = ?")
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("error preparing update statement: %v, rollback error: %v", err, rollbackErr)
		}
		return fmt.Errorf("error preparing update statement: %w", err)
	}
	defer statement.Close()

	_, err = statement.Exec(name, description, id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("error executing update statement: %v, rollback error: %v", err, rollbackErr)
		}
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

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	statement, err := tx.Prepare("INSERT INTO time_entries (name, description, start_time, end_time) VALUES (?, ?, ?, ?)")
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("error preparing statement: %v, rollback error: %v", err, rollbackErr)
		}
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer statement.Close()

	_, err = statement.Exec(name, description, start, end)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("error executing statement: %v, rollback error: %v", err, rollbackErr)
		}
		return fmt.Errorf("error executing statement: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
