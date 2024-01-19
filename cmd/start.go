package cmd

import (
	"context"
	"database/sql"
	"github.com/spf13/cobra"
	godb "go-time/db"
	"log"
)

// StartCmd creates a new start command.
func StartCmd(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "start [task name]",
		Short: "Start a new timer",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			taskName := args[0]
			if err := godb.StartTimer(ctx, db, taskName); err != nil {
				log.Printf("error starting timer: %v", err)
				// Handle the error appropriately
			}
		},
	}
}
