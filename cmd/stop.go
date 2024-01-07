package cmd

import (
	"database/sql"
	"github.com/spf13/cobra"
	"go-time/timer"
)

// StartCmd creates a new start command.
func StopCmd(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "stop [task name]",
		Short: "Stop a the current timer",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			taskName := args[0]
			timer.StopTimer(db, taskName)
		},
	}
}
