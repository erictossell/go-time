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
		start:       key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "start timer")),
		stop:        key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "stop timer")),
		up:          key.NewBinding(key.WithKeys("k"), key.WithHelp("k", "up")),
		down:        key.NewBinding(key.WithKeys("j"), key.WithHelp("j", "down")),
		left:        key.NewBinding(key.WithKeys("h"), key.WithHelp("h", "left")),
		right:       key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "right")),
		list:        key.NewBinding(key.WithKeys("v"), key.WithHelp("v", "views")),
		quit:        key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
		selectTimer: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	}
	return model{
		db:          db,
		currentView: "timers",
		keymap:      keymap,
		help:        help.New(),
		cursor:      0,
	}
}

func (m model) Init() tea.Cmd {
	m.updateTimers()
	m.updateEntries()
	return m.stopwatch.Init()
}
