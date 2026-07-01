package cgroup

import (
	"bufio"
	"os"
	"strings"
)

// MountInfo represents a single cgroup mount entry from /proc/self/cgroup.
type MountInfo struct {
	ID         string
	Controller string
	Path       string
}

// ReadMounts parses /proc/self/cgroup and returns all mount entries.
func ReadMounts() ([]MountInfo, error) {
	f, err := os.Open("/proc/self/cgroup")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var mounts []MountInfo
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), ":", 3)
		if len(parts) < 3 {
			continue
		}
		mounts = append(mounts, MountInfo{
			ID:         parts[0],
			Controller: parts[1],
			Path:       parts[2],
		})
	}
	return mounts, scanner.Err()
}

// IsCgroup2 returns true if the system uses the unified cgroup v2 hierarchy.
func IsCgroup2() bool {
	_, err := os.Stat("/sys/fs/cgroup/cgroup.controllers")
	return err == nil
}
