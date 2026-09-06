package tui

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dev-harshhh19/harshad-tui/internal/commands"
	"github.com/dev-harshhh19/harshad-tui/internal/filesystem"
	"github.com/dev-harshhh19/harshad-tui/internal/navigation"
)

// Model is the main Bubble Tea model for Harshad's terminal portfolio
type Model struct {
	vfs    *filesystem.VFS
	runner *commands.Runner
	nav    *navigation.State
	theme  Theme
	width  int
	height int
	ready  bool

	// View state
	homeCursor   int
	fsCursor     int
	searchCursor int
	searchQuery  string
	cmdInput     string
	cmdOutput    []string
	fileScroll   int
	blogScroll   int

	// Node currently being viewed
	activeNode *filesystem.FSNode
}

// NewModel creates an initialized Bubble Tea model
func NewModel() *Model {
	vfs := filesystem.NewVFS()
	runner := commands.NewRunner(vfs)
	nav := navigation.NewState()
	theme := DefaultTheme()

	return &Model{
		vfs:       vfs,
		runner:    runner,
		nav:       nav,
		theme:     theme,
		width:     80,
		height:    24,
		cmdOutput: []string{},
	}
}

// SetWindowSize sets initial dimensions before first tick
func (m *Model) SetWindowSize(w, h int) {
	m.width = w
	m.height = h
	m.ready = true
}

