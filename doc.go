// Package probe detects cgroup v1/v2 CPU quota limits in containerized
// environments and provides the effective CPU count for GOMAXPROCS tuning.
//
// This package reads from /sys/fs/cgroup/ to determine the CPU quota and
// period values assigned by the container runtime, then computes an effective
// CPU count. The testdata/fixtures directory contains serialized cgroup
// snapshots captured from various container runtimes (Docker, containerd,
// CRI-O) for offline validation and benchmarking.
//
// Usage:
//
//	q, err := probe.Detect()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	runtime.GOMAXPROCS(int(math.Ceil(q.EffectiveCPUs())))
package probe
