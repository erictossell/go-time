package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func EditEntry(ctx context.Context, db *sql.DB, id int, name, description string) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	// Defer a rollback in case of error
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Error rolling back transaction: %v", rbErr)
			}
		}
	}()

	// Prepare the SQL statement with context
	statement, err := tx.PrepareContext(ctx, "UPDATE entries SET name = ?, description = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("error preparing update statement: %w", err)
	}
	defer statement.Close()

	// Execute the statement with context
	if _, err = statement.ExecContext(ctx, name, description, id); err != nil {
		return fmt.Errorf("error executing update statement: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func ListEntries(ctx context.Context, db *sql.DB) ([]Entry, error) {
	const query = "SELECT id, name, description, start_time, end_time FROM entries"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		// Log the error for debugging purposes
		log.Printf("Error querying entries: %v", err)
		return nil, fmt.Errorf("error querying entries: %w", err)
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var entry Entry
		if err := rows.Scan(&entry.ID, &entry.Name, &entry.Description, &entry.StartTime, &entry.EndTime); err != nil {
			log.Printf("Error scanning time entry row: %v", err)
			return nil, fmt.Errorf("error scanning time entry row: %w", err)
		}
		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over time entry rows: %v", err)
		return nil, fmt.Errorf("error iterating over time entry rows: %w", err)
	}

	return entries, nil
}

func SaveTimeEntry(ctx context.Context, tx *sql.Tx, name, description string, start, end time.Time) error {
	// Input validation
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if end.Before(start) {
		return fmt.Errorf("end time cannot be before start time")
	}

	statement, err := tx.PrepareContext(ctx, "INSERT INTO entries (name, description, start_time, end_time) VALUES (?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer statement.Close()

	if _, err = statement.ExecContext(ctx, name, description, start, end); err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}
