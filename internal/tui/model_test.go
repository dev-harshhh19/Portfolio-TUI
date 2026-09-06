package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dev-harshhh19/harshad-tui/internal/navigation"
)

func TestNewModel(t *testing.T) {
	m := NewModel()
	if m == nil {
		t.Fatal("expected non-nil Model")
	}
	if m.vfs == nil {
		t.Fatal("expected non-nil VFS in Model")
	}
	if m.nav.Mode != navigation.ViewHome {
		t.Errorf("expected initial mode ViewHome, got %v", m.nav.Mode)
	}

	m.SetWindowSize(100, 30)
	if m.width != 100 || m.height != 30 || !m.ready {
		t.Errorf("SetWindowSize failed: width=%d, height=%d, ready=%v", m.width, m.height, m.ready)
	}
}

func TestFilterSearchResults(t *testing.T) {
	m := NewModel()

	// Empty query returns all items
	all := m.FilterSearchResults()
	if len(all) == 0 {
		t.Fatal("expected non-empty search index")
	}

	// Filter by project keyword
	m.searchQuery = "freelivo"
	res := m.FilterSearchResults()
	if len(res) == 0 {
		t.Error("expected search results for 'freelivo'")
	}
	found := false
	for _, r := range res {
		if strings.Contains(strings.ToLower(r.Title), "freelivo") || strings.Contains(strings.ToLower(r.Path), "freelivo") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected freelivo to appear in search results")
	}

	// Filter by non-existent query
	m.searchQuery = "xyznonexistentquery123"
	res = m.FilterSearchResults()
	if len(res) != 0 {
		t.Errorf("expected 0 results for non-existent query, got %d", len(res))
	}
}

func TestWrapText(t *testing.T) {
	m := NewModel()

	text := "This is a short sentence for testing text wrapping in terminal."
	wrapped := m.wrapText(text, 20)
	lines := strings.Split(wrapped, "\n")
	for _, line := range lines {
		if len(line) > 20 {
			t.Errorf("wrapped line exceeds maxWidth 20: %q (len=%d)", line, len(line))
		}
	}

	// Edge case: maxWidth <= 0
	if got := m.wrapText(text, 0); got != text {
		t.Errorf("expected original text when maxWidth <= 0, got %q", got)
	}

	// Edge case: empty string
	if got := m.wrapText("", 20); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestModelViewRendering(t *testing.T) {
	m := NewModel()
	m.SetWindowSize(100, 30)

	// Render Home
	view := m.View()
	if !strings.Contains(view, "Harshad") {
		t.Error("expected View() on Home to contain 'Harshad'")
	}

	// Switch to Help
	m.nav.Mode = navigation.ViewHelp
	view = m.View()
	if !strings.Contains(view, "KEYBOARD NAVIGATION") {
		t.Error("expected View() on Help to contain 'KEYBOARD NAVIGATION'")
	}

	// Switch to Search
	m.nav.Mode = navigation.ViewSearch
	view = m.View()
	if !strings.Contains(view, "select result") {
		t.Error("expected View() on Search to contain hints")
	}

	// Switch to Command
	m.nav.Mode = navigation.ViewCommand
	view = m.View()
	if !strings.Contains(view, "VFS COMMAND INTERFACE") {
		t.Error("expected View() on Command to contain 'VFS COMMAND INTERFACE'")
	}

	// Switch to Directory
	m.nav.Mode = navigation.ViewDirectory
	m.nav.CurrentPath = "/projects"
	view = m.View()
	if !strings.Contains(view, "TYPE") || !strings.Contains(view, "NAME") {
		t.Error("expected Directory view to contain table headers")
	}
}

func TestKeyHandling(t *testing.T) {
	m := NewModel()
	m.SetWindowSize(100, 30)

	// Key 'p' navigates to /projects
	m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	if m.nav.CurrentPath != "/projects" {
		t.Errorf("expected current path '/projects', got %s", m.nav.CurrentPath)
	}
	if m.nav.Mode != navigation.ViewDirectory {
		t.Errorf("expected ViewDirectory, got %v", m.nav.Mode)
	}

	// Esc returns to Home
	m.handleKey(tea.KeyMsg{Type: tea.KeyEsc})
	if m.nav.CurrentPath != "/" {
		t.Errorf("expected current path '/', got %s", m.nav.CurrentPath)
	}
	if m.nav.Mode != navigation.ViewHome {
		t.Errorf("expected ViewHome, got %v", m.nav.Mode)
	}

	// ':' enters Command mode
	m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{':'}})
	if m.nav.Mode != navigation.ViewCommand {
		t.Errorf("expected ViewCommand, got %v", m.nav.Mode)
	}

	// Typing 'pwd' then enter in Command mode
	m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'w'}})
	m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m.handleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if len(m.cmdOutput) == 0 {
		t.Fatal("expected command output in buffer")
	}

	// Esc exits command mode
	m.handleKey(tea.KeyMsg{Type: tea.KeyEsc})
	if m.nav.Mode == navigation.ViewCommand {
		t.Error("expected to exit ViewCommand on Esc")
	}
}
