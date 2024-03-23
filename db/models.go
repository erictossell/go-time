package db

import (
	"database/sql"
	"time"
)

// TimeEntry represents a record of time spent on a task.
type Entry struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	StartTime   time.Time      `json:"start_time"`
	EndTime     time.Time      `json:"end_time"`
	Tags        []Tag          `json:"tags"`
}

// TimerState represents the current state of a timer.
type TimerState struct {
	ID        int       `json:"id"`
	IsRunning bool      `json:"is_running"`
	TaskName  string    `json:"task_name"`
	StartTime time.Time `json:"start_time"`
}

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
