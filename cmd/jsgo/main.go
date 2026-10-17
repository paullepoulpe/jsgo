// Command jsgo is the Go port of Closure Compiler's command-line front end:
// `jsgo <args>` behaves like `java -jar closure-compiler.jar <args>`
// (CommandLineRunner.main).
package main

import (
	"os"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func main() {
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
