package entry

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"go-time/pkgs/tag"

	"go-time/pkgs/util"
)

func Form(tags []string) *huh.Form {
	options := util.CreateTagOptions(tags)
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name"),
			huh.NewInput().
				Key("startTime").
				Title("Start Time (YYYY-MM-DD HH:MM:SS)"),
			huh.NewInput().
				Key("endTime").
				Title("End Time (YYYY-MM-DD HH:MM:SS)"),
			huh.NewMultiSelect[string]().
				Key("tags").
				Title("Tags").
				Options(options...).
				Limit(3).
				Value(&tags),
		),
	)
}

func EditForm(entry Entry, tags []string) *huh.Form {
	options := util.CreateTagOptions(tags)
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name").
				Value(&entry.Name),
			huh.NewInput().
				Key("startTime").
				Title("Start Time (YYYY-MM-DD HH:MM:SS)").
				Value(util.TimePtrToStringPtr(&entry.StartTime)), // Convert *time.Time to *string
			huh.NewInput().
				Key("endTime").
				Title("End Time (YYYY-MM-DD HH:MM:SS)").
				Value(util.TimePtrToStringPtr(&entry.EndTime)), // Convert *time.Time to *string
			huh.NewMultiSelect[string]().
				Key("tags").
				Title("Tags").
				Options(options...).
				Limit(3).
				Value(&tags),
		),
	)
}

func HandleForm(ctx context.Context, db *sql.DB) {
	// Fetch the tags from the database
	tagsStr, err := tag.GetTagsAsStrArr(ctx, db)
	if err != nil {
		log.Printf("Error fetching tags: %v", err)
		return
	}

	// Initialize the form with the tags
	form := Form(tagsStr)
	if err = form.Run(); err != nil {
		log.Printf("Error running entry form: %v", err)
		return
	}

	// Extract data from the form
	name := form.GetString("name")
	startTimeStr := form.GetString("startTime")
	endTimeStr := form.GetString("endTime")
	tags, ok := form.Get("tags").([]string)
	if !ok {
		fmt.Println("Error: tags is not of type []string")
		return
	}

	// Convert string times to time.Time objects
	startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr)
	if err != nil {
		fmt.Printf("Error parsing start time: %v\n", err)
		return
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
	if err != nil {
		fmt.Printf("Error parsing end time: %v\n", err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return
	}

	spinner := spinner.New().Title("Saving entry...")
	err = spinner.Action(func() {
		if err := CreateEntry(ctx, tx, name, startTime, endTime, tags); err != nil {
			log.Printf("Error saving entry: %v", err)
			tx.Rollback() // Ensure rollback on error
		} else {
			fmt.Println("Entry saved successfully for:", name)
			tx.Commit() // Commit only on success
		}
	}).Run()

	if err != nil {
		fmt.Println("Error during spinner action: ", err)
	}
}
