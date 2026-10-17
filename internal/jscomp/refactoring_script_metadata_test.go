// Port of com/google/javascript/refactoring/ScriptMetadataTest.java @ closure-compiler v20261015
// (test/com/google/javascript/refactoring/ScriptMetadataTest.java)

/*
 * Copyright 2014 The Closure Compiler Authors.
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

// Java: public final class ScriptMetadataTest extends CompilerTestCase
type ScriptMetadataTest struct {
	CompilerTestCaseBase

	lastMetadata *ScriptMetadata
}

// Java: ScriptMetadataTest()
func NewScriptMetadataTest() *ScriptMetadataTest {
	c := &ScriptMetadataTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestScriptMetadataTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewScriptMetadataTest() })
}

// Java: ScriptMetadataTest.getProcessor(Compiler)
func (c *ScriptMetadataTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
		c.lastMetadata = ScriptMetadata_Create(root.GetOnlyChild(), compiler)
	})
}

// Java: ScriptMetadataTest.testParsesAliases()
func (c *ScriptMetadataTest) TestParsesAliases() {
	c.TestSame(
		"const Foo = goog.require('a.b.Foo');\n" +
			"const Bar = goog.requireType('a.b.Bar');\n" +
			"const Qux = goog.forwardDeclare('a.b.Qux');\n")

	AssertThat(c.lastMetadata.GetAlias("a.b.Foo")).IsEqualTo("Foo")
	AssertThat(c.lastMetadata.GetAlias("a.b.Bar")).IsEqualTo("Bar")
	AssertThat(c.lastMetadata.GetAlias("a.b.Qux")).IsEqualTo("Qux")
}
