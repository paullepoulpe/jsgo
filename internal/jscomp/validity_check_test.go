// Port of com/google/javascript/jscomp/ValidityCheckTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/ValidityCheckTest.java)

/*
 * Copyright 2009 The Closure Compiler Authors.
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

// ValidityCheckTest holds the unit tests of ValidityCheck.
type ValidityCheckTest struct {
	CompilerTestCaseBase

	// May be nil.
	otherPass CompilerPass
}

// Java: ValidityCheckTest()
func NewValidityCheckTest() *ValidityCheckTest {
	c := &ValidityCheckTest{}
	c.InitCompilerTestCaseBase(c)
	c.otherPass = nil
	return c
}

func TestValidityCheckTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewValidityCheckTest() })
}

// Java: ValidityCheckTest.setUp()
func (c *ValidityCheckTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.otherPass = nil
}

// Java: ValidityCheckTest.getProcessor(Compiler)
func (c *ValidityCheckTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
		c.otherPass.Process(externs, root)
		(NewValidityCheck(compiler)).Process(externs, root)
	})
}

// Java: ValidityCheckTest.testUnnormalizeNodeTypes()
func (c *ValidityCheckTest) TestUnnormalizeNodeTypes() {
	c.otherPass =
		CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
			var compiler AbstractCompiler = c.GetLastCompiler()
			script := root.GetFirstChild()
			ifNode := rhino.NewNodeWith2Children(rhino.Token_IF, rhino.NewNode(rhino.Token_TRUE), rhino.NewNode(rhino.Token_EMPTY))
			ifNode.SrcrefTree(script)
			root.GetFirstChild().AddChildToBack(ifNode)
			compiler.ReportChangeToEnclosingScope(script)
		})

	jutil.Try(func() {
		c.Test("var x = 3;", "var x=3;0;0")
		AssertWithMessage("Expected IllegalStateException").Fail()
	}, func(e jutil.IIllegalStateException) {
		AssertThat(e.GetMessage()).Contains("Expected BLOCK but was EMPTY")
	})
}

// Java: ValidityCheckTest.testUnnormalized()
func (c *ValidityCheckTest) TestUnnormalized() {
	c.otherPass =
		CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
			c.GetLastCompiler().SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
		})

	jutil.Try(func() {
		c.TestSame("while(1){}")
		AssertWithMessage("Expected RuntimeException").Fail()
	}, func(e jutil.IRuntimeException) {
		AssertThat(e.GetMessage()).Contains("Normalize constraints violated:\nWHILE node")
	})
}

// Java: ValidityCheckTest.testConstantAnnotationMismatch()
func (c *ValidityCheckTest) TestConstantAnnotationMismatch() {
	c.otherPass =
		CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
			var compiler AbstractCompiler = c.GetLastCompiler()
			script := root.GetFirstChild()
			name := rhino.Node_NewStringWithToken(rhino.Token_NAME, "x")
			name.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
			exprResult := rhino.NewNodeWithChild(rhino.Token_EXPR_RESULT, name).SrcrefTree(script)
			script.AddChildToBack(exprResult)
			compiler.ReportChangeToEnclosingScope(script)
			compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
		})

	jutil.Try(func() {
		c.Test("var x;", "var x; x;")
		AssertWithMessage("Expected RuntimeException").Fail()
	}, func(e jutil.IRuntimeException) {
		AssertThat(e.GetMessage()).
			Contains("The name x is not consistently annotated as constant.")
	})
}

// Java: ValidityCheckTest.testSourceInfoCheck()
func (c *ValidityCheckTest) TestSourceInfoCheck() {
	compiler := NewCompiler()
	script := compiler.ParseTestCode("var x = 3;")
	n := rhino.NewNodeWithChild(rhino.Token_EXPR_RESULT, rhino.Node_NewString("foo")) // Child missing source info
	n.Srcref(script.GetFirstChild())
	script.AddChildToBack(n)

	check := NewSourceInfoCheck(compiler)
	check.Process(nil, rhino.NewNodeWithChild(rhino.Token_ROOT, script))

	AssertThat(compiler.GetErrors()).HasSize(1)
	AssertThat(compiler.GetErrors()[0].Description()).
		Contains("No source location information associated with")
}

// Java: ValidityCheckTest.testErrorPass()
func (c *ValidityCheckTest) TestErrorPass() {
	compiler := NewCompiler()
	testError := DiagnosticType_Error("TEST_ERROR", "Test error description")
	pass := NewErrorPass(compiler, testError)
	script := compiler.ParseTestCode("var x = 1;")
	pass.Process(nil, script)
	AssertThat(compiler.GetErrors()).HasSize(1)
	AssertThat(compiler.GetErrors()[0].Type()).IsEqualTo(testError)
}
