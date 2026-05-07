package gopsutil

import (
	"fmt"

	"github.com/ckinan/cktop/internal/domain"
	"github.com/shirou/gopsutil/v4/mem"
)

type GopsutilMemoryReader struct{}

func (g GopsutilMemoryReader) ReadMemory() (domain.Memory, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return domain.Memory{}, fmt.Errorf("gopsutil VirtualMemory: %w", err)
	}
	return domain.Memory{
		Total:     int64(v.Total),
		Available: int64(v.Available),
		Used:      int64(v.Used),
	}, nil
}
