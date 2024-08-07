package entry

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"go-time/pkgs/tag"
)

type Entry struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	StartTime   time.Time      `json:"start_time"`
	EndTime     time.Time      `json:"end_time"`
	Tags        []tag.Tag      `json:"tags"`
}

func ReadEntries(ctx context.Context, db *sql.DB) ([]Entry, error) {
	const query = "SELECT id, name, description, start_time, end_time FROM entries"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
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

func CreateEntry(ctx context.Context, tx *sql.Tx, name string, start, end time.Time, tags []string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if end.Before(start) {
		return fmt.Errorf("end time cannot be before start time")
	}

	res, err := tx.ExecContext(ctx, "INSERT INTO entries (name,  start_time, end_time) VALUES (?, ?, ?)", name, start, end)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	entryID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %w", err)
	}

	for _, tag := range tags {
		var tagID int
		_, err = tx.ExecContext(ctx, "INSERT OR IGNORE INTO tags (name) VALUES (?)", tag)
		if err != nil {
			return fmt.Errorf("error inserting tag: %w", err)
		}

		err = tx.QueryRowContext(ctx, "SELECT id FROM tags WHERE name = ?", tag).Scan(&tagID)
		if err != nil {
			return fmt.Errorf("error getting tag ID: %w", err)
		}

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

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Error rolling back transaction: %v", rbErr)
			}
		}
	}()

	if _, err = tx.ExecContext(ctx, "UPDATE entries SET name = ?, description = ? WHERE id = ?", name, description, id); err != nil {
		return fmt.Errorf("error executing update statement: %w", err)
	}

	if _, err = tx.ExecContext(ctx, "DELETE FROM entry_tags WHERE entry_id = ?", id); err != nil {
		return fmt.Errorf("error deleting existing tags: %w", err)
	}

	for _, tag := range tags {

		_, err = tx.ExecContext(ctx, "INSERT OR IGNORE INTO tags (name) VALUES (?)", tag)
		if err != nil {
			return fmt.Errorf("error inserting tag: %w", err)
		}

		var tagID int
		err = tx.QueryRowContext(ctx, "SELECT id FROM tags WHERE name = ?", tag).Scan(&tagID)
		if err != nil {
			return fmt.Errorf("error getting tag ID: %w", err)
		}

		_, err = tx.ExecContext(ctx, "INSERT INTO entry_tags (entry_id, tag_id) VALUES (?, ?)", id, tagID)
		if err != nil {
			return fmt.Errorf("error linking tag with entry: %w", err)
		}
	}

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

func GetEntriesByTag(db *sql.DB, tagName string) ([]Entry, error) {
	var entries []Entry
	query := `
    SELECT e.id, e.name, e.description, e.start_time, e.end_time
    FROM entries e
    INNER JOIN entry_tags et ON e.id = et.entry_id
    INNER JOIN tags t ON et.tag_id = t.id
    WHERE t.name = ?`

	rows, err := db.Query(query, tagName)
	if err != nil {
		return nil, fmt.Errorf("error querying entries by tag: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry Entry
		if err := rows.Scan(&entry.ID, &entry.Name, &entry.Description, &entry.StartTime, &entry.EndTime); err != nil {
			return nil, fmt.Errorf("error scanning entry: %w", err)
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return entries, nil
}

func AddTagsToEntry(db *sql.DB, entryID int, tags []string) error {
	for _, tagName := range tags {

		_, err := db.Exec("INSERT OR IGNORE INTO tags (name) VALUES (?)", tagName)
		if err != nil {
			return fmt.Errorf("error inserting tag: %w", err)
		}

		var tagID int
		err = db.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagID)
		if err != nil {
			return fmt.Errorf("error getting tag ID: %w", err)
		}

		_, err = db.Exec("INSERT INTO entry_tags (entry_id, tag_id) VALUES (?, ?)", entryID, tagID)
		if err != nil {
			return fmt.Errorf("error linking tag with entry: %w", err)
		}
	}
	return nil
}
