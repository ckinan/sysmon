package main

import (
	"fmt"

	"github.com/ckinan/system-monitor.go/internal"
)

func main() {
	ram, err := internal.GetRam()
	if err != nil {
		fmt.Println(fmt.Errorf("error: %w", err))
	}
	fmt.Printf("used memory: %d, available memory: %d, total memory: %d\n", ram.MemUsed, ram.MemAvailable, ram.MemTotal)
}
