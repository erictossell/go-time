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
func DelCmd(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "del [id]",
		Short: "Delete an existing time entry",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Invalid ID:", err)
				return
			}

			err = godb.DeleteEntry(ctx, db, id)
			if err != nil {
				fmt.Println("Error editing time entry:", err)
				return
			}
			fmt.Println("Time entry deleted.")
		},
	}
}
