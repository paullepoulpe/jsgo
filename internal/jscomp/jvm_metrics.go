// Port of com/google/javascript/jscomp/JvmMetrics.java @ closure-compiler v20261015

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

// JvmMetrics is a class to report jvm/jmx statistics. Borrowed from:
// http://code.google.com/p/dart/source/browse/trunk/dart/compiler/java/com/google/dart/compiler/metrics/JvmMetrics.java
//
// Go-only: language-adapter — there is no JVM. The java.lang.management
// beans are modelled by jvmMetrics_* adapters that report the Go runtime's
// equivalents: one "Go GC" collector (runtime.MemStats NumGC and
// PauseTotalNs), the heap and non-heap (stack + other system) memory, one
// "Go heap" memory pool, and a "Go" compiler with a total compilation time
// of 0 (Go is compiled ahead of time). The report formats are Java's.
type JvmMetrics struct{}

const (
	jvmMetrics_TABULAR_COLON_POS       = 40
	jvmMetrics_ONE_KILO_BYTE     int64 = 1 << 10
	jvmMetrics_ONE_MEGA_BYTE     int64 = 1 << 20
	jvmMetrics_ONE_GIGA_BYTE     int64 = 1 << 30
)

// jvmMetrics_MemoryUsage is java.lang.management.MemoryUsage.
type jvmMetrics_MemoryUsage struct {
	init, used, committed, max int64
}

func newJvmMetrics_MemoryUsage(init, used, committed, max int64) *jvmMetrics_MemoryUsage {
	return &jvmMetrics_MemoryUsage{init: init, used: used, committed: committed, max: max}
}

func (u *jvmMetrics_MemoryUsage) GetUsed() int64 { return u.used }
func (u *jvmMetrics_MemoryUsage) GetMax() int64  { return u.max }

// jvmMetrics_GarbageCollectorMXBean is java.lang.management.GarbageCollectorMXBean.
type jvmMetrics_GarbageCollectorMXBean struct {
	name            string
	collectionCount int64
	collectionTime  int64
}

func (b *jvmMetrics_GarbageCollectorMXBean) GetName() string           { return b.name }
func (b *jvmMetrics_GarbageCollectorMXBean) GetCollectionCount() int64 { return b.collectionCount }
func (b *jvmMetrics_GarbageCollectorMXBean) GetCollectionTime() int64  { return b.collectionTime }

// jvmMetrics_MemoryPoolMXBean is java.lang.management.MemoryPoolMXBean.
type jvmMetrics_MemoryPoolMXBean struct {
	name      string
	usage     *jvmMetrics_MemoryUsage
	peakUsage *jvmMetrics_MemoryUsage
}

func (b *jvmMetrics_MemoryPoolMXBean) GetName() string                       { return b.name }
func (b *jvmMetrics_MemoryPoolMXBean) GetUsage() *jvmMetrics_MemoryUsage     { return b.usage }
func (b *jvmMetrics_MemoryPoolMXBean) GetPeakUsage() *jvmMetrics_MemoryUsage { return b.peakUsage }

func jvmMetrics_memStats() *runtime.MemStats {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return &ms
}

// ManagementFactory.getGarbageCollectorMXBeans()
func jvmMetrics_getGarbageCollectorMXBeans() []*jvmMetrics_GarbageCollectorMXBean {
	ms := jvmMetrics_memStats()
	return []*jvmMetrics_GarbageCollectorMXBean{{
		name:            "Go GC",
		collectionCount: int64(ms.NumGC),
		collectionTime:  int64(ms.PauseTotalNs / 1e6),
	}}
}

// ManagementFactory.getMemoryPoolMXBeans()
func jvmMetrics_getMemoryPoolMXBeans() []*jvmMetrics_MemoryPoolMXBean {
	ms := jvmMetrics_memStats()
	usage := newJvmMetrics_MemoryUsage(0, int64(ms.HeapAlloc), int64(ms.HeapSys), int64(ms.HeapSys))
	return []*jvmMetrics_MemoryPoolMXBean{{name: "Go heap", usage: usage, peakUsage: usage}}
}

