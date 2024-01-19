package tui

import (
	"database/sql"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	//"go-time/stopwatch"
)

func initialModel(db *sql.DB) model {
	keymap := keymap{
		start: key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "start timer")),
		stop:  key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "stop timer")),
		up:    key.NewBinding(key.WithKeys("k"), key.WithHelp("k", "up")),
		down:  key.NewBinding(key.WithKeys("j"), key.WithHelp("j", "down")),
		list:  key.NewBinding(key.WithKeys("v"), key.WithHelp("v", "views")),
		quit:  key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
	}
	return model{
		db:          db,
		currentView: "timers",
		keymap:      keymap,
		help:        help.New(),
		selected:    make(map[int]struct{}),
		cursor:      0,
	}
}

func (m model) Init() tea.Cmd {
	m.updateTimers()
	m.updateEntries()
	return nil
}
