package db

import "time"

// TimeEntry represents a record of time spent on a task.
type TimeEntry struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

// TimerState represents the current state of a timer.
type TimerState struct {
	ID        int       `json:"id"`
	IsRunning bool      `json:"is_running"`
	TaskName  string    `json:"task_name"`
	StartTime time.Time `json:"start_time"`
}
