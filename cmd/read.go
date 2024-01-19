package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	godb "go-time/db"
	"time"
)

func ReadCmd(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "read [timers|entries]",
		Short: "List all active timers or time entries",
		Long:  `Read command is used to list all active timers or time entries. If no argument is provided, it defaults to listing timers.`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			listType := "timers"
			if len(args) > 0 {
				listType = args[0]
			}

			switch listType {
			case "entries":
				readEntries(ctx, db)
			case "timers":
				readTimers(ctx, db)
			default:
				fmt.Println("Invalid argument. Please specify 'entries' or 'timers'.")
			}
		},
	}
}

func readEntries(ctx context.Context, db *sql.DB) {
	entries, err := godb.ReadEntries(ctx, db)
	if err != nil {
		fmt.Println("Error listing time entries:", err)
		return
	}

	// Determine the width of each column
	idWidth := 4    // Adjust as needed
	nameWidth := 20 // Adjust based on your longest name
	descWidth := 30 // Adjust based on your longest description
	timeWidth := 25 // Adjust to fit the length of the formatted time

	headerFormat := fmt.Sprintf("%%-%ds | %%-%ds | %%-%ds | %%-%ds | %%-%ds\n", idWidth, nameWidth, descWidth, timeWidth, timeWidth)
	rowFormat := fmt.Sprintf("%%-%dd | %%-%ds | %%-%ds | %%-%ds | %%-%ds\n", idWidth, nameWidth, descWidth, timeWidth, timeWidth)

	fmt.Printf(headerFormat, "ID", "Name", "Description", "Start Time", "End Time")

	for _, entry := range entries {
		fmt.Printf(rowFormat, entry.ID, entry.Name, entry.Description, entry.StartTime.Format(time.RFC3339), entry.EndTime.Format(time.RFC3339))
	}
}

func readTimers(ctx context.Context, db *sql.DB) {
	timers, err := godb.ReadTimers(ctx, db) // Implement this function in your timer package
	if err != nil {
		fmt.Println("Error listing active timers:", err)
		return
	}

	// Determine the width of each column
	idWidth := 4    // Adjust as needed
	timeWidth := 25 // Adjust to fit the length of the formatted time
	nameWidth := 30 // Adjust based on your longest task name

	headerFormat := fmt.Sprintf("%%-%ds | %%-%ds | %%-%ds\n", idWidth, timeWidth, nameWidth)
	rowFormat := fmt.Sprintf("%%-%dd | %%-%ds | %%-%ds\n", idWidth, timeWidth, nameWidth)

	fmt.Printf(headerFormat, "ID", "Start Time", "Task Name")

	for _, t := range timers {
		fmt.Printf(rowFormat, t.ID, t.StartTime.Format(time.RFC3339), t.TaskName)
	}
}
