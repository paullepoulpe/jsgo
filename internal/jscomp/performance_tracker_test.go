// Port of com/google/javascript/jscomp/PerformanceTrackerTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/PerformanceTrackerTest.java)

/*
 * Copyright 2013 The Closure Compiler Authors.
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
	"bytes"
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PerformanceTrackerTest holds the unit tests for PerformanceTracker.
type PerformanceTrackerTest struct {
	t               *testing.T
	emptyExternRoot rhino.Node
	emptyJsRoot     rhino.Node
}

func (c *PerformanceTrackerTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *PerformanceTrackerTest) JUnitBefore()           {}
func (c *PerformanceTrackerTest) JUnitAfter()            {}

// Java: PerformanceTrackerTest()
func NewPerformanceTrackerTest() *PerformanceTrackerTest {
	c := &PerformanceTrackerTest{}
	c.emptyExternRoot = rhino.NewNode(rhino.Token_BLOCK)
	c.emptyJsRoot = rhino.NewNode(rhino.Token_BLOCK)
	return c
}

func TestPerformanceTrackerTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewPerformanceTrackerTest() })
}

// Java: PerformanceTrackerTest.testStatsCalculation()
func (c *PerformanceTrackerTest) TestStatsCalculation() {
	tracker :=
		NewPerformanceTracker(c.emptyExternRoot, c.emptyJsRoot, CompilerOptions_TracerMode_ALL)
	handler := tracker.GetCodeChangeHandler()

	// It's sufficient for this test to assume that a single run of any pass
	// takes some fixed amount of time, say 5ms.
	passRuntime := 5

	tracker.RecordPassStart("noloopA", true)
	handler.ReportChange()
	tracker.RecordPassStop("noloopA", int64(passRuntime))

	tracker.RecordPassStart("noloopB", true)
	handler.ReportChange()
	tracker.RecordPassStop("noloopB", int64(passRuntime))

	tracker.RecordPassStart("loopA", false)
	handler.ReportChange()
	tracker.RecordPassStop("loopA", int64(passRuntime))

	tracker.RecordPassStart("loopA", false)
	tracker.RecordPassStop("loopA", int64(passRuntime))

	tracker.RecordPassStart("noloopB", true)
	handler.ReportChange()
	tracker.RecordPassStop("noloopB", int64(passRuntime))

	tracker.RecordPassStart("loopB", false)
	tracker.RecordPassStop("loopB", int64(passRuntime))

	tracker.RecordPassStart("noloopB", true)
	tracker.RecordPassStop("noloopB", int64(passRuntime))

	numRuns := tracker.GetRuns()

	AssertThat(numRuns).IsEqualTo(7)
	AssertThat(numRuns * passRuntime).IsEqualTo(tracker.GetRuntime())
	AssertThat(tracker.GetLoopRuns()).IsEqualTo(3)
	AssertThat(tracker.GetChanges()).IsEqualTo(4) /* reportChange was called 4 times */
	AssertThat(tracker.GetLoopChanges()).IsEqualTo(1)

	stats := tracker.GetStats()
	st := stats.Get("noloopA")
	AssertThat(st.runs).IsEqualTo(1)
	AssertThat(passRuntime).IsEqualTo(st.runtime)
	AssertThat(st.changes).IsEqualTo(1)

	st = stats.Get("noloopB")
	AssertThat(st.runs).IsEqualTo(3)
	AssertThat(3 * passRuntime).IsEqualTo(st.runtime)
	AssertThat(st.changes).IsEqualTo(2)

	st = stats.Get("loopA")
	AssertThat(st.runs).IsEqualTo(2)
	AssertThat(2 * passRuntime).IsEqualTo(st.runtime)
	AssertThat(st.changes).IsEqualTo(1)

	st = stats.Get("loopB")
	AssertThat(st.runs).IsEqualTo(1)
	AssertThat(passRuntime).IsEqualTo(st.runtime)
	AssertThat(st.changes).IsEqualTo(0)
}

