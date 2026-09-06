package commands

import (
	"strings"
	"testing"

	"github.com/dev-harshhh19/harshad-tui/internal/filesystem"
)

func TestRunner(t *testing.T) {
	vfs := filesystem.NewVFS()
	runner := NewRunner(vfs)

	// Test whoami
	res := runner.Execute("/", "whoami")
	if !strings.Contains(res.Output, "dev-harshhh") {
		t.Fatalf("expected whoami to contain dev-harshhh, got: %s", res.Output)
	}

	// Test pwd
	res = runner.Execute("/", "pwd")
	if res.Output != "/" {
		t.Fatalf("expected pwd to be '/', got: %s", res.Output)
	}

	// Test cd projects
	res = runner.Execute("/", "cd projects")
	if res.NewPath != "/projects" {
		t.Fatalf("expected new path /projects, got: %s", res.NewPath)
	}

	// Test cd invalid directory
	res = runner.Execute("/", "cd non_existent_dir")
	if res.Error == nil {
		t.Fatal("expected error when cd into non-existent dir")
	}

	// Test cd file (not a directory)
	res = runner.Execute("/", "cd README.md")
	if res.Error == nil {
		t.Fatal("expected error when cd into a file")
	}

	// Test ls in /projects
	res = runner.Execute("/projects", "ls")
	if !strings.Contains(res.Output, "freelivo") {
		t.Fatalf("expected ls to contain freelivo, got: %s", res.Output)
	}

	// Test ls -l and ls -la
	res = runner.Execute("/", "ls -l")
	if !strings.Contains(res.Output, "total") || !strings.Contains(res.Output, "about/") {
		t.Fatalf("expected long format ls output, got: %s", res.Output)
	}

	res = runner.Execute("/projects", "ls -la")
	if !strings.Contains(res.Output, ".") {
		t.Fatalf("expected ls -la to contain '.', got: %s", res.Output)
	}

	// Test ls on a file
	res = runner.Execute("/", "ls /README.md")
	if res.Output != "README.md" {
		t.Fatalf("expected ls on file to return 'README.md', got: %s", res.Output)
	}

	// Test cat projects/freelivo/README.md
	res = runner.Execute("/", "cat /projects/freelivo/README.md")
	if !strings.Contains(res.Output, "Freelivo") {
		t.Fatalf("expected cat to contain Freelivo, got: %s", res.Output)
	}

	// Test cat errors
	res = runner.Execute("/", "cat")
	if res.Error == nil {
		t.Fatal("expected error for cat with missing operand")
	}

	res = runner.Execute("/", "cat /projects")
	if res.Error == nil {
		t.Fatal("expected error for cat on a directory")
	}

	res = runner.Execute("/", "cat /non-existent-file")
	if res.Error == nil {
		t.Fatal("expected error for cat on non-existent file")
	}

	// Test tree
	res = runner.Execute("/", "tree")
	if !strings.Contains(res.Output, "├── about/") {
		t.Fatalf("expected tree output, got: %s", res.Output)
	}

	// Test echo
	res = runner.Execute("/", "echo hello world")
	if res.Output != "hello world" {
		t.Fatalf("expected 'hello world', got: %s", res.Output)
	}

	// Test history
	res = runner.Execute("/", "history")
	if !strings.Contains(res.Output, "whoami") {
		t.Fatalf("expected history to contain previous commands, got: %s", res.Output)
	}

	// Test clear
	res = runner.Execute("/", "clear")
	if !res.ShouldClear {
		t.Fatal("expected ShouldClear to be true for clear command")
	}

	// Test help
	res = runner.Execute("/", "help")
	if !strings.Contains(res.Output, "VFS Shell Command Help") {
		t.Fatalf("expected help text, got: %s", res.Output)
	}

	// Test exit
	res = runner.Execute("/", "exit")
	if !strings.Contains(res.Output, "Returning") {
		t.Fatalf("expected exit message, got: %s", res.Output)
	}

	// Test unknown command
	res = runner.Execute("/", "invalidcommand123")
	if res.Error == nil {
		t.Fatal("expected error for unknown command")
	}

	// Test empty command
	res = runner.Execute("/", "   ")
	if res.Output != "" || res.Error != nil {
		t.Fatal("expected empty result for empty input")
	}
}
