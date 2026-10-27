// Port of com/google/javascript/refactoring/ScriptMetadataTest.java @ closure-compiler v20261021
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

package refactoring

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/testjunit"
)

func TestScriptMetadataTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return &ScriptMetadataTest{} })
}

type ScriptMetadataTest struct {
	t *testing.T
}

func (c *ScriptMetadataTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *ScriptMetadataTest) JUnitBefore()           {}
func (c *ScriptMetadataTest) JUnitAfter()            {}

// Java: ScriptMetadataTest.testParsesAliases()
func (c *ScriptMetadataTest) TestParsesAliases() {
	src := "const Foo = goog.require('a.b.Foo');\n" +
		"const Bar = goog.requireType('a.b.Bar');\n" +
		"const Qux = goog.forwardDeclare('a.b.Qux');\n"
	compiler := jscomp.NewCompiler()
	compiler.Init(nil, []*jscomp.SourceFile{jscomp.SourceFile_FromCode("test", src)}, jscomp.NewCompilerOptions())
	compiler.Parse()
	metadata := ScriptMetadata_Create(compiler.GetRoot().GetLastChild().GetOnlyChild(), compiler)

	testjunit.AssertThat(metadata.GetAlias("a.b.Foo")).IsEqualTo("Foo")
	testjunit.AssertThat(metadata.GetAlias("a.b.Bar")).IsEqualTo("Bar")
	testjunit.AssertThat(metadata.GetAlias("a.b.Qux")).IsEqualTo("Qux")
}
