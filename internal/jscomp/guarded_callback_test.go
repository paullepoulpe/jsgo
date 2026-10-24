// Port of com/google/javascript/jscomp/GuardedCallbackTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/GuardedCallbackTest.java)

/*
 * Copyright 2020 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// GuardedCallbackTest holds unit tests for GuardedCallback.
//
// Go-only: test-helper — the Java test methods are not named test*; their
// Go methods are Test_<javaName> so that RunJUnit4 finds them.
type GuardedCallbackTest struct {
	CompilerTestCaseBase
}

func NewGuardedCallbackTest() *GuardedCallbackTest {
	c := &GuardedCallbackTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestGuardedCallbackTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewGuardedCallbackTest() })
}

// guardedCallbackTest_pass is the anonymous CompilerPass of getProcessor.
type guardedCallbackTest_pass struct {
	compiler *Compiler
}

// Java: GuardedCallbackTest.getProcessor(Compiler).new CompilerPass() {...}.process(Node, Node)
func (p *guardedCallbackTest_pass) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(p.compiler, root, newGuardedCallbackTest_GuardSwitchingCallback(p.compiler))
}

// Java: GuardedCallbackTest.getProcessor(Compiler)
func (c *GuardedCallbackTest) GetProcessor(compiler *Compiler) CompilerPass {
	return &guardedCallbackTest_pass{compiler: compiler}
}

// guardedCallbackTest_GuardSwitchingCallback replaces all guarded name and
// property references with "GUARDED_NAME" and "GUARDED_PROP",
// respectively.
type guardedCallbackTest_GuardSwitchingCallback struct {
	GuardedCallbackBase[string]
}

// Java: GuardedCallbackTest.GuardSwitchingCallback(AbstractCompiler)
func newGuardedCallbackTest_GuardSwitchingCallback(compiler AbstractCompiler) *guardedCallbackTest_GuardSwitchingCallback {
	cb := &guardedCallbackTest_GuardSwitchingCallback{}
	cb.InitGuardedCallbackBase(cb, compiler)
	return cb
}

// Java: GuardedCallbackTest.GuardSwitchingCallback.visitGuarded(NodeTraversal, Node, Node)
func (cb *guardedCallbackTest_GuardSwitchingCallback) VisitGuarded(traversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsName() && cb.IsGuarded(n.GetString()) {
		n.SetString("GUARDED_NAME")
		traversal.ReportCodeChange()
	} else if n.IsGetProp() || n.IsOptChainGetProp() {
		// prefix guarded resource name with "." to keep properties distinct
		// from names
		if cb.IsGuarded("." + n.GetString()) {
			n.SetString("GUARDED_PROP")
			traversal.ReportCodeChange()
		}
	}
}

// Java: GuardedCallbackTest.unguardedTest()
func (c *GuardedCallbackTest) Test_unguardedTest() {
	externs := NewTestExternsBuilder().AddConsole().AddPromise().Build()
	c.TestSame(c.Externs(externs), c.Srcs("console.log(Promise.allSettled);"))
}

// Java: GuardedCallbackTest.guardedTest()
func (c *GuardedCallbackTest) Test_guardedTest() {
	externs := NewTestExternsBuilder().AddConsole().AddPromise().Build()
	c.Test(
		c.Externs(externs),
		c.Srcs("console.log(Promise && Promise.allSettled && Promise.allSettled);"),
		c.Expected(
			"console.log(GUARDED_NAME && GUARDED_NAME.GUARDED_PROP && GUARDED_NAME.GUARDED_PROP)"))
}

// Java: GuardedCallbackTest.ifGuardTest()
func (c *GuardedCallbackTest) Test_ifGuardTest() {
	externs := NewTestExternsBuilder().AddConsole().AddPromise().Build()
	c.Test(
		c.Externs(externs),
		c.Srcs(
			`if (Promise && Promise.allSettled) {
  Promise.allSettled([]).then(() => console.log('done'));
}
`),
		c.Expected(
			`if (GUARDED_NAME && GUARDED_NAME.GUARDED_PROP) {
  GUARDED_NAME.GUARDED_PROP([]).then(() => console.log('done'));
}
`))
}

// Java: GuardedCallbackTest.ifGuardWithOptChainTest()
func (c *GuardedCallbackTest) Test_ifGuardWithOptChainTest() {
	externs := NewTestExternsBuilder().AddConsole().AddPromise().Build()
	c.Test(
		c.Externs(externs),
		c.Srcs(
			`if (Promise?.allSettled) {
  Promise.allSettled([]).then(() => console.log('done'));
}
`),
		c.Expected(
			`if (GUARDED_NAME?.GUARDED_PROP) {
  GUARDED_NAME.GUARDED_PROP([]).then(() => console.log('done'));
}
`))
	c.Test(
		c.Externs(externs),
		c.Srcs(
			`if (Promise?.allSettled([])) {
  Promise.allSettled([]).then(() => console.log('done'));
}
`),
		c.Expected(
			`if (GUARDED_NAME?.allSettled([])) {
  GUARDED_NAME.allSettled([]).then(() => console.log('done'));
}
`))
}

// Java: GuardedCallbackTest.guardedAndUnguardedTest()
func (c *GuardedCallbackTest) Test_guardedAndUnguardedTest() {
	externs := NewTestExternsBuilder().AddConsole().AddPromise().Build()
	c.Test(
		c.Externs(externs),
		c.Srcs("console.log(Promise && Promise.allSettled && Promise.finally);"),
		// The values of `Promise` and `Promise.allSettled` are only checked
		// for truthiness, so they can be considered guarded. However, the
		// value of `Promise.finally` may end up getting passed to
		// `console.log()`, so it is not guarded.
		c.Expected("console.log(GUARDED_NAME && GUARDED_NAME.GUARDED_PROP && GUARDED_NAME.finally)"))
}

// Java: GuardedCallbackTest.optChainTest()
func (c *GuardedCallbackTest) Test_optChainTest() {
	externs := NewTestExternsBuilder().AddConsole().AddPromise().Build()
	c.TestSame(
		c.Externs(externs),
		// Nothing guarded
		c.Srcs("console.log(Promise.allSettled([Promise.resolve(), x.allSettled]))"))
	c.Test(
		c.Externs(externs),
		c.Srcs("console.log(Promise?.allSettled([Promise.resolve(), x.allSettled]))"),
		c.Expected("console.log(GUARDED_NAME?.allSettled([GUARDED_NAME.resolve(), x.allSettled]))"))
	c.Test(
		c.Externs(externs),
		c.Srcs("console.log(Promise.allSettled?.([Promise.resolve(), x.allSettled]))"),
		c.Expected("console.log(Promise.GUARDED_PROP?.([Promise.resolve(), x.GUARDED_PROP]))"))
	c.Test(
		c.Externs(externs),
		c.Srcs("console.log(Promise?.allSettled?.([Promise.resolve(), x.allSettled]))"),
		c.Expected(
			"console.log(GUARDED_NAME?.GUARDED_PROP?.([GUARDED_NAME.resolve(), x.GUARDED_PROP]))"))
	c.Test(
		c.Externs(externs),
		c.Srcs(
			`console.log(
    Promise?.resolve(Promise)
        .then(x.finally)
        .finally?.(x.finally))
`),
		c.Expected(
			`console.log(
    GUARDED_NAME?.resolve(GUARDED_NAME)
        .then(x.finally)
        .GUARDED_PROP?.(x.GUARDED_PROP))
`))
}
