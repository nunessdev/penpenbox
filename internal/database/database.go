package database

import (
	"database/sql"
	"fmt"

	"github.com/nunessdev/penpenbox/internal/models"
	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	return db, nil
}

func CreateDB(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS games (
		AppID INTEGER NOT NULL,
		Title TEXT NOT NULL,
		Playtime INTEGER,
		Platform TEXT NOT NULL,
		PRIMARY KEY (AppID, Platform)
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("create games table: %w", err)
	}

	return nil
}

func AddGame(db *sql.DB, game models.Game) error {
	query := `INSERT INTO games (AppID, Title, Playtime, Platform) VALUES (?, ?, ?, ?)
			  ON CONFLICT(AppID, Platform) DO UPDATE SET Title = excluded.Title, Playtime = excluded.Playtime`

	_, err := db.Exec(query, game.AppID, game.Title, game.Playtime, game.Platform)
	if err != nil {
		return fmt.Errorf("insert game: %w", err)
	}

	return nil
}

func DeleteGame(db *sql.DB, appid int, platform string) error {
	res, err := db.Exec(`DELETE FROM games WHERE AppID = ? AND Platform = ?`, appid, platform)
	if err != nil {
		return fmt.Errorf("delete game: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete game: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("delete game: no game with appid %d on %s", appid, platform)
	}

	return nil
}

func ListGames(db *sql.DB) ([]models.Game, error) {
	rows, err := db.Query(`SELECT AppID, Title, Playtime, Platform FROM games ORDER BY Title COLLATE NOCASE`)
	if err != nil {
		return nil, fmt.Errorf("list games: %w", err)
	}
	defer rows.Close()

	var games []models.Game
	for rows.Next() {
		var g models.Game
		if err := rows.Scan(&g.AppID, &g.Title, &g.Playtime, &g.Platform); err != nil {
			return nil, fmt.Errorf("scan game row: %w", err)
		}
		games = append(games, g)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate game rows: %w", err)
	}

	return games, nil
}
