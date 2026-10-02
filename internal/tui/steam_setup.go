package tui

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type steamSetupModel struct {
	apiKey  textinput.Model
	steamID textinput.Model
	focused int // 0 = apiKey, 1 = steamID
}

func NewSteamSetupModel() steamSetupModel {
	apiKey := textinput.New()
	apiKey.Placeholder = "Steam API key"
	apiKey.SetWidth(40)
	apiKey.Focus()

	steamID := textinput.New()
	steamID.Placeholder = "Steam User ID"
	steamID.SetWidth(40)

	return steamSetupModel{apiKey: apiKey, steamID: steamID}
}

func (m steamSetupModel) Update(msg tea.Msg) (steamSetupModel, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "esc":
			return m, func() tea.Msg { return backToPlatformMsg{} }

		case "up":
			if m.focused > 0 {
				m.focused = 0
				m.steamID.Blur()
				m.apiKey.Focus()
			}
			return m, nil

		case "down":
			if m.focused < 1 {
				m.focused = 1
				m.apiKey.Blur()
				m.steamID.Focus()
			}
			return m, nil

		case "enter":
			// Submit m.apiKey.Value() and m.steamID.Value()
			apiKey := strings.TrimSpace(m.apiKey.Value())
			steamID := strings.TrimSpace(m.steamID.Value())
			if apiKey == "" || steamID == "" {
				return m, nil // ignore until both fields are filled
			}
			return m, nil
		}
	}

	// Everything else goes to the focused input.
	var cmd tea.Cmd
	if m.focused == 0 {
		m.apiKey, cmd = m.apiKey.Update(msg)
	} else {
		m.steamID, cmd = m.steamID.Update(msg)
	}
	return m, cmd
}

func (m steamSetupModel) View() string {
	return titleStyle.Render("Steam setup") + "\n\n" +
		textStyle.Render("Steam API requires a Steam API key to work. Refer to https://github.com/nunessdev/penpenbox for instructions on how to get one!") + "\n" +
		textStyle.Render("It also requires the Steam User ID of the profile you want to track (make sure your profile is set to public)!") + "\n\n" +
		textStyle.Render(m.apiKey.View()) + "\n" +
		textStyle.Render(m.steamID.View()) + "\n" +
		textStyle.Render("The config will be saved on ~/.config/penpenbox/") + "\n\n" +
		textStyle.Render("up/down: switch field • enter: submit • esc: back") + "\n"
}
