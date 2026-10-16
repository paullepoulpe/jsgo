package jutil

import (
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
