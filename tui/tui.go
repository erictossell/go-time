package tui

import (
	"database/sql"
	"fmt"
	godb "go-time/db"
	"go-time/timer"

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
	db           *sql.DB
	currentView  string
	timeEntries  string
	activeTimers string
	choice       string
	keymap       keymap
	help         help.Model
	// ... other fields as needed
}

func initialModel(db *sql.DB) model {
	keymap := keymap{
		start: key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "start timer")),
		stop:  key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "stop timer")),
		list:  key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "list timers")),
		quit:  key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
	}

	return model{
		db:          db,
		currentView: "mainMenu",
		keymap:      keymap,
		help:        help.New(),
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
			entries, err := ListTimeEntriesForTUI(m.db)
			if err != nil {
				// Handle the error, e.g., set an error message in the model
				return m, nil
			}
			m.timeEntries = entries
			m.currentView = "timeEntries"
		case "2":
			timers, err := ListActiveTimersForTUI(m.db)
			if err != nil {
				// Handle the error
				return m, nil
			}
			m.activeTimers = timers
			m.currentView = "activeTimers"
		case "q":
			return m, tea.Quit	
		}

	}
	return m, nil
}
func (m model) View() string {
	var s string
	switch m.currentView {
	case "mainMenu":
		s = "Main Menu:\n1. View Time Entries\n2. View Active Timers\n" + m.helpView()
	case "timeEntries":
		s = "Time Entries:\n" + m.timeEntries + m.helpView()
	case "activeTimers":
		s = "Active Timers:\n" + m.activeTimers + m.helpView()
		// ... other views
	}
	return s
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
func ListActiveTimersForTUI(db *sql.DB) (string, error) {
	timers, err := timer.GetActiveTimers(db)
	if err != nil {
		return "", err
	}

	var result string
	for _, t := range timers {
		result += fmt.Sprintf("Task: %s, Started at: %s\n", t.TaskName, t.StartTime.Format("2006-01-02 15:04:05"))
	}

	return result, nil
}

func ListTimeEntriesForTUI(db *sql.DB) (string, error) {
	entries, err := godb.ListTimeEntries(db)
	if err != nil {
		return "", err
	}

	var result string
	for _, entry := range entries {
		result += fmt.Sprintf("ID: %d, Name: %s, Description: %s, Start: %s, End: %s\n",
			entry.ID, entry.Name, entry.Description, entry.StartTime.Format("2006-01-02 15:04:05"),
			entry.EndTime.Format("2006-01-02 15:04:05"))
	}

	return result, nil
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
