package cmd

import (
	"context"
	"database/sql"

	"github.com/spf13/cobra"

	"go-time/pkgs/entry"
	"go-time/pkgs/tag"
	"go-time/pkgs/timer"

	"log"
)

func CreateCmd(db *sql.DB) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "create [record type]",
		Short: "Create a new record (entry, timer, or tag)",
		Long:  `Create a new record by specifying the type (entry, timer, or tag).`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			recordType := args[0]
			ctx := context.Background()

			switch recordType {
			case "timer":
				timer.HandleForm(ctx, db)

			case "entry":
				tags, err := tag.GetTagsAsStrArr(ctx, db)
				if err != nil {
					log.Fatal(err)
				}
				entry.HandleForm(ctx, db, tags)

			case "tag":
				tag.HandleForm(ctx, db)

			default:
				log.Println("Invalid record type. Use the --type flag to specify 'entry', 'timer', or 'tag'.")
			}
		},
	}

	return cmd
}
