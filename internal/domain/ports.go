package domain

type MemoryReader interface {
	ReadMemory() (Memory, error)
}

type ProcessReader interface {
	ReadProcesses() ([]Process, error)
}

type CPUReader interface {
	ReadCPU() (float64, error)
}
