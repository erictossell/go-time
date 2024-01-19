package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	godb "go-time/db"
	"time"
)

func ListCmd(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "list [timers|entries]",
		Short: "List all active timers or time entries",
		Long:  `List command is used to list all active timers or time entries. If no argument is provided, it defaults to listing timers.`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			listType := "timers"
			if len(args) > 0 {
				listType = args[0]
			}

			switch listType {
			case "entries":
				listTimeEntries(ctx, db)
			case "timers":
				listActiveTimers(ctx, db)
			default:
				fmt.Println("Invalid argument. Please specify 'entries' or 'timers'.")
			}
		},
	}
}

func listTimeEntries(ctx context.Context, db *sql.DB) {

	entries, err := godb.ListTimeEntries(ctx, db)
	if err != nil {
		fmt.Println("Error listing time entries:", err)
		return
	}
	for _, entry := range entries {
		fmt.Printf("%d | %s | %s | %s | %s\n", entry.ID, entry.Name, entry.Description, entry.StartTime.Format(time.RFC3339), entry.EndTime.Format(time.RFC3339))
	}
}

func listActiveTimers(ctx context.Context, db *sql.DB) {
	timers, err := godb.ListTimers(ctx, db) // Implement this function in your timer package
	if err != nil {
		fmt.Println("Error listing active timers:", err)
		return
	}
	for _, t := range timers {
		fmt.Printf("Task: %s, Started at: %v\n", t.TaskName, t.StartTime.Format(time.RFC3339))
	}
}
