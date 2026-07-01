#!/bin/bash
# benchmark.sh — Run quota probe benchmarks across cgroup versions
set -euo pipefail

ITERATIONS=${1:-1000}
echo "Running $ITERATIONS iterations..."

for v in v1 v2; do
  echo "=== cgroup $v ==="
  time for i in $(seq 1 $ITERATIONS); do
    cat "testdata/fixtures/cpu_${v}_quota.dat" > /dev/null 2>&1 || true
  done
done

echo "Done."
