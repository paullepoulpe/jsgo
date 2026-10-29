// Port of com/google/javascript/jscomp/PerformanceTracker.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"cmp"
	"iter"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PerformanceTracker_StatsMap maps pass names to their cumulative statistics.
type PerformanceTracker_StatsMap map[string]*PerformanceTracker_Stats

// Get returns the cumulative Stats for passName, or nil if absent.
func (m PerformanceTracker_StatsMap) Get(passName string) *PerformanceTracker_Stats {
	return m[passName]
}

// PerformanceTracker collects statistics about the runtime of each pass,
// and how much a pass impacts the size of the compiled output, before and
// after gzip.
type PerformanceTracker struct {
	externsRoot rhino.Node

	mode CompilerOptions_TracerMode

	// Keeps track of AST changes and computes code size estimation
	// if there is any.
	codeChange *RecentChange

	jsRoot rhino.Node

	initAstSize    int
	initCodeSize   int
	initGzCodeSize int

	startTime     int64
	endTime       int64
	passesRuntime int
	maxMem        int
	runs          int
	changes       int
	loopRuns      int
	loopChanges   int

	jsLines       int
	jsSources     int
	externLines   int
	externSources int

	transitiveSummaryLines         int
	transitiveSummarySources       int
	transitiveSummaryLinesParsed   int
	transitiveSummarySourcesParsed int

	prePruningRecorded bool

	// The following fields for tracking size changes are just estimates.
	// They do not take into account preserved license blocks, newline padding,
	// or pretty printing (if enabled), since they don't use CodePrinter.
	// To get exact sizes, call compiler.toSource() for the final generated code.
	astSize    int
	codeSize   int
	gzCodeSize int
	astDiff    int
	diff       int
	gzDiff     int

	// ArrayDeque used as a stack: the top is the last element.
	currentPass []*PerformanceTracker_Stats

	// Cumulative stats for each compiler pass, plus first-seen order.
	passSummary      PerformanceTracker_StatsMap
	passSummaryOrder []*PerformanceTracker_Stats

	astManifest          [256]int32
	astManifestPopulated bool

	// Stats a single run of a compiler pass.
	log []*PerformanceTracker_Stats

	pruningAnalysisSummary string

	disambiguatePropertiesSummary string

	ambiguatePropertiesSummary string
}

const PerformanceTracker_DEFAULT_WHEN_SIZE_UNTRACKED = -1

// Java: PerformanceTracker(Node, Node, TracerMode)
func NewPerformanceTracker(externsRoot rhino.Node, jsRoot rhino.Node, mode CompilerOptions_TracerMode) *PerformanceTracker {
	p := &PerformanceTracker{}
	p.codeChange = NewRecentChange()
	p.initAstSize = PerformanceTracker_DEFAULT_WHEN_SIZE_UNTRACKED
	p.initCodeSize = PerformanceTracker_DEFAULT_WHEN_SIZE_UNTRACKED
	p.initGzCodeSize = PerformanceTracker_DEFAULT_WHEN_SIZE_UNTRACKED
	p.passesRuntime = 0
	p.maxMem = 0
	p.runs = 0
	p.changes = 0
	p.loopRuns = 0
	p.loopChanges = 0
	p.jsLines = 0
	p.jsSources = 0
	p.externLines = 0
	p.externSources = 0
	p.transitiveSummaryLines = 0
	p.transitiveSummarySources = 0
	p.transitiveSummaryLinesParsed = 0
	p.transitiveSummarySourcesParsed = 0
	p.prePruningRecorded = false
	p.astSize = PerformanceTracker_DEFAULT_WHEN_SIZE_UNTRACKED
	p.codeSize = PerformanceTracker_DEFAULT_WHEN_SIZE_UNTRACKED
	p.gzCodeSize = PerformanceTracker_DEFAULT_WHEN_SIZE_UNTRACKED
	p.astDiff = 0
	p.diff = 0
	p.gzDiff = 0
	p.pruningAnalysisSummary = "not executed"
	p.disambiguatePropertiesSummary = "not executed"
	p.ambiguatePropertiesSummary = "not executed"

	jutil.CheckArgumentMsg(mode != CompilerOptions_TracerMode_OFF, "PerformanceTracker can't work without tracer data.")
	p.startTime = time.Now().UnixMilli()
	p.externsRoot = externsRoot
	p.jsRoot = jsRoot
	p.mode = mode
	return p
}

// Java: PerformanceTracker.getCodeChangeHandler()
func (p *PerformanceTracker) GetCodeChangeHandler() CodeChangeHandler {
	return p.codeChange
}

