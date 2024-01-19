package cmd

import (
	"context"
	"database/sql"
	"github.com/spf13/cobra"
	godb "go-time/db"
	"log"
)

// StartCmd creates a new start command.
func StopCmd(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "stop [task name]",
		Short: "Stop a the current timer",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			taskName := args[0]
			description := ""
			if err := godb.StopTimer(ctx, db, taskName, description); err != nil {
				log.Printf("error stopping timer: %v", err)
				// Handle the error appropriately
			}

		},
	}
}
