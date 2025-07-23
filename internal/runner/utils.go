package runner

import (
	"fmt"
	"sort"
	"time"
	"strings"
)

func parseHeaders(headers []string) map[string]string {
	result := make(map[string]string)
	for _, h := range headers {
		parts := splitHeader(h)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

func splitHeader(h string) []string {
	sepIndex := -1
	for i, r := range h {
		if r == ':' {
			sepIndex = i
			break
		}
	}
	if sepIndex == -1 {
		return nil
	}
	key := h[:sepIndex]
	val := h[sepIndex+1:]
	return []string{trim(key), trim(val)}
}

func trim(s string) string {
	return strings.TrimSpace(s)
}

func printSummary(m *Metrics, duration time.Duration) {
	sort.Slice(m.Latencies, func(i, j int) bool {
		return m.Latencies[i] < m.Latencies[j]
	})

	p95 := m.Latencies[int(0.95*float64(len(m.Latencies)))]
	min := m.Latencies[0]
	max := m.Latencies[len(m.Latencies)-1]
	avg := averageLatency(m.Latencies)

	rps := float64(m.TotalRequests) / duration.Seconds()

	fmt.Println("\n--- Benchmark Summary ---")
	fmt.Printf("Total Requests:\t\t%d\n", m.TotalRequests)
	fmt.Printf("Success:\t\t%d\n", m.SuccessCount)
	fmt.Printf("Failed:\t\t\t%d\n", m.ErrorCount)
	fmt.Printf("Test Duration:\t\t%.2fs\n", duration.Seconds())
	fmt.Printf("Requests/sec:\t\t%.2f\n", rps)
	fmt.Printf("Avg Latency:\t\t%s\n", avg)
	fmt.Printf("Min Latency:\t\t%s\n", min)
	fmt.Printf("Max Latency:\t\t%s\n", max)
	fmt.Printf("95th Percentile:\t%s\n", p95)
}

func averageLatency(latencies []time.Duration) time.Duration {
	var sum time.Duration
	for _, l := range latencies {
		sum += l
	}
	return sum / time.Duration(len(latencies))
}
