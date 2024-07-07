package tui

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"time"

	db "go-time/db"
)

const maxWidth = 80

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(indigo).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.
		Foreground(red)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

var (
	tags []string
)

type state int

const (
	statusNormal state = iota
	stateDone
)

type Model struct {
	state state
	form  *huh.Form
}

func timerForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name"),
			huh.NewMultiSelect[string]().
				Key("tags").
				Title("Tags").
				Options(
					huh.NewOption("work", "Work"),
					huh.NewOption("personal", "Personal"),
					huh.NewOption("fun", "Fun"),
				).
				Limit(3).
				Value(&tags),
		),
	)
}

func timerEditForm(timer db.Timer) *huh.Form {
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
				Options(
					huh.NewOption("work", "Work"),
					huh.NewOption("personal", "Personal"),
					huh.NewOption("fun", "Fun"),
				).
				Limit(3).
				Value(&tags),
		),
	)

}

func entryForm() *huh.Form {
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
				Options(
					huh.NewOption("work", "Work"),
					huh.NewOption("personal", "Personal"),
					huh.NewOption("fun", "Fun"),
				).
				Limit(3).
				Value(&tags),
		),
	)
}

func entryEditForm(entry db.Entry) *huh.Form {
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
				Options(
					huh.NewOption("work", "Work"),
					huh.NewOption("personal", "Personal"),
					huh.NewOption("fun", "Fun"),
				).
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
