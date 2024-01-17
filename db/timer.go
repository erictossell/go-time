package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Timer struct {
	ID        int
	TaskName  string
	StartTime time.Time
}

func IsTimerRunning(db *sql.DB, taskName string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM timer_state WHERE is_running = 1 AND task_name = ?", taskName).Scan(&count)
	if err != nil {
		fmt.Println("Error checking timer state:", err)
		return false
	}
	return count > 0
}

func ListTimers(db *sql.DB) ([]Timer, error) {
	query := "SELECT task_name, start_time FROM timer_state WHERE is_running = 1"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying active timers: %w", err)
	}
	defer rows.Close()

	var timers []Timer
	for rows.Next() {
		var timer Timer
		if err := rows.Scan(&timer.TaskName, &timer.StartTime); err != nil {
			return nil, fmt.Errorf("error scanning timer row: %w", err)
		}
		timers = append(timers, timer)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over timer rows: %w", err)
	}

	return timers, nil
}

func StartTimer(db *sql.DB, taskName string) {
	if IsTimerRunning(db, taskName) {
		fmt.Println("Timer is already running for task:", taskName)
		return
	}

	startTime := time.Now()
	_, err := db.Exec("INSERT INTO timer_state (is_running, task_name, start_time) VALUES (?, ?, ?)", true, taskName, startTime)
	if err != nil {
		fmt.Println("Error starting timer:", err)
		return
	}

	fmt.Println("Timer started for task:", taskName)
}

func StopTimer(db *sql.DB, taskName, description string) {
	var startTime time.Time
	err := db.QueryRow("SELECT start_time FROM timer_state WHERE is_running = 1 AND task_name = ?", taskName).Scan(&startTime)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No timer is running for task:", taskName)
		} else {
			fmt.Println("Error fetching running timer:", err)
		}
		return
	}

	endTime := time.Now()
	SaveTimeEntry(db, taskName, description, startTime, endTime) // Use the correct package prefix
	if err != nil {
		fmt.Println("Error saving time entry:", err)
		return
	}

	_, err = db.Exec("UPDATE timer_state SET is_running = 0 WHERE task_name = ?", taskName)
	if err != nil {
		fmt.Println("Error updating timer state:", err)
		return
	}

	fmt.Println("Timer stopped for task:", taskName)
}