// Java: PerformanceTracker.recordPassStart(String, boolean)
func (p *PerformanceTracker) RecordPassStart(passName string, isOneTime bool) {
	p.currentPass = append(p.currentPass, NewPerformanceTracker_Stats(passName, isOneTime))
	// In Compiler, toSource may be called after every pass X. We don't want it
	// to reset the handler, because recordPassStop for pass X has not been
	// called, so we are falsely logging that pass X didn't make changes.
	if passName != "toSource" {
		p.codeChange.Reset()
	}
}

// UpdateAfterDeserialize updates the saved jsRoot and resets the size
// tracking fields accordingly.
//
// Java: PerformanceTracker.updateAfterDeserialize(Node)
func (p *PerformanceTracker) UpdateAfterDeserialize(jsRoot rhino.Node) {
	// TODO(bradfordcsmith): Restore line counts for inputs and externs.
	p.jsRoot = jsRoot
	if !p.TracksAstSize() {
		return
	}
	p.astSize = NodeUtil_CountAstSize(p.jsRoot)
	p.initAstSize = p.astSize
	if !p.TracksSize() {
		return
	}
	estimator := PerformanceTrackerCodeSizeEstimator_Estimate(p.jsRoot, p.TracksGzSize())
	p.codeSize = estimator.GetCodeSize()
	p.initCodeSize = p.codeSize
	if p.TracksGzSize() {
		p.gzCodeSize = estimator.GetZippedCodeSize()
		p.initGzCodeSize = p.gzCodeSize
	}
}

// RecordPassStop collects information about a pass P after P finishes
// running, eg, how much time P took and what was its impact on code size.
//
// passName is the short name of the pass; runtime is the execution time in
// milliseconds.
//
// Java: PerformanceTracker.recordPassStop(String, long)
func (p *PerformanceTracker) RecordPassStop(passName string, runtime int64) {
	allocMem := p.GetAllocatedMegabytes()
	if len(p.currentPass) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.pop()
	}
	logStats := p.currentPass[len(p.currentPass)-1]
	p.currentPass = p.currentPass[:len(p.currentPass)-1]
	jutil.CheckState(passName == logStats.pass)
	p.log = append(p.log, logStats)

	// Update fields that aren't related to code size
	logStats.runtime = runtime
	logStats.allocMem = allocMem
	logStats.runs = 1
	if p.codeChange.HasCodeChanged() {
		logStats.changes = 1
	}
	if passName == PassNames_PARSE_INPUTS {
		p.RecordParsingStop(logStats)
	} else if p.codeChange.HasCodeChanged() && p.TracksAstSize() {
		p.RecordOtherPassStop(logStats)
	}
}

// Java: PerformanceTracker.recordParsingStop(Stats)
func (p *PerformanceTracker) RecordParsingStop(logStats *PerformanceTracker_Stats) {
	if p.TracksAstManifest() {
		p.PopulateAstManifest()
	}
	p.RecordInputCount()
	if !p.TracksAstSize() {
		return
	}
	p.astSize = NodeUtil_CountAstSize(p.jsRoot)
	p.initAstSize = p.astSize
	logStats.astSize = p.initAstSize
	if !p.TracksSize() {
		return
	}
	estimator := PerformanceTrackerCodeSizeEstimator_Estimate(p.jsRoot, p.TracksGzSize())
	p.codeSize = estimator.GetCodeSize()
	p.initCodeSize = p.codeSize
	logStats.size = p.initCodeSize
	if p.TracksGzSize() {
		p.gzCodeSize = estimator.GetZippedCodeSize()
		p.initGzCodeSize = p.gzCodeSize
		logStats.gzSize = p.initGzCodeSize
	}
}

// Java: PerformanceTracker.recordOtherPassStop(Stats)
func (p *PerformanceTracker) RecordOtherPassStop(logStats *PerformanceTracker_Stats) {
	newSize := NodeUtil_CountAstSize(p.jsRoot)
	logStats.astDiff = p.astSize - newSize
	logStats.astSize = newSize
	p.astSize = logStats.astSize
	if !p.TracksSize() {
		return
	}
	estimator := PerformanceTrackerCodeSizeEstimator_Estimate(p.jsRoot, p.TracksGzSize())
	newSize = estimator.GetCodeSize()
	logStats.diff = p.codeSize - newSize
	logStats.size = newSize
	p.codeSize = logStats.size
	if p.TracksGzSize() {
		newSize = estimator.GetZippedCodeSize()
		logStats.gzDiff = p.gzCodeSize - newSize
		logStats.gzSize = newSize
		p.gzCodeSize = logStats.gzSize
	}
}

