package tui

import (
	tea "charm.land/bubbletea/v2"
)

// screen identifies which screen is currently active.
type screen int

const (
	screenPlatform screen = iota // iota makes these 0, 1, 2...
	screenSteamSetup
)

// Messages that screens send to the root to request a switch.
type openSteamSetupMsg struct{}
type backToPlatformMsg struct{}

// Model is the root model. It owns every screen and tracks which is active.
type Model struct {
	screen   screen
	platform platformModel
	steam    steamSetupModel
}

func New() Model {
	return Model{
		screen:   screenPlatform,
		platform: NewPlatformModel(),
		steam:    NewSteamSetupModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Messages meant for the root itself: quitting and screen switches.
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case openSteamSetupMsg:
		m.steam = NewSteamSetupModel() // start with a fresh form each time
		m.screen = screenSteamSetup
		return m, nil

	case backToPlatformMsg:
		m.screen = screenPlatform
		return m, nil
	}

	// Everything else goes to the active screen.
	var cmd tea.Cmd
	switch m.screen {
	case screenPlatform:
		m.platform, cmd = m.platform.Update(msg)
	case screenSteamSetup:
		m.steam, cmd = m.steam.Update(msg)
	}
	return m, cmd
}

func (m Model) View() tea.View {
	var content string

	switch m.screen {
	case screenPlatform:
		content = m.platform.View()
	case screenSteamSetup:
		content = m.steam.View()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}