package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ckinan/cktop/internal/domain"
	"github.com/ckinan/cktop/internal/util"
)

type snapshotMsg domain.Snapshot

func waitForSnapshot(ch <-chan domain.Snapshot) tea.Cmd {
	return func() tea.Msg {
		snap, ok := <-ch
		if !ok {
			return nil
		}
		return snapshotMsg(snap)
	}
}

func calcDir(showDir bool, sortDesc bool) string {
	if !showDir {
		return ""
	}
	if sortDesc == true {
		return " ▼"
	}
	return " ▲"
}

func (m *Model) applySort() {
	// Reserve lines for CPU header (1) + RAM header (1) + blank (1) + [table content] + blank (1) + footer (1) = 5
	// bubbles/table renders its own column header row internally
	cmdW := max(20, m.width-colPIDWidth-colPPIDWidth-colUserWidth-colCPUWidth-colRSSWidth)

	m.table.SetColumns([]table.Column{
		{Title: "PID" + calcDir(m.sortBy == SortByPID, m.sortDesc), Width: colPIDWidth},
		{Title: "PPID", Width: colPPIDWidth},
		{Title: "User", Width: colUserWidth},
		{Title: "CPU%" + calcDir(m.sortBy == SortByCPU, m.sortDesc), Width: colCPUWidth},
		{Title: "RSS" + calcDir(m.sortBy == SortByRSS, m.sortDesc), Width: colRSSWidth},
		{Title: "CmdLine" + calcDir(m.sortBy == SortByCmdLine, m.sortDesc), Width: cmdW},
	})

	var sorted []domain.Process
	switch m.sortBy {
	case SortByRSS:
		sorted = util.SortBy(m.procs, func(p domain.Process) int { return p.Rss }, m.sortDesc)
	case SortByCPU:
		sorted = util.SortBy(m.procs, func(p domain.Process) float64 { return p.CPU }, m.sortDesc)
	case SortByPID:
		sorted = util.SortBy(m.procs, func(p domain.Process) int { return p.Pid }, m.sortDesc)
	case SortByPPID:
		sorted = util.SortBy(m.procs, func(p domain.Process) int { return p.Ppid }, m.sortDesc)
	case SortByCmdLine:
		sorted = util.SortBy(m.procs, func(p domain.Process) string { return p.Cmdline }, m.sortDesc)
	}

	rows := make([]table.Row, len(sorted))
	for i, p := range sorted {
		rows[i] = table.Row{
			fmt.Sprintf("%d", p.Pid),
			fmt.Sprintf("%d", p.Ppid),
			p.Username,
			fmt.Sprintf("%.2f%%", p.CPU),
			util.HumanBytes(p.Rss),
			p.Cmdline,
		}
	}
	m.table.SetRows(rows)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case snapshotMsg:
		snap := domain.Snapshot(msg)
		m.CPU = msg.CPU
		m.memory = msg.Memory
		wasEmpty := len(m.procs) == 0 // first data arrival?
		m.procs = snap.Processes
		m.applySort()
		if wasEmpty {
			m.table.GotoTop()
		}
		return m, waitForSnapshot(m.snapCh)
	case tea.KeyMsg:
		prev := m.sortBy
		isSortKey := true
		switch msg.String() {
		case "M":
			m.sortBy = SortByRSS
		case "C":
			m.sortBy = SortByCPU
		case "P":
			m.sortBy = SortByPID
		case "L":
			m.sortBy = SortByCmdLine
		case "q":
			return m, tea.Quit
		default:
			isSortKey = false
		}
		if isSortKey {
			if m.sortBy == prev {
				// same key: toggle direction
				m.sortDesc = !m.sortDesc
			} else {
				// new field: reset to descending
				m.sortDesc = true
			}
			m.applySort()
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.table.SetHeight(m.height - 5)
		m.applySort()

		return m, nil
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}
