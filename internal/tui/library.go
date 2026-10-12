package tui

import (
	"database/sql"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/nunessdev/penpenbox/internal/config"
	"github.com/nunessdev/penpenbox/internal/database"
	"github.com/nunessdev/penpenbox/internal/steam"
)

type libraryModel struct {
	db      *sql.DB
	syncing bool
	status  string
}

func NewLibraryModel(db *sql.DB) libraryModel {
	return libraryModel{db: db}
}

type gamesSyncedMsg struct {
	count int
	err   error
}

func syncCmd(db *sql.DB) tea.Cmd {
	return func() tea.Msg {
		cfg, err := config.Load()
		if err != nil {
			return gamesSyncedMsg{err: err}
		}

		games, err := steam.GetGames(cfg.APIKey, cfg.SteamID)
		if err != nil {
			return gamesSyncedMsg{err: err}
		}

		for _, g := range games {
			if err := database.AddGame(db, g); err != nil {
				return gamesSyncedMsg{err: err}
			}
		}

		return gamesSyncedMsg{count: len(games)}
	}
}

func (m libraryModel) Update(msg tea.Msg) (libraryModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "r" && !m.syncing {
			m.syncing = true
			m.status = "Syncing..."
			return m, syncCmd(m.db)
		}

	case gamesSyncedMsg:
		m.syncing = false
		if msg.err != nil {
			m.status = "Sync failed: " + msg.err.Error()
		} else {
			m.status = fmt.Sprintf("Synced %d games", msg.count)
		}
	}
	return m, nil
}

func (m libraryModel) View() string {
	s := titleStyle.Render("Library") + "\n\n"
	if m.status != "" {
		s += "\n" + textStyle.Render(m.status) + "\n\n"
	}
	s += textStyle.Render("↑/↓: navigate • r: resync library • q: quit") + "\n"
	return s
}
