package tui

import (
	"strconv"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type addGameModel struct {
	title    textinput.Model
	platform textinput.Model
	playtime textinput.Model
	focused  int // 0 = title, 1 = platform, 2 = playtime
	errMsg   string
}

func NewAddGameModel() addGameModel {
	title := textinput.New()
	title.Placeholder = "Game title"
	title.SetWidth(40)
	title.Focus()

	platform := textinput.New()
	platform.Placeholder = "Platform (e.g. switch)"
	platform.SetWidth(40)

	playtime := textinput.New()
	playtime.Placeholder = "Playtime in hours (optional)"
	playtime.SetWidth(40)

	return addGameModel{title: title, platform: platform, playtime: playtime}
}

// Sent to the root when the form is valid.
type addGameSubmittedMsg struct {
	title    string
	platform string
	playtime *int // minutes, nil if left blank
}

func (m *addGameModel) setFocus(i int) {
	m.focused = i
	m.title.Blur()
	m.platform.Blur()
	m.playtime.Blur()
	switch i {
	case 0:
		m.title.Focus()
	case 1:
		m.platform.Focus()
	case 2:
		m.playtime.Focus()
	}
}

func (m addGameModel) Update(msg tea.Msg) (addGameModel, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "esc":
			return m, func() tea.Msg { return openLibraryMsg{} }

		case "up", "shift+tab":
			if m.focused > 0 {
				m.setFocus(m.focused - 1)
			}
			return m, nil

		case "down", "tab":
			if m.focused < 2 {
				m.setFocus(m.focused + 1)
			}
			return m, nil

		case "enter":
			title := strings.TrimSpace(m.title.Value())
			platform := strings.TrimSpace(m.platform.Value())
			if title == "" || platform == "" {
				m.errMsg = "Title and platform are required."
				return m, nil
			}

			var playtime *int
			if raw := strings.TrimSpace(m.playtime.Value()); raw != "" {
				hours, err := strconv.ParseFloat(raw, 64)
				if err != nil || hours < 0 {
					m.errMsg = "Playtime must be a number of hours, like 12.5."
					return m, nil
				}
				mins := int(hours * 60)
				playtime = &mins
			}

			m.errMsg = ""
			return m, func() tea.Msg {
				return addGameSubmittedMsg{title: title, platform: platform, playtime: playtime}
			}
		}
	}

	// Everything else goes to the focused input.
	var cmd tea.Cmd
	switch m.focused {
	case 0:
		m.title, cmd = m.title.Update(msg)
	case 1:
		m.platform, cmd = m.platform.Update(msg)
	case 2:
		m.playtime, cmd = m.playtime.Update(msg)
	}
	return m, cmd
}

func (m addGameModel) View() string {
	s := titleStyle.Render("Add a game") + "\n\n" +
		textStyle.Render(m.title.View()) + "\n" +
		textStyle.Render(m.platform.View()) + "\n" +
		textStyle.Render(m.playtime.View()) + "\n\n"

	if m.errMsg != "" {
		s += errorStyle.Render(m.errMsg) + "\n\n"
	}

	s += textStyle.Render("up/down: switch field • enter: add • esc: back") + "\n"
	return s
}
