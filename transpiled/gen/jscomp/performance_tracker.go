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
// Port of com/google/javascript/jscomp/PerformanceTracker.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_performance_tracker_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_performance_tracker_2 j2gort.Ext_java_util_function_ToLongFunction

func init() {
	lambdaSingleton_performance_tracker_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_53 any) any {
		entry := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](lp_53)
		_ = entry
		var key string = j2gort.CastString(entry.GetKey())
		_ = key
		var stats *PerformanceTracker_Stats = j2gort.CastPtr[PerformanceTracker_Stats](entry.GetValue())
		_ = stats
		return j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("%s,%d,%d,%d,%d,%d,%d,%d", []any{j2gort.BoxString(key), j2gort.ValueOfLong(stats.runtime), j2gort.ValueOfInt(stats.allocMem), j2gort.ValueOfInt(stats.runs), j2gort.ValueOfInt(stats.changes), j2gort.ValueOfInt(stats.astDiff), j2gort.ValueOfInt(stats.diff), j2gort.ValueOfInt(stats.gzDiff)}))
	})
	lambdaSingleton_performance_tracker_2 = j2gort.NewLambda_Ext_java_util_function_ToLongFunction(func(lp_54 any) int64 {
		e := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](lp_54)
		_ = e
		return j2gort.CastPtr[PerformanceTracker_Stats](e.GetValue()).runtime
	})
}

