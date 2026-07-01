package probe

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetectV2Fixture(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "fixtures", "cpu_v2_max.dat"))
	if err != nil {
		t.Skip("fixture not found")
	}
	parts := strings.Fields(strings.TrimSpace(string(data)))
	if len(parts) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(parts))
	}
	if parts[0] == "max" {
		t.Log("unlimited quota")
	}
}

func TestDetectV1Fixture(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "fixtures", "cpu_v1_quota.dat"))
	if err != nil {
		t.Skip("fixture not found")
	}
	val := strings.TrimSpace(string(data))
	if val == "-1" {
		t.Log("unlimited quota (v1)")
	}
}

func TestEffectiveCPUs(t *testing.T) {
	cases := []struct {
		name   string
		quota  int64
		period int64
		want   float64
	}{
		{"two_cpus", 200000, 100000, 2.0},
		{"half_cpu", 50000, 100000, 0.5},
		{"unlimited", -1, 100000, 0},
		{"zero_period", 100000, 0, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			q := &Quota{QuotaUS: tc.quota, PeriodUS: tc.period}
			got := q.EffectiveCPUs()
			if got != tc.want {
				t.Errorf("got %.2f, want %.2f", got, tc.want)
			}
		})
	}
}
