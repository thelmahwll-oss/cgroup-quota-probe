package cgroup

import (
	"os"
	"strconv"
	"strings"
)

// CPUStats holds parsed CPU throttling statistics from cgroup v2.
type CPUStats struct {
	UsageUSec     int64
	UserUSec      int64
	SystemUSec    int64
	NrPeriods     int64
	NrThrottled   int64
	ThrottledUSec int64
}

// ReadCPUStats parses the cpu.stat file for cgroup v2.
func ReadCPUStats(path string) (*CPUStats, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	stats := &CPUStats{}
	for _, line := range strings.Split(string(data), "\n") {
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}
		val, _ := strconv.ParseInt(parts[1], 10, 64)
		switch parts[0] {
		case "usage_usec":
			stats.UsageUSec = val
		case "user_usec":
			stats.UserUSec = val
		case "system_usec":
			stats.SystemUSec = val
		case "nr_periods":
			stats.NrPeriods = val
		case "nr_throttled":
			stats.NrThrottled = val
		case "throttled_usec":
			stats.ThrottledUSec = val
		}
	}
	return stats, nil
}

// ThrottleRatio returns the ratio of throttled periods to total periods.
func (s *CPUStats) ThrottleRatio() float64 {
	if s.NrPeriods == 0 {
		return 0
	}
	return float64(s.NrThrottled) / float64(s.NrPeriods)
}
