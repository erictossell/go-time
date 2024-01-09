package cmd

import (
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	godb "go-time/db"
	timer "go-time/timer" // Adjust this import path as necessary
	"time"
)

func ListCmd(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "list [timers|entries]",
		Short: "List all active timers or time entries",
		Long:  `List command is used to list all active timers or time entries. If no argument is provided, it defaults to listing timers.`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			listType := "timers"
			if len(args) > 0 {
				listType = args[0]
			}

			switch listType {
			case "entries":
				listTimeEntries(db)
			case "timers":
				listActiveTimers(db)
			default:
				fmt.Println("Invalid argument. Please specify 'entries' or 'timers'.")
			}
		},
	}
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
	timers, err := timer.ListActiveTimers(db) // Implement this function in your timer package
	if err != nil {
		fmt.Println("Error listing active timers:", err)
		return
	}
	for _, t := range timers {
		fmt.Printf("Task: %s, Started at: %v\n", t.TaskName, t.StartTime.Format(time.RFC3339))
	}
}
