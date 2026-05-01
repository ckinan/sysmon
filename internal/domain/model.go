package domain

type Memory struct {
	Total     int
	Available int
	Used      int
}

type Process struct {
	Pid      int
	Ppid     int
	Rss      int
	CPU      float64
	Cmdline  string
	Username string
}

type Snapshot struct {
	CPU       float64
	Memory    Memory
	Processes []Process
}
