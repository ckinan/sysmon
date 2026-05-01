//go:build integration

package gopsutil_test

import (
	"testing"

	"github.com/ckinan/sysmon/internal/adapters/gopsutil"
)

func TestGopsutilMemoryReader_RealSystem(t *testing.T) {
	r := gopsutil.GopsutilMemoryReader{}
	mem, err := r.ReadMemory()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mem.Total == 0 {
		t.Error("Total memory should be > 0 on a real machine")
	}
	if mem.Used == 0 {
		t.Error("Used memory should be > 0 on a real machine")
	}
}

func TestGopsutilCPUReader_RealSystem(t *testing.T) {
	r := gopsutil.GopsutilCPUReader{}
	cpu, err := r.ReadCPU()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cpu < 0 || cpu > 100 {
		t.Errorf("CPU percent should be between 0-100, got %.2f", cpu)
	}
}

func TestGopsutilProcessReader_RealSystem(t *testing.T) {
	r := gopsutil.NewGopsutilProcessReader()
	procs, err := r.ReadProcesses()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(procs) == 0 {
		t.Error("expected at least one process on a real machine")
	}
}
