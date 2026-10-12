package tui

import (
	"database/sql"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/nunessdev/penpenbox/internal/config"
	"github.com/nunessdev/penpenbox/internal/database"
	"github.com/nunessdev/penpenbox/internal/models"
	"github.com/nunessdev/penpenbox/internal/steam"
)

type libraryModel struct {
	db      *sql.DB
	syncing bool
	status  string
	width   int
	height  int
	games   []models.Game
	cursor  int // index of the selected game
	offset  int // index of the first visible game
}

func NewLibraryModel(db *sql.DB) libraryModel {
	return libraryModel{db: db}
}

// Sync games db functions
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

// Load games from db functions
type gamesLoadedMsg struct {
	games []models.Game
	err   error
}

func loadGamesCmd(db *sql.DB) tea.Cmd {
	return func() tea.Msg {
		games, err := database.ListGames(db)
		return gamesLoadedMsg{games: games, err: err}
	}
}

func (m libraryModel) Init() tea.Cmd {
	return loadGamesCmd(m.db)
}

// Update and View
func (m libraryModel) Update(msg tea.Msg) (libraryModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:

		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.games)-1 {
				m.cursor++
			}
		}

		// keep the cursor inside the visible window
		if m.cursor < m.offset {
			m.offset = m.cursor
		}
		if m.cursor >= m.offset+m.visibleRows() {
			m.offset = m.cursor - m.visibleRows() + 1
		}

		if msg.String() == "r" && !m.syncing {
			m.syncing = true
			m.status = "Syncing..."
			return m, syncCmd(m.db)
		}

	case gamesLoadedMsg:
		if msg.err != nil {
			m.status = "Load failed: " + msg.err.Error()
		} else {
			m.games = msg.games
			if m.cursor >= len(m.games) {
				m.cursor = max(0, len(m.games)-1)
			}
			m.status = fmt.Sprintf("%d games loaded", len(m.games))
		}

	case gamesSyncedMsg:
		m.syncing = false
		if msg.err != nil {
			m.status = "Sync failed: " + msg.err.Error()
		} else {
			m.status = fmt.Sprintf("Synced %d games", msg.count)
			return m, loadGamesCmd(m.db)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m libraryModel) visibleRows() int {
	n := m.height - 8
	if n < 1 {
		return 1
	}
	return n
}

func (m libraryModel) View() string {
	s := titleStyle.Render("Library") + "\n\n"

	// rows available: window height minus title (3 lines), blank line, status, footer
	visible := m.height - 8
	if visible < 0 {
		visible = 0
	}
	if visible > len(m.games) {
		visible = len(m.games)
	}

	end := m.offset + m.visibleRows()
	if end > len(m.games) {
		end = len(m.games)
	}

	for i := m.offset; i < end; i++ {
		style := textStyle
		if i == m.cursor {
			style = highlightedStyle
		}
		s += style.Render(m.games[i].Title) + "\n"
	}

	if m.status != "" {
		s += "\n" + textStyle.Render(m.status) + "\n"
	}

	s += textStyle.Render("↑/↓: navigate • r: resync library • q: quit") + "\n"
	return s
}
