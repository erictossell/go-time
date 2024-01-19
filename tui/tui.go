package tui

import (
	"context"
	"database/sql"
	"fmt"
	godb "go-time/db"
	//"go-time/stopwatch"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	db          *sql.DB
	currentView string
	entries     []godb.Entry
	timers      []godb.Timer
	keymap      keymap
	help        help.Model
	cursor      int
	selected    map[int]struct{}
	// ... other fields as needed
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
	m.updateTimers()
	m.updateEntries()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.up):
			if m.cursor > 0 {
				m.cursor--
				return m, nil
			}

		case key.Matches(msg, m.keymap.down):
			if m.cursor < len(m.entries)-1 {
				m.cursor++
				return m, nil
			}

		case key.Matches(msg, m.keymap.list):
			// Switch view
			if m.currentView == "entries" {
				m.currentView = "timers"
			} else {
				m.currentView = "entries"
			}
			return m, nil

		case key.Matches(msg, m.keymap.quit):
			// Quit application
			return m, tea.Quit
		}
	}
	return m, nil
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
	}
	return s
}

func (m *model) updateEntries() error {
	ctx := context.Background()
	entries, err := godb.ListEntries(ctx, m.db)
	if err != nil {
		return err
	}
	m.entries = entries
	return nil
}

func (m *model) updateTimers() error {
	ctx := context.Background()
	timers, err := godb.ListTimers(ctx, m.db)
	if err != nil {
		return err
	}
	m.timers = timers
	return nil
}
