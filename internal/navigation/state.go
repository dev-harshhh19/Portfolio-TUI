package navigation

import (
	"path"
	"strings"
)

// ViewMode represents the active screen in the TUI
type ViewMode int

const (
	ViewHome ViewMode = iota
	ViewDirectory
	ViewFile
	ViewProject
	ViewBlog
	ViewSearch
	ViewCommand
	ViewHelp
)

// State tracks the user's path, history, and active view mode
type State struct {
	CurrentPath string
	History     []string
	Mode        ViewMode
	PrevMode    ViewMode
	CursorPos   map[string]int // remembers cursor position per path
	SearchQuery string
	CommandLine string
}

// NewState initializes navigation state at root
func NewState() *State {
	return &State{
		CurrentPath: "/",
		History:     make([]string, 0),
		Mode:        ViewHome,
		PrevMode:    ViewHome,
		CursorPos:   make(map[string]int),
	}
}

// PushPath navigates to a new path and stores history
func (s *State) PushPath(newPath string) {
	if s.CurrentPath != newPath {
		s.History = append(s.History, s.CurrentPath)
		s.CurrentPath = newPath
	}
}

// PopPath returns to the previous path
func (s *State) PopPath() bool {
	if len(s.History) > 0 {
		prev := s.History[len(s.History)-1]
		s.History = s.History[:len(s.History)-1]
		s.CurrentPath = prev
		return true
	}
	if s.CurrentPath != "/" {
		parent := path.Dir(strings.TrimRight(s.CurrentPath, "/"))
		if parent == "" {
			parent = "/"
		}
		s.CurrentPath = parent
		return true
	}
	return false
}

// GetCursor returns the saved cursor position for the given path
func (s *State) GetCursor(p string) int {
	return s.CursorPos[p]
}

// SetCursor saves the cursor position for the given path
func (s *State) SetCursor(p string, idx int) {
	s.CursorPos[p] = idx
}
