package tui

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	godb "go-time/db"
	"go-time/stopwatch"
	"os"
	"time"
)

type model struct {
	db            *sql.DB
	currentView   string
	entries       []godb.Entry
	timers        []godb.Timer
	tags          []godb.Tag
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
		defer f.Close()
	}
	p := tea.NewProgram(initialModel(db))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				//godb.CreateEntry(context.Background(), tx, name, startTimeParsed, endTimeParsed, tagsParsed)
				if err != nil {
					fmt.Println("Error: ", err)
				}

				m.formActive = false

			case "timers":

				if err != nil {
					fmt.Println("Error: ", err)
				}

				err = godb.CreateTimer(context.Background(), m.db, name, tagsParsed)
				if err != nil {
					fmt.Println("Error: ", err)
				}

				m.formActive = false

			case "tags":
				err := godb.CreateTag(context.Background(), m.db, name)
				if err != nil {
					fmt.Println("Error: ", err)
				}

				m.formActive = false

			}
		} else {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				if msg.Type == tea.KeyEsc {
					m.form = tagForm()
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
			tags, err := godb.GetTags(context.Background(), m.db)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			tagsStr := make([]string, len(tags))
			switch m.currentView {
			case "entries":
				m.form = entryForm(tagsStr)
				m.formActive = true
			case "timers":
				m.form = timerForm(tagsStr)
				m.formActive = true
			case "tags":
				m.form = tagForm()
				m.formActive = true
			}
			return m, nil

		case key.Matches(msg, m.keymap.edit):
			tags, err := godb.GetTags(context.Background(), m.db)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			tagsStr := make([]string, len(tags))

			switch m.currentView {
			case "entries":
				entry := m.entries[m.entriesCursor]
				m.form = entryEditForm(entry, tagsStr)
				m.formActive = true

			case "timers":
				timer := m.timers[m.timersCursor]
				m.form = timerEditForm(timer, tagsStr)
				m.formActive = true

			case "tags":
				tag := m.tags[m.tagsCursor]
				m.form = tagEditForm(tag)
				m.formActive = true

			}

		case key.Matches(msg, m.keymap.delete):
			switch m.currentView {
			case "entries":
				entry := m.entries[m.entriesCursor]
				err := godb.DeleteEntry(context.Background(), m.db, entry.ID)
				if err != nil {
					fmt.Println("Error: ", err)
				}
			case "timers":
				timer := m.timers[m.timersCursor]
				err := godb.DeleteTimer(context.Background(), m.db, timer.ID)

				if err != nil {
					fmt.Println("Error: ", err)
				}

			case "tags":
				tag := m.tags[m.tagsCursor]
				err := godb.DeleteTag(context.Background(), m.db, tag.ID)
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

func (m model) View() string {
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
	entries, err := godb.ReadEntries(ctx, m.db)
	if err != nil {
		return err
	}
	m.entries = entries
	return nil
}

func (m *model) updateTimers() error {
	ctx := context.Background()
	timers, err := godb.ReadTimers(ctx, m.db)
	if err != nil {
		return err
	}
	m.timers = timers
	return nil
}

func (m *model) updateTags() error {
	ctx := context.Background()
	tags, err := godb.GetTags(ctx, m.db)
	if err != nil {
		return err
	}
	m.tags = tags
	return nil
}

func (m *model) startStopwatch(timer godb.Timer) tea.Cmd {
	startTime := timer.StartTime
	elapsedTime := time.Since(startTime)
	m.stopwatch = stopwatch.New()
	m.stopwatch = m.stopwatch.SetElapsedTime(elapsedTime)
	cmd := m.stopwatch.Start()
	return cmd
}

func (m *model) navigateMenu(direction int) {
	menuItems := []string{"entries", "timers", "timer", "tags"}
	currentIndex := indexOf(menuItems, m.currentView)
	if currentIndex != -1 {
		newIndex := (currentIndex + direction + len(menuItems)) % len(menuItems)
		newView := menuItems[newIndex]

		m.currentView = newView
	}
}

func indexOf(slice []string, value string) int {
	for i, item := range slice {
		if item == value {
			return i
		}
	}
	return -1
}
