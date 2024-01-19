package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	godb "go-time/db"
	"strings"
	"time"
)

func ReadCmd(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "read [timers|entries]",
		Short: "List all active timers or time entries",
		Long:  `Read command is used to list all active timers or time entries. If no argument is provided, it defaults to listing timers.`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			listType := "timers"
			if len(args) > 0 {
				listType = args[0]
			}

			switch listType {
			case "entries":
				readEntries(ctx, db)
			case "timers":
				readTimers(ctx, db)
			default:
				fmt.Println("Invalid argument. Please specify 'entries' or 'timers'.")
			}
		},
	}
}

func readEntries(ctx context.Context, db *sql.DB) {
	entries, err := godb.ReadEntries(ctx, db)
	if err != nil {
		fmt.Println("Error listing time entries:", err)
		return
	}

	// Additional function to get tags for each entry
	getTagsForEntry := func(entryID int) ([]string, error) {
		var tags []string
		query := `
        SELECT t.name
        FROM tags t
        INNER JOIN entry_tags et ON t.id = et.tag_id
        WHERE et.entry_id = ?`
		rows, err := db.Query(query, entryID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var tagName string
			if err := rows.Scan(&tagName); err != nil {
				return nil, err
			}
			tags = append(tags, tagName)
		}
		return tags, nil
	}

	// Determine the width of each column
	idWidth := 4
	nameWidth := 20

	timeWidth := 25
	tagsWidth := 20 // Adjust based on expected tag length

	// Adjusted format strings
	headerFormat := fmt.Sprintf("%%-%ds | %%-%ds | %%-%ds | %%-%ds | %%-%ds\n", idWidth, nameWidth, timeWidth, timeWidth, tagsWidth)
	rowFormat := fmt.Sprintf("%%-%dd | %%-%ds | %%-%ds | %%-%ds | %%-%ds\n", idWidth, nameWidth, timeWidth, timeWidth, tagsWidth)

	fmt.Printf(headerFormat, "ID", "Name", "Start Time", "End Time", "Tags")
	for _, entry := range entries {
		tags, err := getTagsForEntry(entry.ID)
		if err != nil {
			fmt.Println("Error fetching tags for entry:", err)
			continue
		}
		tagStr := strings.Join(tags, ", ")
		fmt.Printf(rowFormat, entry.ID, entry.Name, entry.StartTime.Format(time.RFC3339), entry.EndTime.Format(time.RFC3339), tagStr)
	}
}

func readTimers(ctx context.Context, db *sql.DB) {
	timers, err := godb.ReadTimers(ctx, db)
	if err != nil {
		fmt.Println("Error listing time entries:", err)
		return
	}

	// Additional function to get tags for each entry
	getTagsForTimer := func(timerID int) ([]string, error) {
		var tags []string
		query := `
        SELECT t.name
        FROM tags t
        INNER JOIN timer_tags tt ON t.id = tt.tag_id
        WHERE tt.timer_id = ?`
		rows, err := db.Query(query, timerID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var tagName string
			if err := rows.Scan(&tagName); err != nil {
				return nil, err
			}
			tags = append(tags, tagName)
		}
		return tags, nil
	}

	// Determine the width of each column
	idWidth := 4
	nameWidth := 20

	timeWidth := 25
	tagsWidth := 20 // Adjust based on expected tag length

	headerFormat := fmt.Sprintf("%%-%ds | %%-%ds | %%-%ds | %%-%ds \n", idWidth, nameWidth, timeWidth, tagsWidth)
	rowFormat := fmt.Sprintf("%%-%dd | %%-%ds | %%-%ds | %%-%ds \n", idWidth, nameWidth, timeWidth, tagsWidth)

	fmt.Printf(headerFormat, "ID", "Name", "Start Time", "Tags")

	for _, timer := range timers {
		tags, err := getTagsForTimer(timer.ID)
		if err != nil {
			fmt.Println("Error fetching tags for entry:", err)
			continue
		}
		tagStr := strings.Join(tags, ", ")
		fmt.Printf(rowFormat, timer.ID, timer.Name, timer.StartTime.Format(time.RFC3339), tagStr)
	}
}