// Init implements tea.Model
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Global exit
	if key == "ctrl+c" {
		return m, tea.Quit
	}

	// 1. Search Mode
	if m.nav.Mode == navigation.ViewSearch {
		switch key {
		case "esc":
			m.nav.Mode = m.nav.PrevMode
			m.searchQuery = ""
			return m, nil
		case "up", "ctrl+p":
			if m.searchCursor > 0 {
				m.searchCursor--
			}
			return m, nil
		case "down", "ctrl+n":
			results := m.FilterSearchResults()
			if m.searchCursor < len(results)-1 {
				m.searchCursor++
			}
			return m, nil
		case "enter":
			results := m.FilterSearchResults()
			if len(results) > 0 && m.searchCursor < len(results) {
				selected := results[m.searchCursor]
				m.navigateToPath(selected.Path)
			}
			m.searchQuery = ""
			return m, nil
		case "backspace":
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				m.searchCursor = 0
			}
			return m, nil
		default:
			if len(msg.Runes) > 0 {
				m.searchQuery += string(msg.Runes)
				m.searchCursor = 0
			}
			return m, nil
		}
	}

	// 2. Command Mode
	if m.nav.Mode == navigation.ViewCommand {
		switch key {
		case "esc":
			m.nav.Mode = m.nav.PrevMode
			m.cmdInput = ""
			return m, nil
		case "enter":
			if strings.TrimSpace(m.cmdInput) != "" {
				res := m.runner.Execute(m.nav.CurrentPath, m.cmdInput)
				m.cmdOutput = append(m.cmdOutput, fmt.Sprintf("$ %s", m.cmdInput))
				if res.ShouldClear {
					m.cmdOutput = []string{}
				} else if res.Error != nil {
					m.cmdOutput = append(m.cmdOutput, m.theme.ErrorText.Render(res.Error.Error()))
				} else if res.Output != "" {
					m.cmdOutput = append(m.cmdOutput, res.Output)
				}
				if res.NewPath != "" {
					m.nav.PushPath(res.NewPath)
				}
				if res.OpenNode != nil {
					m.openNode(res.OpenNode)
				}
			}
			m.cmdInput = ""
			return m, nil
		case "backspace":
			if len(m.cmdInput) > 0 {
				m.cmdInput = m.cmdInput[:len(m.cmdInput)-1]
			}
			return m, nil
		default:
			if len(msg.Runes) > 0 {
				m.cmdInput += string(msg.Runes)
			}
			return m, nil
		}
	}

	// 3. Help Mode
	if m.nav.Mode == navigation.ViewHelp {
		if key == "esc" || key == "q" || key == "?" || key == "enter" {
			m.nav.Mode = m.nav.PrevMode
		}
		return m, nil
	}

	// 4. File / Blog Scroll Modes
	if m.nav.Mode == navigation.ViewFile {
		switch key {
		case "up", "k":
			if m.fileScroll > 0 {
				m.fileScroll--
			}
			return m, nil
		case "down", "j":
			m.fileScroll++
			return m, nil
		case "pgup":
			m.fileScroll -= 8
			if m.fileScroll < 0 {
				m.fileScroll = 0
			}
			return m, nil
		case "pgdown", "space":
			m.fileScroll += 8
			return m, nil
		case "esc", "q":
			m.goBack()
			return m, nil
		}
	}

	if m.nav.Mode == navigation.ViewBlog {
		switch key {
		case "up", "k":
			if m.blogScroll > 0 {
				m.blogScroll--
			}
			return m, nil
		case "down", "j":
			m.blogScroll++
			return m, nil
		case "pgup":
			m.blogScroll -= 10
			if m.blogScroll < 0 {
				m.blogScroll = 0
			}
			return m, nil
		case "pgdown", "space":
			m.blogScroll += 10
			return m, nil
		case "esc", "q":
			m.goBack()
			return m, nil
		}
	}

	if m.nav.Mode == navigation.ViewProject {
		switch key {
		case "esc", "q":
			m.goBack()
			return m, nil
		}
	}

	// 5. Global Hotkeys
	switch key {
	case "/":
		m.nav.PrevMode = m.nav.Mode
		m.nav.Mode = navigation.ViewSearch
		m.searchQuery = ""
		m.searchCursor = 0
		return m, nil

	case ":":
		m.nav.PrevMode = m.nav.Mode
		m.nav.Mode = navigation.ViewCommand
		m.cmdInput = ""
		return m, nil

	case "?":
		m.nav.PrevMode = m.nav.Mode
		m.nav.Mode = navigation.ViewHelp
		return m, nil

	case "p":
		m.navigateToPath("/projects")
		return m, nil

	case "a":
		m.navigateToPath("/about")
		return m, nil

	case "s":
		m.navigateToPath("/skills")
		return m, nil

	case "e":
		m.navigateToPath("/experience")
		return m, nil

	case "b":
		m.navigateToPath("/blog")
		return m, nil

	case "g":
		m.navigateToPath("/github")
		return m, nil

	case "c":
		m.navigateToPath("/contact")
		return m, nil

	case "esc":
		if m.nav.Mode == navigation.ViewHome {
			return m, tea.Quit
		}
		m.goBack()
		return m, nil

	case "q":
		if m.nav.Mode == navigation.ViewHome {
			return m, tea.Quit
		}
		m.goBack()
		return m, nil
	}

	// 6. Navigation in Home view
	if m.nav.Mode == navigation.ViewHome {
		items := getHomeMenuItems()
		switch key {
		case "up", "k":
			if m.homeCursor > 0 {
				m.homeCursor--
			} else {
				m.homeCursor = len(items) - 1
			}
		case "down", "j":
			if m.homeCursor < len(items)-1 {
				m.homeCursor++
			} else {
				m.homeCursor = 0
			}
		case "enter":
			if m.homeCursor >= 0 && m.homeCursor < len(items) {
				m.navigateToPath(items[m.homeCursor].Path)
			}
		}
		return m, nil
	}

	// 7. Navigation in Directory view
	if m.nav.Mode == navigation.ViewDirectory {
		node, found := m.vfs.GetNode(m.nav.CurrentPath)
		if !found {
			return m, nil
		}

		var children []*filesystem.FSNode
		for _, c := range node.Children {
			children = append(children, c)
		}
		sort.Slice(children, func(i, j int) bool {
			if children[i].Type != children[j].Type {
				return children[i].Type == filesystem.DirNode
			}
			return children[i].Name < children[j].Name
		})

		switch key {
		case "up", "k":
			if m.fsCursor > 0 {
				m.fsCursor--
			} else if len(children) > 0 {
				m.fsCursor = len(children) - 1
			}
		case "down", "j":
			if m.fsCursor < len(children)-1 {
				m.fsCursor++
			} else {
				m.fsCursor = 0
			}
		case "enter":
			if len(children) > 0 && m.fsCursor < len(children) {
				m.openNode(children[m.fsCursor])
			}
		}
		return m, nil
	}

	return m, nil
}

func (m *Model) navigateToPath(targetPath string) {
	node, found := m.vfs.GetNode(targetPath)
	if !found {
		return
	}
	m.openNode(node)
}

