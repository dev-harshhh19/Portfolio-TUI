package navigation

import (
	"testing"
)

func TestNavigationState(t *testing.T) {
	s := NewState()
	if s.CurrentPath != "/" {
		t.Errorf("expected initial path '/', got %s", s.CurrentPath)
	}
	if s.Mode != ViewHome {
		t.Errorf("expected initial mode ViewHome, got %v", s.Mode)
	}

	// Push paths
	s.PushPath("/projects")
	if s.CurrentPath != "/projects" {
		t.Errorf("expected current path '/projects', got %s", s.CurrentPath)
	}
	if len(s.History) != 1 || s.History[0] != "/" {
		t.Errorf("expected history ['/'], got %v", s.History)
	}

	s.PushPath("/projects/freelivo")
	if s.CurrentPath != "/projects/freelivo" {
		t.Errorf("expected current path '/projects/freelivo', got %s", s.CurrentPath)
	}
	if len(s.History) != 2 {
		t.Errorf("expected history length 2, got %d", len(s.History))
	}

	// Pushing the same path should not duplicate history
	s.PushPath("/projects/freelivo")
	if len(s.History) != 2 {
		t.Errorf("expected history length still 2 after pushing same path, got %d", len(s.History))
	}

	// Pop path
	popped := s.PopPath()
	if !popped || s.CurrentPath != "/projects" {
		t.Errorf("expected PopPath() to return true and path '/projects', got %v, %s", popped, s.CurrentPath)
	}

	popped = s.PopPath()
	if !popped || s.CurrentPath != "/" {
		t.Errorf("expected PopPath() to return true and path '/', got %v, %s", popped, s.CurrentPath)
	}

	// Pop when history empty but path is at root
	popped = s.PopPath()
	if popped {
		t.Error("expected PopPath() to return false when at root with empty history")
	}

	// Pop when history empty but path is subpath (fallback to parent)
	s.CurrentPath = "/about/philosophy"
	popped = s.PopPath()
	if !popped || s.CurrentPath != "/about" {
		t.Errorf("expected fallback to parent '/about', got %v, %s", popped, s.CurrentPath)
	}
}

func TestCursorPosition(t *testing.T) {
	s := NewState()

	if pos := s.GetCursor("/projects"); pos != 0 {
		t.Errorf("expected default cursor pos 0, got %d", pos)
	}

	s.SetCursor("/projects", 3)
	if pos := s.GetCursor("/projects"); pos != 3 {
		t.Errorf("expected cursor pos 3, got %d", pos)
	}
}
