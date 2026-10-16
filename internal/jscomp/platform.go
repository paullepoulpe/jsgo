// Port of com/google/javascript/jscomp/Platform.java @ closure-compiler v20260415

package jscomp

import (
	"os"
	"runtime"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Platform is a utility to abstract platform specific logic so it could be
// super-sourced for Web.
type Platform struct{}

// Java: Platform.getFileSeperator()
//
// Go-only: language-adapter — File.separator is os.PathSeparator.
func Platform_GetFileSeperator() string {
	return string(os.PathSeparator)
}

// Java: Platform.isThreadInterrupted()
//
// Go-only: language-adapter — goroutines cannot be interrupted, so the
// interrupted flag is never set.
func Platform_IsThreadInterrupted() bool {
	// Note that this clears interruption. Thread.isInterruped should be
	// preferred instead.
	return false
}

// Java: Platform.freeMemory()
//
// Go-only: language-adapter — Runtime.freeMemory() is the heap memory that
// is reserved but unused (HeapIdle).
func Platform_FreeMemory() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int64(m.HeapIdle)
}

// Java: Platform.totalMemory()
//
// Go-only: language-adapter — Runtime.totalMemory() is the heap memory
// reserved from the OS (HeapSys).
func Platform_TotalMemory() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int64(m.HeapSys)
}

// Java: Platform.formatMessage(String, String...)
func Platform_FormatMessage(message string, arguments ...string) string {
	// Note that MessageFormat is removing single quotes and in many cases
	// intended ones. Consider moving to a simpler formatting version like the
	// Web one.
	args := make([]any, len(arguments))
	for i, a := range arguments {
		args[i] = a
	}
	return jutil.MessageFormat_Format(message, args)
}

// Java: Platform()
func newPlatform() *Platform { return &Platform{} }
