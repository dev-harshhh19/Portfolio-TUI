package content

// Project represents a portfolio project with real metadata
type Project struct {
	ID          string   `json:"id"`
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Tagline     string   `json:"tagline"`
	Year        string   `json:"year"`
	Description string   `json:"description"`
	Detail      string   `json:"detail"`
	Tech        []string `json:"tech"`
	Live        string   `json:"live"`
	GitHub      string   `json:"github"`
}

// SkillCategory groups technical skills
type SkillCategory struct {
	Title  string   `json:"title"`
	Skills []string `json:"skills"`
}

// BlogPost represents real articles published by Harshad
type BlogPost struct {
	Slug          string   `json:"slug"`
	Title         string   `json:"title"`
	ShortTitle    string   `json:"short_title"`
	PublishedDate string   `json:"published_date"`
	ReadTime      string   `json:"read_time"`
	Description   string   `json:"description"`
	Tags          []string `json:"tags"`
	Content       string   `json:"content"`
}

// Principle represents Harshad's engineering philosophy / work style
type Principle struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// HarshadProfile contains the complete personal and professional profile
type HarshadProfile struct {
	Name       string
	Handle     string
	Title      string
	Tagline    string
	Bio        string
	Status     string
	Email      string
	GitHub     string
	LinkedIn   string
	Website    string
	Calendly   string
	SSHHost    string
	Philosophy string
	Projects   []Project
	Skills     []SkillCategory
	Strengths  []Principle
	WorkStyle  []Principle
	BlogPosts  []BlogPost
}

// Profile is an alias for HarshadProfile
type Profile = HarshadProfile

