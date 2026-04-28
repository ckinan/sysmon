package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ckinan/sysmon/internal"
	"github.com/ckinan/sysmon/internal/collector"
)

// Model is the bubbletea model. It holds all UI state
type Model struct {
	snapCh <-chan collector.Snapshot // read-only channel from the collect
	ram    internal.Ram
	procs  []internal.Process
}

func New(ch <-chan collector.Snapshot) Model {
	return Model{snapCh: ch}
}

func (m Model) Init() tea.Cmd {
	return waitForSnapshot(m.snapCh)
}
