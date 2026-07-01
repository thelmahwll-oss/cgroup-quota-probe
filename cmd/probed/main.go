// Command probed is a lightweight daemon that continuously monitors cgroup
// CPU quota parameters and exposes them via a local Unix socket for
// consumption by GOMAXPROCS tuning libraries.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	probe "github.com/thelmahwll-oss/cgroup-quota-probe"
	"github.com/thelmahwll-oss/cgroup-quota-probe/internal/cgroup"
)

var (
	interval = flag.Duration("interval", 5*time.Second, "polling interval")
	logFile  = flag.String("log", "", "optional log file path (default: stderr)")
	version  = flag.Bool("version", false, "print version and exit")
)

const buildVersion = "0.3.1"

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("probed %s\n", buildVersion)
		os.Exit(0)
	}

	if *logFile != "" {
		f, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("open log: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	log.Printf("probed %s starting (interval=%s)", buildVersion, *interval)

	if cgroup.IsCgroup2() {
		log.Println("detected cgroup v2 unified hierarchy")
	} else {
		log.Println("detected cgroup v1 legacy hierarchy")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down")
			return
		case <-ticker.C:
			q, err := probe.Detect()
			if err != nil {
				log.Printf("probe: %v", err)
				continue
			}
			log.Printf("cgroup v%d: quota=%d period=%d effective_cpus=%.2f",
				q.Version, q.QuotaUS, q.PeriodUS, q.EffectiveCPUs())
		}
	}
}
