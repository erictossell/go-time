package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
)

func (m model) topBarView() string {
	view := "---------- Go-Time ---------- \n"
	menuItems := []string{"entries", "timers", "timer"}
	for _, item := range menuItems {
		if m.currentView == item {
			// Highlight the selected menu item. For example, using bold style.
			view += "[" + item + "]"
		} else {
			view += " " + item + " "
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

func (m model) entriesView() string {
	view := m.topBarView()
	err := m.updateEntries()
	if err != nil {
		return "Error: " + err.Error()
	}

	for i, entry := range m.entries {

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

func (m model) timersView() string {
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

func (m model) timerView() string {
	view := m.topBarView()
	timer := m.timers[m.cursor]
	line := fmt.Sprintf("ID: %d, Name: %s, Start: %s",
		timer.ID, timer.TaskName, timer.StartTime.Format("2006-01-02 15:04:05"))
	view += line + "\n"
	view += m.stopwatch.View() + "\n"

	view += m.helpView()
	return view
}
