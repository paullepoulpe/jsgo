// Port of com/google/javascript/refactoring/ApplySuggestedFixesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/refactoring/ApplySuggestedFixesTest.java)

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
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/testjunit"
)

// ApplySuggestedFixesTest is the test case for ApplySuggestedFixes.
//
// Java: public class ApplySuggestedFixesTest
type ApplySuggestedFixesTest struct {
	t *testing.T
}

func (c *ApplySuggestedFixesTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *ApplySuggestedFixesTest) JUnitBefore()           {}
func (c *ApplySuggestedFixesTest) JUnitAfter()            {}

// Java: ApplySuggestedFixesTest()
func NewApplySuggestedFixesTest() *ApplySuggestedFixesTest {
	return &ApplySuggestedFixesTest{}
}

func TestApplySuggestedFixesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewApplySuggestedFixesTest() })
}

// applySuggestedFixesTest_mapOf is ImmutableMap.of(k1, v1, ...).
//
// Go-only: test-helper.
func applySuggestedFixesTest_mapOf(kvs ...string) map[string]string {
	m := make(map[string]string, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		m[kvs[i]] = kvs[i+1]
	}
	return m
}

// Java: ApplySuggestedFixesTest.testApplyCodeReplacements_overlapsAreErrors()
func (c *ApplySuggestedFixesTest) TestApplyCodeReplacements_overlapsAreErrors() {
	replacements :=
		[]*CodeReplacement{CodeReplacement_Create(0, 10, ""), CodeReplacement_Create(5, 3, "")}
	testjunit.AssertThrows[jutil.IIllegalArgumentException](
		func() { ApplySuggestedFixes_ApplyCodeReplacements(replacements, "") })
}

// Java: ApplySuggestedFixesTest.testApplyCodeReplacements_overlapsAreErrors_unlessEqual()
func (c *ApplySuggestedFixesTest) TestApplyCodeReplacements_overlapsAreErrors_unlessEqual() {
	replacements :=
		[]*CodeReplacement{CodeReplacement_Create(0, 3, "A"), CodeReplacement_Create(0, 3, "A")}
	ApplySuggestedFixes_ApplyCodeReplacements(replacements, "abcdef")
}

// Java: ApplySuggestedFixesTest.testApplyCodeReplacements_noOverlapsSucceed()
func (c *ApplySuggestedFixesTest) TestApplyCodeReplacements_noOverlapsSucceed() {
	replacements :=
		[]*CodeReplacement{CodeReplacement_Create(0, 3, ""), CodeReplacement_Create(5, 3, "")}
	ApplySuggestedFixes_ApplyCodeReplacements(replacements, "abcdef")
}

// Java: ApplySuggestedFixesTest.testApplyCodeReplacements()
func (c *ApplySuggestedFixesTest) TestApplyCodeReplacements() {
	replacements :=
		[]*CodeReplacement{CodeReplacement_Create(0, 1, "z"), CodeReplacement_Create(3, 2, "qq")}
	newCode := ApplySuggestedFixes_ApplyCodeReplacements(replacements, "abcdef")
	testjunit.AssertThat(newCode).IsEqualTo("zbcqqf")
}

// Java: ApplySuggestedFixesTest.testApplyCodeReplacements_insertion()
func (c *ApplySuggestedFixesTest) TestApplyCodeReplacements_insertion() {
	replacements :=
		[]*CodeReplacement{CodeReplacement_Create(0, 0, "z")}
	newCode := ApplySuggestedFixes_ApplyCodeReplacements(replacements, "abcdef")
	testjunit.AssertThat(newCode).IsEqualTo("zabcdef")
}

// Java: ApplySuggestedFixesTest.testApplyCodeReplacements_deletion()
func (c *ApplySuggestedFixesTest) TestApplyCodeReplacements_deletion() {
	replacements :=
		[]*CodeReplacement{CodeReplacement_Create(0, 6, "")}
	newCode := ApplySuggestedFixes_ApplyCodeReplacements(replacements, "abcdef")
	testjunit.AssertThat(newCode).IsEmpty()
}

// Java: ApplySuggestedFixesTest.testApplyCodeReplacements_boundaryCases()
func (c *ApplySuggestedFixesTest) TestApplyCodeReplacements_boundaryCases() {
	replacements :=
		[]*CodeReplacement{CodeReplacement_Create(5, 1, "z")}
	newCode := ApplySuggestedFixes_ApplyCodeReplacements(replacements, "abcdef")
	testjunit.AssertThat(newCode).IsEqualTo("abcdez")
}

// Java: ApplySuggestedFixesTest.testApplyCodeReplacements_multipleReplacements()
func (c *ApplySuggestedFixesTest) TestApplyCodeReplacements_multipleReplacements() {
	replacements :=
		[]*CodeReplacement{
			CodeReplacement_Create(0, 2, "z"),
			CodeReplacement_Create(2, 1, "y"),
			CodeReplacement_Create(3, 3, "xwvu")}
	newCode := ApplySuggestedFixes_ApplyCodeReplacements(replacements, "abcdef")
	testjunit.AssertThat(newCode).IsEqualTo("zyxwvu")
}

