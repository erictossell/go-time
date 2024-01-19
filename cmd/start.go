package cmd

import (
	"context"
	"database/sql"
	"github.com/spf13/cobra"
	godb "go-time/db"
	"log"
)

func StartCmd(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "start [task name] [tags]",
		Short: "Start a new timer with optional tags",
		Args:  cobra.MinimumNArgs(1), // At least the task name is required
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			taskName := args[0]
			tags := args[1:] // Remaining arguments are considered as tags

			if err := godb.CreateTimer(ctx, db, taskName, tags); err != nil {
				log.Printf("error starting timer: %v", err)
			}
		},
	}
}
