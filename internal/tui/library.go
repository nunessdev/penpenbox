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
	db            *sql.DB
	syncing       bool
	status        string
	width         int
	height        int
	games         []models.Game
	cursor        int // index of the selected game
	offset        int // index of the first visible game
	confirmDelete bool
}

func NewLibraryModel(db *sql.DB) libraryModel {
	return libraryModel{db: db}
}

// Sync games db functions
type gamesSyncedMsg struct {
	count int
	err   error
}

type gameDeletedMsg struct{ err error }

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

func deleteCmd(db *sql.DB, g models.Game) tea.Cmd {
	return func() tea.Msg {
		return gameDeletedMsg{err: database.DeleteGame(db, g.AppID, g.Platform)}
	}
}

func formatPlaytime(minutes *int) string {
	if minutes == nil || *minutes <= 0 {
		return "never played"
	}
	if *minutes < 60 {
		return fmt.Sprintf("%dm", *minutes)
	}
	return fmt.Sprintf("%.1fh", float64(*minutes)/60)
}

func truncate(s string, width int) string {
	r := []rune(s)
	if len(r) <= width {
		return s
	}
	return string(r[:width-3]) + "..."
}
func (m libraryModel) titleWidth() int {
	const (
		leftPad  = 4
		rightMax = 8 + 2 + 12 // platform + gap + playtime columns
		gap      = 2          // space between title and right part
		margin   = 4          // space kept on the right edge
	)
	w := m.width - leftPad - rightMax - gap - margin
	if w < 10 {
		return 10 // never let the title collapse completely
	}
	return w
}

func (m libraryModel) Init() tea.Cmd {
	return loadGamesCmd(m.db)
}

// Update and View
func (m libraryModel) Update(msg tea.Msg) (libraryModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:

		if m.confirmDelete && msg.String() != "d" {
			m.confirmDelete = false
			m.status = ""
		}

		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.games)-1 {
				m.cursor++
			}
		case "c":
			return m, func() tea.Msg { return backToPlatformMsg{} }

		case "d":
			if len(m.games) == 0 {
				break
			}
			if !m.confirmDelete {
				m.confirmDelete = true
				m.status = "Press d again to delete " + m.games[m.cursor].Title
				break
			}
			m.confirmDelete = false
			return m, deleteCmd(m.db, m.games[m.cursor])

		case "ctrl+c", "q":
			return m, tea.Quit
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

	case gameDeletedMsg:
		if msg.err != nil {
			m.status = "Delete failed: " + msg.err.Error()
		} else {
			m.status = "Game deleted"
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
		g := m.games[i]
		tw := m.titleWidth()
		line := fmt.Sprintf("%-*s  %-8s  %12s", tw, truncate(g.Title, tw), g.Platform, formatPlaytime(g.Playtime))
		s += style.Render(line) + "\n"
	}

	if m.status != "" {
		s += "\n" + textStyle.Render(m.status) + "\n"
	}

	s += textStyle.Render("↑/↓: navigate • r: resync library • c: configure platform • q: quit") + "\n"
	return s
}