// ManagementFactory.getMemoryMXBean().getHeapMemoryUsage() and
// getNonHeapMemoryUsage().
func jvmMetrics_getHeapMemoryUsage() *jvmMetrics_MemoryUsage {
	ms := jvmMetrics_memStats()
	return newJvmMetrics_MemoryUsage(0, int64(ms.HeapAlloc), int64(ms.HeapSys), int64(ms.HeapSys))
}

func jvmMetrics_getNonHeapMemoryUsage() *jvmMetrics_MemoryUsage {
	ms := jvmMetrics_memStats()
	nonHeap := int64(ms.Sys - ms.HeapSys)
	return newJvmMetrics_MemoryUsage(0, int64(ms.StackInuse), nonHeap, nonHeap)
}

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

	if len(st) > 0 {
		for len(st) > 0 {
			types := st[0]
			st = st[1:]
			for _, type_ := range jvmMetrics_tokenize(types, ",") {
				jvmMetrics_writeMetrics(out, type_, verboseMode, prettyMode)
			}
		}
	} else {
		// the default
		jvmMetrics_writeMetrics(out, "all", verboseMode, prettyMode)
	}
}

// jvmMetrics_tokenize is java.util.StringTokenizer(s, delims): the
// non-empty runs between delimiter characters.
func jvmMetrics_tokenize(s string, delims string) []string {
	return strings.FieldsFunc(s, func(r rune) bool { return strings.ContainsRune(delims, r) })
}

func jvmMetrics_println(out io.Writer, s string) {
	io.WriteString(out, s+"\n")
}

// Java: JvmMetrics.writeMetrics(PrintStream, String, boolean, boolean)
func jvmMetrics_writeMetrics(out io.Writer, type_ string, verbose bool, pretty bool) {
	if type_ == "gc" || strings.EqualFold(type_, "all") {
		jvmMetrics_writeGarbageCollectionStats(out, verbose, pretty)
	}
	if type_ == "mem" || strings.EqualFold(type_, "all") {
		jvmMetrics_writeMemoryMetrics(out, verbose, pretty)
	}
	if type_ == "jit" || strings.EqualFold(type_, "all") {
		jvmMetrics_writeJitMetrics(out, verbose, pretty)
	}
}

// Java: JvmMetrics.writeJitMetrics(PrintStream, boolean, boolean)
func jvmMetrics_writeJitMetrics(out io.Writer, verbose bool, pretty bool) {
	// Go-only: language-adapter — CompilationMXBean: no JIT.
	cBeanName := "Go"
	var cBeanTotalCompilationTime int64 = 0
	var name string
	if verbose {
		name = cBeanName
	} else {
		name = "total"
	}

	if pretty {
		jvmMetrics_println(out, "\nJIT Stats")
		io.WriteString(out, jutil.StringFormat("\t%s jit time: %d ms%n", name, cBeanTotalCompilationTime))
	} else {
		jvmMetrics_println(out, jvmMetrics_normalizeTabularColonPos(jutil.StringFormat("%s-jit-time-ms : %d",
			jvmMetrics_normalizeName(name), cBeanTotalCompilationTime)))
	}
}

// Java: JvmMetrics.writeOverallMemoryUsage(PrintStream, MemoryUsage, String, boolean)
func jvmMetrics_writeOverallMemoryUsage(out io.Writer, usage *jvmMetrics_MemoryUsage, prefix string, pretty bool) {
	if pretty {
		io.WriteString(out, jutil.StringFormat("\t%s\n", prefix))
		io.WriteString(out, jutil.StringFormat("\t\tavailable         : %s\n", jvmMetrics_formatBytes(usage.GetMax())))
		io.WriteString(out, jutil.StringFormat("\t\tcurrent           : %s\n", jvmMetrics_formatBytes(usage.GetUsed())))
	} else {
		prefix = jvmMetrics_normalizeName(prefix)
		jvmMetrics_println(out, jvmMetrics_normalizeTabularColonPos(
			jutil.StringFormat("%s-available-bytes : %d", prefix, usage.GetMax())))
		jvmMetrics_println(out, jvmMetrics_normalizeTabularColonPos(
			jutil.StringFormat("%s-current-bytes : %d", prefix, usage.GetUsed())))
	}
}

