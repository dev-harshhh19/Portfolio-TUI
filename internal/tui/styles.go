package tui

import "github.com/charmbracelet/lipgloss"

// Theme holds styling definitions for the portfolio TUI
type Theme struct {
	Primary      lipgloss.AdaptiveColor
	Secondary    lipgloss.AdaptiveColor
	Dim          lipgloss.AdaptiveColor
	Accent       lipgloss.AdaptiveColor
	Highlight    lipgloss.AdaptiveColor
	StatusOpen   lipgloss.AdaptiveColor
	BorderColor  lipgloss.AdaptiveColor
	ActiveBorder lipgloss.AdaptiveColor
	BgDark       lipgloss.AdaptiveColor
	BgSubtle     lipgloss.AdaptiveColor
	ErrorColor   lipgloss.AdaptiveColor

	// Pre-composed lipgloss styles
	AppContainer  lipgloss.Style
	HeaderBox     lipgloss.Style
	Breadcrumb    lipgloss.Style
	PathStyle     lipgloss.Style
	KeyBadge      lipgloss.Style
	KeyText       lipgloss.Style
	Title         lipgloss.Style
	Subtitle      lipgloss.Style
	MetaLabel     lipgloss.Style
	MetaValue     lipgloss.Style
	Body          lipgloss.Style
	DimText       lipgloss.Style
	HighlightRow  lipgloss.Style
	NormalRow     lipgloss.Style
	CursorSymbol  lipgloss.Style
	FooterBar     lipgloss.Style
	SearchBar     lipgloss.Style
	SearchPrompt  lipgloss.Style
	CommandLine   lipgloss.Style
	CommandPrompt lipgloss.Style
	TagBadge      lipgloss.Style
	BoxBorder     lipgloss.Style
	ErrorText     lipgloss.Style
	AccentBold    lipgloss.Style
	SecondaryBold lipgloss.Style
	AccentText    lipgloss.Style
}

// DefaultTheme returns the curated, understated terminal theme
func DefaultTheme() Theme {
	t := Theme{
		Primary:      lipgloss.AdaptiveColor{Light: "#18181B", Dark: "#EDEDED"},
		Secondary:    lipgloss.AdaptiveColor{Light: "#71717A", Dark: "#A1A1AA"},
		Dim:          lipgloss.AdaptiveColor{Light: "#A1A1AA", Dark: "#52525B"},
		Accent:       lipgloss.AdaptiveColor{Light: "#09090B", Dark: "#FAFAFA"},
		Highlight:    lipgloss.AdaptiveColor{Light: "#27272A", Dark: "#FFFFFF"},
		StatusOpen:   lipgloss.AdaptiveColor{Light: "#15803D", Dark: "#34D399"},
		BorderColor:  lipgloss.AdaptiveColor{Light: "#E4E4E7", Dark: "#27272A"},
		ActiveBorder: lipgloss.AdaptiveColor{Light: "#71717A", Dark: "#52525B"},
		BgDark:       lipgloss.AdaptiveColor{Light: "#F4F4F5", Dark: "#0E0E10"},
		BgSubtle:     lipgloss.AdaptiveColor{Light: "#E4E4E7", Dark: "#18181B"},
		ErrorColor:   lipgloss.AdaptiveColor{Light: "#B91C1C", Dark: "#F87171"},
	}

	t.ErrorText = lipgloss.NewStyle().Foreground(t.ErrorColor)
	t.AccentBold = lipgloss.NewStyle().Foreground(t.Accent).Bold(true)
	t.SecondaryBold = lipgloss.NewStyle().Foreground(t.Secondary).Bold(true)
	t.AccentText = lipgloss.NewStyle().Foreground(t.Accent)

	t.AppContainer = lipgloss.NewStyle().
		Padding(1, 2)

	t.HeaderBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderColor).
		Padding(0, 1)

	t.Breadcrumb = lipgloss.NewStyle().
		Foreground(t.Secondary)

	t.PathStyle = lipgloss.NewStyle().
		Foreground(t.Accent).
		Bold(true)

	t.KeyBadge = lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true)

	t.KeyText = lipgloss.NewStyle().
		Foreground(t.Secondary)

	t.Title = lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true)

	t.Subtitle = lipgloss.NewStyle().
		Foreground(t.Secondary)

	t.MetaLabel = lipgloss.NewStyle().
		Foreground(t.Dim).
		Width(12)

	t.MetaValue = lipgloss.NewStyle().
		Foreground(t.Primary)

	t.Body = lipgloss.NewStyle().
		Foreground(t.Primary)

	t.DimText = lipgloss.NewStyle().
		Foreground(t.Dim)

	t.HighlightRow = lipgloss.NewStyle().
		Foreground(t.Accent).
		Background(t.BgSubtle).
		Bold(true)

	t.NormalRow = lipgloss.NewStyle().
		Foreground(t.Primary)

	t.CursorSymbol = lipgloss.NewStyle().
		Foreground(t.Accent).
		Bold(true)

	t.FooterBar = lipgloss.NewStyle().
		Foreground(t.Dim).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(t.BorderColor).
		PaddingTop(0)

	t.SearchBar = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.ActiveBorder).
		Padding(0, 1)

	t.SearchPrompt = lipgloss.NewStyle().
		Foreground(t.Accent).
		Bold(true)

	t.CommandLine = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.ActiveBorder).
		Padding(0, 1)

	t.CommandPrompt = lipgloss.NewStyle().
		Foreground(t.Accent).
		Bold(true)

	t.TagBadge = lipgloss.NewStyle().
		Foreground(t.Secondary).
		Background(t.BgSubtle).
		Padding(0, 1)

	t.BoxBorder = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderColor).
		Padding(1)

	return t
}