// Java: PerformanceTracker.recordInputCount()
func (p *PerformanceTracker) RecordInputCount() {
	for n := p.externsRoot.GetFirstChild(); !n.IsNil(); n = n.GetNext() {
		if NodeUtil_IsFromTypeSummary(n) {
			p.transitiveSummarySources += 1
			p.transitiveSummaryLines += NodeUtil_EstimateNumLines(n)
		} else {
			p.externSources += 1
			p.externLines += NodeUtil_EstimateNumLines(n)
		}
	}

	for n := p.jsRoot.GetFirstChild(); !n.IsNil(); n = n.GetNext() {
		p.jsSources += 1
		p.jsLines += NodeUtil_EstimateNumLines(n)
	}
}

// Java: PerformanceTracker.bytesToMB(long)
func (p *PerformanceTracker) BytesToMB(bytes int64) int {
	return int(int32(bytes / (1024 * 1024)))
}

// Java: PerformanceTracker.getAllocatedMegabytes()
func (p *PerformanceTracker) GetAllocatedMegabytes() int {
	return p.BytesToMB(Platform_TotalMemory() - Platform_FreeMemory())
}

// Java: PerformanceTracker.tracksSize()
func (p *PerformanceTracker) TracksSize() bool {
	return p.mode == CompilerOptions_TracerMode_RAW_SIZE || p.mode == CompilerOptions_TracerMode_ALL
}

// Java: PerformanceTracker.tracksGzSize()
func (p *PerformanceTracker) TracksGzSize() bool {
	return p.mode == CompilerOptions_TracerMode_ALL
}

// Java: PerformanceTracker.tracksAstSize()
func (p *PerformanceTracker) TracksAstSize() bool {
	return p.mode != CompilerOptions_TracerMode_TIMING_ONLY
}

// Java: PerformanceTracker.tracksAstManifest()
func (p *PerformanceTracker) TracksAstManifest() bool {
	return p.mode.IsOn()
}

// Java: PerformanceTracker.getRuntime()
func (p *PerformanceTracker) GetRuntime() int {
	p.CalcTotalStats()
	return p.passesRuntime
}

// Java: PerformanceTracker.getSize()
func (p *PerformanceTracker) GetSize() int {
	p.CalcTotalStats()
	return p.codeSize
}

// Java: PerformanceTracker.getGzSize()
func (p *PerformanceTracker) GetGzSize() int {
	p.CalcTotalStats()
	return p.gzCodeSize
}

// Java: PerformanceTracker.getAstSize()
func (p *PerformanceTracker) GetAstSize() int {
	p.CalcTotalStats()
	return p.astSize
}

// Java: PerformanceTracker.getChanges()
func (p *PerformanceTracker) GetChanges() int {
	p.CalcTotalStats()
	return p.changes
}

// Java: PerformanceTracker.getLoopChanges()
func (p *PerformanceTracker) GetLoopChanges() int {
	p.CalcTotalStats()
	return p.loopChanges
}

// Java: PerformanceTracker.getRuns()
func (p *PerformanceTracker) GetRuns() int {
	p.CalcTotalStats()
	return p.runs
}

// Java: PerformanceTracker.getLoopRuns()
func (p *PerformanceTracker) GetLoopRuns() int {
	p.CalcTotalStats()
	return p.loopRuns
}

// Java: PerformanceTracker.getStats()
func (p *PerformanceTracker) GetStats() PerformanceTracker_StatsMap {
	p.CalcTotalStats()
	return p.passSummary
}

// Java: PerformanceTracker.calcTotalStats()
func (p *PerformanceTracker) CalcTotalStats() {
	// This method only does work the first time it is called
	if p.passSummary != nil {
		return
	}
	p.endTime = time.Now().UnixMilli()

	p.PopulatePassSummary()

	for _, stats := range p.passSummaryOrder {
		p.passesRuntime = int(int32(int64(p.passesRuntime) + stats.runtime))
		p.maxMem = max(p.maxMem, stats.allocMem)
		p.runs += stats.runs
		p.changes += stats.changes
		if !stats.isOneTime {
			p.loopRuns += stats.runs
			p.loopChanges += stats.changes
		}
		p.astDiff += stats.astDiff
		p.diff += stats.diff
		p.gzDiff += stats.gzDiff
	}
	jutil.CheckState(!p.TracksAstSize() || p.initAstSize == p.astDiff+p.astSize)
	jutil.CheckState(!p.TracksSize() || p.initCodeSize == p.diff+p.codeSize)
	jutil.CheckState(!p.TracksGzSize() || p.initGzCodeSize == p.gzDiff+p.gzCodeSize)
}

