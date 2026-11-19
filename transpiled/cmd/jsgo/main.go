// Command jsgo is the mechanically transpiled Go port of Closure Compiler's
// command-line front end (com.google.javascript.jscomp.CommandLineRunner.main).
package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strings"

	"github.com/paullepoulpe/jsgo/transpiled/gen/jscomp"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

func main() {
	nodestore.EnsureCompilationStore()

	if cpuProf := os.Getenv("JSGO_CPU_PROFILE"); cpuProf != "" {
		if f, err := os.Create(cpuProf); err == nil {
			_ = pprof.StartCPUProfile(f)
			stopCPU := func() {
				pprof.StopCPUProfile()
				_ = f.Close()
			}
			javaio.RegisterAtExit(stopCPU)
			defer stopCPU()
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
		javaio.RegisterAtExit(writeMem)
		defer writeMem()
	}

	defer func() {
		r := recover()
		if r == nil {
			javaio.System_Out.Flush()
			javaio.System_Err.Flush()
			return
		}
		if s, ok := r.(string); ok && strings.HasPrefix(s, "UNIMPLEMENTED_EXT:") {
			fmt.Fprintf(os.Stderr, "%s\n%s", s, debug.Stack())
			os.Exit(2)
		}
		t := javalang.RecoverException(r)
		javaio.System_Err.Print("Exception in thread \"main\" ")
		javaio.System_Err.Println(t.Error())
		if os.Getenv("JSGO_STACKTRACE") != "" {
			fmt.Fprintf(os.Stderr, "%s", debug.Stack())
		}
		javaio.System_Out.Flush()
		javaio.System_Err.Flush()
		javaio.RunAtExit()
		os.Exit(1)
	}()

	jscomp.CommandLineRunner_Main(os.Args[1:])
}
