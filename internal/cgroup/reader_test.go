package cgroup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseMountFixture(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "fixtures", "cgroup_mount.dat"))
	if err != nil {
		t.Skip("fixture not found")
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) == 0 {
		t.Fatal("empty fixture")
	}
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 3)
		if len(parts) != 3 {
			t.Errorf("malformed line: %q", line)
		}
	}
}

func TestCgroupBaselineParsing(t *testing.T) {
	entries, err := os.ReadDir(filepath.Join("..", "..", "testdata", "fixtures"))
	if err != nil {
		t.Skip("fixtures dir not found")
	}
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".dat") {
			continue
		}
		t.Run(e.Name(), func(t *testing.T) {
			info, err := e.Info()
			if err != nil {
				t.Fatal(err)
			}
			if info.Size() == 0 {
				t.Error("empty fixture file")
			}
		})
	}
}