// GetProfile returns the canonical, real data for Harshad Nikam
func GetProfile() HarshadProfile {
	return HarshadProfile{
		Name:       "Harshad Nikam",
		Handle:     "dev-harshhh",
		Title:      "Full-Stack Engineer",
		Tagline:    "Full-Stack Engineer who actually ships things.",
		Status:     "Currently open to opportunities",
		Email:      "hi@harshadnikam.me",
		GitHub:     "https://github.com/dev-harshhh19",
		LinkedIn:   "https://www.linkedin.com/in/harshad-nikam-311734281",
		Website:    "https://harshadnikam.me",
		Calendly:   "https://calendly.com/dev-harshhh",
		SSHHost:    "cli.harshadnikam.me",
		Philosophy: "Everything starts at /",
		Bio: `I build full-stack apps end to end. TypeScript, React, and Node.js is mostly what I work with.
I like figuring out how to make things fast, clean, and easy to maintain.

Good code matters to me. So does accessibility, good API design, and making the UI feel right.
I try not to over-engineer things.

Right now I am working on cloud infra, digging into backend performance, and experimenting with
AI features that actually make sense to have.`,
		Projects: []Project{
			{
				ID:          "freelivo",
				Slug:        "freelivo",
				Title:       "Freelivo",
				Tagline:     "Invoicing platform for freelancers",
				Year:        "2026",
				Description: "Invoicing tool built for freelancers. Generate branded invoices, track payments, get analytics. Used by real clients.",
				Detail: `Full-stack SaaS built from scratch: designed the relational data model, built the auth layer with NextAuth, wired up Stripe Connect and Razorpay for multi-tenant payments, set up Cloudflare R2 for file storage, and handled all transactional email with Resend + React Email.

The billing logic alone handled edge cases across multi-currency invoicing, prorated taxes, and webhook settlement events.`,
				Tech:   []string{"Next.js 16", "TypeScript 5", "PostgreSQL", "Prisma ORM", "Cloudflare R2", "Tailwind CSS 4", "NextAuth.js", "Razorpay", "Stripe Connect", "Resend"},
				Live:   "https://freelivo.tech/",
				GitHub: "NA",
			},
			{
				ID:          "discord-bot",
				Slug:        "discord-bot",
				Title:       "Discord Bot",
				Tagline:     "TomMC-SMP Aternos Manager & Discord Bot",
				Year:        "2026",
				Description: "An automation tool for managing Minecraft servers on Aternos with a Discord bot, web dashboard, real-time server monitoring, and automatic recovery.",
				Detail: `TomMC-SMP Aternos Manager & Discord Bot is a full-featured automation system designed to simplify Minecraft server management on Aternos.

It automates server startup and queue handling using Puppeteer, provides real-time player and server status through direct network monitoring, and offers both Discord commands and a web dashboard for easy control.

The system includes role-based permissions, secure authentication, automatic recovery watchdogs for stuck processes, and performance optimizations for low-resource devices such as Raspberry Pi, Termux, and VPS containers.`,
				Tech:   []string{"TypeScript", "Puppeteer", "Node.js", "Express", "Docker"},
				Live:   "NA",
				GitHub: "https://github.com/dev-harshhh19/Discord-BOT",
			},
			{
				ID:          "code-x",
				Slug:        "code-x",
				Title:       "Code X",
				Tagline:     "Browser-side cipher encoder",
				Year:        "2026",
				Description: "Encode and decode text with custom ciphers, all in the browser. No server, no data leaves your machine.",
				Detail: `Built to learn cipher theory and ended up shipping it. Everything runs in the browser — no backend, no telemetry.

The cipher logic is modular so new schemes can be dropped in without touching the UI layer. Zero external tracking, zero server dependencies.`,
				Tech:   []string{"TypeScript", "Tailwind CSS", "Next.js"},
				Live:   "https://code-x-orian.vercel.app/",
				GitHub: "https://github.com/dev-harshhh19/Code-X",
			},
			{
				ID:          "projectflow",
				Slug:        "projectflow",
				Title:       "ProjectFlow",
				Tagline:     "Real-time project management",
				Year:        "2025",
				Description: "Project management tool with live task syncing via Socket.io. Teams track tasks, chat, and see progress.",
				Detail: `The real challenge here was keeping Socket.io state consistent when multiple users edited the same board simultaneously.

Solved with server-authoritative events and optimistic UI updates on the client. PostgreSQL handles the persistence; the socket layer handles the live diff and reconciliation.`,
				Tech:   []string{"React", "Node.js", "Socket.io", "PostgreSQL"},
				Live:   "https://project-flow-dev.vercel.app/",
				GitHub: "https://github.com/dev-harshhh19/ProjectFlow",
			},
			{
				ID:          "gamiex",
				Slug:        "gamiex",
				Title:       "Gamiex",
				Tagline:     "E-commerce for gaming gear",
				Year:        "2025",
				Description: "E-commerce site for gaming gear — product catalog, cart, Stripe checkout, admin panel.",
				Detail: `Built a full e-commerce flow: product listings, a cart with persistent state, Stripe checkout with webhook order confirmation, and a protected admin panel to manage inventory and orders.

Focused on keeping the UI fast and the purchase funnel friction-free.`,
				Tech:   []string{"React", "Node.js", "Express", "MongoDB", "Stripe"},
				Live:   "https://gamiex.vercel.app/",
				GitHub: "https://github.com/dev-harshhh19/Gamiex",
			},
			{
				ID:          "fuelsim",
				Slug:        "fuelsim",
				Title:       "FuelSim",
				Tagline:     "Gas station pump simulator",
				Year:        "2025",
				Description: "A gas station pump simulator built on the Canvas API — physics-accurate, responsive across screen sizes.",
				Detail: `A weekend experiment that grew into a proper project. The physics loop runs on requestAnimationFrame with a fixed timestep so it stays consistent at different refresh rates.

The hardest part was making the nozzle interaction and fluid flow feel natural at small screen sizes.`,
				Tech:   []string{"HTML", "CSS", "JavaScript", "Canvas API"},
				Live:   "https://fuelsim.netlify.app/",
				GitHub: "https://github.com/dev-harshhh19/FuelSim-Gas-Station-Pump-Simulator/",
			},
			{
				ID:          "student-report-card",
				Slug:        "student-report-card",
				Title:       "Student Report Card",
				Tagline:     "Teacher dashboard for grade tracking",
				Year:        "2025",
				Description: "Dashboard for teachers to view student performance. Charts, filters, grade breakdowns. Simple and fast.",
				Detail: `Built for a specific school that needed something lightweight — no SaaS overhead, no login walls.

Teachers can filter by class, subject, or student and get an instant grade breakdown. Charts are rendered with vanilla Canvas to keep the bundle small.`,
				Tech:   []string{"HTML", "CSS", "JavaScript", "Python", "Flask", "Chart.js"},
				Live:   "https://report-card-dashboard.onrender.com/",
				GitHub: "https://github.com/dev-harshhh19/Report-card-Dashboard",
			},
		},
		Skills: []SkillCategory{
			{
				Title:  "Frontend",
				Skills: []string{"React", "Next.js", "TypeScript", "JavaScript", "Tailwind CSS", "Framer Motion", "GSAP"},
			},
			{
				Title:  "Backend",
				Skills: []string{"Node.js", "Express", "REST APIs", "Socket.io", "PostgreSQL", "MongoDB", "Prisma ORM"},
			},
			{
				Title:  "Cloud & DevOps",
				Skills: []string{"Docker", "Vercel", "Cloudflare R2", "Linux", "DigitalOcean", "CI/CD", "GitHub Actions"},
			},
			{
				Title:  "Tools & Services",
				Skills: []string{"Git & GitHub", "Postman", "Resend", "Stripe Connect", "Razorpay", "Supabase", "Firebase", "Zed", "Vim"},
			},
		},
		Strengths: []Principle{
			{
				Title: "Ownership",
				Body:  "When I pick something up, I see it through. I don't wait to be told what needs fixing. If something is broken or unclear, I say so and try to resolve it.",
			},
			{
				Title: "Learning fast",
				Body:  "I've taught myself most of what I know. New languages, new frameworks, new domains — I figure things out quickly and I retain them.",
			},
			{
				Title: "Reliability",
				Body:  "If I say something will be done, it will be done. If something comes up and changes that, I communicate early. Trust is built in small moments, not grand gestures.",
			},
			{
				Title: "Clean code",
				Body:  "I care about the code after it's written. Future me, and future teammates, should be able to read it without confusion. That means naming things well, keeping logic flat, and not being clever when simple works.",
			},
			{
				Title: "Curiosity",
				Body:  "I genuinely enjoy understanding how things work, not just at the surface level. That curiosity has taken me into database internals, compiler behavior, and networking.",
			},
			{
				Title: "Problem solving",
				Body:  "I slow down before jumping to solutions. A poorly framed problem leads to a solution that fixes the wrong thing. I try to understand the root before writing a single line of code.",
			},
		},
		WorkStyle: []Principle{
			{
				Title: "Communicate early",
				Body:  "If something is going wrong or I'm unsure about a direction, I say so before it becomes a problem. Silence is the most expensive thing on a team.",
			},
			{
				Title: "Writing over talking",
				Body:  "I document decisions, write clear PR descriptions, and keep things written down. It's respectful of other people's time and builds a record that helps everyone.",
			},
			{
				Title: "Plan, stay flexible",
				Body:  "I think through what I'm going to build before I build it. But I hold that plan loosely. Requirements change. Rigidly following a plan when reality has shifted is stubbornness, not discipline.",
			},
			{
				Title: "Enjoy the people",
				Body:  "I like reviewing code, pair-debugging hard problems, and talking through architecture. Other people catch things I miss. I don't want to work alone.",
			},
		},
		BlogPosts: []BlogPost{
			{
				Slug:          "turning-old-laptop-into-testing-home-server",
				Title:         "How I Turned My Old Potato Laptop into a Powerful Testing Home Server",
				ShortTitle:    "Potato Laptop to Testing Home Server",
				PublishedDate: "Aug 2026",
				ReadTime:      "8 min read",
				Description:   "A clean, practical guide to converting an old Intel i5-7300U laptop with 16GB DDR4 RAM into a high-performance Proxmox VE testing home server for lightweight LXC containers, Docker, databases, and microservices.",
				Tags:          []string{"Home Server", "Proxmox", "Docker", "LXC", "DevOps"},
				Content: `# How I Turned My Old Potato Laptop into a Powerful Testing Home Server

By Harshad Nikam · Aug 2026 · 8 min read

## The Problem
Old laptops often end up in drawers gathering dust. Mine had an Intel Core i5-7300U (2 cores, 4 threads), 16GB DDR4 RAM, and a cracked screen. Instead of throwing it away, I turned it into a low-power, 24/7 home lab testing server.

## Hardware Specifications
- CPU: Intel Core i5-7300U @ 2.60GHz (Turbo up to 3.50GHz)
- Memory: 16 GB DDR4 RAM
- Storage: 512 GB SATA SSD
- Network: Gigabit Ethernet + 802.11ac Wi-Fi
- Power Draw: ~12-18W idle (integrated battery acts as a built-in UPS!)

## The OS: Proxmox VE 8.x
Proxmox gives you enterprise-grade virtualization (KVM) and lightweight Linux Containers (LXC) on Debian Linux with a web management interface.

LXC containers share the host Linux kernel, consuming virtually zero idle CPU and minimal RAM (a typical Alpine/Debian container uses under 40MB RAM).

## Key Workloads Hosted
1. Docker Host (Ubuntu LXC):
   - Nginx Proxy Manager (SSL certificates & routing)
   - Portainer (container management)
   - PostgreSQL 16 & Redis 7 for local dev testing
2. DNS Sinkhole & Local Resolver (Pi-hole)
3. Automated Backups to Cloudflare R2 via rclone

## Lessons Learned
1. Always disable laptop lid sleep in /etc/systemd/logind.conf:
   HandleLidSwitch=ignore
2. Keep the laptop propped up on small rubber feet for airflow.
3. Clean thermal paste every 2 years — temps dropped by 14°C.
4. The built-in laptop battery provides a free UPS during brief power cuts!`,
			},
			{
				Slug:          "how-to-get-started-with-claude-code",
				Title:         "How to Get Started with Claude Code: Complete Installation Guide for Windows, macOS & Linux",
				ShortTitle:    "How to Get Started with Claude Code",
				PublishedDate: "Aug 2026",
				ReadTime:      "6 min read",
				Description:   "Learn what Claude Code is, why developers love it, and how to install and configure it on Windows, macOS, and Linux with custom endpoints.",
				Tags:          []string{"AI Tools", "Claude Code", "CLI", "Productivity", "Setup Guide"},
				Content: `# How to Get Started with Claude Code

By Harshad Nikam · Aug 2026 · 6 min read

## What is Claude Code?
Claude Code is Anthropic's agentic command-line tool that lives directly in your terminal. It understands your codebase, edits files, runs tests, fixes errors, and creates git commits seamlessly.

## Installation
Claude Code requires Node.js 18 or later.

    npm install -g @anthropic-ai/claude-code

## Quick Start
Navigate to any project directory in your terminal:

    cd my-project
    claude

On first run, it will authenticate via your Anthropic account or API key.

## Best Practices
1. Compact Prompts: Give concrete goals ("Fix the CORS error in src/api/auth.ts").
2. Let it run tests: Ask it to verify changes with your existing test runner.
3. Git Hygiene: Commit clean work before launching big refactorings.

Claude Code bridges the gap between terminal speed and AI assistance.`,
			},
		},
	}
}
