package cmd

import (
	"database/sql"
	"github.com/spf13/cobra"
	godb "go-time/db"
)

// StartCmd creates a new start command.
func StartCmd(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "start [task name]",
		Short: "Start a new timer",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			taskName := args[0]
			godb.StartTimer(db, taskName)
		},
	}
}
