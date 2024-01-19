package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	godb "go-time/db" // Adjust this import path as necessary
	"strconv"
)

// EditCmd creates a new edit command.
func EditCmd(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "edit [id] [name] [description]",
		Short: "Edit an existing time entry",
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Invalid ID:", err)
				return
			}
			name := args[1]
			description := args[2]

			err = godb.EditEntry(ctx, db, id, name, description)
			if err != nil {
				fmt.Println("Error editing time entry:", err)
				return
			}
			fmt.Println("Time entry updated")
		},
	}
}
