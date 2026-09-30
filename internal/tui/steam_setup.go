package tui

import (
	tea "charm.land/bubbletea/v2"
)

type steamSetupModel struct{}

func NewSteamSetupModel() steamSetupModel {
	return steamSetupModel{}
}

func (m steamSetupModel) Update(msg tea.Msg) (steamSetupModel, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "esc" {
		return m, func() tea.Msg { return backToPlatformMsg{} }
	}
	return m, nil
}

func (m steamSetupModel) View() string {
	return titleStyle.Render("Steam setup") + "\n\n" +
		textStyle.Render("(TODO - Implement this :p)") + "\n\n" +
		textStyle.Render("Press esc to go back.") + "\n"
}