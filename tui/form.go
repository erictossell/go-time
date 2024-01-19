package tui

import (
	"github.com/charmbracelet/huh"
)

func addEntryForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name"),
			huh.NewInput().
				Key("description").
				Title("Description"),
		),
	)
}

func editEntryForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name"),
			huh.NewInput().
				Key("description").
				Title("Description"),
			huh.NewInput().
				Key("startTime").
				Title("Start Time (YYYY-MM-DD HH:MM:SS)"),
			huh.NewInput().
				Key("endTime").
				Title("End Time (YYYY-MM-DD HH:MM:SS)"),
		),
	)
}