// Java: PerformanceTrackerTest.testAstSummaryAndFormat()
func (c *PerformanceTrackerTest) TestAstSummaryAndFormat() {
	// Given
	main :=
		rhino.IR_Root(
			rhino.IR_ScriptOfNodes(
				rhino.IR_BlockOfNodes(
					rhino.IR_BlockOfNode(
						rhino.IR_ExprResult( //
							rhino.IR_Hook(rhino.IR_String("a"), rhino.IR_String("b"), rhino.IR_String("c")))),
					rhino.IR_Function(rhino.IR_Name("name"), rhino.IR_ParamList(), rhino.IR_Block()))))
	tracker :=
		NewPerformanceTracker(c.emptyExternRoot, main, CompilerOptions_TracerMode_TIMING_ONLY)

	// When
	tracker.RecordPassStart(PassNames_PARSE_INPUTS, true)
	tracker.RecordPassStop(PassNames_PARSE_INPUTS, 0) // This is what triggers the counting.
	report := PerformanceTrackerTest_ExtractReport(tracker)

	// Then
	AssertThat(report).
		ContainsMatch(
			`Input AST Manifest:
token,count
BLOCK,3
EXPR_RESULT,1
FUNCTION,1
HOOK,1
NAME,1
PARAM_LIST,1
ROOT,1
SCRIPT,1
STRINGLIT,3
`)
}

// Java: PerformanceTrackerTest.testOutputFormat()
//
// Go-only: language-adapter — the regex is compiled by jutil.Pattern, which
// has no DOTALL flag, so the Java pattern's two DOTALL `.*` (which match
// across lines) are written `(?:.|\n)*` (the report is ASCII with "\n" line
// separators).
func (c *PerformanceTrackerTest) TestOutputFormat() {
	tracker :=
		NewPerformanceTracker(c.emptyExternRoot, c.emptyJsRoot, CompilerOptions_TracerMode_ALL)
	report := PerformanceTrackerTest_ExtractReport(tracker)
	p :=
		`(?:.|\n)*TOTAL:
Start time\(ms\): [0-9]+
End time\(ms\): [0-9]+
Wall time\(ms\): [0-9]+
Passes runtime\(ms\): [0-9]+
Max mem usage \(measured after each pass\)\(MB\): -?[0-9]+
#Runs: [0-9]+
#Changing runs: [0-9]+
#Loopable runs: [0-9]+
#Changing loopable runs: [0-9]+
Estimated AST reduction\(#nodes\): [0-9]+
Estimated Reduction\(bytes\): [0-9]+
Estimated GzReduction\(bytes\): [0-9]+
Estimated AST size\(#nodes\): -?[0-9]+
Estimated Size\(bytes\): -?[0-9]+
Estimated GzSize\(bytes\): -?[0-9]+

DisambiguateProperties: not executed
AmbiguateProperties: not executed

Inputs:
JS lines:\s*[0-9]+
JS sources:\s*[0-9]+
Extern lines:\s*[0-9]+
Extern sources:\s*[0-9]+
Type summary lines \(raw input\):\s*[0-9]+
Type summary lines \(post-pruning\):\s*[0-9]+
Type summary sources \(raw input\):\s*[0-9]+
Type summary sources \(post-pruning\):\s*[0-9]+

Dependency pruning analysis: not executed

Summary:
pass,runtime,allocMem,runs,changingRuns,astReduction,reduction,gzReduction

Log:
pass,runtime,allocMem,codeChanged,astReduction,reduction,gzReduction,astSize,size,gzSize
(?:.|\n)*
`

	AssertThat(report).Matches(p)
}

// Java: PerformanceTrackerTest.testDistinguishBetweenSrcExternAndIjs()
func (c *PerformanceTrackerTest) TestDistinguishBetweenSrcExternAndIjs() {
	// Add extern .i.js and .js nodes
	c.emptyExternRoot.AddChildToFront(PerformanceTrackerTest_CreateExtern("externs1.js"))
	c.emptyExternRoot.AddChildToFront(PerformanceTrackerTest_CreateExtern("externs2.js"))
	c.emptyExternRoot.AddChildToFront(PerformanceTrackerTest_CreateExtern("externs3.js"))
	c.emptyExternRoot.AddChildToFront(PerformanceTrackerTest_CreateTypeSummary("type_dep.js.i.js"))
	c.emptyJsRoot.AddChildToFront(PerformanceTrackerTest_CreateScript("input1.js"))
	c.emptyJsRoot.AddChildToFront(PerformanceTrackerTest_CreateScript("input2.js"))
	tracker :=
		NewPerformanceTracker(c.emptyExternRoot, c.emptyJsRoot, CompilerOptions_TracerMode_ALL)

	tracker.RecordPassStart(PassNames_PARSE_INPUTS, true)
	tracker.RecordPassStop(PassNames_PARSE_INPUTS, 0)
	report := PerformanceTrackerTest_ExtractReport(tracker)

	p := `JS lines:\s* 0
JS sources:\s* 2
Extern lines:\s* 0
Extern sources:\s* 3
Type summary lines \(raw input\):\s* 0
Type summary lines \(post-pruning\):\s* 0
Type summary sources \(raw input\):\s* 1
Type summary sources \(post-pruning\):\s* 1
`

	AssertThat(report).ContainsMatch(p)
}

