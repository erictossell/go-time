package cmd

import (
	"context"
	"database/sql"
	"github.com/spf13/cobra"
	godb "go-time/db"
	"log"
)

func StopCmd(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "stop [task name] [tags]",
		Short: "Stop the current timer and add tags",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			taskName := args[0]
			tags := args[1:]  // Remaining arguments are considered as tags
			description := "" // Update or handle this as per your requirement

			if err := godb.StopTimer(ctx, db, taskName, description, tags); err != nil {
				log.Printf("error stopping timer: %v", err)
			}
		},
	}
}