func (m *Model) openNode(node *filesystem.FSNode) {
	m.activeNode = node

	if node.Type == filesystem.DirNode {
		if node.Path == "/" {
			m.nav.CurrentPath = "/"
			m.nav.Mode = navigation.ViewHome
			return
		}
		// If it's a project directory with a ProjectRef, open project view
		if node.ProjectRef != nil {
			m.nav.PushPath(node.Path)
			m.nav.Mode = navigation.ViewProject
			return
		}
		// If it's a blog directory with BlogRef, open blog view
		if node.BlogRef != nil {
			m.nav.PushPath(node.Path)
			m.blogScroll = 0
			m.nav.Mode = navigation.ViewBlog
			return
		}

		m.nav.PushPath(node.Path)
		m.nav.Mode = navigation.ViewDirectory
		m.fsCursor = 0
		return
	}

	// File node
	m.nav.PushPath(node.Path)
	m.fileScroll = 0
	m.nav.Mode = navigation.ViewFile
}

func (m *Model) goBack() {
	if m.nav.PopPath() {
		if m.nav.CurrentPath == "/" {
			m.nav.Mode = navigation.ViewHome
		} else {
			node, found := m.vfs.GetNode(m.nav.CurrentPath)
			if found && node.Type == filesystem.DirNode {
				m.nav.Mode = navigation.ViewDirectory
				m.activeNode = node
			} else {
				m.nav.Mode = navigation.ViewHome
			}
		}
	} else {
		m.nav.Mode = navigation.ViewHome
		m.nav.CurrentPath = "/"
	}
}

// View implements tea.Model
func (m *Model) View() string {
	w := m.width
	h := m.height
	if w <= 0 {
		w = 80
	}
	if h <= 0 {
		h = 24
	}

	var sb strings.Builder

	// 1. Top Bar
	topBar := m.renderTopBar(w)
	sb.WriteString(topBar + "\n")

	// 2. Main Content Area
	var content string
	switch m.nav.Mode {
	case navigation.ViewHome:
		content = m.RenderHome(w, h)
	case navigation.ViewDirectory:
		content = m.RenderDirectory(w, h)
	case navigation.ViewProject:
		if m.activeNode != nil && m.activeNode.ProjectRef != nil {
			content = m.RenderProject(m.activeNode.ProjectRef, w, h)
		} else {
			content = m.RenderDirectory(w, h)
		}
	case navigation.ViewBlog:
		if m.activeNode != nil && m.activeNode.BlogRef != nil {
			content = m.RenderBlog(m.activeNode.BlogRef, w, h)
		} else {
			content = m.RenderDirectory(w, h)
		}
	case navigation.ViewFile:
		content = m.RenderFile(m.activeNode, w, h)
	case navigation.ViewSearch:
		content = m.RenderSearch(w, h)
	case navigation.ViewCommand:
		content = m.RenderCommand(w, h)
	case navigation.ViewHelp:
		content = m.RenderHelp(w, h)
	default:
		content = m.RenderHome(w, h)
	}

	sb.WriteString(content)

	// 3. Bottom Help / Status Bar
	bottomBar := m.renderBottomBar(w)
	sb.WriteString("\n" + bottomBar)

	return sb.String()
}

func (m *Model) renderTopBar(width int) string {
	boxWidth := width - 4
	if boxWidth < 30 {
		boxWidth = 30
	}

	title := m.theme.Title.Render("dev-harshhh@cli")
	loc := m.theme.PathStyle.Render(m.nav.CurrentPath)

	shortcuts := m.theme.DimText.Render("[p] projects  [a] about  [s] skills  [c] contact")

	// Left part: title + path
	left := fmt.Sprintf("%s  %s", title, loc)

	gap := boxWidth - lipgloss.Width(left) - lipgloss.Width(shortcuts) - 4
	if gap < 2 {
		gap = 2
	}

	barContent := fmt.Sprintf("%s%s%s", left, strings.Repeat(" ", gap), shortcuts)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.BorderColor).
		Width(boxWidth).
		Padding(0, 1).
		Render(barContent)
}

func (m *Model) renderBottomBar(width int) string {
	boxWidth := width - 4
	if boxWidth < 30 {
		boxWidth = 30
	}

	hints := "↑↓ navigate   enter select   / search   : command   ? help   q back"
	if m.nav.Mode == navigation.ViewCommand {
		hints = "enter execute   esc exit command mode   help list commands"
	} else if m.nav.Mode == navigation.ViewSearch {
		hints = "↑↓ select result   enter open   esc cancel search"
	}

	renderedHints := m.theme.DimText.Render(hints)
	return lipgloss.PlaceHorizontal(width, lipgloss.Center, renderedHints)
}
