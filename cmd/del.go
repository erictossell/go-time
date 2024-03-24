package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	godb "go-time/db" 
)


func DelCmd(db *sql.DB) *cobra.Command {
	var id int

	cmd := &cobra.Command{
		Use:   "del",
		Short: "Delete an existing time entry",
		Long:  `Delete an existing time entry by specifying its ID.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			err := godb.DeleteEntry(ctx, db, id)
			if err != nil {
				fmt.Println("Error deleting time entry:", err)
				return
			}
			fmt.Println("Time entry deleted.")
		},
	}

	cmd.Flags().IntVarP(&id, "id", "i", 0, "ID of the time entry to delete")
	cmd.MarkFlagRequired("id")

	return cmd
}
