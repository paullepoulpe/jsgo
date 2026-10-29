// Port of com/google/javascript/jscomp/Platform.java @ closure-compiler v20261021

/*
 * Copyright 2020 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"os"
	"runtime"
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
	return parseDiagPattern(message).format(arguments...)
}
