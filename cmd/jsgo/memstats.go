package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/metrics"
	"runtime/pprof"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

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

// Go-only measurement hook: JSGO_PASS_LOG=<file> writes one TSV row per
// pass run and per Compiler tracer section, with wall time, process CPU
// time (user+system, all threads, GC included) and, for passes, the change
// reports and changed scopes they made. Rows are written in end order;
// "depth" is the nesting level. With JSGO_CPU_PROFILE also set, CPU samples
// carry a pprof label pass=<name> (inherited by goroutines started inside
// the pass).
func installPassLog(path string, labels bool) func() {
	start := time.Now()
	type frame struct {
		ev   jscomp.DebugPassEvent
		wall time.Time
		cpu  time.Duration
	}
	var stack []frame
	var rows []string
	cpuNow := func() time.Duration {
		var ru syscall.Rusage
		_ = syscall.Getrusage(syscall.RUSAGE_SELF, &ru)
		return time.Duration(ru.Utime.Nano() + ru.Stime.Nano())
	}
	setLabel := func(name string) {
		if !labels {
			return
		}
		if name == "" {
			pprof.SetGoroutineLabels(context.Background())
			return
		}
		pprof.SetGoroutineLabels(pprof.WithLabels(context.Background(), pprof.Labels("pass", name)))
	}
	jscomp.DebugPassLogHook = func(ev jscomp.DebugPassEvent) {
		if ev.Start {
			stack = append(stack, frame{ev: ev, wall: time.Now(), cpu: cpuNow()})
			setLabel(ev.Name)
			return
		}
		// Pop the matching frame (tracer sections and passes nest).
		i := len(stack) - 1
		for i >= 0 && (stack[i].ev.Name != ev.Name || stack[i].ev.Tracer != ev.Tracer) {
			i--
		}
		if i < 0 {
			return
		}
		f := stack[i]
		stack = stack[:i]
		kind := "pass"
		if ev.Tracer {
			kind = "tracer"
		}
		rows = append(rows, fmt.Sprintf("%s\t%s\t%d\t%d\t%.3f\t%.3f\t%d\t%d\t%d\t%.3f",
			kind, ev.Name, f.ev.Loop, f.ev.Iter,
			float64(time.Since(f.wall).Microseconds())/1000,
			float64((cpuNow()-f.cpu).Microseconds())/1000,
			ev.Reports, ev.ChangedScopes, i,
			float64(f.wall.Sub(start).Microseconds())/1000))
		if len(stack) > 0 {
			setLabel(stack[len(stack)-1].ev.Name)
		} else {
			setLabel("")
		}
	}
	return func() {
		f, err := os.Create(path)
		if err != nil {
			return
		}
		w := bufio.NewWriter(f)
		fmt.Fprintln(w, "kind\tname\tloop\titer\twall_ms\tcpu_ms\treports\tscopes\tdepth\tstart_ms")
		for _, r := range rows {
			fmt.Fprintln(w, r)
		}
		fmt.Fprintf(w, "total\tprocess\t0\t0\t%.3f\t%.3f\t0\t0\t0\t0\n",
			float64(time.Since(start).Microseconds())/1000, float64(cpuNow().Microseconds())/1000)
		_ = w.Flush()
		_ = f.Close()
	}
}
