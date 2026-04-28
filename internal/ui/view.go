package ui

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/ckinan/sysmon/internal"
)

func (m Model) View() string {
	header := fmt.Sprintf("RAM used: %d kB / %d kB\n\n", m.ram.MemUsed, m.ram.MemTotal)
	var buf strings.Builder
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0) // minwidth=0, tabwidth=0, padding=2
	fmt.Fprintln(w, "PID\tNAME\tSTATE\tTHREADS\tRSS(kB)")
	fmt.Fprintln(w, "---\t----\t-----\t-------\t-------")

	// TODO Let's not do hardcoded limits here
	sorted := slices.Clone(m.procs)
	slices.SortFunc(sorted, func(a, b internal.Process) int {
		return cmp.Compare(b.RssKB, a.RssKB)
	})
	limit := 30
	limit = min(limit, len(sorted))
	for _, p := range sorted[:limit] {
		fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%d\n", p.Pid, p.Name, p.State, p.Threads, p.RssKB)
	}

	w.Flush() // must call Flush because tabwriter only writes to buf after seeing all rows
	return header + buf.String() + "\nPress q to quit"
}
