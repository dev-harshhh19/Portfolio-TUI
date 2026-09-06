package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// MenuItem represents a selectable option on the start screen
type MenuItem struct {
	Key         string
	Path        string
	Title       string
	Description string
}

func getHomeMenuItems() []MenuItem {
	return []MenuItem{
		{Key: "p", Path: "/projects", Title: "Projects", Description: "Shipped applications & systems"},
		{Key: "a", Path: "/about", Title: "About", Description: "Bio, story & engineering values"},
		{Key: "s", Path: "/skills", Title: "Skills", Description: "Languages, frameworks & DevOps"},
		{Key: "e", Path: "/experience", Title: "Experience", Description: "Work style, strengths & mindset"},
		{Key: "b", Path: "/blog", Title: "Blog", Description: "Technical guides & articles"},
		{Key: "g", Path: "/github", Title: "GitHub", Description: "Open source repositories"},
		{Key: "c", Path: "/contact", Title: "Contact", Description: "Email, LinkedIn & calendar"},
	}
}

// RenderHome renders the minimalist, centered start screen
func (m *Model) RenderHome(width, height int) string {
	items := getHomeMenuItems()

	var sb strings.Builder

	// Logo / Identity section
	handle := m.theme.DimText.Render(m.vfs.Profile.Handle)
	name := m.theme.Title.Render(m.vfs.Profile.Name)
	title := m.theme.Subtitle.Render(m.vfs.Profile.Title)

	statusDot := lipgloss.NewStyle().Foreground(m.theme.StatusOpen).Render("●")
	statusText := m.theme.DimText.Render(" " + m.vfs.Profile.Status)
	statusBadge := fmt.Sprintf("%s%s", statusDot, statusText)

	// Menu list
	var menuRows []string
	for i, it := range items {
		isSelected := i == m.homeCursor

		cursor := "  "
		if isSelected {
			cursor = m.theme.CursorSymbol.Render("❯ ")
		}

		keyBadge := m.theme.KeyBadge.Render(fmt.Sprintf("[%s]", it.Key))
		itemTitle := it.Title
		if isSelected {
			itemTitle = m.theme.HighlightRow.Render(it.Title)
		} else {
			itemTitle = m.theme.NormalRow.Render(it.Title)
		}

		desc := m.theme.DimText.Render(it.Description)

		row := fmt.Sprintf("%s%-5s  %-14s  %s", cursor, keyBadge, itemTitle, desc)
		menuRows = append(menuRows, row)
	}

	menuBlock := strings.Join(menuRows, "\n")

	tagline := m.theme.DimText.Render(m.vfs.Profile.Philosophy)

	// Assemble central block
	sb.WriteString("\n")
	sb.WriteString(lipgloss.PlaceHorizontal(width, lipgloss.Center, handle))
	sb.WriteString("\n\n")
	sb.WriteString(lipgloss.PlaceHorizontal(width, lipgloss.Center, name))
	sb.WriteString("\n")
	sb.WriteString(lipgloss.PlaceHorizontal(width, lipgloss.Center, title))
	sb.WriteString("\n")
	sb.WriteString(lipgloss.PlaceHorizontal(width, lipgloss.Center, statusBadge))
	sb.WriteString("\n\n\n")

	// Center the menu block
	centeredMenu := lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(menuBlock)
	sb.WriteString(centeredMenu)

	sb.WriteString("\n\n\n")
	sb.WriteString(lipgloss.PlaceHorizontal(width, lipgloss.Center, tagline))
	sb.WriteString("\n")

	return sb.String()
}
