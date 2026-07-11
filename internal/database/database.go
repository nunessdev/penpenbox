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
		return nil, fmt.Errorf("Error opening database: %w\n", err)
	}
	return db, nil
}

func CreateDB(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS games (
		AppID TEXT NOT NULL,
		Title TEXT NOT NULL,
		Playtime INTEGER,
		Platform TEXT NOT NULL,
		PRIMARY KEY (AppID, Platform)
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("Error creating database: %w\n", err)
	}

	return nil
}

func AddGame(db *sql.DB, game models.Game) error {
	query := `INSERT INTO games (AppID, Title, Playtime, Platform) VALUES (?, ?, ?, ?)`

	_, err := db.Exec(query, game.AppID, game.Title, game.Playtime, game.Platform)
	if err != nil {
		return fmt.Errorf("Error inserting game data: %w", err)
	}

	return nil
}

func DeleteGame(db *sql.DB, appid int) error {
	query := `DELETE FROM games WHERE AppID = ?`

	_, err := db.Exec(query, appid)
	if err != nil {
		return fmt.Errorf("Error deleting game entry: %w", err)
	}

	return nil
}

func ListGames(db *sql.DB) error {
	query := `SELECT * FROM games`

	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("Error listing game library: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var appID int
		var title, platform string
		var playtime int

		err = rows.Scan(&appID, &title, &playtime, &platform)
		if err != nil {
			return fmt.Errorf("Error scanning game row: %w\n", err)
		}

		fmt.Printf("[%s] %s (playtime: %d min)\n", platform, title, playtime)
	}

	err = rows.Err()
	if err != nil {
		return fmt.Errorf("Error iterating game rows: %w\n", err)
	}

	return nil
}
