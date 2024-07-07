package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func addTagsToEntry(db *sql.DB, entryID int, tags []string) error {
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

func getEntriesByTag(db *sql.DB, tagName string) ([]Entry, error) {
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

func CreateTag(ctx context.Context, db *sql.DB, name string) error {
	_, err := db.Exec("INSERT INTO tags (name) VALUES (?)", name)
	if err != nil {
		return fmt.Errorf("error inserting tag: %w", err)
	}
	return nil
}

func DeleteTag(ctx context.Context, db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM tags WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting tag: %w", err)
	}
	return nil
}

func GetTags(ctx context.Context, db *sql.DB) ([]Tag, error) {
	var tags []Tag
	rows, err := db.Query("SELECT id, name FROM tags")
	if err != nil {
		return nil, fmt.Errorf("error querying tags: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, fmt.Errorf("error scanning tag: %w", err)
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return tags, nil
}