// Java: PerformanceTracker.populatePassSummary()
func (p *PerformanceTracker) PopulatePassSummary() {
	p.passSummary = make(PerformanceTracker_StatsMap, len(p.log))
	p.passSummaryOrder = make([]*PerformanceTracker_Stats, 0, len(p.log))

	for _, logStat := range p.log {
		passName := logStat.pass
		entry := p.passSummary[passName]
		if entry == nil {
			entry = NewPerformanceTracker_Stats(passName, logStat.isOneTime)
			p.passSummary[passName] = entry
			p.passSummaryOrder = append(p.passSummaryOrder, entry)
		}
		entry.runtime += logStat.runtime
		entry.allocMem = max(entry.allocMem, logStat.allocMem)
		entry.runs++
		entry.changes += logStat.changes
		entry.astDiff += logStat.astDiff
		entry.diff += logStat.diff
		entry.gzDiff += logStat.gzDiff
		// We don't populate the size fields in the passSummary stats.
		// We used to put the size after the last time a pass was run, but that is
		// a pretty meaningless thing to measure.
	}
}

// Java: PerformanceTracker.populateAstManifest()
func (p *PerformanceTracker) PopulateAstManifest() {
	if p.astManifestPopulated {
		return
	}
	p.astManifestPopulated = true
	NodeUtil_VisitPreOrder(p.jsRoot, NodeUtil_VisitorFunc(func(n rhino.Node) {
		p.astManifest[n.GetToken()]++
	}))
}

// Java: PerformanceTracker.setPruningAnalysisSummary(String)
func (p *PerformanceTracker) SetPruningAnalysisSummary(summary string) {
	p.pruningAnalysisSummary = summary
}

// Java: PerformanceTracker.recordPrePruningInputCount(Iterable<CompilerInput>, Iterable<CompilerInput>)
func (p *PerformanceTracker) RecordPrePruningInputCount(
	originalInputs iter.Seq[*CompilerInput], originalExterns iter.Seq[*CompilerInput]) {
	for input := range originalExterns {
		name := input.GetSourceFile().GetName()
		if jutil.StringEndsWith(name, ".i.js") {
			p.transitiveSummarySourcesParsed += 1
			p.transitiveSummaryLinesParsed += input.GetSourceFile().GetNumLines()
		}
	}
	for input := range originalInputs {
		name := input.GetSourceFile().GetName()
		if jutil.StringEndsWith(name, ".i.js") {
			p.transitiveSummarySourcesParsed += 1
			p.transitiveSummaryLinesParsed += input.GetSourceFile().GetNumLines()
		}
	}
	p.prePruningRecorded = true
}

// Java: PerformanceTracker.setDisambiguatePropertiesSummary(String)
func (p *PerformanceTracker) SetDisambiguatePropertiesSummary(summary string) {
	p.disambiguatePropertiesSummary = summary
}

// Java: PerformanceTracker.setAmbiguatePropertiesSummary(String)
func (p *PerformanceTracker) SetAmbiguatePropertiesSummary(summary string) {
	p.ambiguatePropertiesSummary = summary
}

