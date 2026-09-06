package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderHelp displays the compact, elegant help overlay
func (m *Model) RenderHelp(width, height int) string {
	boxWidth := width - 6
	if boxWidth > 74 {
		boxWidth = 74
	}
	if boxWidth < 40 {
		boxWidth = 40
	}

	var sb strings.Builder

	sb.WriteString(m.theme.Title.Render("KEYBOARD NAVIGATION & SHORTCUTS") + "\n")
	sb.WriteString(m.theme.DimText.Render(strings.Repeat("─", boxWidth-4)) + "\n\n")

	sb.WriteString(m.theme.AccentBold.Render("Navigation") + "\n")
	sb.WriteString(fmt.Sprintf("  %-16s %s\n", "↑ / ↓ or j / k", "Move cursor / scroll"))
	sb.WriteString(fmt.Sprintf("  %-16s %s\n", "Enter", "Select / open item"))
	sb.WriteString(fmt.Sprintf("  %-16s %s\n", "Esc / q", "Go back / return to /"))
	sb.WriteString(fmt.Sprintf("  %-16s %s\n", "/ ", "Open instant search"))
	sb.WriteString(fmt.Sprintf("  %-16s %s\n", ": ", "Open command mode (ls, cd, cat...)"))
	sb.WriteString(fmt.Sprintf("  %-16s %s\n\n", "?", "Toggle this help view"))

	sb.WriteString(m.theme.AccentBold.Render("Direct Path Shortcuts") + "\n")
	sb.WriteString(fmt.Sprintf("  %-16s %s\n", "p", "/projects — featured projects & systems"))
	sb.WriteString(fmt.Sprintf("  %-16s %s\n", "a", "/about    — background & philosophy"))
	sb.WriteString(fmt.Sprintf("  %-16s %s\n", "s", "/skills   — technical toolchain & stack"))
	sb.WriteString(fmt.Sprintf("  %-16s %s\n", "e", "/experience — strengths & work style"))
	sb.WriteString(fmt.Sprintf("  %-16s %s\n", "b", "/blog     — engineering articles"))
	sb.WriteString(fmt.Sprintf("  %-16s %s\n", "g", "/github   — open source profile & repos"))
	sb.WriteString(fmt.Sprintf("  %-16s %s\n\n", "c", "/contact  — email, links & calendar"))

	sb.WriteString(m.theme.AccentBold.Render("VFS Commands (in : mode)") + "\n")
	sb.WriteString(fmt.Sprintf("  %-16s %s\n", "pwd", "Print working directory"))
	sb.WriteString(fmt.Sprintf("  %-16s %s\n", "cd <dir>", "Change directory (e.g. cd projects)"))
	sb.WriteString(fmt.Sprintf("  %-16s %s\n", "ls [-l] [-a]", "List contents with details"))
	sb.WriteString(fmt.Sprintf("  %-16s %s\n", "cat <file>", "Read virtual file content"))
	sb.WriteString(fmt.Sprintf("  %-16s %s\n", "tree", "Render filesystem hierarchy"))
	sb.WriteString(fmt.Sprintf("  %-16s %s\n\n", "whoami", "Display user details"))

	sb.WriteString(m.theme.DimText.Render("Press [Esc], [q], or [?] to close help"))

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.ActiveBorder).
		Width(boxWidth).
		Padding(1, 2).
		Render(sb.String())

	return lipgloss.PlaceHorizontal(width, lipgloss.Center, box)
}
