package tui

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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
	stopwatch     stopwatch.Model
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

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.up):
			if m.currentView == "entries" && m.entriesCursor > 0 {
				m.entriesCursor--
			} else if m.currentView == "timers" && m.timersCursor > 0 {
				m.timersCursor--
			}

		case key.Matches(msg, m.keymap.down):
			if m.currentView == "entries" && m.entriesCursor < len(m.entries)-1 {
				m.entriesCursor++
			} else if m.currentView == "timers" && m.timersCursor < len(m.timers)-1 {
				m.timersCursor++
			}
		case key.Matches(msg, m.keymap.selectTimer):
			selectedTimer := m.timers[m.timersCursor]
			cmd := m.startStopwatch(selectedTimer)
			m.currentView = "timer"
			return m, cmd

		case key.Matches(msg, m.keymap.left):
			m.navigateMenu(-1)
			return m, nil

		case key.Matches(msg, m.keymap.right):
			m.navigateMenu(1)
			if m.currentView == "timer" {
				selectedTimer := m.timers[m.timersCursor]
				cmd := m.startStopwatch(selectedTimer)
				return m, cmd

			}
			return m, nil

		case key.Matches(msg, m.keymap.quit):
			// Quit application
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
	startTime := timer.StartTime // Assuming timer has a StartTime field
	elapsedTime := time.Since(startTime)
	m.stopwatch = stopwatch.New()
	m.stopwatch = m.stopwatch.SetElapsedTime(elapsedTime)
	cmd := m.stopwatch.Start() // Start the stopwatch
	return cmd
	// Handle the command generated by the Start method if necessary
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
