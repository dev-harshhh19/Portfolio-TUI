package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dev-harshhh19/harshad-tui/internal/tui"
)

func main() {
	m := tui.NewModel()
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running Harshad TUI: %v\n", err)
		os.Exit(1)
	}
}
