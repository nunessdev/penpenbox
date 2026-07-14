package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

func NewPlatformModel() platformModel {
	return platformModel{
		platforms: []string{"Steam", "GOG (Coming soon??)"},
	}
}

func (m platformModel) Init() tea.Cmd {
	return nil
}

func (m platformModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyPressMsg:

		// What was the actual key pressed?
		switch msg.String() {

		// Exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// Cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// Cursor down
		case "down", "j":
			if m.cursor < len(m.platforms)-1 {
				m.cursor++
			}

			// Enter enters the screen for configuring the respective platform
			//case "enter", "space":
		}
	}
	return m, nil
}

func (m platformModel) View() tea.View {
	s := "Select a platform to configure:\n\n"

	for i, choice := range m.platforms {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	s += "\nPress q to quit.\n"

	v := tea.NewView(s)
	v.AltScreen = true
	return v
}
