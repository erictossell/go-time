package tui

import (
	"database/sql"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type keymap struct {
	start      key.Binding
	stop       key.Binding
	up         key.Binding
	down       key.Binding
	left       key.Binding
	right      key.Binding
	quit       key.Binding
	selectItem key.Binding
}

func initialModel(db *sql.DB) model {
	keymap := keymap{
		start: key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "start timer")),
		stop:  key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "stop timer")),
		up:    key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k", "up")),
		down:  key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j", "down")),
		left:  key.NewBinding(key.WithKeys("h", "left"), key.WithHelp("h", "left")),
		right: key.NewBinding(key.WithKeys("l", "right"), key.WithHelp("l", "right")),

		quit:       key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
		selectItem: key.NewBinding(key.WithKeys("enter", " "), key.WithHelp("enter", "select")),
	}
	return model{
		db:          db,
		currentView: "timers",
		keymap:      keymap,
		help:        help.New(),
		form:        addEntryForm(),
		formActive:  false,
	}
}

func (m model) Init() tea.Cmd {
	err := m.updateTimers()
	if err != nil {
		return tea.Quit
	}
	err = m.updateEntries()
	if err != nil {
		return tea.Quit
	}
	return m.stopwatch.Init()
}
