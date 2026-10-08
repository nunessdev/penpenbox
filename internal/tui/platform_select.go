package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type platformItem struct {
	name     string
	disabled bool
}

type platformModel struct {
	platforms []platformItem // platforms in the list
	cursor    int            // which platform the cursor is pointing at
}

func NewPlatformModel() platformModel {
	return platformModel{
		platforms: []platformItem{
			{name: "Steam", disabled: false},
			{name: "GOG (Coming soon??)", disabled: true},
		},
	}
}

func (m platformModel) Update(msg tea.Msg) (platformModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.platforms)-1 {
				m.cursor++
			}
		case "enter":
			if !m.platforms[m.cursor].disabled {
				return m, func() tea.Msg { return openSteamSetupMsg{} }
			}
		}
	}
	return m, nil
}

func (m platformModel) View() string {
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

	return s
}
