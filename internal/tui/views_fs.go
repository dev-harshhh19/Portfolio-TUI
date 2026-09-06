package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dev-harshhh19/harshad-tui/internal/filesystem"
)

// RenderDirectory renders an interactive directory listing
func (m *Model) RenderDirectory(width, height int) string {
	node, found := m.vfs.GetNode(m.nav.CurrentPath)
	if !found {
		return lipgloss.PlaceHorizontal(width, lipgloss.Center, m.theme.ErrorText.Render("Path not found: "+m.nav.CurrentPath))
	}

	var children []*filesystem.FSNode
	for _, c := range node.Children {
		children = append(children, c)
	}

	// Sort directories first, then alphabetically
	sort.Slice(children, func(i, j int) bool {
		if children[i].Type != children[j].Type {
			return children[i].Type == filesystem.DirNode
		}
		return children[i].Name < children[j].Name
	})

	if len(children) == 0 {
		return lipgloss.PlaceHorizontal(width, lipgloss.Center, m.theme.DimText.Render("Directory is empty"))
	}

	// Ensure cursor is in range
	if m.fsCursor >= len(children) {
		m.fsCursor = len(children) - 1
	}
	if m.fsCursor < 0 {
		m.fsCursor = 0
	}

	var rows []string

	// Header line for columns
	colHeader := fmt.Sprintf("   %-6s  %-24s  %s", "TYPE", "NAME", "DESCRIPTION")
	rows = append(rows, m.theme.DimText.Render(colHeader))
	rows = append(rows, m.theme.DimText.Render("   ──────  ────────────────────────  ──────────────────────────────────"))

	for i, c := range children {
		isSelected := i == m.fsCursor

		cursor := "  "
		if isSelected {
			cursor = m.theme.CursorSymbol.Render("❯ ")
		}

		typeTag := "[dir]"
		if c.Type == filesystem.FileNode {
			typeTag = "[doc]"
		}

		nameStr := c.Name
		if c.Type == filesystem.DirNode {
			nameStr += "/"
		}

		desc := c.Description
		if desc == "" && c.ProjectRef != nil {
			desc = c.ProjectRef.Tagline
		}
		if desc == "" && c.BlogRef != nil {
			desc = c.BlogRef.ShortTitle
		}

		if len(desc) > 36 {
			desc = desc[:33] + "..."
		}

		rowContent := fmt.Sprintf("%-6s  %-24s  %s", typeTag, nameStr, desc)

		if isSelected {
			rows = append(rows, cursor+m.theme.HighlightRow.Render(rowContent))
		} else {
			rows = append(rows, cursor+m.theme.NormalRow.Render(rowContent))
		}
	}

	content := strings.Join(rows, "\n")

	// Frame inside a neat box
	boxWidth := width - 4
	if boxWidth > 82 {
		boxWidth = 82
	}
	if boxWidth < 40 {
		boxWidth = 40
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.BorderColor).
		Width(boxWidth).
		Padding(1, 2).
		Render(content)

	return lipgloss.PlaceHorizontal(width, lipgloss.Center, box)
}
