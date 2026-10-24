// Port of com/google/javascript/jscomp/InlineCostEstimatorTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/InlineCostEstimatorTest.java)

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
	"math"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// InlineCostEstimatorTest is the unit test for InlineCostEstimator.
type InlineCostEstimatorTest struct {
	// Go-only: test-helper — the Go test running the current test method.
	t *testing.T
}

// Java: InlineCostEstimatorTest()
func NewInlineCostEstimatorTest() *InlineCostEstimatorTest {
	return &InlineCostEstimatorTest{}
}

func TestInlineCostEstimatorTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewInlineCostEstimatorTest() })
}

// JUnitSetT gives the instance the Go test running its test method.
func (c *InlineCostEstimatorTest) JUnitSetT(t *testing.T) { c.t = t }

// JUnitBefore: the class has no @Before method.
func (c *InlineCostEstimatorTest) JUnitBefore() {}

// JUnitAfter: the class has no @After method.
func (c *InlineCostEstimatorTest) JUnitAfter() {}

// Java: InlineCostEstimatorTest.parse(String)
func InlineCostEstimatorTest_Parse(js string) rhino.Node {
	compiler := NewCompiler()
	n := compiler.ParseTestCode(js)
	AssertThat(compiler.GetErrors()).IsEmpty()
	return n
}

// Java: InlineCostEstimatorTest.minimize(String)
func InlineCostEstimatorTest_Minimize(js string) string {
	options := NewCompilerOptions()
	options.SetLineLengthThreshold(math.MaxInt32)
	return NewCodePrinter_Builder(InlineCostEstimatorTest_Parse(js)).SetCompilerOptions(options).Build()
}

// Java: InlineCostEstimatorTest.cost(String)
func InlineCostEstimatorTest_Cost(js string) int64 {
	return int64(InlineCostEstimator_GetCost(InlineCostEstimatorTest_Parse(js)))
}

// Java: InlineCostEstimatorTest.testCost()
func (c *InlineCostEstimatorTest) TestCost() {
	c.CheckCost("1", "1")
	c.CheckCost("true", "1")
	c.CheckCost("false", "1")
	c.CheckCost("a", "xx")
	c.CheckCost("a + b", "xx+xx")
	c.CheckCost("foo()", "xx()")
	c.CheckCost("foo(a,b)", "xx(xx,xx)")
	c.CheckCost("10 + foo(a,b)", "0+xx(xx,xx)")
	c.CheckCost("1 + foo(a,b)", "1+xx(xx,xx)")
	c.CheckCost("a ? 1 : 0", "xx?1:0")
	c.CheckCost("a.b", "xx.xx")
	c.CheckCost("new Obj()", "new xx")
	c.CheckCost("function a() {return \"monkey\"}", "function xx(){return\"monkey\"}")
}

// Java: InlineCostEstimatorTest.checkCost(String, String)
func (c *InlineCostEstimatorTest) CheckCost(source string, example string) {

	// The example string should have been minified already.
	AssertThat(example).IsEqualTo(InlineCostEstimatorTest_Minimize(example))

	// cost estimate should be the same as the length of the example string.
	AssertThat(InlineCostEstimatorTest_Cost(source)).IsEqualTo(int64(len(example)))
}
