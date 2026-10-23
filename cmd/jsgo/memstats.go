package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/metrics"
	"runtime/pprof"
	"sync"
	"sync/atomic"
)

// Go-only: JSGO_MEMSTATS=<file> writes allocation and GC totals at exit, as
// JSON, for memory benchmarking (tools/bench/memgate.py). peak_live_bytes is
// the largest heap that survived a GC cycle (sampled once per cycle), which
// is a less noisy measure of the working set than the process's peak RSS.
// JSGO_PEAK_HEAP_PROFILE=<file> also writes a heap profile whenever the live
// heap reaches a new peak (5% above the last one), so the file ends up
// describing the heap at its high-water mark.

var peakLive atomic.Uint64

// trackPeakLive samples the live heap after every GC cycle with a finalizer
// that re-arms itself.
func trackPeakLive(profPath string) {
	var lastProf uint64
	pending := false
	type sentinel struct{ _ [16]byte }
	sample := []metrics.Sample{{Name: "/gc/heap/live:bytes"}}
	var arm func()
	arm = func() {
		s := &sentinel{}
		runtime.SetFinalizer(s, func(*sentinel) {
			metrics.Read(sample)
			// A heap profile describes the cycle before the last one, so
			// the profile of a new peak is written one cycle later.
			if pending && profPath != "" {
				pending = false
				if f, err := os.Create(profPath); err == nil {
					_ = pprof.WriteHeapProfile(f)
					_ = f.Close()
				}
			}
			if v := sample[0].Value.Uint64(); v > peakLive.Load() {
				peakLive.Store(v)
				if v > lastProf+lastProf/20 {
					lastProf = v
					pending = true
				}
			}
			arm()
		})
	}
	arm()
}

func memStatsWriter(path string) func() {
	prof := os.Getenv("JSGO_PEAK_HEAP_PROFILE")
	if prof != "" {
		runtime.MemProfileRate = 64 << 10
	}
	trackPeakLive(prof)
	var once sync.Once
	return func() {
		once.Do(func() {
			names := []string{
				"/gc/heap/allocs:bytes",
				"/gc/heap/allocs:objects",
				"/gc/cycles/total:gc-cycles",
				"/cpu/classes/gc/total:cpu-seconds",
				"/memory/classes/total:bytes",
			}
			samples := make([]metrics.Sample, len(names))
			for i, n := range names {
				samples[i].Name = n
			}
			metrics.Read(samples)
			val := func(i int) float64 {
				switch samples[i].Value.Kind() {
				case metrics.KindUint64:
					return float64(samples[i].Value.Uint64())
				case metrics.KindFloat64:
					return samples[i].Value.Float64()
				}
				return 0
			}
			f, err := os.Create(path)
			if err != nil {
				return
			}
			defer f.Close()
			fmt.Fprintf(f, "{\"alloc_bytes\": %.0f, \"alloc_objects\": %.0f, \"gc_cycles\": %.0f, "+
				"\"gc_cpu_s\": %.3f, \"runtime_mapped_bytes\": %.0f, \"peak_live_bytes\": %d}\n",
				val(0), val(1), val(2), val(3), val(4), peakLive.Load())
		})
	}
}
