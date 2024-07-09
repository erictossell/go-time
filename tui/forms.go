package tui

import (
	"github.com/charmbracelet/huh"
	"time"

	db "go-time/db"
)

func timerForm(tags []string) *huh.Form {

	options := make([]huh.Option[string], len(tags))
	for i, tag := range tags {
		options[i] = huh.NewOption(tag, tag)
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name"),
			huh.NewMultiSelect[string]().
				Key("tags").
				Title("Tags").
				Options(options...).
				Limit(3).
				Value(&tags),
		),
	)
}

func TimerForm(tags []string) *huh.Form {

	options := make([]huh.Option[string], len(tags))
	for i, tag := range tags {
		options[i] = huh.NewOption(tag, tag)
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name"),
			huh.NewMultiSelect[string]().
				Key("tags").
				Title("Tags").
				Options(options...).
				Limit(3).
				Value(&tags),
		),
	)
}

func timerEditForm(timer db.Timer, tags []string) *huh.Form {
	options := make([]huh.Option[string], len(tags))
	for i, tag := range tags {
		options[i] = huh.NewOption(tag, tag)
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name").
				Value(&timer.Name),
			huh.NewInput().
				Key("start_time").
				Title("Start Time").
				Value(timePtrToStringPtr(&timer.StartTime)),
			huh.NewMultiSelect[string]().
				Key("tags").
				Title("Tags").
				Options(options...).
				Limit(3).
				Value(&tags),
		),
	)

}

func entryForm(tags []string) *huh.Form {
	options := make([]huh.Option[string], len(tags))
	for i, tag := range tags {
		options[i] = huh.NewOption(tag, tag)
	}
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

func entryEditForm(entry db.Entry, tags []string) *huh.Form {
	options := make([]huh.Option[string], len(tags))
	for i, tag := range tags {
		options[i] = huh.NewOption(tag, tag)
	}
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name").
				Value(&entry.Name),
			huh.NewInput().
				Key("startTime").
				Title("Start Time (YYYY-MM-DD HH:MM:SS)").
				Value(timePtrToStringPtr(&entry.StartTime)), // Convert *time.Time to *string
			huh.NewInput().
				Key("endTime").
				Title("End Time (YYYY-MM-DD HH:MM:SS)").
				Value(timePtrToStringPtr(&entry.EndTime)), // Convert *time.Time to *string
			huh.NewMultiSelect[string]().
				Key("tags").
				Title("Tags").
				Options(options...).
				Limit(3).
				Value(&tags),
		),
	)
}

func tagForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name"),
		),
	)
}

func tagEditForm(tag db.Tag) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name").
				Value(&tag.Name),
		),
	)
}

func timePtrToStringPtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02 15:04:05")
	return &s
}
