package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderCommand renders the interactive command mode interface
func (m *Model) RenderCommand(width, height int) string {
	boxWidth := width - 6
	if boxWidth > 82 {
		boxWidth = 82
	}
	if boxWidth < 40 {
		boxWidth = 40
	}

	var sb strings.Builder

	// Title
	sb.WriteString(m.theme.Title.Render("VFS COMMAND INTERFACE") + "\n")
	sb.WriteString(m.theme.DimText.Render("Type 'help' for commands (ls, cd, pwd, cat, tree, whoami, clear...) or Esc to exit") + "\n")
	sb.WriteString(m.theme.DimText.Render(strings.Repeat("─", boxWidth-4)) + "\n\n")

	// Output buffer
	bufferLines := m.cmdOutput
	maxLines := height - 12
	if maxLines < 4 {
		maxLines = 4
	}

	if len(bufferLines) > maxLines {
		bufferLines = bufferLines[len(bufferLines)-maxLines:]
	}

	for _, line := range bufferLines {
		sb.WriteString(line + "\n")
	}

	if len(bufferLines) == 0 {
		sb.WriteString(m.theme.DimText.Render("Ready for command input.\n\n"))
	}

	// Command input line
	prompt := m.theme.CommandPrompt.Render(fmt.Sprintf("%s:%s$ ", m.vfs.Profile.Handle, m.nav.CurrentPath))
	cursor := m.theme.CursorSymbol.Render("█")
	inputLine := fmt.Sprintf("%s%s%s", prompt, m.cmdInput, cursor)

	sb.WriteString("\n" + inputLine + "\n")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.ActiveBorder).
		Width(boxWidth).
		Padding(1, 2).
		Render(sb.String())

	return lipgloss.PlaceHorizontal(width, lipgloss.Center, box)
}
