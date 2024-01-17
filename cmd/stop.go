package cmd

import (
	"database/sql"
	"github.com/spf13/cobra"
	godb "go-time/db"
)

// StartCmd creates a new start command.
func StopCmd(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "stop [task name]",
		Short: "Stop a the current timer",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			taskName := args[0]
			description := ""
			godb.StopTimer(db, taskName, description)
		},
	}
}
