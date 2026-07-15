package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

// Styles for different types of text
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			PaddingTop(2).
			PaddingLeft(4)

	textStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			PaddingLeft(4)

	disabledStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BABABA")).
			PaddingLeft(4)

	highlightedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				PaddingLeft(4)
)

type platformItem struct {
	name     string
	disabled bool
}

func NewPlatformModel() platformModel {
	return platformModel{
		platforms: []platformItem{
			{name: "Steam", disabled: false},
			{name: "GOG (Coming soon??)", disabled: true},
		},
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
	s := titleStyle.Render("Select a platform to configure:") + "\n\n"

	for i, choice := range m.platforms {
		cursor := " "
		if m.cursor == i {
			// render cursor on selected item
			cursor = ">"
			if choice.disabled == false {
				// selected option renders with an highlight
				s += highlightedStyle.Render(fmt.Sprintf("%s %s", cursor, choice.name)) + "\n"
			} else {
				// disabled option appears grey when selected
				s += disabledStyle.Render(fmt.Sprintf("%s %s", cursor, choice.name)) + "\n"
			}
		} else {
			// option not selected, render plain text
			s += textStyle.Render(fmt.Sprintf("%s %s", cursor, choice.name)) + "\n"
		}
	}
	s += "\n" + textStyle.Render("Press q to quit.") + "\n"

	v := tea.NewView(s)
	v.AltScreen = true
	return v
}
