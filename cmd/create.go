package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
	godb "go-time/db"
	"go-time/pkgs/tui"
	"log"
)

func CreateCmd(db *sql.DB) *cobra.Command {
	var name string
	var form *huh.Form

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
				tagsStr, err := godb.GetTagsAsStrArr(ctx, db)

				form = tui.TimerForm(tagsStr)
				form.Init()

				err = form.Run()
				if err != nil {
					log.Printf("Error running timer form: %v", err)
					return
				}

				name = form.GetString("name")
				tags := form.Get("tags")
				tagsParsed, ok := tags.([]string)
				if !ok {
					fmt.Println("Error: tags is not of type []string")
				}

				action := func() {
					err := godb.CreateTimer(ctx, db, name, tagsParsed)
					if err != nil {
						fmt.Println("Error: ", err)
					}
				}

				err = spinner.New().
					Title("Creating timer...").
					Action(action).
					Run()
				if err != nil {
					fmt.Println("Error: ", err)
				} else {
					fmt.Println("Timer started for task:", name)
				}

			case "entry":
				// Add logic to create an entry
				log.Println("Entry created.")
			case "tag":
				// Add logic to create a tag
				log.Println("Tag created.")
			default:
				log.Println("Invalid record type. Use the --type flag to specify 'entry', 'timer', or 'tag'.")
			}
		},
	}

	return cmd
}
