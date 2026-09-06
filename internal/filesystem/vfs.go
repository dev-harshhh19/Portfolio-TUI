package filesystem

import (
	"fmt"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/dev-harshhh19/harshad-tui/internal/content"
)

// NodeType designates whether a node is a directory or file
type NodeType int

const (
	DirNode NodeType = iota
	FileNode
)

// FSNode is a node in the virtual filesystem
type FSNode struct {
	Name        string
	Path        string
	Type        NodeType
	Size        int64
	ModTime     time.Time
	Permissions string
	Owner       string
	Group       string
	Description string
	Content     string
	Children    map[string]*FSNode
	Parent      *FSNode
	ProjectRef  *content.Project
	BlogRef     *content.BlogPost
}

// VFS is the virtual filesystem
type VFS struct {
	Root    *FSNode
	Profile content.Profile
}

// NewVFS constructs the full Unix-like virtual filesystem for Harshad Nikam
func NewVFS() *VFS {
	p := content.GetProfile()

	baseTime := time.Date(2026, 8, 8, 12, 0, 0, 0, time.UTC)

	root := &FSNode{
		Name:        "/",
		Path:        "/",
		Type:        DirNode,
		Permissions: "drwxr-xr-x",
		Owner:       "dev-harshhh",
		Group:       "staff",
		ModTime:     baseTime,
		Children:    make(map[string]*FSNode),
	}

	vfs := &VFS{
		Root:    root,
		Profile: p,
	}

	// Helper to add directories
	addDir := func(parent *FSNode, name, desc string) *FSNode {
		pStr := "/" + name
		if parent.Path != "/" {
			pStr = parent.Path + "/" + name
		}
		dir := &FSNode{
			Name:        name,
			Path:        pStr,
			Type:        DirNode,
			Permissions: "drwxr-xr-x",
			Owner:       "dev-harshhh",
			Group:       "staff",
			Description: desc,
			ModTime:     baseTime,
			Children:    make(map[string]*FSNode),
			Parent:      parent,
		}
		parent.Children[name] = dir
		return dir
	}

	// Helper to add files
	addFile := func(parent *FSNode, name, contentStr, desc string) *FSNode {
		pStr := "/" + name
		if parent.Path != "/" {
			pStr = parent.Path + "/" + name
		}
		f := &FSNode{
			Name:        name,
			Path:        pStr,
			Type:        FileNode,
			Permissions: "-rw-r--r--",
			Owner:       "dev-harshhh",
			Group:       "staff",
			Description: desc,
			ModTime:     baseTime,
			Size:        int64(len(contentStr)),
			Content:     contentStr,
			Parent:      parent,
		}
		parent.Children[name] = f
		return f
	}

	// /README.md
	rootReadme := fmt.Sprintf(`# Harshad Nikam — Portfolio Terminal
==================================================
Handle:      dev-harshhh
Title:       Full-Stack Engineer
Core Motto:  %s
SSH Gateway: ssh %s

Welcome to my interactive terminal portfolio.
Everything starts at /.

Quick Navigation:
  [p] /projects     Things I've built & shipped
  [a] /about        Bio, background & philosophy
  [s] /skills       Frontend, backend, cloud & tools
  [e] /experience   Work style & engineering strengths
  [b] /blog         Technical guides & articles
  [g] /github       Open source repositories
  [c] /contact      Direct reach & links
  [:] Command Mode  Execute Unix commands (ls, cd, cat...)
  [/] Search        Instant global search
`, p.Philosophy, p.SSHHost)
	addFile(root, "README.md", rootReadme, "Terminal welcome & overview")

	// 1. /about
	aboutDir := addDir(root, "about", "Personal bio, background and engineering principles")
	aboutReadme := fmt.Sprintf(`# About Harshad Nikam
==================================================
Role:   Full-Stack Engineer
Status: %s
HQ:     Remote / India

%s

Engineering Philosophy
----------------------
- I prioritize shipping real software that users rely on.
- Clean code matters: flat logic, readable naming, zero cleverness for the sake of cleverness.
- Real-time reactivity, bulletproof API boundaries, and low-latency client performance.
- Modern stack: TypeScript, Next.js, React, Node.js, PostgreSQL, Docker, Linux.

Everything starts at /
`, p.Status, p.Bio)
	addFile(aboutDir, "README.md", aboutReadme, "Harshad's full story & bio")
	addFile(aboutDir, "philosophy.txt", fmt.Sprintf(`CORE PHILOSOPHY
--------------------------------------------------
"%s"

1. Ownership: When I pick something up, I see it through to production.
2. Flat & Clean: Don't introduce abstractions until you need them twice.
3. Fast & Accessible: High-contrast typography, zero lag, minimal dependencies.
`, p.Philosophy), "Core engineering values")

	// 2. /projects
	projectsDir := addDir(root, "projects", "Shipped production applications and technical tools")
	var projIndex strings.Builder
	projIndex.WriteString("# Projects Directory\n==================================================\n\n")
	projIndex.WriteString("Select any project folder to view architecture and case studies:\n\n")

	for _, prj := range p.Projects {
		pDir := addDir(projectsDir, prj.Slug, prj.Tagline)
		pDir.ProjectRef = &prj

		projIndex.WriteString(fmt.Sprintf("- %s/ (%s) — %s\n", prj.Slug, prj.Year, prj.Tagline))

		readme := fmt.Sprintf(`# %s
==================================================
Tagline: %s
Year:    %s
Stack:   %s

Overview
--------
%s

Deep Dive & Architecture
------------------------
%s

Links
-----
Live Demo:   %s
Source Code: %s
`, prj.Title, prj.Tagline, prj.Year, strings.Join(prj.Tech, ", "), prj.Description, prj.Detail, prj.Live, prj.GitHub)

		addFile(pDir, "README.md", readme, fmt.Sprintf("Case study for %s", prj.Title))
		addFile(pDir, "stack.txt", fmt.Sprintf("TECHNOLOGY STACK:\n%s\n", strings.Join(prj.Tech, "\n")), "Tech stack list")
		addFile(pDir, "links.txt", fmt.Sprintf("Live URL:   %s\nGitHub URL: %s\n", prj.Live, prj.GitHub), "External links")
	}
	addFile(projectsDir, "README.md", projIndex.String(), "List of all featured projects")

	// 3. /skills
	skillsDir := addDir(root, "skills", "Technical toolchain, frameworks and languages")
	var skillsTxt strings.Builder
	skillsTxt.WriteString("# Technical Skills & Toolchain\n==================================================\n\n")
	for _, cat := range p.Skills {
		skillsTxt.WriteString(fmt.Sprintf("## %s\n", cat.Title))
		for _, s := range cat.Skills {
			skillsTxt.WriteString(fmt.Sprintf("  • %s\n", s))
		}
		skillsTxt.WriteString("\n")

		// Create individual category files
		slug := strings.ToLower(strings.ReplaceAll(cat.Title, " & ", "-"))
		slug = strings.ReplaceAll(slug, " ", "-")
		addFile(skillsDir, slug+".txt", strings.Join(cat.Skills, "\n"), fmt.Sprintf("%s skills", cat.Title))
	}
	addFile(skillsDir, "stack.txt", skillsTxt.String(), "Complete skills matrix")

	// 4. /experience
	expDir := addDir(root, "experience", "Strengths, engineering habits and teamwork principles")
	var strTxt strings.Builder
	strTxt.WriteString("# Core Strengths\n==================================================\n\n")
	for _, s := range p.Strengths {
		strTxt.WriteString(fmt.Sprintf("• %s\n  %s\n\n", s.Title, s.Body))
	}
	addFile(expDir, "strengths.txt", strTxt.String(), "Core strengths")

	var wsTxt strings.Builder
	wsTxt.WriteString("# Work Style & Principles\n==================================================\n\n")
	for _, w := range p.WorkStyle {
		wsTxt.WriteString(fmt.Sprintf("• %s\n  %s\n\n", w.Title, w.Body))
	}
	addFile(expDir, "work-style.txt", wsTxt.String(), "Collaboration principles")

	// 5. /blog
	blogDir := addDir(root, "blog", "Engineering blog posts and guides")
	for _, post := range p.BlogPosts {
		bDir := addDir(blogDir, post.Slug, post.ShortTitle)
		bDir.BlogRef = &post
		addFile(bDir, "article.md", post.Content, post.Title)
		addFile(bDir, "meta.txt", fmt.Sprintf("Title:     %s\nPublished: %s\nRead Time: %s\nTags:      %s\n\n%s\n",
			post.Title, post.PublishedDate, post.ReadTime, strings.Join(post.Tags, ", "), post.Description), "Post metadata")
	}
	var blogIndex strings.Builder
	blogIndex.WriteString("# Harshad's Technical Articles\n==================================================\n\n")
	for _, post := range p.BlogPosts {
		blogIndex.WriteString(fmt.Sprintf("- /blog/%s\n  %s (%s · %s)\n  Tags: %s\n\n",
			post.Slug, post.Title, post.PublishedDate, post.ReadTime, strings.Join(post.Tags, ", ")))
	}
	addFile(blogDir, "README.md", blogIndex.String(), "Articles index")

	// 6. /github
	ghDir := addDir(root, "github", "GitHub profile and repository links")
	ghContent := fmt.Sprintf(`# GitHub Profile
==================================================
User:    %s
URL:     %s
Profile: Harshad Nikam (dev-harshhh19)

Featured Public Repositories:
- Discord-BOT:            https://github.com/dev-harshhh19/Discord-BOT
- Code-X:                 https://github.com/dev-harshhh19/Code-X
- ProjectFlow:            https://github.com/dev-harshhh19/ProjectFlow
- Gamiex:                 https://github.com/dev-harshhh19/Gamiex
- FuelSim:                https://github.com/dev-harshhh19/FuelSim-Gas-Station-Pump-Simulator
- Report-card-Dashboard:  https://github.com/dev-harshhh19/Report-card-Dashboard
- Harshad-Portfolio:      https://github.com/dev-harshhh19/Harshad-Portfolio
`, p.Handle, p.GitHub)
	addFile(ghDir, "profile.txt", ghContent, "GitHub profile and repo list")

	// 7. /contact
	contactDir := addDir(root, "contact", "Contact information and connection channels")
	contactContent := fmt.Sprintf(`# Contact & Connect
==================================================
Email:       %s
LinkedIn:    %s
GitHub:      %s
Website:     %s
Calendly:    %s
SSH Gateway: ssh %s

Available for full-time engineering roles, high-impact contract work,
and technical discussions.
`, p.Email, p.LinkedIn, p.GitHub, p.Website, p.Calendly, p.SSHHost)
	addFile(contactDir, "links.txt", contactContent, "Contact methods and direct links")

	// 8. /resume
	resumeDir := addDir(root, "resume", "Curriculum vitae and technical experience summary")
	resumeContent := fmt.Sprintf(`# Harshad Nikam — Full-Stack Engineer Resume
==================================================
Email:   %s
Website: %s
GitHub:  %s

Summary:
Full-Stack Engineer specialized in TypeScript, Next.js, React, and Node.js.
Passionate about low-latency web architecture, distributed systems, clean APIs,
and developer tooling.

Primary Technologies:
- Languages:  TypeScript, JavaScript, Go, Python, SQL, HTML/CSS
- Frontend:   React, Next.js, Tailwind CSS, GSAP, Framer Motion
- Backend:    Node.js, Express, Socket.io, REST APIs, GraphQL
- Databases:  PostgreSQL, MongoDB, Prisma ORM
- Cloud/Ops:  Docker, Linux, Cloudflare R2, Vercel, CI/CD, Git

Key Achievements:
- Freelivo: Production SaaS invoicing tool with Stripe/Razorpay and Cloudflare R2.
- TomMC-SMP Aternos Manager: Automated Minecraft server management with telemetry.
- ProjectFlow: Real-time board with Socket.io diff reconciliation on PostgreSQL.
- Testing Home Server: Built multi-container LXC/Docker environment on Proxmox VE.

Download / PDF:
To fetch the PDF version directly from your terminal:
  curl -sL https://harshadnikam.me/resume.pdf -o harshad_nikam_resume.pdf
`, p.Email, p.Website, p.GitHub)
	addFile(resumeDir, "resume.txt", resumeContent, "Full text resume summary")

	return vfs
}

