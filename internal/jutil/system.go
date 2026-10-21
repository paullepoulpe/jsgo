package jutil

import (
	"io"
	"os"
	"sync"
)

// Go-only: language-adapter — java.lang.System's standard streams
// (System.out, System.err, System.setOut, System.setErr), as PrintStreams
// over the process's stdout and stderr.

var (
	system_mu  sync.RWMutex
	system_out = NewPrintStream(os.Stdout)
	system_err = NewPrintStream(os.Stderr)
)

// System_Out is System.out.
func System_Out() *PrintStream {
	system_mu.RLock()
	defer system_mu.RUnlock()
	return system_out
}

// System_Err is System.err.
func System_Err() *PrintStream {
	system_mu.RLock()
	defer system_mu.RUnlock()
	return system_err
}

// System_SetOut is System.setOut(PrintStream).
func System_SetOut(out *PrintStream) {
	system_mu.Lock()
	defer system_mu.Unlock()
	system_out = out
}

// System_SetErr is System.setErr(PrintStream).
func System_SetErr(err *PrintStream) {
	system_mu.Lock()
	defer system_mu.Unlock()
	system_err = err
}

var system_atExit []func()

// System_AtExit registers a callback to run before System_Exit calls os.Exit.
func System_AtExit(fn func()) {
	system_mu.Lock()
	defer system_mu.Unlock()
	system_atExit = append(system_atExit, fn)
}

// System_Exit is System.exit(int): terminates the process with the given
// status (the shell sees its low 8 bits).
func System_Exit(status int) {
	system_mu.Lock()
	hooks := system_atExit
	system_atExit = nil
	system_mu.Unlock()
	for i := len(hooks) - 1; i >= 0; i-- {
		hooks[i]()
	}
	os.Exit(status)
}

// System_In is System.in: the process's standard input.
func System_In() io.Reader {
	return os.Stdin
}
