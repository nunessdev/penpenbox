package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	// "github.com/nunessdev/penpenbox/internal/database"
	// "github.com/nunessdev/penpenbox/internal/steam"
	"github.com/nunessdev/penpenbox/internal/database"
	"github.com/nunessdev/penpenbox/internal/tui"
)

func main() {
	db, err := database.Open()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer db.Close()

	p := tea.NewProgram(tui.New(db))
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
