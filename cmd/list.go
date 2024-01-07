package cmd

import (
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	godb "go-time/db"
	timer "go-time/timer" // Adjust this import path as necessary
	"time"
)

// ListCmd creates a new list command.
func ListCmd(db *sql.DB) *cobra.Command {
	var listType string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all time entries or active timers",
		Run: func(cmd *cobra.Command, args []string) {
			switch listType {
			case "entries":
				listTimeEntries(db)
			case "timers":
				listActiveTimers(db)
			default:
				fmt.Println("Invalid list type. Please specify 'entries' or 'timers'.")
			}
		},
	}

	cmd.Flags().StringVarP(&listType, "type", "t", "timers", "Specify what to list: 'entries' or 'timers'")
	return cmd
}

func listTimeEntries(db *sql.DB) {
	entries, err := godb.ListTimeEntries(db)
	if err != nil {
		fmt.Println("Error listing time entries:", err)
		return
	}
	for _, entry := range entries {
		fmt.Printf("%d | %s | %s | %s | %s\n", entry.ID, entry.Name, entry.Description, entry.StartTime.Format(time.RFC3339), entry.EndTime.Format(time.RFC3339))
	}
}

func listActiveTimers(db *sql.DB) {
	timers, err := timer.GetActiveTimers(db) // Implement this function in your timer package
	if err != nil {
		fmt.Println("Error listing active timers:", err)
		return
	}
	for _, t := range timers {
		fmt.Printf("Task: %s, Started at: %v\n", t.TaskName, t.StartTime.Format(time.RFC3339))
	}
}
