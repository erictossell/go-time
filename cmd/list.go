package cmd

import (
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	godb "go-time/db"
	"time"
)

// ListCmd creates a new list command.
func ListCmd(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all time entries",
		Run: func(cmd *cobra.Command, args []string) {
			entries, err := godb.ListTimeEntries(db)
			if err != nil {
				fmt.Println("Error listing time entries:", err)
				return
			}
			for _, entry := range entries {
				fmt.Printf("%d | %s | %s | %s | %s\n", entry.ID, entry.Name, entry.Description, entry.StartTime.Format(time.RFC3339), entry.EndTime.Format(time.RFC3339))
			}
		},
	}
}
