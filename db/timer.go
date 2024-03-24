package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Timer struct {
	ID        int
	Name      string
	StartTime time.Time
	Tags      []string
}

type TimerState struct {
	ID        int       `json:"id"`
	IsRunning bool      `json:"is_running"`
	TaskName  string    `json:"task_name"`
	StartTime time.Time `json:"start_time"`
}

func IsTimerRunning(ctx context.Context, db *sql.DB, taskName string) (bool, error) {
	var count int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM timer_state WHERE is_running = 1 AND task_name = ?", taskName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking timer state: %w", err)
	}
	return count > 0, nil
}

func ReadTimers(ctx context.Context, db *sql.DB) ([]Timer, error) {
	query := "SELECT id, task_name, start_time FROM timer_state WHERE is_running = 1"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying active timers: %w", err)
	}
	defer rows.Close()

	var timers []Timer
	for rows.Next() {
		var timer Timer
		if err := rows.Scan(&timer.ID, &timer.Name, &timer.StartTime); err != nil {
			return nil, fmt.Errorf("error scanning timer row: %w", err)
		}
		timers = append(timers, timer)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over timer rows: %w", err)
	}
	return timers, nil
}

func CreateTimer(ctx context.Context, db *sql.DB, taskName string, tags []string) error {
	isRunning, err := IsTimerRunning(ctx, db, taskName)
	if err != nil {
		return fmt.Errorf("error checking if timer is running: %w", err)
	}
	if isRunning {
		return fmt.Errorf("timer is already running for task: %s", taskName)
	}

	startTime := time.Now()
	res, err := db.ExecContext(ctx, "INSERT INTO timer_state (is_running, task_name, start_time) VALUES (?, ?, ?)", true, taskName, startTime)
	if err != nil {
		return fmt.Errorf("error starting timer: %w", err)
	}

	timerID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %w", err)
	}

	for _, tag := range tags {
		var tagID int
		
		_, err = db.ExecContext(ctx, "INSERT OR IGNORE INTO tags (name) VALUES (?)", tag)
		if err != nil {
			return fmt.Errorf("error inserting tag: %w", err)
		}

		
		err = db.QueryRowContext(ctx, "SELECT id FROM tags WHERE name = ?", tag).Scan(&tagID)
		if err != nil {
			return fmt.Errorf("error getting tag ID: %w", err)
		}

		
		_, err = db.ExecContext(ctx, "INSERT INTO timer_tags (timer_id, tag_id) VALUES (?, ?)", timerID, tagID)
		if err != nil {
			return fmt.Errorf("error linking tag with timer: %w", err)
		}
	}

	return nil
}

func StopTimer(ctx context.Context, db *sql.DB, taskName string) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil && rbErr != sql.ErrTxDone {
			log.Printf("transaction rollback error: %v", rbErr)
		}
	}()

	var startTime time.Time
	var timerID int
	err = tx.QueryRowContext(ctx, "SELECT id, start_time FROM timer_state WHERE is_running = 1 AND task_name = ?", taskName).Scan(&timerID, &startTime)
	if err != nil {
		return fmt.Errorf("error fetching running timer: %w", err)
	}

	
	tags, err := fetchTagsForTimer(ctx, tx, timerID)
	if err != nil {
		return fmt.Errorf("error fetching tags for timer: %w", err)
	}

	endTime := time.Now()
	if err = InsertTimeEntry(ctx, tx, taskName, startTime, endTime, tags); err != nil {
		return fmt.Errorf("error saving time entry: %w", err)
	}

	if _, err = tx.ExecContext(ctx, "UPDATE timer_state SET is_running = 0 WHERE task_name = ?", taskName); err != nil {
		return fmt.Errorf("error updating timer state: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func fetchTagsForTimer(ctx context.Context, tx *sql.Tx, timerID int) ([]string, error) {
	var tags []string
	query := `
    SELECT t.name 
    FROM tags t 
    INNER JOIN timer_tags tt ON t.id = tt.tag_id 
    WHERE tt.timer_id = ?`

	rows, err := tx.QueryContext(ctx, query, timerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}
