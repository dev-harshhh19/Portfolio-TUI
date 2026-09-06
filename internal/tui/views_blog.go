package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dev-harshhh19/harshad-tui/internal/content"
)

// RenderBlog renders the terminal article reader
func (m *Model) RenderBlog(post *content.BlogPost, width, height int) string {
	if post == nil {
		return "Article not found"
	}

	boxWidth := width - 6
	if boxWidth > 80 {
		boxWidth = 80
	}
	if boxWidth < 40 {
		boxWidth = 40
	}

	contentLines := strings.Split(post.Content, "\n")
	var formattedLines []string

	for _, line := range contentLines {
		if strings.HasPrefix(line, "# ") {
			formattedLines = append(formattedLines, m.theme.Title.Render(strings.TrimPrefix(line, "# ")))
		} else if strings.HasPrefix(line, "## ") {
			formattedLines = append(formattedLines, "\n"+m.theme.AccentBold.Render(strings.TrimPrefix(line, "## ")))
		} else if strings.HasPrefix(line, "### ") {
			formattedLines = append(formattedLines, "\n"+m.theme.SecondaryBold.Render(strings.TrimPrefix(line, "### ")))
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

	// Clamp scroll offset
	maxOffset := totalLines - visibleHeight
	if maxOffset < 0 {
		maxOffset = 0
	}
	if m.blogScroll > maxOffset {
		m.blogScroll = maxOffset
	}
	if m.blogScroll < 0 {
		m.blogScroll = 0
	}

	endIdx := m.blogScroll + visibleHeight
	if endIdx > totalLines {
		endIdx = totalLines
	}

	slice := formattedLines[m.blogScroll:endIdx]

	var sb strings.Builder

	// Top meta bar
	metaStr := fmt.Sprintf("%s  ·  %s  ·  %s", post.PublishedDate, post.ReadTime, strings.Join(post.Tags, ", "))
	sb.WriteString(m.theme.DimText.Render(metaStr) + "\n")
	sb.WriteString(m.theme.DimText.Render(strings.Repeat("─", boxWidth-4)) + "\n\n")

	sb.WriteString(strings.Join(slice, "\n"))

	// Scroll progress footer
	percent := 100
	if totalLines > visibleHeight {
		percent = int(float64(m.blogScroll) / float64(maxOffset) * 100)
	}

	scrollInfo := fmt.Sprintf("\n\n%s  [↑↓/j/k scroll, Esc/q back]", m.theme.DimText.Render(fmt.Sprintf("[%d%%]", percent)))
	sb.WriteString(scrollInfo)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.BorderColor).
		Width(boxWidth).
		Padding(1, 2).
		Render(sb.String())

	return lipgloss.PlaceHorizontal(width, lipgloss.Center, box)
}
