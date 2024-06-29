package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	godb "go-time/db" 
)


func EditCmd(db *sql.DB) *cobra.Command {
	var id int
	var name, description string
	var tags []string

	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit an existing time entry",
		Long:  `Edit an existing time entry by specifying its ID, name, and description.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			err := godb.EditEntry(ctx, db, id, name, description, tags)
			if err != nil {
				fmt.Println("Error editing time entry:", err)
				return
			}
			fmt.Println("Time entry updated")
		},
	}

	cmd.Flags().IntVarP(&id, "id", "i", 0, "ID of the time entry to edit")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Name of the time entry")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Description of the time entry")
	cmd.MarkFlagRequired("description")
	cmd.Flags().StringArrayVarP(&tags, "tags", "t", nil, "Tags for the time entry")

	return cmd
}
