package main

import (
	"context"
	"log/slog"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ckinan/sysmon/internal/collector"
	"github.com/ckinan/sysmon/internal/ui"
)

func main() {
	// context.WithCancel gives us a cancel function to stop the collector cleanly.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // ensure goroutine is always stopped when main() exits

	// start collector - returns a channel immediately, goroutine runs in background
	snapshotCh := collector.Start(ctx, 2*time.Second)

	// Read a few snapshots then exit (for now)
	p := tea.NewProgram(ui.New(snapshotCh), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		slog.Error("error running TUI", "err", err)
	}
}
