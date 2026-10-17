// Port of com/google/javascript/jscomp/SubstituteEs6SyntaxTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/SubstituteEs6SyntaxTest.java)

/*
 * Copyright 2015 The Closure Compiler Authors.
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
)

type SubstituteEs6SyntaxTest struct {
	CompilerTestCaseBase
}

// Java: SubstituteEs6SyntaxTest()
func NewSubstituteEs6SyntaxTest() *SubstituteEs6SyntaxTest {
	c := &SubstituteEs6SyntaxTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestSubstituteEs6SyntaxTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewSubstituteEs6SyntaxTest() })
}

// Java: SubstituteEs6SyntaxTest.getProcessor(Compiler)
func (c *SubstituteEs6SyntaxTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewSubstituteEs6Syntax(compiler)
}

// Java: SubstituteEs6SyntaxTest.nullishCoalesceFunctions()
func (c *SubstituteEs6SyntaxTest) TestNullishCoalesceFunctions() {
	c.Test("()=>{ return x ?? y; }", "()=> x??y")
	c.TestSame("(function() { this.x = a??b; })")
}

// Java: SubstituteEs6SyntaxTest.nullishCoalesceObjectPattern()
func (c *SubstituteEs6SyntaxTest) TestNullishCoalesceObjectPattern() {
	// Tree comparisons don't fail on node property differences, so compare as strings instead.
	c.DisableCompareAsTree()
	c.Test("const {x:x}=a??b", "const {x}=a??b")
	c.TestSame("const {x:y}=a??b")
}

// Java: SubstituteEs6SyntaxTest.testArrowFunctions()
func (c *SubstituteEs6SyntaxTest) TestArrowFunctions() {
	c.TestSame("function f() {}")
	c.TestSame("(function() { this.x = 5; })")
	c.TestSame("(function() { return arguments[0]; })")
	c.TestSame("(function() { return ()=>this; })")
	c.TestSame("(function() { return ()=>arguments[0]; })")
	c.TestSame("(function() { x++; return 5; })")
	c.Test("()=>{ return 5; }", "()=>5")
	c.Test("()=>{ return; }", "()=>undefined")
	c.Test("(x)=>{ return x+1 }", "(x) => x+1")
	c.Test("(x)=>{ return x++,5; }", "(x) => (x++,5)")
}

// Java: SubstituteEs6SyntaxTest.testObjectPattern()
func (c *SubstituteEs6SyntaxTest) TestObjectPattern() {
	// Tree comparisons don't fail on node property differences, so compare as strings instead.
	c.DisableCompareAsTree()
	c.Test("const {x:x}=y", "const {x}=y")
	c.TestSame("const {x:y}=z")
	c.TestSame("const {[x]:x}=y")
	c.TestSame("const {[\"x\"]:x}=y")
}

// Java: SubstituteEs6SyntaxTest.testObjectLiteral()
func (c *SubstituteEs6SyntaxTest) TestObjectLiteral() {
	// Tree comparisons don't fail on node property differences, so compare as strings instead.
	c.DisableCompareAsTree()
	c.TestSame("const o={x:y}")
	c.TestSame("const o={[x]:x}")
	c.TestSame("const o={[\"x\"]:x}")
	c.Test("const o={x:x}", "const o={x}")
}
