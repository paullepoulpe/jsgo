// Port of com/google/javascript/jscomp/PhaseOptimizerTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/PhaseOptimizerTest.java)

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
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PhaseOptimizerTest holds the tests for PhaseOptimizer.
type PhaseOptimizerTest struct {
	t           *testing.T
	passesRun   []string
	dummyRoot   rhino.Node
	dummyScript rhino.Node
	optimizer   *PhaseOptimizer
	compiler    *Compiler
	tracker     *PerformanceTracker
}

func (c *PhaseOptimizerTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *PhaseOptimizerTest) JUnitBefore()           { c.SetUp() }
func (c *PhaseOptimizerTest) JUnitAfter()            {}

func TestPhaseOptimizerTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return &PhaseOptimizerTest{} })
}

// Java: PhaseOptimizerTest.setUp()
func (c *PhaseOptimizerTest) SetUp() {
	c.passesRun = nil
	dummyExternsRoot := rhino.NewNode(rhino.Token_ROOT)
	c.dummyScript = rhino.IR_Script()
	c.dummyRoot = rhino.IR_Root(c.dummyScript)
	c.compiler = NewCompiler()
	c.compiler.InitCompilerOptionsIfTesting()
	c.tracker = NewPerformanceTracker(dummyExternsRoot, c.dummyRoot, CompilerOptions_TracerMode_TIMING_ONLY)
	c.optimizer = NewPhaseOptimizer(c.compiler, c.tracker)
}

// Java: PhaseOptimizerTest.testOneRun()
func (c *PhaseOptimizerTest) TestOneRun() {
	c.AddOneTimePass("x")
	c.AssertPasses("x")
}

// Java: PhaseOptimizerTest.testLoop1()
func (c *PhaseOptimizerTest) TestLoop1() {
	loop := c.optimizer.AddFixedPointLoop()
	c.AddLoopedPass(loop, "x", 0)
	c.AssertPasses("x")
}

// Java: PhaseOptimizerTest.testLoop2()
func (c *PhaseOptimizerTest) TestLoop2() {
	loop := c.optimizer.AddFixedPointLoop()
	c.AddLoopedPass(loop, "x", 3)
	c.AssertPasses("x", "x", "x", "x")
}

// Java: PhaseOptimizerTest.testSchedulingOfLoopablePasses()
func (c *PhaseOptimizerTest) TestSchedulingOfLoopablePasses() {
	loop := c.optimizer.AddFixedPointLoop()
	c.AddLoopedPass(loop, "x", 3)
	c.AddLoopedPass(loop, "y", 1)
	// The pass iterations can be grouped as: [x y] [x y] [x] [x] [y]
	c.AssertPasses("x", "y", "x", "y", "x", "x", "y")
}

// Java: PhaseOptimizerTest.testCapLoopIterations()
func (c *PhaseOptimizerTest) TestCapLoopIterations() {
	options := c.compiler.GetOptions()
	options.SetMaxOptimizationLoopIterations(1)
	c.optimizer = NewPhaseOptimizer(c.compiler, c.tracker)
	loop := c.optimizer.AddFixedPointLoop()
	c.AddLoopedPass(loop, PassNames_PEEPHOLE_OPTIMIZATIONS, 2)
	c.AssertPasses(PassNames_PEEPHOLE_OPTIMIZATIONS)
}

// Java: PhaseOptimizerTest.testNotInfiniteLoop()
func (c *PhaseOptimizerTest) TestNotInfiniteLoop() {
	loop := c.optimizer.AddFixedPointLoop()
	c.AddLoopedPass(loop, "x", PhaseOptimizer_MAX_LOOPS-2)
	c.optimizer.Process(rhino.NilNode, c.dummyRoot)
	AssertWithMessage("There should be no errors.").That(c.compiler.GetErrorCount()).IsEqualTo(0)
}

// Java: PhaseOptimizerTest.testInfiniteLoop()
func (c *PhaseOptimizerTest) TestInfiniteLoop() {
	loop := c.optimizer.AddFixedPointLoop()
	c.AddLoopedPass(loop, "x", PhaseOptimizer_MAX_LOOPS+1)
	jutil.Try(func() {
		c.optimizer.Process(rhino.NilNode, c.dummyRoot)
		AssertWithMessage("Expected RuntimeException").Fail()
	}, func(e jutil.IRuntimeException) {
		AssertWithMessage(e.GetMessage()).
			That(strings.Contains(e.GetMessage(), PhaseOptimizer_OPTIMIZE_LOOP_ERROR)).
			IsTrue()
	})
}

// Java: PhaseOptimizerTest.testSchedulingOfAnyKindOfPasses1()
func (c *PhaseOptimizerTest) TestSchedulingOfAnyKindOfPasses1() {
	c.AddOneTimePass("a")
	loop := c.optimizer.AddFixedPointLoop()
	c.AddLoopedPass(loop, "x", 3)
	c.AddLoopedPass(loop, "y", 1)
	c.AddOneTimePass("z")
	c.AssertPasses("a", "x", "y", "x", "y", "x", "x", "y", "z")
}

