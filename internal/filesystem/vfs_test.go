package filesystem

import (
	"strings"
	"testing"
)

func TestNewVFS(t *testing.T) {
	vfs := NewVFS()
	if vfs == nil || vfs.Root == nil {
		t.Fatal("expected non-nil VFS and Root")
	}

	expectedDirs := []string{
		"about",
		"projects",
		"skills",
		"experience",
		"blog",
		"github",
		"contact",
		"resume",
	}

	for _, dirName := range expectedDirs {
		child, exists := vfs.Root.Children[dirName]
		if !exists {
			t.Errorf("expected directory %s to exist in root", dirName)
			continue
		}
		if child.Type != DirNode {
			t.Errorf("expected %s to be DirNode, got %v", dirName, child.Type)
		}
	}

	// Verify root README.md exists
	readme, exists := vfs.Root.Children["README.md"]
	if !exists || readme.Type != FileNode {
		t.Error("expected /README.md file to exist")
	}
}

func TestResolvePath(t *testing.T) {
	vfs := NewVFS()

	cases := []struct {
		current  string
		target   string
		expected string
	}{
		{"/", "", "/"},
		{"/", ".", "/"},
		{"/projects", ".", "/projects"},
		{"/projects", "..", "/"},
		{"/projects/freelivo", "..", "/projects"},
		{"/", "~", "/"},
		{"/", "~/skills", "/skills"},
		{"/", "about", "/about"},
		{"/about", "/projects", "/projects"},
		{"/projects", "freelivo", "/projects/freelivo"},
	}

	for _, tc := range cases {
		got := vfs.ResolvePath(tc.current, tc.target)
		if got != tc.expected {
			t.Errorf("ResolvePath(%q, %q) = %q, want %q", tc.current, tc.target, got, tc.expected)
		}
	}
}

func TestGetNodeAndListDirectory(t *testing.T) {
	vfs := NewVFS()

	// Valid directory
	node, found := vfs.GetNode("/about")
	if !found || node == nil {
		t.Fatal("expected /about to be found")
	}
	if node.Type != DirNode {
		t.Errorf("expected DirNode, got %v", node.Type)
	}

	// Non-existent path
	_, found = vfs.GetNode("/non-existent-path")
	if found {
		t.Error("expected non-existent path to return false")
	}

	// List directory
	items, err := vfs.ListDirectory("/about")
	if err != nil {
		t.Fatalf("unexpected error listing /about: %v", err)
	}
	if len(items) == 0 {
		t.Error("expected items in /about directory")
	}

	// Ensure sorted: directories first, then alphabetical
	for i := 1; i < len(items); i++ {
		prev := items[i-1]
		curr := items[i]
		if prev.Type == curr.Type && prev.Name > curr.Name {
			t.Errorf("expected sorted list, but %s > %s", prev.Name, curr.Name)
		}
	}

	// Error when listing a file as directory
	_, err = vfs.ListDirectory("/README.md")
	if err == nil {
		t.Error("expected error when calling ListDirectory on a file")
	}

	// Error when listing non-existent directory
	_, err = vfs.ListDirectory("/does-not-exist")
	if err == nil {
		t.Error("expected error when calling ListDirectory on non-existent path")
	}
}

func TestTree(t *testing.T) {
	vfs := NewVFS()

	treeStr, err := vfs.Tree("/")
	if err != nil {
		t.Fatalf("unexpected error generating tree: %v", err)
	}

	if !strings.HasPrefix(treeStr, "/\n") {
		t.Errorf("expected tree to start with '/\\n', got %q", treeStr[:min(len(treeStr), 10)])
	}
	if !strings.Contains(treeStr, "├── about/") {
		t.Errorf("expected tree output to contain '├── about/' due to deterministic sorting")
	}

	// Test error on invalid tree path
	_, err = vfs.Tree("/invalid-path")
	if err == nil {
		t.Error("expected error for invalid path in Tree")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
