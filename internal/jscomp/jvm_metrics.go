// Port of com/google/javascript/jscomp/JvmMetrics.java @ closure-compiler v20261028

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
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

const (
	jvmTabularColonPos       = 40
	jvmOneKiloByte     int64 = 1 << 10
	jvmOneMegaByte     int64 = 1 << 20
	jvmOneGigaByte     int64 = 1 << 30
)

// Java: JvmMetrics.maybeWriteJvmMetrics(PrintStream, String)
func JvmMetrics_MaybeWriteJvmMetrics(out io.Writer, options base.OptString) {
	opts, ok := options.Get()
	if !ok {
		return
	}
	JvmMetrics_WriteJvmMetrics(out, opts)
}

// JvmMetrics_MaybeWriteJvmMetricsOpt writes JVM metrics when options is present.
func JvmMetrics_MaybeWriteJvmMetricsOpt(out io.Writer, options base.OptString) {
	JvmMetrics_MaybeWriteJvmMetrics(out, options)
}

// JvmMetrics_WriteJvmMetrics writes JVM metrics for non-null options.
func JvmMetrics_WriteJvmMetrics(out io.Writer, options string) {
	verboseMode := false
	prettyMode := false

	st := jvmTokenize(options, ":")

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
		for _, types := range st {
			for _, kind := range jvmTokenize(types, ",") {
				jvmWriteMetrics(out, kind, verboseMode, prettyMode, getMemStats)
			}
		}
	} else {
		jvmWriteMetrics(out, "all", verboseMode, prettyMode, getMemStats)
	}
}

func jvmTokenize(s string, delims string) []string {
	return strings.FieldsFunc(s, func(r rune) bool { return strings.ContainsRune(delims, r) })
}

func jvmPrintln(out io.Writer, s string) {
	_, _ = io.WriteString(out, s+"\n")
}

func jvmWriteTabularMetric(out io.Writer, key string, val int64) {
	jvmPrintln(out, jvmNormalizeTabularColonPos(key+" : "+strconv.FormatInt(val, 10)))
}

func jvmWriteMetrics(out io.Writer, kind string, verbose, pretty bool, getMemStats func() *runtime.MemStats) {
	if kind == "gc" || strings.EqualFold(kind, "all") {
		jvmWriteGarbageCollectionStats(out, verbose, pretty, getMemStats())
	}
	if kind == "mem" || strings.EqualFold(kind, "all") {
		jvmWriteMemoryMetrics(out, verbose, pretty, getMemStats())
	}
	if kind == "jit" || strings.EqualFold(kind, "all") {
		jvmWriteJitMetrics(out, verbose, pretty)
	}
}

func jvmWriteJitMetrics(out io.Writer, verbose, pretty bool) {
	name := "total"
	if verbose {
		name = "Go"
	}
	if pretty {
		jvmPrintln(out, "\nJIT Stats")
		_, _ = fmt.Fprintf(out, "\t%s jit time: %d ms\n", name, int64(0))
	} else {
		jvmWriteTabularMetric(out, jvmNormalizeName(name)+"-jit-time-ms", 0)
	}
}

func jvmWriteOverallMemoryUsage(out io.Writer, used, maxBytes int64, prefix string, pretty bool) {
	if pretty {
		_, _ = fmt.Fprintf(
			out,
			"\t%s\n\t\tavailable         : %s\n\t\tcurrent           : %s\n",
			prefix,
			jvmFormatBytes(maxBytes),
			jvmFormatBytes(used))
	} else {
		prefix = jvmNormalizeName(prefix)
		jvmWriteTabularMetric(out, prefix+"-available-bytes", maxBytes)
		jvmWriteTabularMetric(out, prefix+"-current-bytes", used)
	}
}

func jvmWritePoolMemoryUsage(out io.Writer, used, peakUsed, maxBytes int64, prefix string, pretty bool) {
	if pretty {
		_, _ = fmt.Fprintf(
			out,
			"\t\tavailable         : %s\n\t\tpeak              : %s\n\t\tcurrent           : %s\n",
			jvmFormatBytes(maxBytes),
			jvmFormatBytes(peakUsed),
			jvmFormatBytes(used))
	} else {
		jvmWriteTabularMetric(out, prefix+"-available-bytes", maxBytes)
		jvmWriteTabularMetric(out, prefix+"-peak-bytes", peakUsed)
		jvmWriteTabularMetric(out, prefix+"-current-bytes", used)
	}
}

