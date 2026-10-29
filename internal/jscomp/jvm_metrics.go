// Port of com/google/javascript/jscomp/JvmMetrics.java @ closure-compiler v20261021

/*
 * Copyright 2012 The Closure Compiler Authors.
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
	"io"
	"runtime"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

const (
	jvmMetrics_TABULAR_COLON_POS       = 40
	jvmMetrics_ONE_KILO_BYTE     int64 = 1 << 10
	jvmMetrics_ONE_MEGA_BYTE     int64 = 1 << 20
	jvmMetrics_ONE_GIGA_BYTE     int64 = 1 << 30
)

// Java: JvmMetrics.maybeWriteJvmMetrics(PrintStream, String)
func JvmMetrics_MaybeWriteJvmMetrics(out io.Writer, options *string) {
	if options == nil {
		return
	}

	verboseMode := false
	prettyMode := false

	st := jvmMetrics_tokenize(*options, ":")

	// options are grouped in order 'detail:format:types'
	if len(st) > 0 {
		mode := st[0]
		st = st[1:]
		if strings.EqualFold(mode, "verbose") {
			verboseMode = true
		}
	}

	if len(st) > 0 {
		mode := st[0]
		st = st[1:]
		if strings.EqualFold(mode, "pretty") {
			prettyMode = true
		}
	}

	var ms runtime.MemStats
	var msLoaded bool
	getMemStats := func() *runtime.MemStats {
		if !msLoaded {
			runtime.ReadMemStats(&ms)
			msLoaded = true
		}
		return &ms
	}

	if len(st) > 0 {
		for len(st) > 0 {
			types := st[0]
			st = st[1:]
			for _, type_ := range jvmMetrics_tokenize(types, ",") {
				jvmMetrics_writeMetrics(out, type_, verboseMode, prettyMode, getMemStats)
			}
		}
	} else {
		jvmMetrics_writeMetrics(out, "all", verboseMode, prettyMode, getMemStats)
	}
}

func jvmMetrics_tokenize(s string, delims string) []string {
	return strings.FieldsFunc(s, func(r rune) bool { return strings.ContainsRune(delims, r) })
}

func jvmMetrics_println(out io.Writer, s string) {
	io.WriteString(out, s+"\n")
}

func jvmMetrics_writeMetrics(out io.Writer, type_ string, verbose, pretty bool, getMemStats func() *runtime.MemStats) {
	if type_ == "gc" || strings.EqualFold(type_, "all") {
		jvmMetrics_writeGarbageCollectionStats(out, verbose, pretty, getMemStats())
	}
	if type_ == "mem" || strings.EqualFold(type_, "all") {
		jvmMetrics_writeMemoryMetrics(out, verbose, pretty, getMemStats())
	}
	if type_ == "jit" || strings.EqualFold(type_, "all") {
		jvmMetrics_writeJitMetrics(out, verbose, pretty)
	}
}

func jvmMetrics_writeJitMetrics(out io.Writer, verbose, pretty bool) {
	name := "total"
	if verbose {
		name = "Go"
	}
	if pretty {
		jvmMetrics_println(out, "\nJIT Stats")
		io.WriteString(out, jutil.StringFormat("\t%s jit time: %d ms%n", name, int64(0)))
	} else {
		jvmMetrics_println(out, jvmMetrics_normalizeTabularColonPos(jutil.StringFormat("%s-jit-time-ms : %d",
			jvmMetrics_normalizeName(name), int64(0))))
	}
}

func jvmMetrics_writeOverallMemoryUsage(out io.Writer, used, maxBytes int64, prefix string, pretty bool) {
	if pretty {
		io.WriteString(out, jutil.StringFormat("\t%s\n", prefix))
		io.WriteString(out, jutil.StringFormat("\t\tavailable         : %s\n", jvmMetrics_formatBytes(maxBytes)))
		io.WriteString(out, jutil.StringFormat("\t\tcurrent           : %s\n", jvmMetrics_formatBytes(used)))
	} else {
		prefix = jvmMetrics_normalizeName(prefix)
		jvmMetrics_println(out, jvmMetrics_normalizeTabularColonPos(
			jutil.StringFormat("%s-available-bytes : %d", prefix, maxBytes)))
		jvmMetrics_println(out, jvmMetrics_normalizeTabularColonPos(
			jutil.StringFormat("%s-current-bytes : %d", prefix, used)))
	}
}

func jvmMetrics_writePoolMemoryUsage(out io.Writer, used, peakUsed, maxBytes int64, prefix *string, pretty bool) {
	if pretty {
		io.WriteString(out, jutil.StringFormat("\t\tavailable         : %s\n", jvmMetrics_formatBytes(maxBytes)))
		io.WriteString(out, jutil.StringFormat("\t\tpeak              : %s\n", jvmMetrics_formatBytes(peakUsed)))
		io.WriteString(out, jutil.StringFormat("\t\tcurrent           : %s\n", jvmMetrics_formatBytes(used)))
	} else {
		var p any = prefix
		jvmMetrics_println(out, jvmMetrics_normalizeTabularColonPos(
			jutil.StringFormat("%s-available-bytes : %d", p, maxBytes)))
		jvmMetrics_println(out, jvmMetrics_normalizeTabularColonPos(
			jutil.StringFormat("%s-peak-bytes : %d", p, peakUsed)))
		jvmMetrics_println(out, jvmMetrics_normalizeTabularColonPos(
			jutil.StringFormat("%s-current-bytes : %d", p, used)))
	}
}

func jvmMetrics_writeMemoryMetrics(out io.Writer, verbose, pretty bool, ms *runtime.MemStats) {
	if pretty {
		jvmMetrics_println(out, "\nMemory usage")
	}

	heapUsed := int64(ms.HeapAlloc)
	heapMax := int64(ms.HeapSys)

	if verbose {
		jvmMetrics_writeOverallMemoryUsage(out, heapUsed, heapMax, "Heap", pretty)
		nonHeap := int64(ms.Sys - ms.HeapSys)
		jvmMetrics_writeOverallMemoryUsage(out, int64(ms.StackInuse), nonHeap, "Non-heap", pretty)

		const poolName = "Go heap"
		if pretty {
			jvmMetrics_println(out, "\tPool "+poolName)
			jvmMetrics_writePoolMemoryUsage(out, heapUsed, heapUsed, heapMax, nil, true)
		} else {
			prefix := "mem-pool-" + jvmMetrics_normalizeName(poolName)
			jvmMetrics_writePoolMemoryUsage(out, heapUsed, heapUsed, heapMax, &prefix, false)
		}
	} else {
		if pretty {
			io.WriteString(out, jutil.StringFormat("\tAggregate of %d memory pools\n", 1))
			jvmMetrics_writePoolMemoryUsage(out, heapUsed, heapUsed, heapMax, nil, true)
		} else {
			prefix := "mem"
			jvmMetrics_writePoolMemoryUsage(out, heapUsed, heapUsed, heapMax, &prefix, false)
		}
	}
}

func jvmMetrics_writeGarbageCollectionStats(out io.Writer, verbose, pretty bool, ms *runtime.MemStats) {
	collectionCount := int64(ms.NumGC)
	collectionTime := int64(ms.PauseTotalNs / 1e6)

	if verbose {
		const gcName = "Go GC"
		if pretty {
			jvmMetrics_println(out, "\nGarbage collection stats")
			jvmMetrics_println(out, "\tCollector "+gcName)
			io.WriteString(out, jutil.StringFormat("\t\tcollection count   : %d\n", collectionCount))
			io.WriteString(out, jutil.StringFormat("\t\tcollection time    : %d ms\n", collectionTime))
		} else {
			name := jvmMetrics_normalizeName(gcName)
			jvmMetrics_println(out, jvmMetrics_normalizeTabularColonPos(
				jutil.StringFormat("gc-"+name+"-collection-count : %d", collectionCount)))
			jvmMetrics_println(out, jvmMetrics_normalizeTabularColonPos(
				jutil.StringFormat("gc-"+name+"-collection-time-ms : %d", collectionTime)))
		}
	} else {
		if pretty {
			jvmMetrics_println(out, "\nGarbage collection stats")
			io.WriteString(out, jutil.StringFormat("\tAggregate of %d collectors\n", 1))
			io.WriteString(out, jutil.StringFormat("\t\tcollection count   : %d\n", collectionCount))
			io.WriteString(out, jutil.StringFormat("\t\tcollection time    : %d ms\n", collectionTime))
		} else {
			name := jvmMetrics_normalizeName("aggregate")
			jvmMetrics_println(out, jvmMetrics_normalizeTabularColonPos(
				jutil.StringFormat("gc-%s-collection-count : %d", name, collectionCount)))
			jvmMetrics_println(out, jvmMetrics_normalizeTabularColonPos(
				jutil.StringFormat("gc-%s-collection-time-ms : %d", name, collectionTime)))
		}
	}
}

func jvmMetrics_normalizeName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "_"))
}

func jvmMetrics_normalizeTabularColonPos(s string) string {
	index := jutil.StringIndexOf(s, ":")
	if index < 0 {
		panic(jutil.NewStringIndexOutOfBoundsException(index, jutil.StringLength(s)))
	}
	if index >= jvmMetrics_TABULAR_COLON_POS {
		return s
	}
	return jutil.StringSubstring(s, 0, index) + strings.Repeat(" ", jvmMetrics_TABULAR_COLON_POS-index) + jutil.StringSubstringFrom(s, index)
}

func jvmMetrics_formatBytes(numBytes int64) string {
	if numBytes < jvmMetrics_ONE_KILO_BYTE {
		return jutil.StringFormat("%d B", numBytes)
	} else if numBytes < jvmMetrics_ONE_MEGA_BYTE {
		return jutil.StringFormat("%d KB", numBytes/jvmMetrics_ONE_KILO_BYTE)
	} else if numBytes < jvmMetrics_ONE_GIGA_BYTE {
		return jutil.StringFormat("%d MB", numBytes/jvmMetrics_ONE_MEGA_BYTE)
	}
	return jutil.StringFormat("%d GB", numBytes/jvmMetrics_ONE_GIGA_BYTE)
}
