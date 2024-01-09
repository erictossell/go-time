package tui

import (
	"database/sql"
	"fmt"
	godb "go-time/db"
	"go-time/timer"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"
)

type tickMsg struct{}

type keymap struct {
	start key.Binding
	stop  key.Binding
	list  key.Binding
	quit  key.Binding
}

type model struct {
	db            *sql.DB
	currentView   string
	timeEntries   string
	activeTimers  string
	choice        string
	keymap        keymap
	help          help.Model
	selectedEntry int
	// ... other fields as needed
}

type TimeEntry struct {
	ID          int
	Name        string
	Description string
	StartTime   time.Time
	EndTime     time.Time
}

func initialModel(db *sql.DB) model {
	keymap := keymap{
		start: key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "start timer")),
		stop:  key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "stop timer")),
		list:  key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "lists")),
		quit:  key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
	}

	return model{
		db:            db,
		currentView:   "mainMenu",
		keymap:        keymap,
		help:          help.New(),
		selectedEntry: -1,
	}
}

func (m model) Init() tea.Cmd {
	return nil // Return `nil`, which means "no I/O right now, please."
}

func Main(db *sql.DB) {
	p := tea.NewProgram(initialModel(db))
	if _, err := p.Run(); err != nil {
		// handle error
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "l":
			if m.currentView == "timeEntries" {
				m.currentView = "activeTimers"
				// Load active timers
			} else {
				m.currentView = "timeEntries"
				// Load time entries
			}
		case "q":
			return m, tea.Quit
		}
		switch m.currentView {
		case "timeEntries":
			switch msg.String() {
			case "j", "down": // Move selection down
				m.selectedEntry = (m.selectedEntry + 1) % len(m.timeEntries)
			case "k", "up": // Move selection up
				if m.selectedEntry > 0 {
					m.selectedEntry--
				} else {
					m.selectedEntry = len(m.timeEntries) - 1
				}
			case "enter": // Edit selected entry
				// Add logic to edit the selected entry
			}
		}

	}
	return m, nil
}
func (m model) View() string {
	var s string
	switch m.currentView {
	case "timeEntries":
		return m.timeEntriesView()
	case "activeTimers":
		s = "Active Timers:\n" + m.activeTimers + m.helpView()
		// ... other views

	case "mainMenu":
		return m.mainMenuView()
	}
	return s
}

func (m model) timeEntriesView() string {
	var menu string
	entries, err := ListTimeEntriesForTUI(m.db)
	if err != nil {
		return "Error: " + err.Error()
	}

	for i, entry := range entries {
		line := fmt.Sprintf("ID: %d, Name: %s, Description: %s, Start: %s, End: %s",
			entry.ID, entry.Name, entry.Description, entry.StartTime.Format("2006-01-02 15:04:05"),
			entry.EndTime.Format("2006-01-02 15:04:05"))

		if i == m.selectedEntry {
			menu += "> " + line + "\n"
		} else {
			menu += "  " + line + "\n"
		}
	}
	menu += m.helpView()
	return menu
}

func (m model) activeTimersView() string {
	var menu string
	entries, err := ListActiveTimersForTUI(m.db)
	if err != nil {
		return "Error: " + err.Error()
	}

	for i, timer := range entries {
		line := fmt.Sprintf("ID: %d, Name: %s, Description: %s, Start: %s, End: %s",
			timer.TaskName, timer.StartTime.Format("2006-01-02 15:04:05"))

		if i == m.selectedEntry {
			menu += "> " + line + "\n"
		} else {
			menu += "  " + line + "\n"
		}
	}
	menu += m.helpView()
	return menu
}

func (m model) mainMenuView() string {
	var menu string
	menu += "Go-Time Tracker Main Menu\n"
	menu += "-------------------------\n"
	menu += "[s] Start a new timer\n"
	menu += "[t] Stop the current timer\n"
	menu += "[l] List time entries\n"
	menu += "[a] List active timers\n"
	menu += "[q] Quit\n"
	menu += "\nSelect an option: "

	return menu
}

func (m model) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.list,
		m.keymap.quit,
	})
}

func (m *keymap) bindings() []key.Binding {
	return []key.Binding{
		m.start,
		m.stop,
		m.list,
		m.quit,
	}
}
func ListActiveTimersForTUI(db *sql.DB) ([]timer.ActiveTimer, error) {
	return timer.ListActiveTimers(db)
}

func ListTimeEntriesForTUI(db *sql.DB) ([]godb.TimeEntry, error) {
	return godb.ListTimeEntries(db)
}

func StartTimerForTUI(db *sql.DB, taskName string) string {
	if timer.IsTimerRunning(db, taskName) {
		return "Timer is already running for task: " + taskName
	}

	timer.StartTimer(db, taskName)
	return "Timer started for task: " + taskName
}

func StopTimerForTUI(db *sql.DB, taskName, description string) string {
	timer.StopTimer(db, taskName, description)
	return "Timer stopped for task: " + taskName
}
