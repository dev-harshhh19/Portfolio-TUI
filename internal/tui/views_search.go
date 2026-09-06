package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// SearchItem is an indexable item
type SearchItem struct {
	Path        string
	Title       string
	Description string
	Category    string
}

func (m *Model) getSearchIndex() []SearchItem {
	var index []SearchItem

	// Sections
	index = append(index,
		SearchItem{Path: "/about", Title: "About Harshad Nikam", Description: "Bio, story & principles", Category: "Section"},
		SearchItem{Path: "/skills", Title: "Technical Skills", Description: "Frontend, backend, cloud & tools", Category: "Section"},
		SearchItem{Path: "/experience", Title: "Experience & Strengths", Description: "Ownership, clean code, learning fast", Category: "Section"},
		SearchItem{Path: "/blog", Title: "Engineering Blog", Description: "Technical articles & deep dives", Category: "Section"},
		SearchItem{Path: "/github", Title: "GitHub Repositories", Description: "Open source profile & repos", Category: "Section"},
		SearchItem{Path: "/contact", Title: "Contact & Links", Description: "Email, LinkedIn, Calendly", Category: "Section"},
		SearchItem{Path: "/resume", Title: "Resume", Description: "CV & professional experience", Category: "Section"},
	)

	// Projects
	for _, p := range m.vfs.Profile.Projects {
		index = append(index, SearchItem{
			Path:        "/projects/" + p.Slug,
			Title:       p.Title,
			Description: fmt.Sprintf("%s · %s", p.Tagline, strings.Join(p.Tech, ", ")),
			Category:    "Project",
		})
	}

	// Blog Posts
	for _, b := range m.vfs.Profile.BlogPosts {
		index = append(index, SearchItem{
			Path:        "/blog/" + b.Slug,
			Title:       b.Title,
			Description: fmt.Sprintf("%s · %s", b.PublishedDate, b.Description),
			Category:    "Article",
		})
	}

	// Skills
	for _, c := range m.vfs.Profile.Skills {
		for _, s := range c.Skills {
			index = append(index, SearchItem{
				Path:        "/skills",
				Title:       s,
				Description: fmt.Sprintf("Technical Skill · %s", c.Title),
				Category:    "Skill",
			})
		}
	}

	// Strengths
	for _, st := range m.vfs.Profile.Strengths {
		index = append(index, SearchItem{
			Path:        "/experience",
			Title:       st.Title,
			Description: st.Body,
			Category:    "Principle",
		})
	}

	return index
}

// FilterSearchResults filters items matching the search query
func (m *Model) FilterSearchResults() []SearchItem {
	query := strings.TrimSpace(strings.ToLower(m.searchQuery))
	index := m.getSearchIndex()
	if query == "" {
		return index
	}

	var results []SearchItem
	for _, it := range index {
		contentMatch := strings.Contains(strings.ToLower(it.Title), query) ||
			strings.Contains(strings.ToLower(it.Path), query) ||
			strings.Contains(strings.ToLower(it.Description), query) ||
			strings.Contains(strings.ToLower(it.Category), query)

		if contentMatch {
			results = append(results, it)
		}
	}
	return results
}

// RenderSearch renders the search interface overlay
func (m *Model) RenderSearch(width, height int) string {
	boxWidth := width - 6
	if boxWidth > 76 {
		boxWidth = 76
	}
	if boxWidth < 40 {
		boxWidth = 40
	}

	var sb strings.Builder

	// Prompt
	prompt := m.theme.SearchPrompt.Render("/ ")
	queryText := m.searchQuery
	cursor := m.theme.CursorSymbol.Render("█")
	sb.WriteString(fmt.Sprintf("%s%s%s\n", prompt, queryText, cursor))
	sb.WriteString(m.theme.DimText.Render(strings.Repeat("─", boxWidth-4)) + "\n")

	results := m.FilterSearchResults()

	if len(results) == 0 {
		sb.WriteString("\n  " + m.theme.DimText.Render("No matching paths found for '"+m.searchQuery+"'"))
	} else {
		maxResults := 8
		if len(results) > maxResults {
			results = results[:maxResults]
		}

		if m.searchCursor >= len(results) {
			m.searchCursor = len(results) - 1
		}
		if m.searchCursor < 0 {
			m.searchCursor = 0
		}

		for i, it := range results {
			isSelected := i == m.searchCursor

			c := "  "
			if isSelected {
				c = m.theme.CursorSymbol.Render("❯ ")
			}

			catBadge := m.theme.TagBadge.Render(it.Category)
			pathText := it.Path
			if isSelected {
				pathText = m.theme.HighlightRow.Render(it.Path)
			} else {
				pathText = m.theme.PathStyle.Render(it.Path)
			}

			row := fmt.Sprintf("%s%-10s  %-26s  %s", c, catBadge, pathText, it.Title)
			sb.WriteString(row + "\n")
		}
	}

	sb.WriteString("\n" + m.theme.DimText.Render("↑↓ navigate  ·  Enter select  ·  Esc cancel"))

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.ActiveBorder).
		Width(boxWidth).
		Padding(1, 2).
		Render(sb.String())

	return lipgloss.PlaceHorizontal(width, lipgloss.Center, box)
}
