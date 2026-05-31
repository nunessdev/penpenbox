package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateDB(db *sql.DB) error {
	const create string = `
	CREATE TABLE IF NOT EXISTS games (
		GameID INTEGER NOT NULL PRIMARY KEY,
		Title TEXT NOT NULL,
		Playtime INTEGER,
		Platform TEXT
	);`

	_, err := db.Exec(create)
	if err != nil {
		return err
	}

	return nil
}
