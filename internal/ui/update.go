package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ckinan/sysmon/internal/collector"
)

type snapshotMsg collector.Snapshot

func waitForSnapshot(ch <-chan collector.Snapshot) tea.Cmd {
	return func() tea.Msg {
		snap, ok := <-ch
		if !ok {
			return nil
		}
		return snapshotMsg(snap)
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case snapshotMsg:
		m.ram = msg.Ram
		m.procs = msg.Processes
		return m, waitForSnapshot(m.snapCh)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}
