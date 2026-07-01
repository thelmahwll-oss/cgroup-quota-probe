// Package probe provides cgroup v1/v2 CPU quota detection for containerized
// Go services. It reads kernel-exported cgroup parameters and exposes them
// through a simple API that higher-level libraries (e.g. automaxprocs) can
// consume to set GOMAXPROCS correctly.
package probe

import (
	"os"
	"strconv"
	"strings"
)

// Quota holds the detected CPU quota parameters.
type Quota struct {
	QuotaUS  int64
	PeriodUS int64
	Max      string
	Version  int
}

// Detect reads the cgroup CPU quota from the filesystem.
// It returns the detected quota and the cgroup version (1 or 2).
func Detect() (*Quota, error) {
	if q, err := detectV2(); err == nil {
		return q, nil
	}
	return detectV1()
}

func detectV2() (*Quota, error) {
	data, err := os.ReadFile("/sys/fs/cgroup/cpu.max")
	if err != nil {
		return nil, err
	}
	parts := strings.Fields(strings.TrimSpace(string(data)))
	q := &Quota{Version: 2, Max: parts[0]}
	if len(parts) >= 2 {
		q.PeriodUS, _ = strconv.ParseInt(parts[1], 10, 64)
	}
	if parts[0] != "max" {
		q.QuotaUS, _ = strconv.ParseInt(parts[0], 10, 64)
	}
	return q, nil
}

func detectV1() (*Quota, error) {
	q := &Quota{Version: 1}
	quota, err := os.ReadFile("/sys/fs/cgroup/cpu/cpu.cfs_quota_us")
	if err != nil {
		return nil, err
	}
	q.QuotaUS, _ = strconv.ParseInt(strings.TrimSpace(string(quota)), 10, 64)

	period, err := os.ReadFile("/sys/fs/cgroup/cpu/cpu.cfs_period_us")
	if err != nil {
		return nil, err
	}
	q.PeriodUS, _ = strconv.ParseInt(strings.TrimSpace(string(period)), 10, 64)
	return q, nil
}

// EffectiveCPUs returns the number of effective CPUs available to the process
// based on the detected quota and period.
func (q *Quota) EffectiveCPUs() float64 {
	if q == nil || q.PeriodUS <= 0 || q.QuotaUS <= 0 {
		return 0
	}
	return float64(q.QuotaUS) / float64(q.PeriodUS)
}
