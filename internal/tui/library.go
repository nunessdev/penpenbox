package tui

import tea "charm.land/bubbletea/v2"

type libraryModel struct{}

func NewLibraryModel() libraryModel {
	return libraryModel{}
}

func (m libraryModel) Update(msg tea.Msg) (libraryModel, tea.Cmd) {
	return m, nil
}

func (m libraryModel) View() string {
	return titleStyle.Render("Library") + "\n"
}
