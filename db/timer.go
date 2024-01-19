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
	TaskName  string
	StartTime time.Time
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
		if err := rows.Scan(&timer.ID, &timer.TaskName, &timer.StartTime); err != nil {
			return nil, fmt.Errorf("error scanning timer row: %w", err)
		}
		timers = append(timers, timer)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over timer rows: %w", err)
	}
	return timers, nil
}

func StartTimer(ctx context.Context, db *sql.DB, taskName string) error {
	isRunning, err := IsTimerRunning(ctx, db, taskName)
	if err != nil {
		return fmt.Errorf("error checking if timer is running: %w", err)
	}
	if isRunning {
		return fmt.Errorf("timer is already running for task: %s", taskName)
	}

	startTime := time.Now()
	_, err = db.ExecContext(ctx, "INSERT INTO timer_state (is_running, task_name, start_time) VALUES (?, ?, ?)", true, taskName, startTime)
	if err != nil {
		return fmt.Errorf("error starting timer: %w", err)
	}
	return nil
}

func StopTimer(ctx context.Context, db *sql.DB, taskName, description string) error {
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
	err = tx.QueryRowContext(ctx, "SELECT start_time FROM timer_state WHERE is_running = 1 AND task_name = ?", taskName).Scan(&startTime)
	if err != nil {
		return fmt.Errorf("error fetching running timer: %w", err)
	}

	endTime := time.Now()
	if err = InsertTimeEntry(ctx, tx, taskName, description, startTime, endTime); err != nil {
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
