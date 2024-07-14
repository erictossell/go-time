package tag

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

func Form() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name"),
		),
	)
}

func EditForm(tag Tag) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name").
				Value(&tag.Name),
		),
	)
}

func HandleForm(ctx context.Context, db *sql.DB) {
	// Initialize the form
	form := Form()
	if err := form.Run(); err != nil {
		log.Printf("Error running tag form: %v", err)
		return
	}

	// Extract data from the form
	tagName := form.GetString("name")

	// Perform the action with a spinner
	spinner := spinner.New().Title("Adding new tag...")
	err := spinner.Action(func() {
		err := CreateTag(ctx, db, tagName)
		if err != nil {
			log.Printf("Error adding new tag: %v", err)
		} else {
			fmt.Println("New tag added successfully:", tagName)
		}
	}).Run()

	if err != nil {
		fmt.Println("Error during spinner action: ", err)
	}
}
