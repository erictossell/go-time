package tui

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"

	"go-time/pkgs/entry"
	"go-time/pkgs/tag"
	"go-time/pkgs/timer"

	"go-time/pkgs/stopwatch"
	"go-time/pkgs/util"
	"os"
	"time"
)

type model struct {
	db            *sql.DB
	currentView   string
	entries       []entry.Entry
	timers        []timer.Timer
	tags          []tag.Tag
	keymap        keymap
	help          help.Model
	entriesCursor int
	timersCursor  int
	tagsCursor    int
	menuCursor    int
	stopwatch     stopwatch.Model
	form          *huh.Form
	formActive    bool
}

func Main(db *sql.DB) {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {

			}
		}(f)
	}
	p := tea.NewProgram(initialModel(db))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	err := m.updateTimers()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = m.updateEntries()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = m.updateTags()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	if m.formActive {

		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
			cmds = append(cmds, cmd)
		}
		if m.form.State == huh.StateCompleted {
			name := m.form.GetString("name")

			tags := m.form.Get("tags")

			// Convert tags to []string
			tagsParsed, ok := tags.([]string)
			if !ok {
				fmt.Println("Error: tags is not of type []string")
			}

			switch m.currentView {

			case "entries":

				if err != nil {
					fmt.Println("Error: ", err)
				}

				//tx, err := m.db.Begin()
				//entry.CreateEntry(context.Background(), tx, name, startTimeParsed, endTimeParsed, tagsParsed)
				if err != nil {
					fmt.Println("Error: ", err)
				}

				m.formActive = false

			case "timers":

				if err != nil {
					fmt.Println("Error: ", err)
				}
				action := func() {
					err := timer.CreateTimer(context.Background(), m.db, name, tagsParsed)
					if err != nil {
						fmt.Println("Error: ", err)
					}
					time.Sleep(1 * time.Second)
				}
				err := spinner.New().
					Title("Creating timer...").
					Action(action).
					Run()
				if err != nil {
					fmt.Println("Error: ", err)
				}

				m.formActive = false

			case "tags":
				err := tag.CreateTag(context.Background(), m.db, name)
				if err != nil {
					fmt.Println("Error: ", err)
				}

				m.formActive = false

			}
		} else {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				if msg.Type == tea.KeyEsc {
					m.form = tag.Form()
					m.formActive = false
				}
			}
		}

		return m, tea.Batch(cmds...)
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.add):
			tags, err := tag.GetTags(context.Background(), m.db)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			tagsStr := util.Map(tags, func(tag tag.Tag) string {
				return tag.Name
			})
			switch m.currentView {
			case "entries":
				m.form = entry.Form(tagsStr)

				m.formActive = true
			case "timers":
				m.form = timer.Form(tagsStr)

				m.formActive = true
			case "tags":
				m.form = tag.Form()

				m.formActive = true
			}
			return m, nil

		case key.Matches(msg, m.keymap.edit):
			tagsStr, err := tag.GetTagsAsStrArr(context.Background(), m.db)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			switch m.currentView {
			case "entries":
				e := m.entries[m.entriesCursor]
				m.form = entry.EditForm(e, tagsStr)
				m.formActive = true

			case "timers":
				t := m.timers[m.timersCursor]
				m.form = timer.EditForm(t, tagsStr)
				m.formActive = true

			case "tags":
				t := m.tags[m.tagsCursor]
				m.form = tag.EditForm(t)
				m.formActive = true
			}

		case key.Matches(msg, m.keymap.delete):
			switch m.currentView {
			case "entries":
				e := m.entries[m.entriesCursor]
				err := entry.DeleteEntry(context.Background(), m.db, e.ID)
				if err != nil {
					fmt.Println("Error: ", err)
				}
			case "timers":
				t := m.timers[m.timersCursor]
				err := timer.DeleteTimer(context.Background(), m.db, t.ID)

				if err != nil {
					fmt.Println("Error: ", err)
				}

			case "tags":
				t := m.tags[m.tagsCursor]
				err := tag.DeleteTag(context.Background(), m.db, t.ID)
				if err != nil {
					fmt.Println("Error: ", err)
				}
			}

		case key.Matches(msg, m.keymap.up):
			switch m.currentView {
			case "entries":
				if m.entriesCursor > 0 {
					m.entriesCursor--
				}
			case "timers", "timer":
				if m.timersCursor > 0 {
					m.timersCursor--
					cmd := m.startStopwatch(m.timers[m.timersCursor])
					return m, cmd
				}
			case "tags":
				if m.tagsCursor > 0 {
					m.tagsCursor--
				}
			}

		case key.Matches(msg, m.keymap.down):
			switch m.currentView {
			case "entries":
				if m.entriesCursor < len(m.entries)-1 {
					m.entriesCursor++
				}
			case "timers", "timer":
				if m.timersCursor < len(m.timers)-1 {
					m.timersCursor++
					cmd := m.startStopwatch(m.timers[m.timersCursor])
					return m, cmd
				}
			case "tags":
				if m.tagsCursor < len(m.tags)-1 {
					m.tagsCursor++
				}
			}

		case key.Matches(msg, m.keymap.left):
			m.navigateMenu(-1)
			return m, nil

		case key.Matches(msg, m.keymap.right):
			m.navigateMenu(1)
			if m.currentView == "timer" && len(m.timers) > 0 {
				cmd := m.startStopwatch(m.timers[m.timersCursor])
				return m, cmd
			}
			return m, nil

		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	var s string
	var err error
	if m.formActive {
		m.form.Init()
		return m.form.View()
	}
	switch m.currentView {
	case "entries":
		err = m.updateEntries()
		if err != nil {
			s += "Error: " + err.Error()
		} else {
			s += m.entriesView()
		}
	case "timers":
		err = m.updateTimers()
		if err != nil {
			s += "Error: " + err.Error()
		} else {
			s += m.timersView()
		}
	case "tags":
		err = m.updateTags()
		if err != nil {
			s += "Error: " + err.Error()
		} else {
			s += m.tagsView()
		}
	case "timer":
		s += m.timerView()
	}

	return s
}

func (m *model) updateEntries() error {
	ctx := context.Background()
	entries, err := entry.ReadEntries(ctx, m.db)
	if err != nil {
		return err
	}
	m.entries = entries
	return nil
}

func (m *model) updateTimers() error {
	ctx := context.Background()
	timers, err := timer.ReadTimers(ctx, m.db)
	if err != nil {
		return err
	}
	m.timers = timers
	return nil
}

func (m *model) updateTags() error {
	ctx := context.Background()
	tags, err := tag.GetTags(ctx, m.db)
	if err != nil {
		return err
	}
	m.tags = tags
	return nil
}

func (m *model) startStopwatch(timer timer.Timer) tea.Cmd {
	startTime := timer.StartTime
	elapsedTime := time.Since(startTime)
	m.stopwatch = stopwatch.New()
	m.stopwatch = m.stopwatch.SetElapsedTime(elapsedTime)
	cmd := m.stopwatch.Start()
	return cmd
}

func (m *model) navigateMenu(direction int) {
	menuItems := []string{"entries", "timers", "timer", "tags"}
	currentIndex := util.IndexOf(menuItems, m.currentView)
	if currentIndex != -1 {
		newIndex := (currentIndex + direction + len(menuItems)) % len(menuItems)
		newView := menuItems[newIndex]

		m.currentView = newView
	}
}
