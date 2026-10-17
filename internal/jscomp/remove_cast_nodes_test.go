// Port of com/google/javascript/jscomp/RemoveCastNodesTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/RemoveCastNodesTest.java)

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

import (
	"testing"
)

type RemoveCastNodesTest struct {
	CompilerTestCaseBase
}

// Java: RemoveCastNodesTest()
func NewRemoveCastNodesTest() *RemoveCastNodesTest {
	c := &RemoveCastNodesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestRemoveCastNodesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRemoveCastNodesTest() })
}

// Java: RemoveCastNodesTest.getProcessor(Compiler)
func (c *RemoveCastNodesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewRemoveCastNodes(compiler)
}

// Java: RemoveCastNodesTest.removesCastInExpression()
func (c *RemoveCastNodesTest) TestRemovesCastInExpression() {
	c.Test(c.Srcs("/** @type {string} */ (0);"), c.Expected("0"))
}

// Java: RemoveCastNodesTest.keepsTypeAnnotations()
func (c *RemoveCastNodesTest) TestKeepsTypeAnnotations() {
	c.TestSame(c.Srcs("/** @type {string} */ const x = 0;"))
}
