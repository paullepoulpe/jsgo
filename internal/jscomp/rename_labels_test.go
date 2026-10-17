// Port of com/google/javascript/jscomp/RenameLabelsTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/RenameLabelsTest.java)

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
	"testing"
)

type RenameLabelsTest struct {
	CompilerTestCaseBase
}

// Java: RenameLabelsTest()
func NewRenameLabelsTest() *RenameLabelsTest {
	c := &RenameLabelsTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestRenameLabelsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRenameLabelsTest() })
}

// Java: RenameLabelsTest.getProcessor(Compiler)
func (c *RenameLabelsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewRenameLabels(compiler)
}

// Java: RenameLabelsTest.testRenameInFunction()
func (c *RenameLabelsTest) TestRenameInFunction() {
	c.Test("function x(){ Foo:a(); }", "function x(){ a(); }")

	c.Test("function x(){ Foo:{ a(); break Foo; } }", "function x(){ a:{ a(); break a; } }")

	c.Test(
		"function x() {\n"+
			"Foo:{\n"+
			"function goo() {\n"+
			"Foo: {\n"+
			"a();\n"+
			"break Foo;\n"+
			"}\n"+
			"}\n"+
			"}\n"+
			"}\n",
		"function x(){{function goo(){a:{ a(); break a; }}}}")

	c.Test(
		"function x() {\n"+
			"Foo:{\n"+
			"function goo() {\n"+
			"Foo: {\n"+
			"a();\n"+
			"break Foo;\n"+
			"}\n"+
			"}\n"+
			"break Foo;\n"+
			"}\n"+
			"}\n",
		"function x(){a:{function goo(){a:{ a(); break a; }} break a;}}")
}

// Java: RenameLabelsTest.testRenameForArrowFunction()
func (c *RenameLabelsTest) TestRenameForArrowFunction() {
	// remove label that is not referenced
	c.Test("() => { Foo:a(); } ", "() => {     a(); }")

	c.Test("Foo:() => { a(); }", "    () => { a(); }")

	// label is referenced
	c.Test("() => { Foo:{ a(); break Foo; } }", "() => {   a:{ a(); break   a; } }")
}

// Java: RenameLabelsTest.testRenameForOf()
func (c *RenameLabelsTest) TestRenameForOf() {
	c.Test(
		"loop:\n"+
			"for (let x of [1, 2, 3]) {\n"+
			"  if (x > 2) {\n"+
			"    break loop;\n"+
			"  }\n"+
			"}\n",
		"a:\n"+
			"for (let x of [1, 2, 3]) {\n"+
			"  if (x > 2) {\n"+
			"    break a;\n"+
			"  }\n"+
			"}\n")
}

// Java: RenameLabelsTest.testRenameGlobals()
func (c *RenameLabelsTest) TestRenameGlobals() {
	c.Test("Foo:{a();}", "a();")
	c.Test("Foo:{a(); break Foo;}", "a:{a(); break a;}")
	c.Test("Foo:{Goo:a(); break Foo;}", "a:{a(); break a;}")
	c.Test(
		"Foo:{Goo:while(1){a(); continue Goo; break Foo;}}",
		"a:{b:while(1){a(); continue b;break a;}}")
	c.Test(
		"Foo:Goo:while(1){a(); continue Goo; break Foo;}",
		"a:b:while(1){a(); continue b;break a;}")

	c.Test("Foo:Bar:X:{ break Bar; }", "b:{ break b; }")
	c.Test("Foo:Bar:X:{ break Bar; break X; }", "b:c:{ break b; break c;}")
	c.Test("Foo:Bar:X:{ break Bar; break Foo; }", "a:b:{ break b; break a;}")

	c.Test("Foo:while (1){a(); break;}", "while (1){a(); break;}")

	// Remove label that is not referenced.
	c.Test("Foo:{a(); while (1) break;}", "a(); while (1) break;")
}

// Java: RenameLabelsTest.testRenameReused()
func (c *RenameLabelsTest) TestRenameReused() {
	c.Test("foo:{break foo}; foo:{break foo}", "a:{break a};a:{break a}")
}
