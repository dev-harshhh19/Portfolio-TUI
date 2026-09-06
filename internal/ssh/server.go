package ssh

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/dev-harshhh19/harshad-tui/internal/tui"
)

// Config holds configuration for the Wish SSH daemon
type Config struct {
	Host        string
	Port        int
	KeyPath     string
	IdleTimeout time.Duration
	MaxTimeout  time.Duration
}

// DefaultConfig returns production-ready defaults
func DefaultConfig() Config {
	return Config{
		Host:        "0.0.0.0",
		Port:        2222,
		KeyPath:     filepath.Join(".ssh", "harshad_ed25519"),
		IdleTimeout: 15 * time.Minute,
		MaxTimeout:  2 * time.Hour,
	}
}

// StartServer starts the SSH listener and blocks until context cancellation or interrupt signal
func StartServer(ctx context.Context, cfg Config) error {
	// Ensure directory for host key exists
	if dir := filepath.Dir(cfg.KeyPath); dir != "." && dir != "" {
		_ = os.MkdirAll(dir, 0700)
	}

	addr := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))

	s, err := wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath(cfg.KeyPath),
		wish.WithIdleTimeout(cfg.IdleTimeout),
		wish.WithMaxTimeout(cfg.MaxTimeout),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			logging.Middleware(),
		),
	)
	if err != nil {
		return fmt.Errorf("could not create wish ssh server: %w", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Starting Harshad Nikam SSH Portfolio Server", "host", cfg.Host, "port", cfg.Port, "addr", addr)
	log.Info("Access via: ssh -p " + fmt.Sprintf("%d", cfg.Port) + " <host> or ssh cli.harshadnikam.me")

	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	select {
	case <-done:
	case <-ctx.Done():
	}

	log.Info("Stopping SSH server gracefully...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		return fmt.Errorf("could not stop ssh server: %w", err)
	}
	log.Info("SSH server stopped cleanly")
	return nil
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, active := s.Pty()
	if !active {
		wish.Fatalln(s, "No active PTY requested. Please run ssh with terminal allocation enabled: ssh -t <host>")
		return nil, nil
	}

	m := tui.NewModel()
	m.SetWindowSize(pty.Window.Width, pty.Window.Height)

	return m, []tea.ProgramOption{
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	}
}
