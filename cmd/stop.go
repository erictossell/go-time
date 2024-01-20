package cmd

import (
	"context"
	"database/sql"
	"github.com/spf13/cobra"
	godb "go-time/db"
	"log"
)

func StopCmd(db *sql.DB) *cobra.Command {
	var taskName string

	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the current timer for a task",
		Long:  `Stop the current timer for a task. Specify the task name using the --name flag.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			if taskName == "" {
				log.Println("Task name is required. Use the --name flag to specify the task name.")
				return
			}

			if err := godb.StopTimer(ctx, db, taskName); err != nil {
				log.Printf("Error stopping timer: %v", err)
			} else {
				log.Println("Timer stopped for task:", taskName)
			}
		},
	}

	cmd.Flags().StringVarP(&taskName, "name", "n", "", "Name of the task to stop")
	cmd.MarkFlagRequired("name")

	return cmd
}
