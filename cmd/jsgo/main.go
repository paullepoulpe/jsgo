// Command jsgo is the Go port of Closure Compiler's command-line front end:
// `jsgo <args>` behaves like `java -jar closure-compiler.jar <args>`
// (CommandLineRunner.main).
package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func main() {
	if cpuProf := os.Getenv("JSGO_CPU_PROFILE"); cpuProf != "" {
		if f, err := os.Create(cpuProf); err == nil {
			_ = pprof.StartCPUProfile(f)
			stopCPU := func() {
				pprof.StopCPUProfile()
				_ = f.Close()
			}
			jutil.System_AtExit(stopCPU)
			defer stopCPU()
		}
	}
	if logPath := os.Getenv("JSGO_PASS_LOG"); logPath != "" {
		writeLog := installPassLog(logPath, os.Getenv("JSGO_CPU_PROFILE") != "")
		jutil.System_AtExit(writeLog)
		defer writeLog()
	}
	if statsPath := os.Getenv("JSGO_MEMSTATS"); statsPath != "" {
		writeStats := memStatsWriter(statsPath)
		jutil.System_AtExit(writeStats)
		defer writeStats()
	}
	if spec := os.Getenv("JSGO_HEAP_PROFILE_AFTER_PASS"); spec != "" {
		// Go-only: "<passName>:<file>" writes a heap profile (after a forced
		// GC) right after the named pass finishes.
		if pass, file, ok := strings.Cut(spec, ":"); ok {
			jscomp.DebugAfterPassHook = func(c *jscomp.Compiler, passName string) {
				if passName != pass {
					return
				}
				runtime.GC()
				runtime.GC()
				if f, err := os.Create(file); err == nil {
					_ = pprof.WriteHeapProfile(f)
					_ = f.Close()
				}
				if c.IsTypeRegistryCleared() {
					fmt.Fprintf(os.Stderr, "[debug] after %s: type registry released: %v\n", passName, c.TypeRegistryReleased())
				}
			}
		}
	}
	if memProf := os.Getenv("JSGO_MEM_PROFILE"); memProf != "" {
		writeMem := func() {
			runtime.GC()
			if f, err := os.Create(memProf); err == nil {
				_ = pprof.WriteHeapProfile(f)
				_ = f.Close()
			}
		}
		jutil.System_AtExit(writeMem)
		defer writeMem()
	}
	// Go-only: language-adapter — the JVM's startup and shutdown around
	// CommandLineRunner.main:
	//   - System.err encodes text with native.encoding (the locale's
	//     codeset; US-ASCII in the C/POSIX locale). System.out is left as
	//     UTF-8: the runner prints only ASCII text to it and writes
	//     everything else as bytes it has already encoded.
	//   - An exception that escapes main is reported by the default
	//     uncaught-exception handler as `Exception in thread "main" ` and
	//     the stack trace on System.err, and the JVM exits with status 1.
	//   - main returning normally exits with status 0.
	jutil.System_SetErr(jutil.NewPrintStreamWithCharset(os.Stderr, jutil.System_NativeCharset()))
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		t, ok := r.(jutil.Throwable)
		if !ok {
			t = jutil.AsThrowable(r)
		}
		err := jutil.System_Err()
		err.Print("Exception in thread \"main\" ")
		jutil.Throwable_PrintStackTrace(t, err)
		err.Flush()
		os.Exit(1)
	}()
	jscomp.CommandLineRunner_Main(os.Args[1:])
}
