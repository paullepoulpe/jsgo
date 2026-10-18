// Command gowarm repeats the real CommandLineRunner in one process, using
// the same job-file/iteration protocol as tools/bench/java/CliBatch.java.
// This is a benchmark driver; the normal jsgo CLI remains unchanged.
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
)

func main() {
	cache := flag.Bool("extern-cache", false, "enable the opt-in built-in-extern parse cache")
	flag.Parse()
	if flag.NArg() != 4 {
		panic("usage: gowarm [-extern-cache] WARMUP ITERS JOBFILE OUTDIR")
	}
	warmup, err := strconv.Atoi(flag.Arg(0))
	must(err)
	iters, err := strconv.Atoi(flag.Arg(1))
	must(err)
	if warmup < 0 || iters < 1 {
		panic("invalid iteration counts")
	}
	f, err := os.Open(flag.Arg(2))
	must(err)
	var jobs [][]string
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 4096), 8<<20)
	for scanner.Scan() {
		if scanner.Text() != "" {
			jobs = append(jobs, strings.Split(scanner.Text(), "\t"))
		}
	}
	must(scanner.Err())
	must(f.Close())
	outdir := flag.Arg(3)
	must(os.MkdirAll(outdir, 0o755))
	parsing.ParserRunner_EnableBuiltinExternCache(*cache)
	jscomp.PhaseOptimizer_logger.SetLevel(jutil.Level_OFF)
	fmt.Println("iter\tkind\twall_ms\tprocess_cpu_ms\tgc_count\tgc_ms\trss_mb\tstart_epoch_ms\tend_epoch_ms")
	for iteration := 1; iteration <= warmup+iters; iteration++ {
		var before, after runtime.MemStats
		runtime.ReadMemStats(&before)
		cpu0 := processCPU()
		start := time.Now()
		for j, args := range jobs {
			out, stderr, status := run(args)
			if iteration == warmup+iters {
				prefix := filepath.Join(outdir, strconv.Itoa(j))
				must(os.WriteFile(prefix+".stdout", out, 0o644))
				must(os.WriteFile(prefix+".stderr", stderr, 0o644))
				must(os.WriteFile(prefix+".exit", []byte(fmt.Sprintf("%d\n", status)), 0o644))
			}
		}
		end := time.Now()
		cpu := processCPU() - cpu0
		runtime.ReadMemStats(&after)
		kind, index := "measure", iteration-warmup
		if iteration <= warmup {
			kind, index = "warmup", iteration-warmup-1
		}
		// gc_ms is Go STW pause time, not Java's collector time; compare
		// whole-process CPU/wall rather than treating those columns alike.
		fmt.Printf("%d\t%s\t%.1f\t%.1f\t%d\t%.3f\t%.0f\t%d\t%d\n", index, kind,
			float64(end.Sub(start))/1e6, float64(cpu)/1e6, after.NumGC-before.NumGC,
			float64(after.PauseTotalNs-before.PauseTotalNs)/1e6, statusKB("VmRSS")/1024,
			start.UnixMilli(), end.UnixMilli())
	}
	fmt.Printf("# peak_rss_mb=%.0f\n", statusKB("VmHWM")/1024)
}

func run(args []string) (out, stderr []byte, status int) {
	var stdout, errbuf bytes.Buffer
	outStream := jutil.NewPrintStream(&stdout)
	errStream := jutil.NewPrintStreamWithCharset(&errbuf, jutil.System_NativeCharset())
	oldOut, oldErr := jutil.System_Out(), jutil.System_Err()
	jutil.System_SetOut(outStream)
	jutil.System_SetErr(errStream)
	defer func() {
		if r := recover(); r != nil {
			t, ok := r.(jutil.Throwable)
			if !ok {
				t = jutil.AsThrowable(r)
			}
			errStream.Print("Exception in thread \"main\" ")
			jutil.Throwable_PrintStackTrace(t, errStream)
			status = 1
		}
		outStream.Flush()
		errStream.Flush()
		out, stderr = stdout.Bytes(), errbuf.Bytes()
		jutil.System_SetOut(oldOut)
		jutil.System_SetErr(oldErr)
	}()
	runner := jscomp.NewCommandLineRunnerWithStreams(args, bytes.NewReader(nil), outStream, errStream)
	if runner.ShouldRunCompiler() {
		runner.SetExitCodeReceiver(func(code int) { status = int(uint8(code)) })
		runner.Run()
	} else if runner.HasErrors() {
		status = 255
	}
	return
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
func processCPU() int64 {
	var r syscall.Rusage
	must(syscall.Getrusage(syscall.RUSAGE_SELF, &r))
	return r.Utime.Nano() + r.Stime.Nano()
}
func statusKB(field string) float64 {
	b, err := os.ReadFile("/proc/self/status")
	must(err)
	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(line, field+":") {
			v, err := strconv.ParseFloat(strings.Fields(line)[1], 64)
			must(err)
			return v
		}
	}
	return -1
}
