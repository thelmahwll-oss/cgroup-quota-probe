# cgroup-quota-probe

Lightweight cgroup v1/v2 CPU quota detection library for containerized Go services.

Reads kernel-exported cgroup parameters (`cpu.max`, `cpu.cfs_quota_us`,
`cpu.cfs_period_us`) and exposes them through a simple API that higher-level
libraries like [automaxprocs](https://pkg.go.dev/go.uber.org/automaxprocs) can
consume to set `GOMAXPROCS` correctly.

## Features

- Supports both cgroup v1 and v2 hierarchies
- Zero dependencies
- Automatic version detection
- Lightweight — no CGO required

## Usage

```go
import probe "github.com/thelmahwll-oss/cgroup-quota-probe"

q, err := probe.Detect()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Effective CPUs: %.1f (cgroup v%d)\n", q.EffectiveCPUs(), q.Version)
```

## Daemon

A standalone monitoring daemon is included in `cmd/probed`:

```bash
go build -o probed ./cmd/probed
./probed -interval 10s
```

## Test Fixtures

The `testdata/fixtures/` directory contains serialized cgroup snapshots and
baseline data captured from real container environments (Docker, containerd,
CRI-O). These are used by CI and benchmark scripts to validate quota detection
across both cgroup versions and various runtime configurations.

## License

MIT