// OutputTracerReport prints a summary, which contains aggregate stats for
// all runs of each pass and a log, which contains stats for each individual
// run.
//
// Java: PerformanceTracker.outputTracerReport(PrintStream)
func (p *PerformanceTracker) OutputTracerReport(output *jutil.PrintStream) {
	jvmMetricsOptions := "verbose:pretty:all"
	JvmMetrics_MaybeWriteJvmMetrics(output, &jvmMetricsOptions)
	p.CalcTotalStats()

	output.Println(
		strings.Join([]string{
			"",
			"TOTAL:",
			"Start time(ms): " + jutil.ToJavaString(p.startTime),
			"End time(ms): " + jutil.ToJavaString(p.endTime),
			"Wall time(ms): " + jutil.ToJavaString(p.endTime-p.startTime),
			"Passes runtime(ms): " + jutil.ToJavaString(p.passesRuntime),
			"Max mem usage (measured after each pass)(MB): " + jutil.ToJavaString(p.maxMem),
			"#Runs: " + jutil.ToJavaString(p.runs),
			"#Changing runs: " + jutil.ToJavaString(p.changes),
			"#Loopable runs: " + jutil.ToJavaString(p.loopRuns),
			"#Changing loopable runs: " + jutil.ToJavaString(p.loopChanges),
			"Estimated AST reduction(#nodes): " + jutil.ToJavaString(p.astDiff),
			"Estimated Reduction(bytes): " + jutil.ToJavaString(p.diff),
			"Estimated GzReduction(bytes): " + jutil.ToJavaString(p.gzDiff),
			"Estimated AST size(#nodes): " + jutil.ToJavaString(p.astSize),
			"Estimated Size(bytes): " + jutil.ToJavaString(p.codeSize),
			"Estimated GzSize(bytes): " + jutil.ToJavaString(p.gzCodeSize),
			"",
			"DisambiguateProperties: " + p.disambiguatePropertiesSummary,
			"AmbiguateProperties: " + p.ambiguatePropertiesSummary,
		}, "\n"))

	summaryLinesRaw := p.transitiveSummaryLines
	summarySourcesRaw := p.transitiveSummarySources
	if p.prePruningRecorded {
		summaryLinesRaw = p.transitiveSummaryLinesParsed
		summarySourcesRaw = p.transitiveSummarySourcesParsed
	}
	output.Println(
		strings.Join([]string{
			"",
			"Inputs:",
			"JS lines:   " + jutil.ToJavaString(p.jsLines),
			"JS sources: " + jutil.ToJavaString(p.jsSources),
			"Extern lines:   " + jutil.ToJavaString(p.externLines),
			"Extern sources: " + jutil.ToJavaString(p.externSources),
			"Type summary lines (raw input):   " +
				jutil.ToJavaString(summaryLinesRaw),
			"Type summary lines (post-pruning):   " + jutil.ToJavaString(p.transitiveSummaryLines),
			"Type summary sources (raw input):   " +
				jutil.ToJavaString(summarySourcesRaw),
			"Type summary sources (post-pruning):   " + jutil.ToJavaString(p.transitiveSummarySources),
			"",
			"Dependency pruning analysis: " + p.pruningAnalysisSummary,
		}, "\n"))

	output.Println("\nSummary:\npass,runtime,allocMem,runs,changingRuns,astReduction,reduction,gzReduction")
	entries := slices.Clone(p.passSummaryOrder)
	// Stream.sorted is stable for an ordered stream.
	slices.SortStableFunc(entries, func(a, b *PerformanceTracker_Stats) int {
		return cmp.Compare(a.runtime, b.runtime)
	})
	for _, stats := range entries {
		output.Println(jutil.StringFormat(
			"%s,%d,%d,%d,%d,%d,%d,%d",
			stats.pass,
			stats.runtime,
			stats.allocMem,
			stats.runs,
			stats.changes,
			stats.astDiff,
			stats.diff,
			stats.gzDiff))
	}

	output.Println("\nLog:\npass,runtime,allocMem,codeChanged,astReduction,reduction,gzReduction,astSize,size,gzSize")
	for _, stats := range p.log {
		output.Printf(
			"%s,%d,%d,%b,%d,%d,%d,%d,%d,%d\n",
			stats.pass,
			stats.runtime,
			stats.allocMem,
			stats.changes == 1,
			stats.astDiff,
			stats.diff,
			stats.gzDiff,
			stats.astSize,
			stats.size,
			stats.gzSize)
	}

	if p.astManifestPopulated {
		output.Println("\nInput AST Manifest:\ntoken,count")
		var lines []string
		for tok, count := range p.astManifest {
			if count > 0 {
				lines = append(lines, rhino.Token(tok).String()+","+strconv.Itoa(int(count)))
			}
		}
		slices.SortFunc(lines, jutil.StringCompareTo)
		for _, line := range lines {
			output.Println(line)
		}
	}

	output.PrintlnEmpty()
	// this.output can be System.out, so don't close it to not lose subsequent
	// error messages. Flush to ensure that you will see the tracer report.
	output.Flush()
}

// PerformanceTracker_Stats contains statistics about a pass run, such as
// running time, size changes, etc
type PerformanceTracker_Stats struct {
	pass      string
	isOneTime bool
	runtime   int64
	allocMem  int
	runs      int
	changes   int
	diff      int
	gzDiff    int
	size      int
	gzSize    int
	astDiff   int
	astSize   int
}

// Java: PerformanceTracker.Stats(String, boolean)
func NewPerformanceTracker_Stats(pass string, iot bool) *PerformanceTracker_Stats {
	s := &PerformanceTracker_Stats{}
	s.runtime = 0
	s.allocMem = 0
	s.runs = 0
	s.changes = 0
	s.diff = 0
	s.gzDiff = 0
	s.size = 0
	s.gzSize = 0
	s.astDiff = 0
	s.astSize = 0
	s.pass = pass
	s.isOneTime = iot
	return s
}
