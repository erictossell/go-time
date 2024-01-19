package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func ReadEntries(ctx context.Context, db *sql.DB) ([]Entry, error) {
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

func InsertTimeEntry(ctx context.Context, tx *sql.Tx, name string, start, end time.Time, tags []string) error {
	// Input validation
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if end.Before(start) {
		return fmt.Errorf("end time cannot be before start time")
	}

	// Insert entry
	res, err := tx.ExecContext(ctx, "INSERT INTO entries (name,  start_time, end_time) VALUES (?, ?, ?)", name, start, end)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	// Get the last inserted ID for the entry
	entryID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %w", err)
	}

	// Insert and link tags for the entry
	for _, tag := range tags {
		var tagID int
		// Insert tag, ignore if exists
		_, err = tx.ExecContext(ctx, "INSERT OR IGNORE INTO tags (name) VALUES (?)", tag)
		if err != nil {
			return fmt.Errorf("error inserting tag: %w", err)
		}

		// Get tag ID
		err = tx.QueryRowContext(ctx, "SELECT id FROM tags WHERE name = ?", tag).Scan(&tagID)
		if err != nil {
			return fmt.Errorf("error getting tag ID: %w", err)
		}

		// Link tag with entry
		_, err = tx.ExecContext(ctx, "INSERT INTO entry_tags (entry_id, tag_id) VALUES (?, ?)", entryID, tagID)
		if err != nil {
			return fmt.Errorf("error linking tag with entry: %w", err)
		}
	}

	return nil
}
func EditEntry(ctx context.Context, db *sql.DB, id int, name, description string, tags []string) error {
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
	// Update entry
	if _, err = tx.ExecContext(ctx, "UPDATE entries SET name = ?, description = ? WHERE id = ?", name, description, id); err != nil {
		return fmt.Errorf("error executing update statement: %w", err)
	}

	// Delete existing tags associations
	if _, err = tx.ExecContext(ctx, "DELETE FROM entry_tags WHERE entry_id = ?", id); err != nil {
		return fmt.Errorf("error deleting existing tags: %w", err)
	}

	// Insert and link new tags
	for _, tag := range tags {
		// Insert tag, ignore if exists
		_, err = tx.ExecContext(ctx, "INSERT OR IGNORE INTO tags (name) VALUES (?)", tag)
		if err != nil {
			return fmt.Errorf("error inserting tag: %w", err)
		}

		// Get tag ID
		var tagID int
		err = tx.QueryRowContext(ctx, "SELECT id FROM tags WHERE name = ?", tag).Scan(&tagID)
		if err != nil {
			return fmt.Errorf("error getting tag ID: %w", err)
		}

		// Link tag with entry
		_, err = tx.ExecContext(ctx, "INSERT INTO entry_tags (entry_id, tag_id) VALUES (?, ?)", id, tagID)
		if err != nil {
			return fmt.Errorf("error linking tag with entry: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func DeleteEntry(ctx context.Context, db *sql.DB, id int) error {
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

	statement, err := tx.PrepareContext(ctx, "DELETE FROM entries WHERE id = ?")
	if err != nil {
		return fmt.Errorf("error preparing delete statement: %w", err)
	}
	defer statement.Close()

	if _, err = statement.ExecContext(ctx, id); err != nil {
		return fmt.Errorf("error executing delete statement: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