// CleanPath normalizes a virtual path
func CleanPath(p string) string {
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	cleaned := path.Clean(p)
	if cleaned == "" {
		return "/"
	}
	return cleaned
}

// ResolvePath takes current directory and target argument, returning normalized absolute path
func (v *VFS) ResolvePath(currentDir, target string) string {
	target = strings.TrimSpace(target)
	if target == "" || target == "." {
		return currentDir
	}
	if target == "~" {
		return "/"
	}
	if strings.HasPrefix(target, "~/") {
		target = "/" + target[2:]
	}
	if strings.HasPrefix(target, "/") {
		return CleanPath(target)
	}
	return CleanPath(path.Join(currentDir, target))
}

// GetNode looks up an FSNode by absolute or relative path
func (v *VFS) GetNode(absPath string) (*FSNode, bool) {
	absPath = CleanPath(absPath)
	if absPath == "/" {
		return v.Root, true
	}

	parts := strings.Split(strings.Trim(absPath, "/"), "/")
	curr := v.Root
	for _, part := range parts {
		if part == "" {
			continue
		}
		child, exists := curr.Children[part]
		if !exists {
			return nil, false
		}
		curr = child
	}
	return curr, true
}

// ListDirectory returns sorted children of a directory
func (v *VFS) ListDirectory(absPath string) ([]*FSNode, error) {
	node, found := v.GetNode(absPath)
	if !found {
		return nil, fmt.Errorf("directory not found: %s", absPath)
	}
	if node.Type != DirNode {
		return nil, fmt.Errorf("not a directory: %s", absPath)
	}

	var items []*FSNode
	for _, c := range node.Children {
		items = append(items, c)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Type != items[j].Type {
			return items[i].Type == DirNode
		}
		return items[i].Name < items[j].Name
	})
	return items, nil
}

// Tree generates an ASCII representation of the virtual filesystem
func (v *VFS) Tree(startPath string) (string, error) {
	node, found := v.GetNode(startPath)
	if !found {
		return "", fmt.Errorf("path not found: %s", startPath)
	}

	var sb strings.Builder
	sb.WriteString(node.Path + "\n")
	v.buildTree(&sb, node, "")
	return sb.String(), nil
}

func (v *VFS) buildTree(sb *strings.Builder, current *FSNode, prefix string) {
	keys := make([]string, 0, len(current.Children))
	for k := range current.Children {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		isLast := i == len(keys)-1
		child := current.Children[k]

		connector := "├── "
		childPrefix := prefix + "│   "
		if isLast {
			connector = "└── "
			childPrefix = prefix + "    "
		}

		display := child.Name
		if child.Type == DirNode {
			display += "/"
		}
		sb.WriteString(prefix + connector + display + "\n")

		if child.Type == DirNode {
			v.buildTree(sb, child, childPrefix)
		}
	}
}
