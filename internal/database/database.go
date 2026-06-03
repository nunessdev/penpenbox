package database

import (
	"database/sql"
	"fmt"

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

func AddGame(db *sql.DB, appid int, title string, playtime int, platform string) error {
	query := `INSERT INTO games (AppID, Title, Playtime, Platform) VALUES (?, ?, ?, ?)`

	_, err := db.Exec(query, appid, title, playtime, platform)
	if err != nil {
		return fmt.Errorf("Error inserting game data: %w\n", err)
	}

	return nil
}
