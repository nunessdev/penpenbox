package tui

import (
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/nunessdev/penpenbox/internal/config"
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

// screen identifies which screen is currently active.
type screen int

const (
	screenPlatform screen = iota // iota makes these 0, 1, 2...
	screenSteamSetup
	screenLibrary
)

// Messages that screens send to the root to request a switch.
type openSteamSetupMsg struct{}
type backToPlatformMsg struct{}
type configSubmittedMsg struct{ apiKey, steamID string }
type configSavedMsg struct{ err error }

// Model is the root model. It owns every screen and tracks which is active.
type Model struct {
	screen   screen
	platform platformModel
	steam    steamSetupModel
	library  libraryModel
}

func New() Model {
	m := Model{
		screen:   screenPlatform,
		platform: NewPlatformModel(),
		steam:    NewSteamSetupModel(),
		library:  NewLibraryModel(),
	}

	// If a valid config already exists, skip to library
	cfg, err := config.Load()
	if err == nil && cfg.HasSteam() {
		m.screen = screenLibrary
	}

	return m
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

	// Steam Setup
	case openSteamSetupMsg:
		m.steam = NewSteamSetupModel() // start with a fresh form each time
		m.screen = screenSteamSetup
		return m, nil

	// Back
	case backToPlatformMsg:
		m.screen = screenPlatform
		return m, nil

	// User submitted a config for Steam
	case configSubmittedMsg:
		return m, func() tea.Msg {
			cfg := config.Config{
				APIKey:  msg.apiKey,
				SteamID: msg.steamID,
			}
			err := config.Save(cfg)
			return configSavedMsg{err: err}
		}

	case configSavedMsg:
		if msg.err != nil {
			return m, nil // save failed: TODO - Display an error message
		}
		m.screen = screenLibrary
		return m, nil
	}

	// Everything else goes to the active screen.
	var cmd tea.Cmd
	switch m.screen {
	case screenPlatform:
		m.platform, cmd = m.platform.Update(msg)

	case screenSteamSetup:
		m.steam, cmd = m.steam.Update(msg)

	case screenLibrary:
		m.library, cmd = m.library.Update(msg)
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

	case screenLibrary:
		content = m.library.View()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}