// Java: JvmMetrics.writePoolMemoryUsage(PrintStream, MemoryUsage, MemoryUsage, String, boolean)
func jvmMetrics_writePoolMemoryUsage(out io.Writer, usage *jvmMetrics_MemoryUsage, peakUsage *jvmMetrics_MemoryUsage, prefix *string, pretty bool) {
	if pretty {
		io.WriteString(out, jutil.StringFormat("\t\tavailable         : %s\n",
			jvmMetrics_formatBytes(usage.GetMax())))
		io.WriteString(out, jutil.StringFormat("\t\tpeak              : %s\n",
			jvmMetrics_formatBytes(peakUsage.GetUsed())))
		io.WriteString(out, jutil.StringFormat("\t\tcurrent           : %s\n",
			jvmMetrics_formatBytes(usage.GetUsed())))
	} else {
		var p any = prefix
		jvmMetrics_println(out, jvmMetrics_normalizeTabularColonPos(
			jutil.StringFormat("%s-available-bytes : %d", p, usage.GetMax())))
		jvmMetrics_println(out, jvmMetrics_normalizeTabularColonPos(
			jutil.StringFormat("%s-peak-bytes : %d", p, peakUsage.GetUsed())))
		jvmMetrics_println(out, jvmMetrics_normalizeTabularColonPos(
			jutil.StringFormat("%s-current-bytes : %d", p, usage.GetUsed())))
	}
}

// Java: JvmMetrics.writeMemoryMetrics(PrintStream, boolean, boolean)
func jvmMetrics_writeMemoryMetrics(out io.Writer, verbose bool, pretty bool) {
	if pretty {
		jvmMetrics_println(out, "\nMemory usage")
	}

	// only show overall stats in verbose mode
	if verbose {
		usage := jvmMetrics_getHeapMemoryUsage()
		jvmMetrics_writeOverallMemoryUsage(out, usage, "Heap", pretty)

		usage = jvmMetrics_getNonHeapMemoryUsage()
		jvmMetrics_writeOverallMemoryUsage(out, usage, "Non-heap", pretty)
	}

	if verbose {
		mpBeans := jvmMetrics_getMemoryPoolMXBeans()
		for _, mpBean := range mpBeans {
			currentUsage := mpBean.GetUsage()
			peakUsage := mpBean.GetPeakUsage()
			if pretty {
				jvmMetrics_println(out, "\tPool "+mpBean.GetName())
				jvmMetrics_writePoolMemoryUsage(out, currentUsage, peakUsage, nil, true)
			} else {
				prefix := "mem-pool-" + jvmMetrics_normalizeName(mpBean.GetName())
				jvmMetrics_writePoolMemoryUsage(out, currentUsage, peakUsage, &prefix, false)
			}
		}
	} else {
		var available int64 = 0
		var current int64 = 0
		var peak int64 = 0
		mpBeans := jvmMetrics_getMemoryPoolMXBeans()
		for _, mpBean := range mpBeans {
			currentUsage := mpBean.GetUsage()
			available += currentUsage.GetMax()
			current += currentUsage.GetUsed()
			peakUsage := mpBean.GetPeakUsage()
			peak += peakUsage.GetUsed()
		}
		summaryUsage := newJvmMetrics_MemoryUsage(0, current, current, available)
		summaryPeakUsage := newJvmMetrics_MemoryUsage(0, peak, peak, peak)
		if pretty {
			io.WriteString(out, jutil.StringFormat("\tAggregate of %d memory pools\n", len(mpBeans)))
			jvmMetrics_writePoolMemoryUsage(out, summaryUsage, summaryPeakUsage, nil, true)
		} else {
			prefix := "mem"
			jvmMetrics_writePoolMemoryUsage(out, summaryUsage, summaryPeakUsage, &prefix, false)
		}
	}
}

