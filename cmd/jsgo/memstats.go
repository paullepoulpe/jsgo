package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/metrics"
	"runtime/pprof"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
)

// Go-only: JSGO_MEMSTATS=<file> writes allocation and GC totals at exit, as
// JSON, for memory benchmarking (tools/bench/memgate.py). peak_live_bytes is
// the largest heap that survived a GC cycle (sampled once per cycle), which
// is a less noisy measure of the working set than the process's peak RSS.
// JSGO_PEAK_HEAP_PROFILE=<file> also writes a heap profile whenever the live
// heap reaches a new peak (5% above the last one), so the file ends up
// describing the heap at its high-water mark. peak_scan_bytes and
// peak_objects are the scannable (pointerful) heap bytes and the number of
// live objects at that same cycle: the GC's mark work per cycle grows with
// the scannable bytes, not with pointer-free ones.

var peakLive, peakScan, peakObjects, liveAfterColors atomic.Uint64

// recordLiveAfterColors forces a collection and records the live heap, for
// live_after_colors_bytes: what survives the types-to-colors conversion.
func recordLiveAfterColors() {
	runtime.GC()
	sample := []metrics.Sample{{Name: "/gc/heap/live:bytes"}}
	metrics.Read(sample)
	liveAfterColors.Store(sample[0].Value.Uint64())
}

// memClassNames are the /memory/classes metrics sampled after every GC
// cycle; peakClasses holds them at the cycle whose total was the largest
// (the runtime's view of the peak RSS: live and not-yet-swept objects,
// unused span slack, free and released spans, stacks and metadata).
var memClassNames = []string{
	"/memory/classes/total:bytes",
	"/memory/classes/heap/objects:bytes",
	"/memory/classes/heap/unused:bytes",
	"/memory/classes/heap/free:bytes",
	"/memory/classes/heap/released:bytes",
	"/memory/classes/heap/stacks:bytes",
	"/memory/classes/metadata/mspan/inuse:bytes",
	"/memory/classes/metadata/other:bytes",
	"/memory/classes/other:bytes",
}

var (
	peakClassesMu sync.Mutex
	peakClasses   []uint64
)

// trackPeakLive samples the live heap after every GC cycle with a finalizer
// that re-arms itself.
func trackPeakLive(profPath string) {
	var lastProf uint64
	pending := false
	type sentinel struct{ _ [16]byte }
	sample := []metrics.Sample{
		{Name: "/gc/heap/live:bytes"},
		{Name: "/gc/scan/heap:bytes"},
		{Name: "/gc/heap/objects:objects"},
	}
	for _, n := range memClassNames {
		sample = append(sample, metrics.Sample{Name: n})
	}
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
			peakClassesMu.Lock()
			if total := sample[3].Value.Uint64(); len(peakClasses) == 0 || total > peakClasses[0] {
				peakClasses = peakClasses[:0]
				for _, s := range sample[3:] {
					peakClasses = append(peakClasses, s.Value.Uint64())
				}
			}
			peakClassesMu.Unlock()
			if v := sample[0].Value.Uint64(); v > peakLive.Load() {
				peakLive.Store(v)
				peakScan.Store(sample[1].Value.Uint64())
				peakObjects.Store(sample[2].Value.Uint64())
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
	jscomp.AfterTypesToColorsHook = recordLiveAfterColors
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
				"\"gc_cpu_s\": %.3f, \"runtime_mapped_bytes\": %.0f, \"peak_live_bytes\": %d, "+
				"\"peak_scan_bytes\": %d, \"peak_objects\": %d, \"live_after_colors_bytes\": %d",
				val(0), val(1), val(2), val(3), val(4), peakLive.Load(), peakScan.Load(), peakObjects.Load(), liveAfterColors.Load())
			peakClassesMu.Lock()
			for i, v := range peakClasses {
				name := strings.TrimSuffix(strings.TrimPrefix(memClassNames[i], "/memory/classes/"), ":bytes")
				fmt.Fprintf(f, ", \"peak_mapped/%s\": %d", name, v)
			}
			peakClassesMu.Unlock()
			fmt.Fprintf(f, "}\n")
		})
	}
}
