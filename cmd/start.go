package cmd

import (
	"database/sql"
	"github.com/spf13/cobra"
	"go-time/timer" // Adjust this import path as necessary
)

// StartCmd creates a new start command.
func StartCmd(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "start [task name]",
		Short: "Start a new timer",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			taskName := args[0]
			timer.StartTimer(db, taskName)
		},
	}
}
