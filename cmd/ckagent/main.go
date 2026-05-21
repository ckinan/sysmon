package main

import (
	"flag"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/ckinan/cktop/internal/adapters/proc"
	"github.com/ckinan/cktop/internal/domain"
)

type memCache struct {
	mu  sync.RWMutex
	mem domain.Memory
}

func (c *memCache) refresh(r proc.ProcMemoryReader) {
	m, err := r.ReadMemory()
	if err != nil {
		log.Printf("read memory: %v", err)
		return
	}
	c.mu.Lock()
	c.mem = m
	c.mu.Unlock()
}

func (c *memCache) get() domain.Memory {
	// RLock blocks if refresh() is currently writing, ensuring we never read
	// a partially-updated struct. Multiple concurrent reads are allowed.
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.mem
}

type cpuCache struct {
	mu    sync.RWMutex
	value float64
}

func (c *cpuCache) start(interval time.Duration, r *proc.ProcCPUReader) {
	// Seed the baseline sample so the first tick produces a real delta.
	if _, err := r.ReadCPU(); err != nil {
		log.Printf("seed cpu baseline: %v", err)
	}
	go func() {
		t := time.NewTicker(interval)
		for range t.C {
			v, err := r.ReadCPU()
			if err != nil {
				log.Printf("read cpu: %v", err)
				continue
			}
			c.mu.Lock()
			c.value = v
			c.mu.Unlock()
		}
	}()
}

func (c *cpuCache) get() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

func main() {
	addr := flag.String("addr", ":9110", "HTTP listen address for /metrics")
	flag.Parse()

	memReader := proc.ProcMemoryReader{}
	mem := &memCache{}

	metrics.NewGauge(`ckagent_memory_total_bytes`, func() float64 { return float64(mem.get().Total) })
	metrics.NewGauge(`ckagent_memory_free_bytes`, func() float64 { return float64(mem.get().Free) })
	metrics.NewGauge(`ckagent_memory_available_bytes`, func() float64 { return float64(mem.get().Available) })
	metrics.NewGauge(`ckagent_memory_used_bytes`, func() float64 { return float64(mem.get().Used) })
	metrics.NewGauge(`ckagent_memory_buffers_bytes`, func() float64 { return float64(mem.get().Buffers) })
	metrics.NewGauge(`ckagent_memory_cached_bytes`, func() float64 { return float64(mem.get().Cached) })
	metrics.NewGauge(`ckagent_memory_shmem_bytes`, func() float64 { return float64(mem.get().Shmem) })
	metrics.NewGauge(`ckagent_swap_total_bytes`, func() float64 { return float64(mem.get().SwapTotal) })
	metrics.NewGauge(`ckagent_swap_free_bytes`, func() float64 { return float64(mem.get().SwapFree) })
	metrics.NewGauge(`ckagent_swap_used_bytes`, func() float64 { return float64(mem.get().SwapUsed) })

	cpu := &cpuCache{}
	cpu.start(5*time.Second, proc.NewProcCPUReader())
	metrics.NewGauge(`ckagent_cpu_usage_percent`, func() float64 { return cpu.get() })

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		mem.refresh(memReader)
		metrics.WritePrometheus(w, false)
	})

	log.Printf("ckagent listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
