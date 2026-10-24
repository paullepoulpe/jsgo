// Port of com/google/javascript/jscomp/ChangeVerifierTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/ChangeVerifierTest.java)

/*
 * Copyright 2017 The Closure Compiler Authors.
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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ChangeVerifierTest is the test class for ChangeVerifier.java. Test cases
// in this class roughly do the following:
//
// 1. Creates a compiler instance and parses a script using it.
//
// 2. Creates a new ChangeVerifier class with the same compiler.
//
// 3. Saves a clean snapshot of the parsed script from step#1 using the
// ChangeVerifier instance from step#2
//
// 4. Changes the script from step #1.
//
// 5. Either records/skips recording change to the script from #1 into the
// compiler.
//
// 6. Invokes changeVerifier (which already contains the clean snapshot) and
// gives it the changed script from step#1 to compare nodes.
type ChangeVerifierTest struct {
	t        *testing.T
	compiler *Compiler
}

func (c *ChangeVerifierTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *ChangeVerifierTest) JUnitBefore()           { c.Setup() }
func (c *ChangeVerifierTest) JUnitAfter()            {}

func TestChangeVerifierTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return &ChangeVerifierTest{} })
}

// Java: ChangeVerifierTest.setup()
func (c *ChangeVerifierTest) Setup() {
	c.compiler = nil
}

// Java: ChangeVerifierTest.testCorrectValidationOfScriptWithChangeAfterFunction()
func (c *ChangeVerifierTest) TestCorrectValidationOfScriptWithChangeAfterFunction() {
	script := c.parse("function A() {} if (0) { A(); }")
	jutil.CheckState(script.IsScript())

	verifier := NewChangeVerifier(c.compiler).Snapshot(script)

	changeVerifierTest_getCallNode(script).Detach()
	c.compiler.ReportChangeToChangeScope(script)

	// checks that a change was made and reported.
	verifier.CheckRecordedChangesWithPassName("test1", script)
}

// Java: ChangeVerifierTest.testChangeToScriptNotReported_newChild()
func (c *ChangeVerifierTest) TestChangeToScriptNotReported_newChild() {
	script := c.parse("function A() {} if (0) { A(); }")
	jutil.CheckState(script.IsScript())

	verifier := NewChangeVerifier(c.compiler).Snapshot(script)

	// no change, no problem
	verifier.CheckRecordedChangesWithPassName("test1", script)

	// add a statement, but don't report the change.
	script.AddChildToBack(rhino.IR_ExprResult(rhino.IR_NullNode()))

	e := AssertThrows[*jutil.IllegalStateException](func() { verifier.CheckRecordedChangesWithPassName("test2", script) })
	AssertThat(e.GetMessage()).Contains("changed scope not marked as changed")
	AssertThat(e.GetMessage()).Contains("differing child count")
}

// Java: ChangeVerifierTest.testChangeToScriptNotReported_changeToChild()
func (c *ChangeVerifierTest) TestChangeToScriptNotReported_changeToChild() {
	script := c.parse("function A() {} if (0) { A(); }")
	jutil.CheckState(script.IsScript())

	verifier := NewChangeVerifier(c.compiler).Snapshot(script)

	// no change, no problem
	verifier.CheckRecordedChangesWithPassName("test1", script)

	// modify the if condition, but don't report the change.
	ifStatement := script.GetLastChild()
	condition := NodeUtil_GetConditionExpression(ifStatement)
	condition.ReplaceWith(rhino.IR_Number(1))

	e := AssertThrows[*jutil.IllegalStateException](func() { verifier.CheckRecordedChangesWithPassName("test2", script) })
	AssertThat(e.GetMessage()).
		Contains(
			`test2: changed scope not marked as changed: SCRIPT: testcode.
shallow inequivalence
Before: NUMBER 0.0 1:20  [length: 1] [source_file: testcode]
After:  NUMBER 1.0 1:20  [length: 1] [source_file: testcode]

Ancestor nodes:
SCRIPT 1:0  [length: 31] [source_file: testcode] [input_id: InputId: testcode] [feature_set: []]
  IF 1:16  [length: 15] [source_file: testcode]
    NUMBER 1.0 1:20  [length: 1] [source_file: testcode]
`)
}

// Java: ChangeVerifierTest.testChangeToScriptNotReported_changeToFunctionName()
func (c *ChangeVerifierTest) TestChangeToScriptNotReported_changeToFunctionName() {
	script := c.parse("function A() {} if (0) { A(); }")
	jutil.CheckState(script.IsScript())

	verifier := NewChangeVerifier(c.compiler).Snapshot(script)

	// no change, no problem
	verifier.CheckRecordedChangesWithPassName("test1", script)

	// modify the if condition, but don't report the change.
	functionNode := script.GetFirstChild()
	functionNameNode := NodeUtil_GetNameNode(functionNode)
	functionNameNode.ReplaceWith(rhino.IR_Name("B"))

	e := AssertThrows[*jutil.IllegalStateException](func() { verifier.CheckRecordedChangesWithPassName("test2", script) })
	AssertThat(e.GetMessage()).Contains("changed scope not marked as changed")
	AssertThat(e.GetMessage()).
		Contains(
			`function name changed
Before: A
After:  B
`)
}

// Java: ChangeVerifierTest.testChangeToFunction_notReported()
func (c *ChangeVerifierTest) TestChangeToFunction_notReported() {
	script := c.parse("function A() {}")
	jutil.CheckState(script.IsScript())
	function := script.GetFirstChild()
	jutil.CheckState(function.IsFunction())

	verifier := NewChangeVerifier(c.compiler).Snapshot(script)

	// no change, no problem.
	verifier.CheckRecordedChangesWithPassName("test1", script)

	// add a statement, but don't report the change.
	function.AddChildToBack(rhino.IR_ExprResult(rhino.IR_NullNode()))

	e := AssertThrows[*jutil.IllegalStateException](func() { verifier.CheckRecordedChangesWithPassName("test2", script) })
	AssertThat(e.GetMessage()).Contains("changed scope not marked as changed")
}

// Java: ChangeVerifierTest.testChangeToFunction_newlyUnusedParameter()
func (c *ChangeVerifierTest) TestChangeToFunction_newlyUnusedParameter() {
	script := c.parse("function f(x) {}")
	jutil.CheckStateMsg1(script.IsScript(), "%s", script)
	function := script.GetFirstChild()
	jutil.CheckStateMsg1(function.IsFunction(), "%s", function)
	xParam := NodeUtil_GetFunctionParameters(function).GetOnlyChild()
	jutil.CheckStateMsg1(xParam.MatchesName("x"), "%s", xParam)

	verifier := NewChangeVerifier(c.compiler).Snapshot(script)

	// no change, no problem.
	verifier.CheckRecordedChangesWithPassName("test1", script)

	// mark parameter x as unused
	xParam.SetUnusedParameter(true)

	e := AssertThrows[*jutil.IllegalStateException](func() { verifier.CheckRecordedChangesWithPassName("test2", script) })
	AssertThat(e.GetMessage()).Contains("changed scope not marked as changed")
	AssertThat(e.GetMessage()).
		Contains(
			`shallow inequivalence
Before: NAME x 1:11  [length: 1] [source_file: testcode]
After:  NAME x 1:11  [length: 1] [source_file: testcode] [is_unused_parameter: 1]

Ancestor nodes:
FUNCTION f 1:0  [length: 16] [source_file: testcode]
  PARAM_LIST 1:10  [length: 3] [source_file: testcode]
    NAME x 1:11  [length: 1] [source_file: testcode] [is_unused_parameter: 1]
`)
}

// Java: ChangeVerifierTest.testChangeToArrowFunction_notReported()
func (c *ChangeVerifierTest) TestChangeToArrowFunction_notReported() {
	script := c.parse("() => {}")
	jutil.CheckState(script.IsScript())
	function := script.GetFirstFirstChild()
	jutil.CheckState(function.IsArrowFunction())

	verifier := NewChangeVerifier(c.compiler).Snapshot(script)

	// no change, no problem.
	verifier.CheckRecordedChangesWithPassName("test1", script)

	// add a statement, but don't report the change.
	function.AddChildToBack(rhino.IR_ExprResult(rhino.IR_NullNode()))
	e := AssertThrows[*jutil.IllegalStateException](func() { verifier.CheckRecordedChangesWithPassName("test2", script) })
	AssertThat(e.GetMessage()).Contains("changed scope not marked as changed")
}

// Java: ChangeVerifierTest.testChangeToArrowFunction_correctlyReportedChange()
func (c *ChangeVerifierTest) TestChangeToArrowFunction_correctlyReportedChange() {
	script := c.parse("() => {}")
	jutil.CheckState(script.IsScript())

	function := script.GetFirstFirstChild()
	jutil.CheckState(function.IsArrowFunction())

	verifier := NewChangeVerifier(c.compiler).Snapshot(script)

	// no change, no problem.
	verifier.CheckRecordedChangesWithPassName("test1", script)

	// add a statement, and report the change.
	function.AddChildToBack(rhino.IR_ExprResult(rhino.IR_NullNode()))
	c.compiler.ReportChangeToChangeScope(function)

	// checks that a change was made and recorded.
	verifier.CheckRecordedChangesWithPassName("test2", script)
}

// Java: ChangeVerifierTest.testDeletedFunction()
func (c *ChangeVerifierTest) TestDeletedFunction() {
	script := c.parse("function A() {}")

	jutil.CheckState(script.IsScript())

	verifier := NewChangeVerifier(c.compiler).Snapshot(script)

	// no change
	verifier.CheckRecordedChangesWithPassName("test1", script)

	// remove the function. report the change in the script but not the
	// function deletion.
	fnNode := script.GetFirstChild()
	fnNode.Detach()
	c.compiler.ReportChangeToChangeScope(script)

	e := AssertThrows[*jutil.IllegalStateException](func() { verifier.CheckRecordedChangesWithPassName("test2", script) })
	AssertThat(e.GetMessage()).Contains("deleted scope was not reported")

	// now try again after reporting the function deletion.
	c.compiler.ReportFunctionDeleted(fnNode)

	// no longer throws an exception.
	verifier.CheckRecordedChangesWithPassName("test2", script)
}

// Java: ChangeVerifierTest.testNotDeletedFunction()
func (c *ChangeVerifierTest) TestNotDeletedFunction() {
	script := c.parse("function A() {}")

	jutil.CheckState(script.IsScript())

	verifier := NewChangeVerifier(c.compiler).Snapshot(script)

	// no change
	verifier.CheckRecordedChangesWithPassName("test1", script)

	// mark the function deleted even though it's alive.
	fnNode := script.GetFirstChild()
	c.compiler.ReportFunctionDeleted(fnNode)

	e := AssertThrows[*jutil.IllegalStateException](func() { verifier.CheckRecordedChangesWithPassName("test2", script) })
	AssertThat(e.GetMessage()).Contains("existing scope is improperly marked as deleted")
}

// Java: ChangeVerifierTest.testChangeVerification()
func (c *ChangeVerifierTest) TestChangeVerification() {
	mainScript := c.parse("")

	verifier := NewChangeVerifier(c.compiler).Snapshot(mainScript)

	verifier.CheckRecordedChanges(mainScript)

	mainScript.AddChildToFront(rhino.IR_Function(rhino.IR_Name("A"), rhino.IR_ParamList(), rhino.IR_Block()))
	c.compiler.ReportChangeToChangeScope(mainScript)

	e := AssertThrows[*jutil.IllegalStateException](func() { verifier.CheckRecordedChangesWithPassName("test2", mainScript) })
	// ensure that e was thrown from the right code-path
	// especially important if it's something as frequent
	// as an IllegalArgumentException, etc.
	AssertThat(e.GetMessage()).Contains("new scope not explicitly marked as changed:")

	// works fine when the newly created function scope is marked as changed.
	fnNode := mainScript.GetFirstChild()
	c.compiler.ReportChangeToChangeScope(fnNode)
	verifier.CheckRecordedChanges(mainScript)
}

// parse initializes a new compiler, parses the script using it and returns
// the script node
//
// Java: ChangeVerifierTest.parse(String)
func (c *ChangeVerifierTest) parse(js string) rhino.Node {
	c.compiler = NewCompiler()
	c.compiler.InitCompilerOptionsIfTesting()
	n := c.compiler.ParseTestCode(js)
	AssertThat(c.compiler.GetErrors()).IsEmpty()
	return n
}

// changeVerifierTest_getCallNode performs a depth first search and returns
// the first call node it finds
//
// Java: ChangeVerifierTest.getCallNode(Node)
func changeVerifierTest_getCallNode(n rhino.Node) rhino.Node {
	if n.IsCall() {
		return n
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		result := changeVerifierTest_getCallNode(c)
		if !result.IsNil() {
			return result
		}
	}
	return rhino.Node{}
}