// Java class: com.google.javascript.jscomp.PerformanceTracker
type PerformanceTracker struct {
	codeChange                     *RecentChange                                     `gson:"codeChange"`
	currentPass                    j2gort.Deque[*PerformanceTracker_Stats]           `gson:"currentPass"`
	passSummary                    j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"passSummary"`
	astManifest                    []int32                                           `gson:"astManifest"`
	log                            []*PerformanceTracker_Stats                       `gson:"log"`
	pruningAnalysisSummary         string                                            `gson:"pruningAnalysisSummary"`
	disambiguatePropertiesSummary  string                                            `gson:"disambiguatePropertiesSummary"`
	ambiguatePropertiesSummary     string                                            `gson:"ambiguatePropertiesSummary"`
	startTime                      int64                                             `gson:"startTime"`
	endTime                        int64                                             `gson:"endTime"`
	externsRoot                    pkg_rhino.Node                                    `gson:"externsRoot"`
	jsRoot                         pkg_rhino.Node                                    `gson:"jsRoot"`
	initAstSize                    int32                                             `gson:"initAstSize"`
	initCodeSize                   int32                                             `gson:"initCodeSize"`
	initGzCodeSize                 int32                                             `gson:"initGzCodeSize"`
	passesRuntime                  int32                                             `gson:"passesRuntime"`
	maxMem                         int32                                             `gson:"maxMem"`
	runs                           int32                                             `gson:"runs"`
	changes                        int32                                             `gson:"changes"`
	loopRuns                       int32                                             `gson:"loopRuns"`
	loopChanges                    int32                                             `gson:"loopChanges"`
	jsLines                        int32                                             `gson:"jsLines"`
	jsSources                      int32                                             `gson:"jsSources"`
	externLines                    int32                                             `gson:"externLines"`
	externSources                  int32                                             `gson:"externSources"`
	transitiveSummaryLines         int32                                             `gson:"transitiveSummaryLines"`
	transitiveSummarySources       int32                                             `gson:"transitiveSummarySources"`
	transitiveSummaryLinesParsed   int32                                             `gson:"transitiveSummaryLinesParsed"`
	transitiveSummarySourcesParsed int32                                             `gson:"transitiveSummarySourcesParsed"`
	astSize                        int32                                             `gson:"astSize"`
	codeSize                       int32                                             `gson:"codeSize"`
	gzCodeSize                     int32                                             `gson:"gzCodeSize"`
	astDiff                        int32                                             `gson:"astDiff"`
	diff                           int32                                             `gson:"diff"`
	gzDiff                         int32                                             `gson:"gzDiff"`
	idHash_                        atomic.Int32
	mode                           CompilerOptions_TracerMode `gson:"mode"`
	prePruningRecorded             bool                       `gson:"prePruningRecorded"`
}

func (this_ *PerformanceTracker) BindSelf_(s any) {
	_ = s
}
func (this_ *PerformanceTracker) IsJavaNil_() bool { return this_ == nil }
func (*PerformanceTracker) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PerformanceTracker", "PerformanceTracker"
}
func (this_ *PerformanceTracker) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PerformanceTracker) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PerformanceTracker", "com.google.javascript.jscomp.PerformanceTracker", func(v_ any) bool {
		return j2gort.InstanceOf[*PerformanceTracker](v_)
	})
}

func StaticGet_PerformanceTracker_DEFAULT_WHEN_SIZE_UNTRACKED() int32 { return int32(-1) }

var static_PerformanceTracker_TOKEN_VALUES []pkg_rhino.Token

func StaticGet_PerformanceTracker_TOKEN_VALUES() []pkg_rhino.Token {
	clinit_PerformanceTracker()
	return static_PerformanceTracker_TOKEN_VALUES
}
func StaticSet_PerformanceTracker_TOKEN_VALUES(v []pkg_rhino.Token) {
	clinit_PerformanceTracker()
	static_PerformanceTracker_TOKEN_VALUES = v
}

var clinitDone_PerformanceTracker uint32
var clinitState_PerformanceTracker j2gort.ClinitState

func clinit_PerformanceTracker() {
	if atomic.LoadUint32(&clinitDone_PerformanceTracker) != 0 {
		return
	}
	clinitSlow_PerformanceTracker()
}

func clinitSlow_PerformanceTracker() {
	if !j2gort.EnterClinit(&clinitState_PerformanceTracker, &clinitDone_PerformanceTracker) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_PerformanceTracker, &clinitDone_PerformanceTracker)
	static_PerformanceTracker_TOKEN_VALUES = pkg_rhino.Token_Values()
}

func (this_ *PerformanceTracker) initFields_PerformanceTracker() {
	_ = this_
	this_.codeChange = NewRecentChange()
	this_.initAstSize = int32(-1)
	this_.initCodeSize = int32(-1)
	this_.initGzCodeSize = int32(-1)

	this_.astSize = int32(-1)
	this_.codeSize = int32(-1)
	this_.gzCodeSize = int32(-1)

	this_.currentPass = j2gort.Deque[*PerformanceTracker_Stats]{}
	this_.log = []*PerformanceTracker_Stats{}
	this_.pruningAnalysisSummary = "not executed"
	this_.disambiguatePropertiesSummary = "not executed"
	this_.ambiguatePropertiesSummary = "not executed"
}

// Java constructor: PerformanceTracker(com.google.javascript.rhino.Node externsRoot, com.google.javascript.rhino.Node jsRoot, com.google.javascript.jscomp.CompilerOptions.TracerMode mode)
func NewPerformanceTracker(externsRoot pkg_rhino.Node, jsRoot pkg_rhino.Node, mode CompilerOptions_TracerMode) *PerformanceTracker {
	clinit_PerformanceTracker()
	this_ := &PerformanceTracker{}
	InitPerformanceTracker(this_, externsRoot, jsRoot, mode)
	return this_
}

func InitPerformanceTracker(this_ *PerformanceTracker, externsRoot pkg_rhino.Node, jsRoot pkg_rhino.Node, mode CompilerOptions_TracerMode) {
	_ = externsRoot
	_ = jsRoot
	_ = mode
	this_.initFields_PerformanceTracker()
	if !(mode != StaticGet_CompilerOptions_TracerMode_OFF()) {
		j2gort.CheckArgumentStr(false, "PerformanceTracker can't work without tracer data.")
	}
	this_.startTime = javalang.System_CurrentTimeMillis()
	this_.externsRoot = externsRoot
	this_.jsRoot = jsRoot
	this_.mode = mode
}

// Java method: PerformanceTracker.getCodeChangeHandler()
func (this_ *PerformanceTracker) GetCodeChangeHandler() CodeChangeHandler {
	return ToIface_CodeChangeHandler(this_.codeChange)
}

// Java method: PerformanceTracker.recordPassStart(java.lang.String passName, boolean isOneTime)
func (this_ *PerformanceTracker) RecordPassStart(passName string, isOneTime bool) {
	_ = passName
	_ = isOneTime
	this_.currentPass.Push(NewPerformanceTracker_Stats(passName, isOneTime))
	if !(j2gort.CheckNotNullString(passName) == "toSource") {
		this_.codeChange.Reset()
	}
}

// Java method: PerformanceTracker.updateAfterDeserialize(com.google.javascript.rhino.Node jsRoot)
func (this_ *PerformanceTracker) UpdateAfterDeserialize(jsRoot pkg_rhino.Node) {
	_ = jsRoot
	this_.jsRoot = jsRoot
	if !(this_.TracksAstSize()) {
		return
	}
	this_.initAstSize = func() int32 {
		var asgn_1 int32 = NodeUtil_CountAstSize(this_.jsRoot)
		this_.astSize = asgn_1
		return asgn_1
	}()
	if !(this_.TracksSize()) {
		return
	}
	var estimator *PerformanceTrackerCodeSizeEstimator = PerformanceTrackerCodeSizeEstimator_Estimate(this_.jsRoot, this_.TracksGzSize())
	_ = estimator
	this_.initCodeSize = func() int32 { var asgn_2 int32 = estimator.GetCodeSize(); this_.codeSize = asgn_2; return asgn_2 }()
	if this_.TracksGzSize() {
		this_.initGzCodeSize = func() int32 {
			var asgn_3 int32 = estimator.GetZippedCodeSize()
			this_.gzCodeSize = asgn_3
			return asgn_3
		}()
	}
}

// Java method: PerformanceTracker.recordPassStop(java.lang.String passName, long runtime)
func (this_ *PerformanceTracker) RecordPassStop(passName string, runtime int64) {
	_ = passName
	_ = runtime
	var allocMem int32 = this_.priv_PerformanceTracker_getAllocatedMegabytes_4308()
	_ = allocMem
	var logStats *PerformanceTracker_Stats = this_.currentPass.Pop()
	_ = logStats
	{
		var iifeCond_1 bool
		{

			sr_4 := passName
			se_5 := logStats.pass
			iifeCond_1 = j2gort.CheckNotNullString(sr_4) == se_5
		}
		if !(iifeCond_1) {
			j2gort.CheckState(false)
		}
	}
	this_.log = append(this_.log, logStats)
	{
		r_6 := logStats
		r_6.runtime = runtime
	}
	{
		r_7 := logStats
		r_7.allocMem = allocMem
	}
	{
		r_8 := logStats
		r_8.runs = int32(1)
	}
	if this_.codeChange.HasCodeChanged() {
		{
			r_9 := logStats
			r_9.changes = int32(1)
		}
	}
	{
		var iifeCond_1 bool
		{

			sr_10 := passName
			se_11 := "parseInputs"
			iifeCond_1 = j2gort.CheckNotNullString(sr_10) == se_11
		}
		if iifeCond_1 {
			this_.priv_PerformanceTracker_recordParsingStop_4304(logStats)
		} else if this_.codeChange.HasCodeChanged() && this_.TracksAstSize() {
			this_.priv_PerformanceTracker_recordOtherPassStop_4305(logStats)
		}
	}
}

// Java method: PerformanceTracker.recordParsingStop(com.google.javascript.jscomp.PerformanceTracker.Stats logStats)
func (this_ *PerformanceTracker) priv_PerformanceTracker_recordParsingStop_4304(logStats *PerformanceTracker_Stats) {
	_ = logStats
	if this_.TracksAstManifest() {
		this_.priv_PerformanceTracker_populateAstManifest_4324()
	}
	this_.priv_PerformanceTracker_recordInputCount_4306()
	if !(this_.TracksAstSize()) {
		return
	}
	{
		r_14 := logStats
		r_14.astSize = func() int32 {
			var asgn_12 int32
			{
				var asgn_13 int32 = NodeUtil_CountAstSize(this_.jsRoot)
				this_.astSize = asgn_13
				asgn_12 = asgn_13
			}

			this_.initAstSize = asgn_12
			return asgn_12
		}()
	}
	if !(this_.TracksSize()) {
		return
	}
	var estimator *PerformanceTrackerCodeSizeEstimator = PerformanceTrackerCodeSizeEstimator_Estimate(this_.jsRoot, this_.TracksGzSize())
	_ = estimator
	{
		r_17 := logStats
		r_17.size = func() int32 {
			var asgn_15 int32
			{
				var asgn_16 int32 = estimator.GetCodeSize()
				this_.codeSize = asgn_16
				asgn_15 = asgn_16
			}
			this_.initCodeSize = asgn_15
			return asgn_15
		}()
	}
	if this_.TracksGzSize() {
		{
			r_20 := logStats
			r_20.gzSize = func() int32 {
				var asgn_18 int32
				{
					var asgn_19 int32 = estimator.GetZippedCodeSize()
					this_.gzCodeSize = asgn_19
					asgn_18 = asgn_19
				}

				this_.initGzCodeSize = asgn_18
				return asgn_18
			}()
		}
	}
}

// Java method: PerformanceTracker.recordOtherPassStop(com.google.javascript.jscomp.PerformanceTracker.Stats logStats)
func (this_ *PerformanceTracker) priv_PerformanceTracker_recordOtherPassStop_4305(logStats *PerformanceTracker_Stats) {
	_ = logStats
	var newSize int32 = NodeUtil_CountAstSize(this_.jsRoot)
	_ = newSize
	{
		r_21 := logStats
		r_21.astDiff = (this_.astSize - newSize)
	}
	this_.astSize = func() int32 {
		var asgn_22 int32 = newSize
		{
			r_23 := logStats
			r_23.astSize = asgn_22
		}
		return asgn_22
	}()
	if !(this_.TracksSize()) {
		return
	}
	var estimator *PerformanceTrackerCodeSizeEstimator = PerformanceTrackerCodeSizeEstimator_Estimate(this_.jsRoot, this_.TracksGzSize())
	_ = estimator
	newSize = estimator.GetCodeSize()
	{
		r_24 := logStats
		r_24.diff = (this_.codeSize - newSize)
	}
	this_.codeSize = func() int32 {
		var asgn_25 int32 = newSize
		{
			r_26 := logStats
			r_26.size = asgn_25
		}
		return asgn_25
	}()
	if this_.TracksGzSize() {
		newSize = estimator.GetZippedCodeSize()
		{
			r_27 := logStats
			r_27.gzDiff = (this_.gzCodeSize - newSize)
		}
		this_.gzCodeSize = func() int32 {
			var asgn_28 int32 = newSize
			{
				r_29 := logStats
				r_29.gzSize = asgn_28
			}
			return asgn_28
		}()
	}
}

// Java method: PerformanceTracker.recordInputCount()
func (this_ *PerformanceTracker) priv_PerformanceTracker_recordInputCount_4306() {
	{
		var n pkg_rhino.Node = this_.externsRoot.GetFirstChild()
		_ = n
		for ; n != 0; n = n.GetNext() {
			if NodeUtil_IsFromTypeSummary(n) {
				this_.transitiveSummarySources = (this_.transitiveSummarySources + int32(1))
				this_.transitiveSummaryLines = (this_.transitiveSummaryLines + NodeUtil_EstimateNumLines(n))
			} else {
				this_.externSources = (this_.externSources + int32(1))
				this_.externLines = (this_.externLines + NodeUtil_EstimateNumLines(n))
			}
		}
	}
	{
		var n_2 pkg_rhino.Node = this_.jsRoot.GetFirstChild()
		_ = n_2
		for ; n_2 != 0; n_2 = n_2.GetNext() {
			this_.jsSources = (this_.jsSources + int32(1))
			this_.jsLines = (this_.jsLines + NodeUtil_EstimateNumLines(n_2))
		}
	}
}

// Java method: PerformanceTracker.bytesToMB(long bytes)
func (this_ *PerformanceTracker) priv_PerformanceTracker_bytesToMB_4307(bytes int64) int32 {
	_ = bytes
	return int32((bytes / int64((int32(1024) * int32(1024)))))
}

// Java method: PerformanceTracker.getAllocatedMegabytes()
func (this_ *PerformanceTracker) priv_PerformanceTracker_getAllocatedMegabytes_4308() int32 {
	return this_.priv_PerformanceTracker_bytesToMB_4307((Platform_TotalMemory() - Platform_FreeMemory()))
}

// Java method: PerformanceTracker.tracksSize()
func (this_ *PerformanceTracker) TracksSize() bool {
	return ((this_.mode == StaticGet_CompilerOptions_TracerMode_RAW_SIZE()) || (this_.mode == StaticGet_CompilerOptions_TracerMode_ALL()))
}

// Java method: PerformanceTracker.tracksGzSize()
func (this_ *PerformanceTracker) TracksGzSize() bool {
	return (this_.mode == StaticGet_CompilerOptions_TracerMode_ALL())
}

// Java method: PerformanceTracker.tracksAstSize()
func (this_ *PerformanceTracker) TracksAstSize() bool {
	return (this_.mode != StaticGet_CompilerOptions_TracerMode_TIMING_ONLY())
}

// Java method: PerformanceTracker.tracksAstManifest()
func (this_ *PerformanceTracker) TracksAstManifest() bool {
	return this_.mode.IsOn()
}

// Java method: PerformanceTracker.getRuntime()
func (this_ *PerformanceTracker) GetRuntime() int32 {
	this_.priv_PerformanceTracker_calcTotalStats_4322()
	return this_.passesRuntime
}

// Java method: PerformanceTracker.getSize()
func (this_ *PerformanceTracker) GetSize() int32 {
	this_.priv_PerformanceTracker_calcTotalStats_4322()
	return this_.codeSize
}

// Java method: PerformanceTracker.getGzSize()
func (this_ *PerformanceTracker) GetGzSize() int32 {
	this_.priv_PerformanceTracker_calcTotalStats_4322()
	return this_.gzCodeSize
}

// Java method: PerformanceTracker.getAstSize()
func (this_ *PerformanceTracker) GetAstSize() int32 {
	this_.priv_PerformanceTracker_calcTotalStats_4322()
	return this_.astSize
}

// Java method: PerformanceTracker.getChanges()
func (this_ *PerformanceTracker) GetChanges() int32 {
	this_.priv_PerformanceTracker_calcTotalStats_4322()
	return this_.changes
}

// Java method: PerformanceTracker.getLoopChanges()
func (this_ *PerformanceTracker) GetLoopChanges() int32 {
	this_.priv_PerformanceTracker_calcTotalStats_4322()
	return this_.loopChanges
}

// Java method: PerformanceTracker.getRuns()
func (this_ *PerformanceTracker) GetRuns() int32 {
	this_.priv_PerformanceTracker_calcTotalStats_4322()
	return this_.runs
}

// Java method: PerformanceTracker.getLoopRuns()
func (this_ *PerformanceTracker) GetLoopRuns() int32 {
	this_.priv_PerformanceTracker_calcTotalStats_4322()
	return this_.loopRuns
}

// Java method: PerformanceTracker.getStats()
func (this_ *PerformanceTracker) GetStats() j2gort.Ext_com_google_common_collect_ImmutableMap {
	this_.priv_PerformanceTracker_calcTotalStats_4322()
	return this_.passSummary
}

// Java method: PerformanceTracker.calcTotalStats()
func (this_ *PerformanceTracker) priv_PerformanceTracker_calcTotalStats_4322() {
	if this_.passSummary != nil {
		return
	}
	this_.endTime = javalang.System_CurrentTimeMillis()
	this_.priv_PerformanceTracker_populatePassSummary_4323()
	{
		cur_30 := j2gort.NewMapValCursor(this_.passSummary)
		for cur_30.HasNext() {
			rawElem_31 := cur_30.Next()
			stats := j2gort.CastPtr[PerformanceTracker_Stats](rawElem_31)
			_ = stats
			this_.passesRuntime = int32((int64(this_.passesRuntime) + stats.runtime))
			this_.maxMem = max(this_.maxMem, stats.allocMem)
			this_.runs = (this_.runs + stats.runs)
			this_.changes = (this_.changes + stats.changes)
			if !(stats.isOneTime) {
				this_.loopRuns = (this_.loopRuns + stats.runs)
				this_.loopChanges = (this_.loopChanges + stats.changes)
			}
			this_.astDiff = (this_.astDiff + stats.astDiff)
			this_.diff = (this_.diff + stats.diff)
			this_.gzDiff = (this_.gzDiff + stats.gzDiff)
		}
	}
	if !(!(this_.TracksAstSize()) || (this_.initAstSize == (this_.astDiff + this_.astSize))) {
		j2gort.CheckState(false)
	}
	if !(!(this_.TracksSize()) || (this_.initCodeSize == (this_.diff + this_.codeSize))) {
		j2gort.CheckState(false)
	}
	if !(!(this_.TracksGzSize()) || (this_.initGzCodeSize == (this_.gzDiff + this_.gzCodeSize))) {
		j2gort.CheckState(false)
	}
}

// Java method: PerformanceTracker.populatePassSummary()
func (this_ *PerformanceTracker) priv_PerformanceTracker_populatePassSummary_4323() {
	var tmpPassSummary j2gort.OrderedMap[string, *PerformanceTracker_Stats] = j2gort.NewOrderedMap[string, *PerformanceTracker_Stats]()
	_ = tmpPassSummary
	for _, logStat := range this_.log {
		_ = logStat
		var passName string = logStat.pass
		_ = passName
		var entry *PerformanceTracker_Stats = j2gort.OrderedMapComputeIfAbsent(&tmpPassSummary, passName, nil, func(lp_38 string) *PerformanceTracker_Stats {
			k := lp_38
			_ = k
			return NewPerformanceTracker_Stats(k, logStat.isOneTime)
		})
		_ = entry
		r_39 := entry
		r_39.runtime = (r_39.runtime + logStat.runtime)
		{
			r_40 := entry
			r_40.allocMem = max(entry.allocMem, logStat.allocMem)
		}
		r_41 := entry
		r_41.runs = r_41.runs + 1
		r_42 := entry
		r_42.changes = (r_42.changes + logStat.changes)
		r_43 := entry
		r_43.astDiff = (r_43.astDiff + logStat.astDiff)
		r_44 := entry
		r_44.diff = (r_44.diff + logStat.diff)
		r_45 := entry
		r_45.gzDiff = (r_45.gzDiff + logStat.gzDiff)
	}

	this_.passSummary = tmpPassSummary.ToImmutableMap()
}

// Java method: PerformanceTracker.populateAstManifest()
func (this_ *PerformanceTracker) priv_PerformanceTracker_populateAstManifest_4324() {
	if this_.astManifest != nil {
		return
	}
	var counts []int32 = make([]int32, int32(len(StaticGet_PerformanceTracker_TOKEN_VALUES())))
	_ = counts
	PerformanceTracker_CountAstTokens(this_.jsRoot, counts)
	this_.astManifest = counts
}

// Java static method: PerformanceTracker.countAstTokens(com.google.javascript.rhino.Node n, int[] counts)
func PerformanceTracker_CountAstTokens(n pkg_rhino.Node, counts []int32) {
	clinit_PerformanceTracker()
	_ = n
	_ = counts
	a_46 := counts
	i_47 := (n.GetToken()).Ordinal()
	a_46[i_47] = a_46[i_47] + 1
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			PerformanceTracker_CountAstTokens(c, counts)
		}
	}
}

// Java method: PerformanceTracker.setPruningAnalysisSummary(java.lang.String summary)
func (this_ *PerformanceTracker) SetPruningAnalysisSummary(summary string) {
	_ = summary
	this_.pruningAnalysisSummary = summary
}

// Java method: PerformanceTracker.recordPrePruningInputCount(java.lang.Iterable<com.google.javascript.jscomp.CompilerInput> originalInputs, java.lang.Iterable<com.google.javascript.jscomp.CompilerInput> originalExterns)
func (this_ *PerformanceTracker) RecordPrePruningInputCount(originalInputs j2gort.Ext_java_lang_Iterable, originalExterns j2gort.Ext_java_lang_Iterable) {
	_ = originalInputs
	_ = originalExterns
	{
		cur_48 := j2gort.NewIterCursor(j2gort.CheckNotNull(originalExterns))
		for cur_48.HasNext() {
			rawElem_49 := cur_48.Next()
			input := j2gort.CastPtr[CompilerInput](rawElem_49)
			_ = input
			var name string = input.GetSourceFile().GetName()
			_ = name
			if javalang.StringEndsWith(j2gort.CheckNotNullString(name), ".i.js") {
				this_.transitiveSummarySourcesParsed = (this_.transitiveSummarySourcesParsed + int32(1))
				this_.transitiveSummaryLinesParsed = (this_.transitiveSummaryLinesParsed + input.GetSourceFile().GetNumLines())
			}
		}
	}
	{
		cur_50 := j2gort.NewIterCursor(j2gort.CheckNotNull(originalInputs))
		for cur_50.HasNext() {
			rawElem_51 := cur_50.Next()
			input_2 := j2gort.CastPtr[CompilerInput](rawElem_51)
			_ = input_2
			var name_2 string = input_2.GetSourceFile().GetName()
			_ = name_2
			if javalang.StringEndsWith(j2gort.CheckNotNullString(name_2), ".i.js") {
				this_.transitiveSummarySourcesParsed = (this_.transitiveSummarySourcesParsed + int32(1))
				this_.transitiveSummaryLinesParsed = (this_.transitiveSummaryLinesParsed + input_2.GetSourceFile().GetNumLines())
			}
		}
	}
	this_.prePruningRecorded = true
}

// Java method: PerformanceTracker.setDisambiguatePropertiesSummary(java.lang.String summary)
func (this_ *PerformanceTracker) SetDisambiguatePropertiesSummary(summary string) {
	_ = summary
	this_.disambiguatePropertiesSummary = summary
}

// Java method: PerformanceTracker.setAmbiguatePropertiesSummary(java.lang.String summary)
func (this_ *PerformanceTracker) SetAmbiguatePropertiesSummary(summary string) {
	_ = summary
	this_.ambiguatePropertiesSummary = summary
}

// Java method: PerformanceTracker.outputTracerReport(java.io.PrintStream output)
func (this_ *PerformanceTracker) OutputTracerReport(output *j2gort.Ext_java_io_PrintStream) {
	_ = output
	JvmMetrics_MaybeWriteJvmMetrics(output, "verbose:pretty:all")
	this_.priv_PerformanceTracker_calcTotalStats_4322()
	output.Println_1_string(pkg_base.JSCompStrings_Lines_1_slice_string([]string{"", "TOTAL:", ("Start time(ms): " + j2gort.StrValInt64(this_.startTime)), ("End time(ms): " + j2gort.StrValInt64(this_.endTime)), ("Wall time(ms): " + j2gort.StrValInt64((this_.endTime - this_.startTime))), ("Passes runtime(ms): " + j2gort.StrValInt32(this_.passesRuntime)), ("Max mem usage (measured after each pass)(MB): " + j2gort.StrValInt32(this_.maxMem)), ("#Runs: " + j2gort.StrValInt32(this_.runs)), ("#Changing runs: " + j2gort.StrValInt32(this_.changes)), ("#Loopable runs: " + j2gort.StrValInt32(this_.loopRuns)), ("#Changing loopable runs: " + j2gort.StrValInt32(this_.loopChanges)), ("Estimated AST reduction(#nodes): " + j2gort.StrValInt32(this_.astDiff)), ("Estimated Reduction(bytes): " + j2gort.StrValInt32(this_.diff)), ("Estimated GzReduction(bytes): " + j2gort.StrValInt32(this_.gzDiff)), ("Estimated AST size(#nodes): " + j2gort.StrValInt32(this_.astSize)), ("Estimated Size(bytes): " + j2gort.StrValInt32(this_.codeSize)), ("Estimated GzSize(bytes): " + j2gort.StrValInt32(this_.gzCodeSize)), "", ("DisambiguateProperties: " + j2gort.StrValString(this_.disambiguatePropertiesSummary)), ("AmbiguateProperties: " + j2gort.StrValString(this_.ambiguatePropertiesSummary))}))
	output.Println_1_string(pkg_base.JSCompStrings_Lines_1_slice_string([]string{"", "Inputs:", ("JS lines:   " + j2gort.StrValInt32(this_.jsLines)), ("JS sources: " + j2gort.StrValInt32(this_.jsSources)), ("Extern lines:   " + j2gort.StrValInt32(this_.externLines)), ("Extern sources: " + j2gort.StrValInt32(this_.externSources)), ("Type summary lines (raw input):   " + j2gort.StrValInt32(func() int32 {
		if this_.prePruningRecorded {
			return this_.transitiveSummaryLinesParsed
		}
		return this_.transitiveSummaryLines
	}())), ("Type summary lines (post-pruning):   " + j2gort.StrValInt32(this_.transitiveSummaryLines)), ("Type summary sources (raw input):   " + j2gort.StrValInt32(func() int32 {
		if this_.prePruningRecorded {
			return this_.transitiveSummarySourcesParsed
		}
		return this_.transitiveSummarySources
	}())), ("Type summary sources (post-pruning):   " + j2gort.StrValInt32(this_.transitiveSummarySources)), "", ("Dependency pruning analysis: " + j2gort.StrValString(this_.pruningAnalysisSummary))}))
	output.Println_1_string(pkg_base.JSCompStrings_Lines_1_slice_string([]string{"", "Summary:", "pass,runtime,allocMem,runs,changingRuns,astReduction,reduction,gzReduction"}))
	j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](this_.passSummary.EntrySet()).Stream().Sorted_1(j2gort.Ext_java_util_Comparator_ComparingLong(lambdaSingleton_performance_tracker_2)).Map(lambdaSingleton_performance_tracker_1).ForEach(func() j2gort.Ext_java_util_function_Consumer {
		bound_52 := output
		return j2gort.NewLambda_Ext_java_util_function_Consumer(func(mrp0_ any) { bound_52.Println_1_string(j2gort.CastString(mrp0_)) })
	}())
	output.Println_1_string(pkg_base.JSCompStrings_Lines_1_slice_string([]string{"", "Log:", "pass,runtime,allocMem,codeChanged,astReduction,reduction,gzReduction,astSize,size,gzSize"}))
	for _, stats := range this_.log {
		_ = stats
		output.Printf_2("%s,%d,%d,%b,%d,%d,%d,%d,%d,%d\n", []any{j2gort.BoxString(stats.pass), j2gort.ValueOfLong(stats.runtime), j2gort.ValueOfInt(stats.allocMem), j2gort.ValueOfBoolean((stats.changes == int32(1))), j2gort.ValueOfInt(stats.astDiff), j2gort.ValueOfInt(stats.diff), j2gort.ValueOfInt(stats.gzDiff), j2gort.ValueOfInt(stats.astSize), j2gort.ValueOfInt(stats.size), j2gort.ValueOfInt(stats.gzSize)})
	}

	if this_.astManifest != nil {
		output.Println_1_string(pkg_base.JSCompStrings_Lines_1_slice_string([]string{"", "Input AST Manifest:", "token,count"}))
		var manifestLines []string = []string{}
		_ = manifestLines
		{
			var i int32 = int32(0)
			_ = i
			for ; i < int32(len(StaticGet_PerformanceTracker_TOKEN_VALUES())); i = i + 1 {
				var count int32 = this_.astManifest[i]
				_ = count
				if count > int32(0) {
					manifestLines = append(manifestLines, j2gort.StrVal(StaticGet_PerformanceTracker_TOKEN_VALUES()[i])+","+j2gort.StrValInt32(count))
				}
			}
		}
		j2gort.SliceSortComparator(manifestLines, nil)
		{
			var i_2 int32 = int32(0)
			_ = i_2
			var len_ int32 = int32(len(manifestLines))
			_ = len_
			for ; i_2 < len_; i_2 = i_2 + 1 {
				output.Println_1_string(manifestLines[i_2])
			}
		}
	}
	output.Println_0()
	output.Flush()
}

// Java class: com.google.javascript.jscomp.PerformanceTracker.Stats
type PerformanceTracker_Stats struct {
	pass      string `gson:"pass"`
	runtime   int64  `gson:"runtime"`
	allocMem  int32  `gson:"allocMem"`
	runs      int32  `gson:"runs"`
	changes   int32  `gson:"changes"`
	diff      int32  `gson:"diff"`
	gzDiff    int32  `gson:"gzDiff"`
	size      int32  `gson:"size"`
	gzSize    int32  `gson:"gzSize"`
	astDiff   int32  `gson:"astDiff"`
	astSize   int32  `gson:"astSize"`
	idHash_   atomic.Int32
	isOneTime bool `gson:"isOneTime"`
}

func (this_ *PerformanceTracker_Stats) BindSelf_(s any) {
	_ = s
}
func (this_ *PerformanceTracker_Stats) IsJavaNil_() bool { return this_ == nil }
func (*PerformanceTracker_Stats) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PerformanceTracker$Stats", "Stats"
}
func (this_ *PerformanceTracker_Stats) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PerformanceTracker_Stats) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PerformanceTracker.Stats", "com.google.javascript.jscomp.PerformanceTracker$Stats", func(v_ any) bool {
		return j2gort.InstanceOf[*PerformanceTracker_Stats](v_)
	})
}

func (this_ *PerformanceTracker_Stats) initFields_PerformanceTracker_Stats() {
	_ = this_
	this_.pass = j2gort.NullString

}

// Java constructor: Stats(java.lang.String pass, boolean iot)
func NewPerformanceTracker_Stats(pass string, iot bool) *PerformanceTracker_Stats {
	this_ := &PerformanceTracker_Stats{}
	InitPerformanceTracker_Stats(this_, pass, iot)
	return this_
}

func InitPerformanceTracker_Stats(this_ *PerformanceTracker_Stats, pass string, iot bool) {
	_ = pass
	_ = iot
	this_.initFields_PerformanceTracker_Stats()
	this_.pass = pass
	this_.isOneTime = iot
}
