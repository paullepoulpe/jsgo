// Command parsebench benchmarks the ported front end over a corpus, K times in
// one process. It is the Go counterpart of tools/bench/java/ParseBench.java and
// prints the same columns.
//
//	--mode=parser (default): internal/parser Parser.ParseProgram
//	--mode=runner: internal/parsing ParserRunner_Parse (Parser + IRFactory +
//	               JsDocInfoParser) with the compiler's parsing.Config
//
//	go run ./tools/bench/parsebench [--mode=runner] <iters> <root> @<listfile> <per-file-tsv>
//
// All files are read into memory first. Each iteration parses every file once
// on one goroutine, with the configuration the compiler uses for
// --language_in=ECMASCRIPT_NEXT. "thread_cpu_ms" is not measurable per
// goroutine in Go, so it repeats process CPU time (which, as on the JVM side,
// includes the garbage collector's background workers). Set
// PARSEBENCH_CPUPROFILE=<file> to write a pprof CPU profile.
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type rhinoCountingReporter struct{ errors int }

func (r *rhinoCountingReporter) Warning(message, sourceName string, line, lineOffset int) {}
func (r *rhinoCountingReporter) Error(message, sourceName string, line, lineOffset int) {
	r.errors++
}

type countingReporter struct {
	parser.ErrorReporterBase
	errors int
}

func (r *countingReporter) ReportError(location parser.SourcePosition, message string) {
	r.errors++
}

func (r *countingReporter) ReportWarning(location parser.SourcePosition, message string) {}

func processCPU() time.Duration {
	var ru syscall.Rusage
	_ = syscall.Getrusage(syscall.RUSAGE_SELF, &ru)
	return time.Duration(ru.Utime.Nano() + ru.Stime.Nano())
}

func main() {
	args := os.Args[1:]
	runner := false
	if len(args) > 0 && strings.HasPrefix(args[0], "--mode=") {
		runner = args[0] == "--mode=runner"
		args = args[1:]
	}
	iters, _ := strconv.Atoi(args[0])
	root := args[1]
	listFile := strings.TrimPrefix(args[2], "@")
	perFile := args[3]

	f, err := os.Open(listFile)
	if err != nil {
		panic(err)
	}
	var names []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if l := sc.Text(); l != "" {
			names = append(names, l)
		}
	}
	f.Close()
	contents := make([]string, len(names))
	totalChars := 0
	for i, n := range names {
		b, err := os.ReadFile(filepath.Join(root, n))
		if err != nil {
			panic(err)
		}
		contents[i] = string(b)
		totalChars += jutil.StringLength(contents[i])
	}
	if prof := os.Getenv("PARSEBENCH_CPUPROFILE"); prof != "" {
		pf, err := os.Create(prof)
		if err != nil {
			panic(err)
		}
		_ = pprof.StartCPUProfile(pf)
		defer pprof.StopCPUProfile()
	}
	runnerConfig := parsing.ParserRunner_CreateConfigWithModes(parsing.Config_LanguageMode_ES_NEXT,
		parsing.Config_JsDocParsing_TYPES_ONLY, parsing.Config_RunMode_STOP_AFTER_ERROR, nil, true,
		parsing.Config_StrictMode_STRICT)
	config := parser.NewParser_ConfigWithMode(parser.Parser_Config_Mode_ES8_OR_GREATER, true)

	mode := "parser"
	if runner {
		mode = "runner"
	}
	fmt.Printf("# go mode=%s files=%d utf16_chars=%d GOMAXPROCS=%d GOGC=%s\n", mode, len(names), totalChars, runtime.GOMAXPROCS(0), os.Getenv("GOGC"))
	fmt.Println("iter\twall_ms\tthread_cpu_ms\tprocess_cpu_ms\talloc_mb\tgc_count\tgc_ms")
	last := make([]time.Duration, len(names))
	errors := 0
	var ms runtime.MemStats
	for it := 1; it <= iters; it++ {
		runtime.ReadMemStats(&ms)
		alloc0, gc0, pause0 := ms.TotalAlloc, ms.NumGC, ms.PauseTotalNs
		cpu0 := processCPU()
		wall0 := time.Now()
		errors = 0
		for i := range names {
			t0 := time.Now()
			if runner {
				rr := &rhinoCountingReporter{}
				parsing.ParserRunner_Parse(rhino.NewSimpleSourceFile(names[i], rhino.StaticSourceFile_SourceKind_STRONG),
					contents[i], runnerConfig, rr)
				errors += rr.errors
				last[i] = time.Since(t0)
				continue
			}
			file := parser.NewSourceFile(names[i], contents[i])
			r := &countingReporter{}
			r.InitErrorReporterBase(r)
			parser.NewParser(config, r, file).ParseProgram()
			errors += r.errors
			last[i] = time.Since(t0)
		}
		wall := time.Since(wall0)
		cpu := processCPU() - cpu0
		runtime.ReadMemStats(&ms)
		fmt.Printf("%d\t%.0f\t%.0f\t%.0f\t%.0f\t%d\t%.0f\n", it,
			float64(wall)/1e6, float64(cpu)/1e6, float64(cpu)/1e6,
			float64(ms.TotalAlloc-alloc0)/1048576, ms.NumGC-gc0, float64(ms.PauseTotalNs-pause0)/1e6)
	}
	fmt.Printf("# files_with_errors_reported=%d\n", errors)
	var sb strings.Builder
	for i, n := range names {
		fmt.Fprintf(&sb, "%s\t%d\t%d\n", n, jutil.StringLength(contents[i]), last[i].Nanoseconds())
	}
	if err := os.WriteFile(perFile, []byte(sb.String()), 0o644); err != nil {
		panic(err)
	}
}
