package timer

import (
	"database/sql"
	"fmt"
	godb "go-time/db" // Adjust this import path as necessary
	"time"
)

func IsTimerRunning(db *sql.DB) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM timer_state WHERE is_running = 1").Scan(&count)
	if err != nil {
		fmt.Println("Error checking timer state:", err)
		return false
	}
	return count > 0
}

func StartTimer(db *sql.DB, taskName string) {
	if IsTimerRunning(db) {
		fmt.Println("Timer is already running")
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

func StopTimer(db *sql.DB, description string) {
	var taskName string
	var startTime time.Time
	err := db.QueryRow("SELECT task_name, start_time FROM timer_state WHERE is_running = 1").Scan(&taskName, &startTime)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No timer is running")
		} else {
			fmt.Println("Error fetching running timer:", err)
		}
		return
	}

	endTime := time.Now()
	godb.SaveTimeEntry(db, taskName, description, startTime, endTime) // Use the correct package prefix
	if err != nil {
		fmt.Println("Error saving time entry:", err)
		return
	}

	fmt.Println("Timer stopped for task:", taskName)
}
