package gopsutil

import "github.com/shirou/gopsutil/v4/cpu"

type GopsutilCPUReader struct{}

func (g GopsutilCPUReader) ReadCPU() (float64, error) {
	cpuPcts, err := cpu.Percent(0, false)
	if err != nil || len(cpuPcts) == 0 {
		return 0, err
	}
	return cpuPcts[0], nil
}