func jvmWriteMemoryMetrics(out io.Writer, verbose, pretty bool, ms *runtime.MemStats) {
	if pretty {
		jvmPrintln(out, "\nMemory usage")
	}

	heapUsed := int64(ms.HeapAlloc)
	heapMax := int64(ms.HeapSys)

	if verbose {
		jvmWriteOverallMemoryUsage(out, heapUsed, heapMax, "Heap", pretty)
		nonHeap := int64(ms.Sys - ms.HeapSys)
		jvmWriteOverallMemoryUsage(out, int64(ms.StackInuse), nonHeap, "Non-heap", pretty)

		const poolName = "Go heap"
		if pretty {
			jvmPrintln(out, "\tPool "+poolName)
			jvmWritePoolMemoryUsage(out, heapUsed, heapUsed, heapMax, "", true)
		} else {
			prefix := "mem-pool-" + jvmNormalizeName(poolName)
			jvmWritePoolMemoryUsage(out, heapUsed, heapUsed, heapMax, prefix, false)
		}
	} else {
		if pretty {
			_, _ = io.WriteString(out, "\tAggregate of 1 memory pools\n")
			jvmWritePoolMemoryUsage(out, heapUsed, heapUsed, heapMax, "", true)
		} else {
			jvmWritePoolMemoryUsage(out, heapUsed, heapUsed, heapMax, "mem", false)
		}
	}
}

func jvmWriteGarbageCollectionStats(out io.Writer, verbose, pretty bool, ms *runtime.MemStats) {
	collectionCount := int64(ms.NumGC)
	collectionTime := int64(ms.PauseTotalNs / 1e6)

	if verbose {
		const gcName = "Go GC"
		if pretty {
			jvmPrintln(out, "\nGarbage collection stats")
			jvmPrintln(out, "\tCollector "+gcName)
			_, _ = fmt.Fprintf(
				out,
				"\t\tcollection count   : %d\n\t\tcollection time    : %d ms\n",
				collectionCount,
				collectionTime)
		} else {
			name := jvmNormalizeName(gcName)
			jvmWriteTabularMetric(out, "gc-"+name+"-collection-count", collectionCount)
			jvmWriteTabularMetric(out, "gc-"+name+"-collection-time-ms", collectionTime)
		}
	} else {
		if pretty {
			jvmPrintln(out, "\nGarbage collection stats")
			_, _ = io.WriteString(out, "\tAggregate of 1 collectors\n")
			_, _ = fmt.Fprintf(
				out,
				"\t\tcollection count   : %d\n\t\tcollection time    : %d ms\n",
				collectionCount,
				collectionTime)
		} else {
			name := jvmNormalizeName("aggregate")
			jvmWriteTabularMetric(out, "gc-"+name+"-collection-count", collectionCount)
			jvmWriteTabularMetric(out, "gc-"+name+"-collection-time-ms", collectionTime)
		}
	}
}

func jvmNormalizeName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "_"))
}

func jvmNormalizeTabularColonPos(s string) string {
	index := strings.IndexByte(s, ':')
	if index < 0 {
		panic(jutil.NewStringIndexOutOfBoundsException(index, len(s)))
	}
	if index >= jvmTabularColonPos {
		return s
	}
	return s[:index] + strings.Repeat(" ", jvmTabularColonPos-index) + s[index:]
}

func jvmFormatBytes(numBytes int64) string {
	switch {
	case numBytes < jvmOneKiloByte:
		return strconv.FormatInt(numBytes, 10) + " B"
	case numBytes < jvmOneMegaByte:
		return strconv.FormatInt(numBytes/jvmOneKiloByte, 10) + " KB"
	case numBytes < jvmOneGigaByte:
		return strconv.FormatInt(numBytes/jvmOneMegaByte, 10) + " MB"
	default:
		return strconv.FormatInt(numBytes/jvmOneGigaByte, 10) + " GB"
	}
}
