package main

import (
	"flag"
	"log"
	"net/http"
	"sync"

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

func main() {
	addr := flag.String("addr", ":9110", "HTTP listen address for /metrics")
	flag.Parse()

	reader := proc.ProcMemoryReader{}
	cache := &memCache{}

	metrics.NewGauge(`ckagent_memory_total_bytes`, func() float64 { return float64(cache.get().Total) })
	metrics.NewGauge(`ckagent_memory_free_bytes`, func() float64 { return float64(cache.get().Free) })
	metrics.NewGauge(`ckagent_memory_available_bytes`, func() float64 { return float64(cache.get().Available) })
	metrics.NewGauge(`ckagent_memory_used_bytes`, func() float64 { return float64(cache.get().Used) })
	metrics.NewGauge(`ckagent_memory_buffers_bytes`, func() float64 { return float64(cache.get().Buffers) })
	metrics.NewGauge(`ckagent_memory_cached_bytes`, func() float64 { return float64(cache.get().Cached) })
	metrics.NewGauge(`ckagent_memory_shmem_bytes`, func() float64 { return float64(cache.get().Shmem) })
	metrics.NewGauge(`ckagent_swap_total_bytes`, func() float64 { return float64(cache.get().SwapTotal) })
	metrics.NewGauge(`ckagent_swap_free_bytes`, func() float64 { return float64(cache.get().SwapFree) })
	metrics.NewGauge(`ckagent_swap_used_bytes`, func() float64 { return float64(cache.get().SwapUsed) })

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		cache.refresh(reader)
		metrics.WritePrometheus(w, false)
	})

	log.Printf("ckagent listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
