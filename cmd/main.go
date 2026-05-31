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
		return
	}

	// Create database
	err = database.CreateDB(db)
	if err != nil {
		return
	}

	// Just a test for steam games fetching
	games, err := steam.GetGames(os.Getenv("STEAM_API_KEY"), os.Getenv("STEAM_ID"))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for i := 0; i < len(games); i++ {
		fmt.Println(games[i].Name)
	}
}
