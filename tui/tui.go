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
	keymap        keymap
	help          help.Model
	entriesCursor int
	timersCursor  int
	menuCursor    int
	stopwatch     stopwatch.Model

	form       *huh.Form
	formActive bool
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
	err := m.updateTimers()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	err = m.updateEntries()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	if m.formActive {
		var cmd tea.Cmd
		updatedForm, cmd := m.form.Update(msg)
		if updatedForm, ok := updatedForm.(*huh.Form); ok {
			m.form = updatedForm
		}
		if m.form.State == huh.StateCompleted {
			name := m.form.GetString("name")

			err := godb.CreateTimer(context.Background(), m.db, name, []string{})
			if err != nil {
				fmt.Println("Error: ", err)
			}
			m.form = addEntryForm()
			m.formActive = false
			m.currentView = "timers"
		} else {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				if msg.Type == tea.KeyEsc {
					m.form = nil
					m.formActive = false
					return m, nil
				}
			}
		}

		return m, cmd
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.up):
			if m.currentView == "entries" && m.entriesCursor > 0 {
				m.entriesCursor--
			} else if (m.currentView == "timers" || m.currentView == "timer") && m.timersCursor > 0 {
				m.timersCursor--
				cmd := m.startStopwatch(m.timers[m.timersCursor])
				return m, cmd
			}

		case key.Matches(msg, m.keymap.down):
			if m.currentView == "entries" && m.entriesCursor < len(m.entries)-1 {
				m.entriesCursor++
			} else if (m.currentView == "timers" || m.currentView == "timer") && m.timersCursor < len(m.timers)-1 {
				m.timersCursor++
				cmd := m.startStopwatch(m.timers[m.timersCursor])
				return m, cmd
			}

		case key.Matches(msg, m.keymap.selectItem):
			switch m.currentView {
			case "entries":
				m.form = editEntryForm()
				m.formActive = true
			case "timers":
				m.form = addEntryForm()
				m.formActive = true
			}
			return m, nil

		case key.Matches(msg, m.keymap.left):
			m.navigateMenu(-1)
			return m, nil

		case key.Matches(msg, m.keymap.right):
			m.navigateMenu(1)
			if m.currentView == "timer" {
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

func (m *model) startStopwatch(timer godb.Timer) tea.Cmd {
	startTime := timer.StartTime
	elapsedTime := time.Since(startTime)
	m.stopwatch = stopwatch.New()
	m.stopwatch = m.stopwatch.SetElapsedTime(elapsedTime)
	cmd := m.stopwatch.Start()
	return cmd
}

func (m *model) navigateMenu(direction int) {
	menuItems := []string{"entries", "timers", "timer"}
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