// Java: PhaseOptimizerTest.testSchedulingOfAnyKindOfPasses2()
func (c *PhaseOptimizerTest) TestSchedulingOfAnyKindOfPasses2() {
	c.optimizer.Consume(
		[]*PassFactory{
			c.CreatePassFactory("a", 0, true),
			c.CreatePassFactory("b", 1, false),
			c.CreatePassFactory("c", 2, false),
			c.CreatePassFactory("d", 1, false),
			c.CreatePassFactory("e", 1, true),
			c.CreatePassFactory("f", 0, true),
		})
	// The pass iterations can be grouped as:
	// 1. [a]  2. [b c d]  3. [b c d]  4. [c]  5. [b]  6. [e]  7. [f]
	// In loop #3, "b" is run and does not make changes, then "c" is run and makes changes, then "d"
	// is run and does not make changes.
	// In loop #4, "c" is run and does not make changes.
	// In loop #5, "b" is run but "d" is not run. This is because the AST changed after the last run
	// of "b" but has not changed after the last run of "d".
	c.AssertPasses("a", "b", "c", "d", "b", "c", "d", "c", "b", "e", "f")
}

// Java: PhaseOptimizerTest.testSchedulingOfAnyKindOfPasses3()
func (c *PhaseOptimizerTest) TestSchedulingOfAnyKindOfPasses3() {
	c.optimizer.Consume(
		[]*PassFactory{
			c.CreatePassFactory("a", 2, false),
			c.CreatePassFactory("b", 1, true),
			c.CreatePassFactory("c", 1, false),
		})
	c.AssertPasses("a", "a", "a", "b", "c", "c")
}

// Java: PhaseOptimizerTest.testSchedulingOfAnyKindOfPasses4()
func (c *PhaseOptimizerTest) TestSchedulingOfAnyKindOfPasses4() {
	c.optimizer.Consume(
		[]*PassFactory{
			c.CreatePassFactory("a", 2, true),
			c.CreatePassFactory("b", 0, false),
			c.CreatePassFactory("c", 0, false),
		})
	c.AssertPasses("a", "b", "c")
}

// Java: PhaseOptimizerTest.testDuplicateLoop()
func (c *PhaseOptimizerTest) TestDuplicateLoop() {
	loop := c.optimizer.AddFixedPointLoop()
	c.AddLoopedPass(loop, "x", 1)
	jutil.Try(func() {
		c.AddLoopedPass(loop, "x", 1)
		AssertWithMessage("Expected exception").Fail()
	}, func(e jutil.IIllegalArgumentException) {
	})
}

// Java: PhaseOptimizerTest.preconditionCheck_success()
func (c *PhaseOptimizerTest) TestPreconditionCheck_success() {
	passFactory :=
		PassFactory_Builder_().
			SetName("myPass").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass { return c.CreatePass("myPass", 0) }).
			SetPreconditionCheck(func(options *CompilerOptions) *PassFactory_PreconditionResult {
				return NewPassFactory_PreconditionResult(true, nil)
			}).
			Build()
	c.optimizer.AddOneTimePass(passFactory)
	c.AssertPasses("myPass")
}

// Java: PhaseOptimizerTest.preconditionCheck_failure()
func (c *PhaseOptimizerTest) TestPreconditionCheck_failure() {
	passFactory :=
		PassFactory_Builder_().
			SetName("myPass").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass { return c.CreatePass("myPass", 0) }).
			SetPreconditionCheck(func(options *CompilerOptions) *PassFactory_PreconditionResult {
				message := "message"
				return NewPassFactory_PreconditionResult(false, &message)
			}).
			Build()
	c.optimizer.AddOneTimePass(passFactory)

	ex :=
		AssertThrows[jutil.IIllegalArgumentException](func() { c.AssertPasses("myPass") })

	AssertThat(ex.GetMessage()).IsEqualTo("Precondition for pass myPass failed: message")
}

// Java: PhaseOptimizerTest.assertPasses(String...)
func (c *PhaseOptimizerTest) AssertPasses(names ...string) {
	c.optimizer.Process(rhino.NilNode, c.dummyRoot)
	AssertThat(c.passesRun).IsEqualTo(names)
}

// Java: PhaseOptimizerTest.addOneTimePass(String)
func (c *PhaseOptimizerTest) AddOneTimePass(name string) {
	c.optimizer.AddOneTimePass(c.CreatePassFactory(name, 0, true))
}

// Java: PhaseOptimizerTest.addLoopedPass(Loop, String, int)
func (c *PhaseOptimizerTest) AddLoopedPass(loop *PhaseOptimizer_Loop, name string, numChanges int) {
	loop.AddLoopedPass(c.CreatePassFactory(name, numChanges, false))
}

// Java: PhaseOptimizerTest.createPassFactory(String, int, boolean)
func (c *PhaseOptimizerTest) CreatePassFactory(name string, numChanges int, isOneTime bool) *PassFactory {
	return c.CreatePassFactoryWithPass(name, c.CreatePass(name, numChanges), isOneTime)
}

// Java: PhaseOptimizerTest.createPassFactory(String, CompilerPass, boolean)
func (c *PhaseOptimizerTest) CreatePassFactoryWithPass(name string, pass CompilerPass, isOneTime bool) *PassFactory {
	return PassFactory_Builder_().
		SetName(name).
		SetRunInFixedPointLoop(!isOneTime).
		SetInternalFactory(func(compiler AbstractCompiler) CompilerPass { return pass }).
		Build()
}

// Java: PhaseOptimizerTest.createPass(String, int)
func (c *PhaseOptimizerTest) CreatePass(name string, numChanges int) CompilerPass {
	self := c
	numChangesClosure := []int{numChanges}
	return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
		c.passesRun = append(c.passesRun, name)
		if numChangesClosure[0] > 0 {
			numChangesClosure[0] = numChangesClosure[0] - 1
			c.compiler.ReportChangeToEnclosingScope(self.dummyScript)
		}
	})
}