// Java: ApplySuggestedFixesTest.testApplySuggestedFixes()
func (c *ApplySuggestedFixesTest) TestApplySuggestedFixes() {
	code := "var someNode;"
	compiler := applySuggestedFixesTest_getCompiler(code)
	root := applySuggestedFixesTest_compileToScriptRoot(compiler)
	fixes :=
		[]*SuggestedFix{NewSuggestedFix_Builder().Delete(root).Build()}
	codeMap := applySuggestedFixesTest_mapOf("test", code)
	newCodeMap :=
		ApplySuggestedFixes_ApplySuggestedFixesToCode(slices.Values(fixes), codeMap)
	testjunit.AssertThat(newCodeMap).ContainsExactlyEntriesIn(applySuggestedFixesTest_mapOf("test", ""))
}

// Java: ApplySuggestedFixesTest.testApplySuggestedFixes_insideJSDoc()
func (c *ApplySuggestedFixesTest) TestApplySuggestedFixes_insideJSDoc() {
	code := "/** @type {Foo} */\nvar foo = new Foo()"
	compiler := applySuggestedFixesTest_getCompiler(code)
	root := applySuggestedFixesTest_compileToScriptRoot(compiler)
	varNode := root.GetFirstChild()
	jsdocRoot :=
		applySuggestedFixesTest_getOnlyElement(varNode.GetJSDocInfo().GetTypeNodes())
	fix := NewSuggestedFix_Builder().
		InsertBefore(jsdocRoot, "!").
		Build()
	fixes := []*SuggestedFix{fix}
	codeMap := applySuggestedFixesTest_mapOf("test", code)
	newCodeMap :=
		ApplySuggestedFixes_ApplySuggestedFixesToCode(slices.Values(fixes), codeMap)
	testjunit.AssertThat(newCodeMap).ContainsExactlyEntriesIn(
		applySuggestedFixesTest_mapOf("test", "/** @type {!Foo} */\nvar foo = new Foo()"))
}

// Java: ApplySuggestedFixesTest.testApplySuggestedFixes_multipleFixesInJsdoc()
func (c *ApplySuggestedFixesTest) TestApplySuggestedFixes_multipleFixesInJsdoc() {
	code := "/** @type {Array<Foo>} */\nvar arr = [new Foo()];"
	compiler := applySuggestedFixesTest_getCompiler(code)
	root := applySuggestedFixesTest_compileToScriptRoot(compiler)
	varNode := root.GetFirstChild()
	jsdocRoot :=
		applySuggestedFixesTest_getOnlyElement(varNode.GetJSDocInfo().GetTypeNodes())
	fix1 := NewSuggestedFix_Builder().
		InsertBefore(jsdocRoot, "!").
		Build()
	foo := jsdocRoot.GetFirstFirstChild()
	fix2 := NewSuggestedFix_Builder().
		InsertBefore(foo, "!").
		Build()
	fixes := []*SuggestedFix{fix1, fix2}
	codeMap := applySuggestedFixesTest_mapOf("test", code)
	newCodeMap :=
		ApplySuggestedFixes_ApplySuggestedFixesToCode(slices.Values(fixes), codeMap)
	testjunit.AssertThat(newCodeMap).
		ContainsExactlyEntriesIn(applySuggestedFixesTest_mapOf("test", "/** @type {!Array<!Foo>} */\nvar arr = [new Foo()];"))
}

// Java: ApplySuggestedFixesTest.testApplySuggestedFixes_noFixes()
func (c *ApplySuggestedFixesTest) TestApplySuggestedFixes_noFixes() {
	codeMap :=
		applySuggestedFixesTest_mapOf(
			"file1", "abcdef",
			"file2", "abcdef")
	expectedNewCodeMap := applySuggestedFixesTest_mapOf()
	fixes := []*SuggestedFix{}
	newCodeMap :=
		ApplySuggestedFixes_ApplySuggestedFixesToCode(slices.Values(fixes), codeMap)
	testjunit.AssertThat(newCodeMap).IsEqualTo(expectedNewCodeMap)
}

// Java: ApplySuggestedFixesTest.testApplySuggestedFixes_missingCodeForFile()
func (c *ApplySuggestedFixesTest) TestApplySuggestedFixes_missingCodeForFile() {
	codeMap := applySuggestedFixesTest_mapOf()
	code := "var someNode;"
	compiler := applySuggestedFixesTest_getCompiler(code)
	root := applySuggestedFixesTest_compileToScriptRoot(compiler)
	fixes :=
		[]*SuggestedFix{NewSuggestedFix_Builder().Delete(root).Build()}
	testjunit.AssertThrows[jutil.IIllegalArgumentException](
		func() { ApplySuggestedFixes_ApplySuggestedFixesToCode(slices.Values(fixes), codeMap) })
}

