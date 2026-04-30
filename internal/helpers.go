package internal

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

func HumanBytes(b int) string {
	switch {
	case b >= 1<<30: // >= 1GiB
		return fmt.Sprintf("%.2f GiB", float64(b)/float64(1<<30))
	case b >= 20: // >= 1 MiB
		return fmt.Sprintf("%.2f MiB", float64(b)/float64(1<<20))
	default:
		return fmt.Sprintf("%d B", b)
	}
}

// TODO remove this unused function
func extractFieldFromLine(line string) (string, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return "", fmt.Errorf("invalid line, expected at least 2 fields, got %v, line: %s", len(fields), line)
	}
	return fields[1], nil
}

func SortBy[T any, K cmp.Ordered](items []T, key func(T) K, desc bool) []T {
	out := slices.Clone(items)
	slices.SortFunc(out, func(a, b T) int {
		if desc {
			return cmp.Compare(key(b), key(a))
		}
		return cmp.Compare(key(a), key(b))
	})
	return out
}
