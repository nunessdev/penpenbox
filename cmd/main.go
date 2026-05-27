package main

import (
	"fmt"
	"os"

	"github.com/nunessdev/penpenbox/internal/steam"
)

func main() {
	games, err := steam.GetGames(os.Getenv("STEAM_API_KEY"), os.Getenv("STEAM_ID"))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for i := 0; i < len(games); i++ {
		fmt.Println(games[i].Name)
	}
}
