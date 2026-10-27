package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
)

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
