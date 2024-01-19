package tui

import (
	"context"
	"database/sql"
	"fmt"
	godb "go-time/db"

	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"
)

type keymap struct {
	start key.Binding
	stop  key.Binding
	up    key.Binding
	down  key.Binding
	list  key.Binding
	quit  key.Binding
}

type model struct {
	db          *sql.DB
	currentView string
	timeEntries []godb.TimeEntry
	timers      []godb.Timer
	keymap      keymap
	help        help.Model
	cursor      int
	selected    map[int]struct{}
	// ... other fields as needed
}

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
		currentView: "mainMenu",
		keymap:      keymap,
		help:        help.New(),
		selected:    make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil // Return `nil`, which means "no I/O right now, please."
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor > len(m.timeEntries)-1 {
				m.cursor++
			}
		case "enter", " ": // Edit selected entry
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}

		case "v":
			if m.currentView == "timeEntries" {
				m.currentView = "activeTimers"
				// Load active timers
			} else {
				m.currentView = "timeEntries"
				// Load time entries
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	var s string
	var err error
	switch m.currentView {
	case "timeEntries":
		err = m.updateTimeEntries()
		if err != nil {
			s += "Error: " + err.Error()
		} else {
			s += m.timeEntriesView()
		}
	case "activeTimers":
		err = m.updateTimers()
		if err != nil {
			s += "Error: " + err.Error()
		} else {
			s += m.activeTimersView()
		}
	case "mainMenu":
		s += m.mainMenuView()
	}
	return s

}

func (m model) timeEntriesView() string {
	view := m.topBarView()
	err := m.updateTimeEntries()
	if err != nil {
		return "Error: " + err.Error()
	}

	for i, entry := range m.timeEntries {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}
		line := fmt.Sprintf("%s ID: %d, Name: %s, Description: %s, Start: %s, End: %s",
			cursor, entry.ID, entry.Name, entry.Description, entry.StartTime.Format("2006-01-02 15:04:05"),
			entry.EndTime.Format("2006-01-02 15:04:05"))

		view += line + "\n"
	}
	view += m.helpView()
	return view
}

func (m model) activeTimersView() string {
	view := m.topBarView()
	err := m.updateTimers()
	if err != nil {
		return "Error: " + err.Error()
	}

	for i, timer := range m.timers {
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}
		line := fmt.Sprintf("%s ID: %d, Name: %s, Start: %s",
			cursor, timer.ID, timer.TaskName, timer.StartTime.Format("2006-01-02 15:04:05"))

		view += line + "\n"

	}
	view += m.helpView()
	return view
}

func (m model) mainMenuView() string {
	view := m.topBarView()
	view += "-------------------------\n"
	view += "[s] Start a new timer\n"
	view += "[t] Stop the current timer\n"
	view += "[v] Change Views \n"
	view += "[q] Quit\n"
	view += "\nSelect an option: "

	return view
}

func (m model) topBarView() string {
	view := "Go-Time - Your CLI Time Tracker\n"
	view += strconv.Itoa(m.cursor) + " " + strconv.Itoa(len(m.timeEntries)) + "\n"
	menuItems := []string{"timeEntries", "activeTimers", "mainMenu"}
	for _, item := range menuItems {
		if m.currentView == item {
			// Highlight the selected menu item. For example, using bold style.
			view += item + " "
		} else {
			view += item + " "
		}
	}
	return view + "\n"
}

func (m model) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.up,
		m.keymap.down,
		m.keymap.list,
		m.keymap.quit,
	})
}

func (m *model) updateTimeEntries() error {
	ctx := context.Background()
	entries, err := godb.ListTimeEntries(ctx, m.db)
	if err != nil {
		return err
	}
	m.timeEntries = entries
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
