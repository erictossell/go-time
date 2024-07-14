package timer

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"go-time/pkgs/tag"

	"go-time/pkgs/util"
)

func Form(tags []string) *huh.Form {
	options := util.CreateTagOptions(tags)
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Key("name").Title("Name"),
			huh.NewMultiSelect[string]().Key("tags").Title("Tags").Options(options...).Limit(3).Value(&tags),
		),
	)
}

func EditForm(timer Timer, tags []string) *huh.Form {
	options := util.CreateTagOptions(tags)
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Key("name").Title("Name").Value(&timer.Name),
			huh.NewInput().Key("start_time").Title("Start Time").Value(util.TimePtrToStringPtr(&timer.StartTime)),
			huh.NewMultiSelect[string]().Key("tags").Title("Tags").Options(options...).Limit(3).Value(&tags),
		),
	)
}

func HandleForm(ctx context.Context, db *sql.DB) {
	tagsStr, err := tag.GetTagsAsStrArr(ctx, db)
	if err != nil {
		log.Printf("Error fetching tags: %v", err)
		return
	}

	form := Form(tagsStr)
	if err = form.Run(); err != nil {
		log.Printf("Error running timer form: %v", err)
		return
	}

	name := form.GetString("name")
	tags, ok := form.Get("tags").([]string)
	if !ok {
		fmt.Println("Error: tags is not of type []string")
		return
	}

	spinner := spinner.New().Title("Creating timer...")
	err = spinner.Action(func() {
		err := CreateTimer(ctx, db, name, tags)
		if err != nil {
			log.Printf("Error creating timer: %v", err)
		} else {
			fmt.Println("Timer started for task:", name)
		}
	}).Run()

	if err != nil {
		fmt.Println("Error: ", err)
	}
}
