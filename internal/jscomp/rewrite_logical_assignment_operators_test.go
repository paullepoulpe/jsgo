// Port of com/google/javascript/jscomp/RewriteLogicalAssignmentOperatorsTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/RewriteLogicalAssignmentOperatorsTest.java)

/*
 * Copyright 2021 The Closure Compiler Authors.
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

import "testing"

// RewriteLogicalAssignmentOperatorsTest has the test cases for the
// transpilation pass that replaces the logical assignment operators (`||=`,
// `&&=`, `??=`).
type RewriteLogicalAssignmentOperatorsTest struct {
	CompilerTestCaseBase
}

// Java: RewriteLogicalAssignmentOperatorsTest()
func NewRewriteLogicalAssignmentOperatorsTest() *RewriteLogicalAssignmentOperatorsTest {
	c := &RewriteLogicalAssignmentOperatorsTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestRewriteLogicalAssignmentOperatorsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRewriteLogicalAssignmentOperatorsTest() })
}

// Java: RewriteLogicalAssignmentOperatorsTest.setUp()
func (c *RewriteLogicalAssignmentOperatorsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeInfoValidation()
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
}

// Java: RewriteLogicalAssignmentOperatorsTest.getProcessor(Compiler)
func (c *RewriteLogicalAssignmentOperatorsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewRewriteLogicalAssignmentOperatorsPass(compiler)
}

// Java: RewriteLogicalAssignmentOperatorsTest.testSimple()
func (c *RewriteLogicalAssignmentOperatorsTest) TestSimple() {
	c.Test(c.Srcs("a ||= b"), c.Expected("a || (a = b);"))
	c.Test(c.Srcs("a &&= b"), c.Expected("a && (a = b);"))
	c.Test(c.Srcs("a ??= b"), c.Expected("a ?? (a = b);"))
}

// Java: RewriteLogicalAssignmentOperatorsTest.testPropertyReference()
func (c *RewriteLogicalAssignmentOperatorsTest) TestPropertyReference() {
	c.Test(
		c.Srcs("a.x ||= b"),
		c.Expected(
			`let $jscomp$logical$assign$tmpm1146332801$0;
($jscomp$logical$assign$tmpm1146332801$0 = a).x
   ||
($jscomp$logical$assign$tmpm1146332801$0.x = b);
`))
	c.Test(
		c.Srcs("a.foo &&= null"),
		c.Expected(
			`let $jscomp$logical$assign$tmpm1146332801$0;
($jscomp$logical$assign$tmpm1146332801$0 = a).foo
   &&
($jscomp$logical$assign$tmpm1146332801$0.foo = null);
`))
	c.Test(
		c.Srcs(
			`foo().x = null;
foo().x ??= y
`),
		c.Expected(
			`foo().x = null;
let $jscomp$logical$assign$tmpm1146332801$0;
($jscomp$logical$assign$tmpm1146332801$0 = foo()).x
   ??
($jscomp$logical$assign$tmpm1146332801$0.x = y);
`))
}

// Java: RewriteLogicalAssignmentOperatorsTest.testPropertyReferenceElement()
func (c *RewriteLogicalAssignmentOperatorsTest) TestPropertyReferenceElement() {
	c.Test(
		c.Srcs("a[x] ||= b"),
		c.Expected(
			`let $jscomp$logical$assign$tmpm1146332801$0;
let $jscomp$logical$assign$tmpindexm1146332801$0;
($jscomp$logical$assign$tmpm1146332801$0 = a)
[$jscomp$logical$assign$tmpindexm1146332801$0 = x]
   ||
($jscomp$logical$assign$tmpm1146332801$0
[$jscomp$logical$assign$tmpindexm1146332801$0] = b);
`))
	c.Test(
		c.Srcs("a[x + 5 + 's'] &&= b"),
		c.Expected(
			`let $jscomp$logical$assign$tmpm1146332801$0;
let $jscomp$logical$assign$tmpindexm1146332801$0;
($jscomp$logical$assign$tmpm1146332801$0 = a)
[$jscomp$logical$assign$tmpindexm1146332801$0 = (x + 5 + 's')]
   &&
($jscomp$logical$assign$tmpm1146332801$0
[$jscomp$logical$assign$tmpindexm1146332801$0] = b);
`))
	c.Test(
		c.Srcs("foo[x] ??= bar[y]"),
		c.Expected(
			`let $jscomp$logical$assign$tmpm1146332801$0;
let $jscomp$logical$assign$tmpindexm1146332801$0;
($jscomp$logical$assign$tmpm1146332801$0 = foo)
[$jscomp$logical$assign$tmpindexm1146332801$0 = x]
   ??
($jscomp$logical$assign$tmpm1146332801$0
[$jscomp$logical$assign$tmpindexm1146332801$0] = bar[y]);
`))
}

// Java: RewriteLogicalAssignmentOperatorsTest.testInExpressions()
func (c *RewriteLogicalAssignmentOperatorsTest) TestInExpressions() {
	c.Test(c.Srcs("a ?? (b ??= 1) ?? (c ??= 2)"), c.Expected("a ?? (b ?? (b = 1)) ?? (c ?? (c = 2));"))
	c.Test(
		c.Srcs("for (let x = foo(); !x.b; x.b ||= {}) {}"),
		c.Expected(
			`let $jscomp$logical$assign$tmpm1146332801$0;
for (let x = foo(); !x.b;
($jscomp$logical$assign$tmpm1146332801$0 = x).b
   ||
($jscomp$logical$assign$tmpm1146332801$0.b = {})
) {}
`))
	c.Test(
		c.Srcs("for (let x = foo; !x[b]; x[b] &&= {}) {}"),
		c.Expected(
			`let $jscomp$logical$assign$tmpm1146332801$0;
let $jscomp$logical$assign$tmpindexm1146332801$0;
for (let x = foo; !x[b];
($jscomp$logical$assign$tmpm1146332801$0 = x)
[$jscomp$logical$assign$tmpindexm1146332801$0 = b]
   &&
($jscomp$logical$assign$tmpm1146332801$0
[$jscomp$logical$assign$tmpindexm1146332801$0] = {})
) {}
`))
	c.Test(
		c.Srcs(
			`a = 4;
do {
  a = a - 1;
} while (a ||= null);
`),
		c.Expected(
			`a = 4;
do {
  a = a - 1;
} while (a || (a = null));
`))
	c.Test(
		c.Srcs(
			`foo().x = 4;
do {
  foo().x = foo().x - 1;
} while (foo().x ||= null);
`),
		c.Expected(
			`foo().x = 4;
let $jscomp$logical$assign$tmpm1146332801$0;
do {
  foo().x = foo().x - 1;
} while (($jscomp$logical$assign$tmpm1146332801$0 = foo()).x
          ||
         ($jscomp$logical$assign$tmpm1146332801$0.x = null));
`))
	c.Test(
		c.Srcs(
			`foo[x] = 4;
do {
  foo[x] = foo[x] - 1;
} while (foo[x] &&= null);
`),
		c.Expected(
			`foo[x] = 4;
let $jscomp$logical$assign$tmpm1146332801$0;
let $jscomp$logical$assign$tmpindexm1146332801$0;
do {
  foo[x] = foo[x] - 1;
} while (($jscomp$logical$assign$tmpm1146332801$0 = foo)
         [$jscomp$logical$assign$tmpindexm1146332801$0 = x]
          &&
         ($jscomp$logical$assign$tmpm1146332801$0
         [$jscomp$logical$assign$tmpindexm1146332801$0] = null));
`))
}

// Java: RewriteLogicalAssignmentOperatorsTest.testNestedPropertyReference()
func (c *RewriteLogicalAssignmentOperatorsTest) TestNestedPropertyReference() {
	c.Test(
		c.Srcs("a ||= (b &&= (c ??= d));"), //
		c.Expected("a || (a = (b && (b = (c ?? (c = d)))));"))

	c.Test(
		c.Srcs("a.x ??= (b.x ||= {});"),
		c.Expected(
			`let $jscomp$logical$assign$tmpm1146332801$0;
let $jscomp$logical$assign$tmpm1146332801$1;
(($jscomp$logical$assign$tmpm1146332801$1 = a).x
   ??
($jscomp$logical$assign$tmpm1146332801$1.x =
($jscomp$logical$assign$tmpm1146332801$0 = b).x
   ||
($jscomp$logical$assign$tmpm1146332801$0.x = {})));
`))

	c.Test(
		c.Srcs("a[x] &&= (b[x] ??= a[x]);"),
		c.Expected(
			`let $jscomp$logical$assign$tmpm1146332801$0;
let $jscomp$logical$assign$tmpindexm1146332801$0;
let $jscomp$logical$assign$tmpm1146332801$1;
let $jscomp$logical$assign$tmpindexm1146332801$1;
($jscomp$logical$assign$tmpm1146332801$1 = a)
[$jscomp$logical$assign$tmpindexm1146332801$1 = x]
   &&
($jscomp$logical$assign$tmpm1146332801$1
[$jscomp$logical$assign$tmpindexm1146332801$1] =
($jscomp$logical$assign$tmpm1146332801$0 = b)
[$jscomp$logical$assign$tmpindexm1146332801$0 = x]
   ??
($jscomp$logical$assign$tmpm1146332801$0
[$jscomp$logical$assign$tmpindexm1146332801$0] = a[x]));
`))
}
