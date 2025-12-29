package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/edellingham/nirimatic/internal/tui"
)

func main() {
	app := tui.NewApp()

	p := tea.NewProgram(
		app,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running nirimatic: %v\n", err)
		os.Exit(1)
	}
}
