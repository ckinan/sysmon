package internal

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Ram struct {
	MemTotal     int
	MemAvailable int
	MemUsed      int
}

func GetRam() (Ram, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return Ram{}, err
	}
	defer file.Close()

	var memTotal, memAvailable string
	var memTotalInt, memAvailableInt, memUsedInt int
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "MemTotal:") {
			memTotal = line
		} else if strings.HasPrefix(line, "MemAvailable:") {
			memAvailable = line
		}

		if memTotal != "" && memAvailable != "" {
			break
		}
	}

	memTotalInt, err = extractFieldFromLine(memTotal)
	if err != nil {
		return Ram{}, err
	}
	memAvailableInt, err = extractFieldFromLine(memAvailable)
	if err != nil {
		return Ram{}, err
	}
	memUsedInt = memTotalInt - memAvailableInt

	return Ram{
		MemTotal:     memTotalInt,
		MemAvailable: memAvailableInt,
		MemUsed:      memUsedInt,
	}, nil
}

func extractFieldFromLine(line string) (int, error) {
	fields := strings.Fields(line)
	intval, err := strconv.Atoi(fields[1])
	if err != nil {
		return -1, err
	}
	return intval, nil
}
