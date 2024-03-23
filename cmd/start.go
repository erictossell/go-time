package cmd

import (
	"context"
	"database/sql"
	"github.com/spf13/cobra"
	godb "go-time/db"
	"log"
)

func StartCmd(db *sql.DB) *cobra.Command {
	var taskName string
	var tags []string

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start a new timer with optional tags",
		Long:  `Start a new timer for a task with optional tags. Specify the task name and tags using flags.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			if taskName == "" {
				log.Println("Task name is required. Use the --name flag to specify the task name.")
				return
			}

			if err := godb.CreateTimer(ctx, db, taskName, tags); err != nil {
				log.Printf("Error starting timer: %v", err)
			} else {
				log.Println("Timer started for task:", taskName)
			}
		},
	}

	cmd.Flags().StringVarP(&taskName, "name", "n", "", "Name of the task")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringArrayVarP(&tags, "tags", "t", nil, "Tags for the timer")

	return cmd
}
