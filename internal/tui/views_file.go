package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dev-harshhh19/harshad-tui/internal/filesystem"
)

// RenderFile renders a text file view with clean formatting and scrolling
func (m *Model) RenderFile(node *filesystem.FSNode, width, height int) string {
	if node == nil {
		return "File not found"
	}

	boxWidth := width - 6
	if boxWidth > 82 {
		boxWidth = 82
	}
	if boxWidth < 40 {
		boxWidth = 40
	}

	rawLines := strings.Split(node.Content, "\n")
	var formattedLines []string

	for _, line := range rawLines {
		if strings.HasPrefix(line, "# ") {
			formattedLines = append(formattedLines, m.theme.Title.Render(strings.TrimPrefix(line, "# ")))
		} else if strings.HasPrefix(line, "## ") {
			formattedLines = append(formattedLines, "\n"+m.theme.AccentBold.Render(strings.TrimPrefix(line, "## ")))
		} else if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "===") {
			formattedLines = append(formattedLines, m.theme.DimText.Render(strings.Repeat("─", boxWidth-6)))
		} else if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "• ") {
			formattedLines = append(formattedLines, "  • "+m.wrapText(line[2:], boxWidth-8))
		} else {
			formattedLines = append(formattedLines, m.wrapText(line, boxWidth-6))
		}
	}

	totalLines := len(formattedLines)
	visibleHeight := height - 10
	if visibleHeight < 8 {
		visibleHeight = 8
	}

	// Clamp scroll
	maxOffset := totalLines - visibleHeight
	if maxOffset < 0 {
		maxOffset = 0
	}
	if m.fileScroll > maxOffset {
		m.fileScroll = maxOffset
	}
	if m.fileScroll < 0 {
		m.fileScroll = 0
	}

	endIdx := m.fileScroll + visibleHeight
	if endIdx > totalLines {
		endIdx = totalLines
	}

	slice := formattedLines[m.fileScroll:endIdx]

	var sb strings.Builder
	sb.WriteString(m.theme.DimText.Render(fmt.Sprintf("FILE: %s  (%d bytes)", node.Path, node.Size)) + "\n")
	sb.WriteString(m.theme.DimText.Render(strings.Repeat("─", boxWidth-4)) + "\n\n")

	sb.WriteString(strings.Join(slice, "\n"))

	if totalLines > visibleHeight {
		percent := int(float64(m.fileScroll) / float64(maxOffset) * 100)
		sb.WriteString(fmt.Sprintf("\n\n%s  [↑↓ scroll, Esc/q back]", m.theme.DimText.Render(fmt.Sprintf("[%d%%]", percent))))
	} else {
		sb.WriteString("\n\n" + m.theme.DimText.Render("[Press Esc or q to go back]"))
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.BorderColor).
		Width(boxWidth).
		Padding(1, 2).
		Render(sb.String())

	return lipgloss.PlaceHorizontal(width, lipgloss.Center, box)
}