// Java: PerformanceTrackerTest.testRecordPrePruningInputCount_distinguishesPrunedAndNonPrunedIjs()
func (c *PerformanceTrackerTest) TestRecordPrePruningInputCount_distinguishesPrunedAndNonPrunedIjs() {
	// Create three type summary nodes and inputs, one of which will be pruned.
	typeSummary1SourceFile := SourceFile_FromCode("type1.js.i.js", "")
	typeSummary2SourceFile := SourceFile_FromCode("type2.js.i.js", "")
	typeSummaryPrunedSourceFile :=
		SourceFile_FromCode("type_depPruned.js.i.js", "a \n b \n c")
	typeSummary1Input := NewCompilerInputWithInputId(typeSummary1SourceFile, rhino.NewInputId("type1.js.i.js"))
	typeSummary2Input := NewCompilerInputWithInputId(typeSummary2SourceFile, rhino.NewInputId("type2.js.i.js"))
	typeSummaryPrunedInput :=
		NewCompilerInputWithInputId(typeSummaryPrunedSourceFile, rhino.NewInputId("pruned.js.i.js"))
	c.emptyExternRoot.AddChildToFront(PerformanceTrackerTest_CreateTypeSummary("type1.js.i.js"))
	c.emptyExternRoot.AddChildToFront(PerformanceTrackerTest_CreateTypeSummary("type2.js.i.js"))
	tracker :=
		NewPerformanceTracker(c.emptyExternRoot, c.emptyJsRoot, CompilerOptions_TracerMode_ALL)

	tracker.RecordPassStart(PassNames_PARSE_INPUTS, true)
	tracker.RecordPrePruningInputCount(
		slices.Values([]*CompilerInput{typeSummary1Input, typeSummaryPrunedInput}),
		slices.Values([]*CompilerInput{typeSummary2Input}))
	tracker.RecordPassStop(PassNames_PARSE_INPUTS, 0)
	report := PerformanceTrackerTest_ExtractReport(tracker)

	p := `Type summary lines \(raw input\):\s* 5
Type summary lines \(post-pruning\):\s* 0
Type summary sources \(raw input\):\s* 3
Type summary sources \(post-pruning\):\s* 2
`

	AssertThat(report).ContainsMatch(p)
}

// Java: PerformanceTrackerTest.createScript(String)
func PerformanceTrackerTest_CreateScript(sourceFileName string) rhino.Node {
	sourceFile := SourceFile_FromCode(sourceFileName, "")
	script := rhino.IR_Script()
	script.SetStaticSourceFile(sourceFile)
	return script
}

// Java: PerformanceTrackerTest.createTypeSummary(String)
func PerformanceTrackerTest_CreateTypeSummary(sourceFileName string) rhino.Node {
	script := PerformanceTrackerTest_CreateScript(sourceFileName)
	scriptInfo := rhino.JSDocInfo_Builder_()
	scriptInfo.RecordTypeSummary()
	script.SetJSDocInfo(scriptInfo.Build())
	return script
}

// Java: PerformanceTrackerTest.createExtern(String)
func PerformanceTrackerTest_CreateExtern(sourceFileName string) rhino.Node {
	script := PerformanceTrackerTest_CreateScript(sourceFileName)
	scriptInfo := rhino.JSDocInfo_Builder_()
	scriptInfo.RecordExterns()
	script.SetJSDocInfo(scriptInfo.Build())
	return script
}

// Java: PerformanceTrackerTest.extractReport(PerformanceTracker)
func PerformanceTrackerTest_ExtractReport(tracker *PerformanceTracker) string {
	var output bytes.Buffer
	outstream := jutil.NewPrintStream(&output)
	tracker.OutputTracerReport(outstream)
	outstream.Close()
	return jutil.NewStringFromUTF8Bytes(output.Bytes())
}
