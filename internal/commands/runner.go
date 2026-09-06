package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dev-harshhh19/harshad-tui/internal/filesystem"
)

// CommandResult represents the output of a command execution
type CommandResult struct {
	Output      string
	NewPath     string
	ShouldClear bool
	OpenNode    *filesystem.FSNode
	Error       error
}

// Runner handles virtual shell command parsing and execution
type Runner struct {
	VFS     *filesystem.VFS
	History []string
}

// NewRunner creates a command runner backed by the virtual filesystem
func NewRunner(vfs *filesystem.VFS) *Runner {
	return &Runner{
		VFS:     vfs,
		History: make([]string, 0),
	}
}

// Execute executes a command string safely within the sandboxed VFS
func (r *Runner) Execute(currentDir, input string) CommandResult {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return CommandResult{}
	}

	r.History = append(r.History, trimmed)

	args := parseArgs(trimmed)
	if len(args) == 0 {
		return CommandResult{}
	}

	cmd := strings.ToLower(args[0])
	cmdArgs := args[1:]

	switch cmd {
	case "pwd":
		return CommandResult{Output: currentDir}

	case "whoami":
		return CommandResult{Output: r.VFS.Profile.Handle + " (" + r.VFS.Profile.Name + " — " + r.VFS.Profile.Title + ")"}

	case "clear", "cls":
		return CommandResult{ShouldClear: true}

	case "history":
		var sb strings.Builder
		for i, h := range r.History {
			sb.WriteString(fmt.Sprintf("  %3d  %s\n", i+1, h))
		}
		return CommandResult{Output: strings.TrimRight(sb.String(), "\n")}

	case "cd":
		target := "/"
		if len(cmdArgs) > 0 {
			target = cmdArgs[0]
		}
		newPath := r.VFS.ResolvePath(currentDir, target)
		node, found := r.VFS.GetNode(newPath)
		if !found {
			return CommandResult{Error: fmt.Errorf("cd: no such directory: %s", target)}
		}
		if node.Type != filesystem.DirNode {
			return CommandResult{Error: fmt.Errorf("cd: not a directory: %s", target)}
		}
		return CommandResult{NewPath: node.Path}

	case "ls":
		target := currentDir
		showAll := false
		longFormat := false

		for _, arg := range cmdArgs {
			if strings.HasPrefix(arg, "-") {
				if strings.Contains(arg, "a") {
					showAll = true
				}
				if strings.Contains(arg, "l") {
					longFormat = true
				}
			} else {
				target = r.VFS.ResolvePath(currentDir, arg)
			}
		}

		node, found := r.VFS.GetNode(target)
		if !found {
			return CommandResult{Error: fmt.Errorf("ls: cannot access '%s': No such file or directory", target)}
		}

		if node.Type == filesystem.FileNode {
			if longFormat {
				return CommandResult{Output: formatLong(node)}
			}
			return CommandResult{Output: node.Name}
		}

		var children []*filesystem.FSNode
		for _, c := range node.Children {
			children = append(children, c)
		}
		sort.Slice(children, func(i, j int) bool {
			return children[i].Name < children[j].Name
		})

		var sb strings.Builder
		if longFormat {
			sb.WriteString(fmt.Sprintf("total %d\n", len(children)))
			if showAll {
				sb.WriteString(fmt.Sprintf("%s  1 dev-harshhh staff %6d Aug  8 12:00 .\n", node.Permissions, 4096))
				if node.Parent != nil {
					sb.WriteString(fmt.Sprintf("%s  1 dev-harshhh staff %6d Aug  8 12:00 ..\n", node.Parent.Permissions, 4096))
				}
			}
			for _, c := range children {
				sb.WriteString(formatLong(c) + "\n")
			}
		} else {
			for _, c := range children {
				suffix := ""
				if c.Type == filesystem.DirNode {
					suffix = "/"
				}
				sb.WriteString(fmt.Sprintf("%-18s ", c.Name+suffix))
			}
		}
		return CommandResult{Output: strings.TrimRight(sb.String(), "\n ")}

	case "cat":
		if len(cmdArgs) == 0 {
			return CommandResult{Error: fmt.Errorf("cat: missing operand")}
		}
		target := r.VFS.ResolvePath(currentDir, cmdArgs[0])
		node, found := r.VFS.GetNode(target)
		if !found {
			return CommandResult{Error: fmt.Errorf("cat: %s: No such file or directory", cmdArgs[0])}
		}
		if node.Type == filesystem.DirNode {
			return CommandResult{Error: fmt.Errorf("cat: %s: Is a directory", cmdArgs[0])}
		}
		return CommandResult{Output: node.Content, OpenNode: node}

	case "tree":
		target := currentDir
		if len(cmdArgs) > 0 {
			target = r.VFS.ResolvePath(currentDir, cmdArgs[0])
		}
		treeStr, err := r.VFS.Tree(target)
		if err != nil {
			return CommandResult{Error: err}
		}
		return CommandResult{Output: treeStr}

	case "echo":
		return CommandResult{Output: strings.Join(cmdArgs, " ")}

	case "help":
		helpText := `VFS Shell Command Help:
  pwd                 Print current working directory
  cd <dir>            Change directory (supports cd .., cd /, cd projects)
  ls [-l] [-a] [dir]  List directory contents
  cat <file>          Display file contents
  tree [dir]          Render directory structure as an ASCII tree
  whoami              Print current user identity
  history             List recently executed commands
  clear               Clear command output buffer
  help                Display this reference guide
  exit / q            Exit command mode and return to interactive TUI
`
		return CommandResult{Output: helpText}

	case "exit", "quit":
		return CommandResult{Output: "Returning to interactive view."}

	default:
		return CommandResult{Error: fmt.Errorf("command not found: %s. Type 'help' for available commands", cmd)}
	}
}

func formatLong(node *filesystem.FSNode) string {
	name := node.Name
	if node.Type == filesystem.DirNode {
		name += "/"
	}
	size := node.Size
	if node.Type == filesystem.DirNode {
		size = 4096
	}
	return fmt.Sprintf("%s  1 %-10s %-6s %6d Aug  8 12:00 %s",
		node.Permissions, node.Owner, node.Group, size, name)
}

func parseArgs(input string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false
	quoteChar := rune(0)

	for _, r := range input {
		if inQuotes {
			if r == quoteChar {
				inQuotes = false
			} else {
				current.WriteRune(r)
			}
		} else {
			if r == '"' || r == '\'' {
				inQuotes = true
				quoteChar = r
			} else if r == ' ' || r == '\t' {
				if current.Len() > 0 {
					args = append(args, current.String())
					current.Reset()
				}
			} else {
				current.WriteRune(r)
			}
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}