// Java: ApplySuggestedFixesTest.testApplySuggestedFixes_withOveralppingEqualParts_areAllApplied()
func (c *ApplySuggestedFixesTest) TestApplySuggestedFixes_withOveralppingEqualParts_areAllApplied() {
	code := "var first, second, shared;"
	compiler := applySuggestedFixesTest_getCompiler(code)
	root := applySuggestedFixesTest_compileToScriptRoot(compiler)
	fmt.Println(root.ToStringTree())
	var_ := root.GetFirstChild()

	fixes :=
		[]*SuggestedFix{
			NewSuggestedFix_Builder().
				Rename(var_.GetLastChild(), "newShared").
				Rename(var_.GetFirstChild(), "newFirst").
				Build(),
			NewSuggestedFix_Builder().
				Rename(var_.GetLastChild(), "newShared").
				Rename(var_.GetSecondChild(), "newSecond").
				Build()}

	newCodeMap :=
		ApplySuggestedFixes_ApplySuggestedFixesToCode(slices.Values(fixes), applySuggestedFixesTest_mapOf("test", code))
	testjunit.AssertThat(newCodeMap).ContainsExactlyEntriesIn(applySuggestedFixesTest_mapOf("test", "var newFirst, newSecond, newShared;"))
}

// applySuggestedFixesTest_compileToScriptRoot returns the root script node
// produced from the compiled JS input.
//
// Java: ApplySuggestedFixesTest.compileToScriptRoot(Compiler)
func applySuggestedFixesTest_compileToScriptRoot(compiler *jscomp.Compiler) rhino.Node {
	root := compiler.GetRoot()
	// The last child of the compiler root is a Block node, and the first
	// child of that is the Script node.
	return root.GetLastChild().GetFirstChild()
}

// Java: ApplySuggestedFixesTest.getCompiler(String)
func applySuggestedFixesTest_getCompiler(jsInput string) *jscomp.Compiler {
	compiler := jscomp.NewCompiler()
	options := RefactoringDriver_GetCompilerOptions()
	compiler.Init(
		[]*jscomp.SourceFile{}, // Externs
		[]*jscomp.SourceFile{jscomp.SourceFile_FromCode("test", jsInput)},
		options)
	compiler.Parse()
	return compiler
}

// applySuggestedFixesTest_getOnlyElement is Iterables.getOnlyElement.
//
// Go-only: test-helper.
func applySuggestedFixesTest_getOnlyElement(nodes []rhino.Node) rhino.Node {
	jutil.CheckArgumentMsg1(len(nodes) == 1, "expected one element but was: %s", len(nodes))
	return nodes[0]
}

func singleUseSeq[T any](items []T) iter.Seq[T] {
	consumed := false
	return func(yield func(T) bool) {
		if consumed {
			panic("single-use iter.Seq iterated more than once")
		}
		consumed = true
		for _, item := range items {
			if !yield(item) {
				return
			}
		}
	}
}

func (c *ApplySuggestedFixesTest) TestApplySuggestedFixes_singleUseIterSeq() {
	tmpFile := filepath.Join(c.t.TempDir(), "input.js")
	code := "var oldName;"
	if err := os.WriteFile(tmpFile, []byte(code), 0644); err != nil {
		c.t.Fatal(err)
	}
	compiler := jscomp.NewCompiler()
	options := RefactoringDriver_GetCompilerOptions()
	compiler.Init(
		[]*jscomp.SourceFile{},
		[]*jscomp.SourceFile{jscomp.SourceFile_FromCode(tmpFile, code)},
		options)
	compiler.Parse()
	root := applySuggestedFixesTest_compileToScriptRoot(compiler)
	nameNode := root.GetFirstFirstChild()

	fix := NewSuggestedFix_Builder().Rename(nameNode, "newName").Build()
	ApplySuggestedFixes_ApplySuggestedFixesToFiles(singleUseSeq([]*SuggestedFix{fix}))
	gotBytes, err := os.ReadFile(tmpFile)
	if err != nil {
		c.t.Fatal(err)
	}
	testjunit.AssertThat(string(gotBytes)).IsEqualTo("var newName;")

	alt2 := NewSuggestedFix_Builder().Rename(nameNode, "choiceB").Build()
	withAlts := NewSuggestedFix_Builder().Rename(nameNode, "choiceA").AddAlternative(alt2).Build()
	choices := ApplySuggestedFixes_ApplyAllSuggestedFixChoicesToCode(
		singleUseSeq([]*SuggestedFix{withAlts}),
		applySuggestedFixesTest_mapOf(tmpFile, code))
	testjunit.AssertThat(len(choices)).IsEqualTo(2)
	testjunit.AssertThat(choices[0][tmpFile]).IsEqualTo("var choiceA;")
	testjunit.AssertThat(choices[1][tmpFile]).IsEqualTo("var choiceB;")
}
