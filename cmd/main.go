package main

import (
	"fmt"
	"os"

	"github.com/nunessdev/penpenbox/internal/database"
	"github.com/nunessdev/penpenbox/internal/steam"
)

func main() {

	// Get database path
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	path := home + "/.local/share/penpenbox/penpenbox.db"

	// Connect to database
	err = os.MkdirAll(home+"/.local/share/penpenbox", os.ModePerm)
	db, err := database.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	// Create database
	err = database.CreateDB(db)
	if err != nil {
		fmt.Println(err)
	}

	// Fetch and insert steam games in the database
	games, err := steam.GetGames(os.Getenv("STEAM_API_KEY"), os.Getenv("STEAM_ID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := range games {
		fmt.Printf("Inserting %s...\n", games[i].Title)
		err := database.AddGame(db, games[i])
		if err != nil {
			fmt.Println(err)
		}
	}

	// Test listing the games
	err = database.ListGames(db)
	if err != nil {
		fmt.Println(err)
	}
}