// Java: JvmMetrics.writeGarbageCollectionStats(PrintStream, boolean, boolean)
func jvmMetrics_writeGarbageCollectionStats(out io.Writer, verbose bool, pretty bool) {
	gcBeans := jvmMetrics_getGarbageCollectorMXBeans()

	if verbose {
		if pretty {
			jvmMetrics_println(out, "\nGarbage collection stats")
			for _, gcBean := range gcBeans {
				jvmMetrics_println(out, "\tCollector "+gcBean.GetName())
				io.WriteString(out, jutil.StringFormat(
					"\t\tcollection count   : %d\n", gcBean.GetCollectionCount()))
				io.WriteString(out, jutil.StringFormat(
					"\t\tcollection time    : %d ms\n", gcBean.GetCollectionTime()))
			}
		} else {
			for _, gcBean := range gcBeans {
				name := jvmMetrics_normalizeName(gcBean.GetName())
				jvmMetrics_println(out, jvmMetrics_normalizeTabularColonPos(jutil.StringFormat("gc-"+name+
					"-collection-count : %d", gcBean.GetCollectionCount())))
				jvmMetrics_println(out, jvmMetrics_normalizeTabularColonPos(jutil.StringFormat("gc-"+name+
					"-collection-time-ms : %d", gcBean.GetCollectionTime())))
			}
		}
	} else {
		var collectionCount int64 = 0
		var collectionTime int64 = 0
		collectorCount := len(gcBeans)
		for _, gcBean := range gcBeans {
			collectionCount += gcBean.GetCollectionCount()
			collectionTime += gcBean.GetCollectionTime()
		}
		if pretty {
			jvmMetrics_println(out, "\nGarbage collection stats")
			io.WriteString(out, jutil.StringFormat("\tAggregate of %d collectors\n", collectorCount))
			io.WriteString(out, jutil.StringFormat("\t\tcollection count   : %d\n", collectionCount))
			io.WriteString(out, jutil.StringFormat("\t\tcollection time    : %d ms\n", collectionTime))
		} else {
			name := jvmMetrics_normalizeName("aggregate")
			jvmMetrics_println(out, jvmMetrics_normalizeTabularColonPos(
				jutil.StringFormat("gc-%s-collection-count : %d", name,
					collectionCount)))
			jvmMetrics_println(out, jvmMetrics_normalizeTabularColonPos(
				jutil.StringFormat("gc-%s-collection-time-ms : %d", name,
					collectionTime)))
		}
	}
}

// Java: JvmMetrics.normalizeName(String)
func jvmMetrics_normalizeName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "_"))
}

// Java: JvmMetrics.normalizeTabularColonPos(String)
func jvmMetrics_normalizeTabularColonPos(s string) string {
	sb := jutil.NewStringBuilderFrom(s)
	index := jutil.StringIndexOf(s, ":")
	for ; index < jvmMetrics_TABULAR_COLON_POS; index++ {
		sb.InsertChar(index, ' ')
	}
	return sb.String()
}

// Java: JvmMetrics.formatBytes(long)
func jvmMetrics_formatBytes(numBytes int64) string {
	if numBytes < jvmMetrics_ONE_KILO_BYTE {
		return jutil.StringFormat("%d B", numBytes)
	} else if numBytes < jvmMetrics_ONE_MEGA_BYTE {
		return jutil.StringFormat("%d KB", numBytes/jvmMetrics_ONE_KILO_BYTE)
	} else if numBytes < jvmMetrics_ONE_GIGA_BYTE {
		return jutil.StringFormat("%d MB", numBytes/jvmMetrics_ONE_MEGA_BYTE)
	} else {
		return jutil.StringFormat("%d GB", numBytes/jvmMetrics_ONE_GIGA_BYTE)
	}
}
