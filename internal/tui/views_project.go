package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dev-harshhh19/harshad-tui/internal/content"
)

// RenderProject renders a compact, elegant project case study
func (m *Model) RenderProject(p *content.Project, width, height int) string {
	if p == nil {
		return "Project not found"
	}

	boxWidth := width - 6
	if boxWidth > 78 {
		boxWidth = 78
	}
	if boxWidth < 42 {
		boxWidth = 42
	}

	var sb strings.Builder

	// Title & Tagline
	titleStr := m.theme.Title.Render(p.Title)
	yearStr := m.theme.DimText.Render("(" + p.Year + ")")
	taglineStr := m.theme.Subtitle.Render(p.Tagline)

	sb.WriteString(fmt.Sprintf("%s  %s\n", titleStr, yearStr))
	sb.WriteString(taglineStr + "\n")
	sb.WriteString(m.theme.DimText.Render(strings.Repeat("─", boxWidth-4)) + "\n\n")

	// Overview
	sb.WriteString(m.theme.MetaLabel.Render("OVERVIEW:") + "\n")
	sb.WriteString(m.wrapText(p.Description, boxWidth-6) + "\n\n")

	// Architecture & Detail
	sb.WriteString(m.theme.MetaLabel.Render("DEEP DIVE & ARCHITECTURE:") + "\n")
	sb.WriteString(m.wrapText(p.Detail, boxWidth-6) + "\n\n")

	// Tech Stack
	sb.WriteString(m.theme.MetaLabel.Render("STACK:") + "\n")
	var badges []string
	for _, tech := range p.Tech {
		badges = append(badges, m.theme.TagBadge.Render(tech))
	}
	sb.WriteString(strings.Join(badges, " ") + "\n\n")

	// Links
	sb.WriteString(m.theme.MetaLabel.Render("ACCESS:") + "\n")
	if p.Live != "" && p.Live != "NA" {
		sb.WriteString(fmt.Sprintf("  • Live:   %s\n", m.theme.AccentText.Render(p.Live)))
	}
	if p.GitHub != "" && p.GitHub != "NA" {
		sb.WriteString(fmt.Sprintf("  • GitHub: %s\n", m.theme.AccentText.Render(p.GitHub)))
	}
	if (p.Live == "" || p.Live == "NA") && (p.GitHub == "" || p.GitHub == "NA") {
		sb.WriteString("  • Internal / Proprietary\n")
	}

	sb.WriteString("\n" + m.theme.DimText.Render("Press [Esc] or [q] to return to /projects"))

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.BorderColor).
		Width(boxWidth).
		Padding(1, 2).
		Render(sb.String())

	return lipgloss.PlaceHorizontal(width, lipgloss.Center, box)
}

func (m *Model) wrapText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return text
	}
	paragraphs := strings.Split(text, "\n")
	var wrappedParagraphs []string

	for _, p := range paragraphs {
		words := strings.Fields(p)
		if len(words) == 0 {
			wrappedParagraphs = append(wrappedParagraphs, "")
			continue
		}

		var line strings.Builder
		lineLen := 0

		for _, w := range words {
			wLen := len(w)
			if lineLen == 0 {
				line.WriteString(w)
				lineLen = wLen
			} else if lineLen+1+wLen <= maxWidth {
				line.WriteString(" " + w)
				lineLen += 1 + wLen
			} else {
				wrappedParagraphs = append(wrappedParagraphs, line.String())
				line.Reset()
				line.WriteString(w)
				lineLen = wLen
			}
		}
		if line.Len() > 0 {
			wrappedParagraphs = append(wrappedParagraphs, line.String())
		}
	}
	return strings.Join(wrappedParagraphs, "\n")
}
