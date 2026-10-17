// Port of com/google/javascript/jscomp/DataFlowAnalysisTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/DataFlowAnalysisTest.java): the tests for computeEscaped.

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

// Go-only: test-helper — DataFlowAnalysisTest's computeEscaped tests, kept
// in a file of their own (they go through LiveVariablesAnalysis and
// GoogleCodingConvention).

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// DataFlowAnalysisEscapedTest holds DataFlowAnalysisTest's computeEscaped
// tests.
type DataFlowAnalysisEscapedTest struct{ DataFlowAnalysisTest }

func TestDataFlowAnalysisEscapedTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return &DataFlowAnalysisEscapedTest{} })
}

// tests for computeEscaped method

// Java: DataFlowAnalysisTest.testEscaped()
func (c *DataFlowAnalysisEscapedTest) TestEscaped() {
	AssertThat(
		c.computeEscapedLocals(
			`function f() {
    var x = 0;
    setTimeout(function() { x++; });
    alert(x);
}
`).Size()).
		IsEqualTo(1)
	AssertThat(c.computeEscapedLocals("function f() {var _x}").Size()).IsEqualTo(1)
	AssertThat(c.computeEscapedLocals("function f() {try{} catch(e){}}").Size()).IsEqualTo(1)
}

// Java: DataFlowAnalysisTest.testEscapedFunctionLayered()
func (c *DataFlowAnalysisEscapedTest) TestEscapedFunctionLayered() {
	AssertThat(
		c.computeEscapedLocals(
			`function f() {
    function ff() {
        var x = 0;
        setTimeout(function() { x++; });
        alert(x);
    }
}
`).IsEmpty()).
		IsTrue()
}

// Java: DataFlowAnalysisTest.testEscapedLetConstSimple()
func (c *DataFlowAnalysisEscapedTest) TestEscapedLetConstSimple() {
	AssertThat(c.computeEscapedLocals("function f() { let x = 0; x ++; x; }").IsEmpty()).IsTrue()
}

// Java: DataFlowAnalysisTest.testEscapedFunctionAssignment()
func (c *DataFlowAnalysisEscapedTest) TestEscapedFunctionAssignment() {
	AssertThat(c.computeEscapedLocals("function f() {var x = function () { return 1; }; }").IsEmpty()).
		IsTrue()
	AssertThat(c.computeEscapedLocals("function f() {var x = function (y) { return y; }; }").IsEmpty()).
		IsTrue()
	AssertThat(c.computeEscapedLocals("function f() {let x = function () { return 1; }; }").IsEmpty()).
		IsTrue()
}

// Java: DataFlowAnalysisTest.testEscapedArrowFunction()
func (c *DataFlowAnalysisEscapedTest) TestEscapedArrowFunction() {
	// When the body of the arrow fn is analyzed, x is considered an escaped
	// var. When the outer block containing "const value ..." is analyzed, 'x'
	// is not considered an escaped var
	AssertThat(
		c.computeEscapedLocals(
			`function f() {const value = () => {
    var x = 0;
    setTimeout(function() { x++; });
    alert(x);
 };}
`).IsEmpty()).
		IsTrue()
}

// Java: DataFlowAnalysisTest.testEscapedInstanceClassField()
func (c *DataFlowAnalysisEscapedTest) TestEscapedInstanceClassField() {
	AssertThat(
		c.computeEscapedLocals(
			`function f() {
  let s = 1;
  class Box {
    html = s;
  }
  return new Box();
}
`).Size()).
		IsEqualTo(1)
}

// Java: DataFlowAnalysisTest.testEscapedComputedInstanceClassField()
func (c *DataFlowAnalysisEscapedTest) TestEscapedComputedInstanceClassField() {
	AssertThat(
		c.computeEscapedLocals(
			`function f() {
  let s = 1;
  class Box {
    ['html'] = s;
  }
  return new Box();
}
`).Size()).
		IsEqualTo(1)
}

// Java: DataFlowAnalysisTest.testNotEscapedStaticClassField()
func (c *DataFlowAnalysisEscapedTest) TestNotEscapedStaticClassField() {
	AssertThat(
		c.computeEscapedLocals(
			`function f() {
  let s = 1;
  class Box {
    static html = s;
  }
  return Box;
}
`).IsEmpty()).
		IsTrue()
}

// Java: DataFlowAnalysisTest.testNotEscapedStaticComputedClassField()
func (c *DataFlowAnalysisEscapedTest) TestNotEscapedStaticComputedClassField() {
	AssertThat(
		c.computeEscapedLocals(
			`function f() {
  let s = 1;
  class Box {
    static ['html'] = s;
  }
  return Box;
}
`).IsEmpty()).
		IsTrue()
}

// Java: DataFlowAnalysisTest.testNotEscapedComputedClassFieldKey()
func (c *DataFlowAnalysisEscapedTest) TestNotEscapedComputedClassFieldKey() {
	AssertThat(
		c.computeEscapedLocals(
			`function f() {
  let s = 'key';
  class Box {
    [s] = 1;
  }
  return new Box();
}
`).IsEmpty()).
		IsTrue()
}

// computeEscapedLocals: test computeEscaped helper method that returns the
// liveness analysis performed by the LiveVariablesAnalysis class
//
// Java: DataFlowAnalysisTest.computeEscapedLocals(String)
func (c *DataFlowAnalysisEscapedTest) computeEscapedLocals(src string) *jutil.ObjectLinkedHashSet[*Var] {
	// Set up compiler
	compiler := NewCompiler()
	options := NewCompilerOptions()
	options.SetCodingConvention(NewGoogleCodingConvention())
	compiler.InitOptions(options)
	compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)

	n := compiler.ParseTestCode(src).RemoveFirstChild()
	script := rhino.NewNodeWithChild(rhino.Token_SCRIPT, n)
	script.SetInputId(rhino.NewInputId("test"))
	AssertThat(compiler.GetErrors()).IsEmpty()

	// Create scopes
	scopeCreator := NewSyntacticScopeCreator(compiler)
	scope := scopeCreator.CreateScopeOfScope(n, Scope_CreateGlobalScope(script))
	var childScope *Scope
	if script.GetFirstChild().IsFunction() {
		childScope = scopeCreator.CreateScopeOfScope(NodeUtil_GetFunctionBody(n), scope)
	} else {
		childScope = nil
	}

	// Control flow graph
	cfg := ControlFlowAnalysis_Builder_().
		SetCompiler(compiler).
		SetCfgRoot(script).
		SetIncludeEdgeAnnotations(true).
		ComputeCfg()

	// All variables declared in function
	allVarsDeclaredInFunction := NodeUtil_GetAllVarsDeclaredInFunction(compiler, scopeCreator, scope)

	// Compute liveness of variables
	analysis := NewLiveVariablesAnalysis(
		cfg, scope, childScope, compiler, scopeCreator, allVarsDeclaredInFunction)
	analysis.Analyze()
	return analysis.GetEscapedLocals()
}
