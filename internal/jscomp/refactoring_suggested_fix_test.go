// Port of com/google/javascript/refactoring/SuggestedFixTest.java @ closure-compiler v20261021
// (test/com/google/javascript/refactoring/SuggestedFixTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// SuggestedFixTest holds the unit tests for JsFlume SuggestedFix.
//
// Java: public class SuggestedFixTest
type SuggestedFixTest struct {
	t *testing.T
}

func (c *SuggestedFixTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *SuggestedFixTest) JUnitBefore()           {}
func (c *SuggestedFixTest) JUnitAfter()            {}

// Java: SuggestedFixTest()
func NewSuggestedFixTest() *SuggestedFixTest {
	return &SuggestedFixTest{}
}

func TestSuggestedFixTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewSuggestedFixTest() })
}

// Java: SuggestedFixTest.testInsertBefore()
func (c *SuggestedFixTest) TestInsertBefore() {
	before := "var someRandomCode = {};"
	after := "/** some comment */\ngoog.foo();"
	compiler := suggestedFixTest_getCompiler(before + after)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	newNode :=
		rhino.IR_ExprResult(rhino.IR_Call(rhino.IR_Getprop(rhino.IR_Name("goog2"), "get"), rhino.IR_String("service")))
	fix :=
		NewSuggestedFix_Builder().
			InsertBeforeNode(root.GetLastChild().GetFirstChild(), newNode, compiler).
			Build()
	replacement :=
		CodeReplacement_Create(jutil.StringLength(before), 0, "goog2.get('service');\n")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testInsertBeforeWithNonJSDocBlockComment()
func (c *SuggestedFixTest) TestInsertBeforeWithNonJSDocBlockComment() {
	before := "var someRandomCode = {};\n"
	after := "/* some comment */\ngoog.foo();"
	compiler := suggestedFixTest_getCompiler(before + after)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	newNode :=
		rhino.IR_ExprResult(rhino.IR_Call(rhino.IR_Getprop(rhino.IR_Name("goog2"), "get"), rhino.IR_String("service")))
	fix :=
		NewSuggestedFix_Builder().InsertBeforeNode(root.GetLastChild(), newNode, compiler).Build()
	replacement :=
		CodeReplacement_Create(jutil.StringLength(before), 0, "goog2.get('service');\n")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testInsertBeforeWithNonJSDocLineComment()
func (c *SuggestedFixTest) TestInsertBeforeWithNonJSDocLineComment() {
	before := "var someRandomCode = {};\n"
	after := "// some comment \ngoog.foo();"
	compiler := suggestedFixTest_getCompiler(before + after)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	newNode :=
		rhino.IR_ExprResult(rhino.IR_Call(rhino.IR_Getprop(rhino.IR_Name("goog2"), "get"), rhino.IR_String("service")))
	fix :=
		NewSuggestedFix_Builder().InsertBeforeNode(root.GetLastChild(), newNode, compiler).Build()
	replacement :=
		CodeReplacement_Create(jutil.StringLength(before), 0, "goog2.get('service');\n")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testReplaceText()
func (c *SuggestedFixTest) TestReplaceText() {
	input := "var foo = new Bar();"
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().ReplaceText(root.GetFirstFirstChild(), 3, "quux").Build()
	replacement := CodeReplacement_Create(4, 3, "quux")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testDelete()
func (c *SuggestedFixTest) TestDelete() {
	input := "var foo = new Bar();"
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix := NewSuggestedFix_Builder().Delete(root.GetFirstChild()).Build()
	replacement := CodeReplacement_Create(0, jutil.StringLength(input), "")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testDelete_spaceBeforeNode()
func (c *SuggestedFixTest) TestDelete_spaceBeforeNode() {
	before := "var foo = new Bar();"
	after := "\n\nvar baz = new Baz();"
	input := before + after
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix := NewSuggestedFix_Builder().Delete(root.GetLastChild()).Build()
	replacement := CodeReplacement_Create(jutil.StringLength(before), jutil.StringLength(after), "")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testDelete_NonJSDocLineCommentBeforeNode()
func (c *SuggestedFixTest) TestDelete_NonJSDocLineCommentBeforeNode() {
	before := "var foo = new Bar();"
	after := "// blah \n var baz = new Baz();"
	input := before + after
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix := NewSuggestedFix_Builder().Delete(root.GetLastChild()).Build()
	replacement := CodeReplacement_Create(jutil.StringLength(before), jutil.StringLength(after), "")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testDelete_NonJSDocBlockCommentBeforeNode()
func (c *SuggestedFixTest) TestDelete_NonJSDocBlockCommentBeforeNode() {
	before := "var foo = new Bar();"
	after := "/* blah */ var baz = new Baz();"
	input := before + after
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix := NewSuggestedFix_Builder().Delete(root.GetLastChild()).Build()
	replacement := CodeReplacement_Create(jutil.StringLength(before), jutil.StringLength(after), "")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testDelete_dontDeleteSpaceBeforeNode()
func (c *SuggestedFixTest) TestDelete_dontDeleteSpaceBeforeNode() {
	before := "var foo = new Bar();\n\n"
	after := "var baz = new Baz();"
	input := before + after
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().
			DeleteWithoutRemovingWhitespaceBefore(root.GetLastChild()).
			Build()
	replacement := CodeReplacement_Create(jutil.StringLength(before), jutil.StringLength(after), "")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testDelete_multipleVarDeclaration()
func (c *SuggestedFixTest) TestDelete_multipleVarDeclaration() {
	input := "var foo = 3, bar, baz;"
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)

	// Delete the 1st variable on the line. Make sure the deletion includes the assignment and the
	// trailing comma.
	fix := NewSuggestedFix_Builder().Delete(root.GetFirstFirstChild()).Build()
	replacement := CodeReplacement_Create(4, jutil.StringLength("foo = 3, "), "")
	SuggestedFixes_AssertReplacement(fix, replacement)

	// Delete the 2nd variable.
	fix = NewSuggestedFix_Builder().Delete(root.GetFirstChild().GetSecondChild()).Build()
	replacement = CodeReplacement_Create(13, jutil.StringLength("bar, "), "")
	SuggestedFixes_AssertReplacement(fix, replacement)

	// Delete the last variable. Make sure it removes the leading comma.
	fix = NewSuggestedFix_Builder().Delete(root.GetFirstChild().GetLastChild()).Build()
	replacement = CodeReplacement_Create(16, jutil.StringLength(", baz"), "")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testDelete_multipleVarDeclarationWithJSDocComments()
func (c *SuggestedFixTest) TestDelete_multipleVarDeclarationWithJSDocComments() {
	input := "var /** fooComment*/foo = 3, bar, baz;"
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)

	// Delete the 1st variable on the line. Make sure the deletion includes the assignment and the
	// trailing comma.
	fix := NewSuggestedFix_Builder().Delete(root.GetFirstFirstChild()).Build()
	replacement :=
		CodeReplacement_Create(4, jutil.StringLength("/** fooComment*/foo = 3, "), "")
	SuggestedFixes_AssertReplacement(fix, replacement)

	// Delete the 2nd variable.
	fix = NewSuggestedFix_Builder().Delete(root.GetFirstChild().GetSecondChild()).Build()
	replacement = CodeReplacement_Create(jutil.StringIndexOf(input, "bar, "), jutil.StringLength("bar, "), "")
	SuggestedFixes_AssertReplacement(fix, replacement)

	// Delete the last variable. Make sure it removes the leading comma.
	fix = NewSuggestedFix_Builder().Delete(root.GetFirstChild().GetLastChild()).Build()
	replacement = CodeReplacement_Create(jutil.StringIndexOf(input, ", baz"), jutil.StringLength(", baz"), "")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testDelete_multipleVarDeclarationWithNonJSDocComments()
func (c *SuggestedFixTest) TestDelete_multipleVarDeclarationWithNonJSDocComments() {
	input := "var /* fooComment*/ foo = 3, bar, baz;"
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)

	// Delete the 1st variable on the line. Make sure the deletion includes the assignment and the
	// trailing comma.
	fix := NewSuggestedFix_Builder().Delete(root.GetFirstFirstChild()).Build()

	replacement :=
		CodeReplacement_Create(4, jutil.StringLength("/* fooComment*/ foo = 3, "), "")
	SuggestedFixes_AssertReplacement(fix, replacement)

	// Delete the 2nd variable.
	fix = NewSuggestedFix_Builder().Delete(root.GetFirstChild().GetSecondChild()).Build()
	replacement = CodeReplacement_Create(jutil.StringIndexOf(input, "bar, "), jutil.StringLength("bar, "), "")
	SuggestedFixes_AssertReplacement(fix, replacement)

	// Delete the last variable. Make sure it removes the leading comma.
	fix = NewSuggestedFix_Builder().Delete(root.GetFirstChild().GetLastChild()).Build()
	replacement = CodeReplacement_Create(jutil.StringIndexOf(input, ", baz"), jutil.StringLength(", baz"), "")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testDelete_multipleLetDeclaration()
func (c *SuggestedFixTest) TestDelete_multipleLetDeclaration() {
	input := "let foo = 3, bar, baz;"
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)

	// Delete the 1st variable on the line. Make sure the deletion includes the assignment and the
	// trailing comma.
	fix := NewSuggestedFix_Builder().Delete(root.GetFirstFirstChild()).Build()
	replacement := CodeReplacement_Create(4, jutil.StringLength("foo = 3, "), "")
	SuggestedFixes_AssertReplacement(fix, replacement)

	// Delete the 2nd variable.
	fix = NewSuggestedFix_Builder().Delete(root.GetFirstChild().GetSecondChild()).Build()
	replacement = CodeReplacement_Create(13, jutil.StringLength("bar, "), "")
	SuggestedFixes_AssertReplacement(fix, replacement)

	// Delete the last variable. Make sure it removes the leading comma.
	fix = NewSuggestedFix_Builder().Delete(root.GetFirstChild().GetLastChild()).Build()
	replacement = CodeReplacement_Create(16, jutil.StringLength(", baz"), "")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testDelete_multipleLetDeclarationWithNonJSDocComment()
func (c *SuggestedFixTest) TestDelete_multipleLetDeclarationWithNonJSDocComment() {
	input := "let foo = 3, /* blah */ bar, baz;"
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)

	// Delete the 1st variable on the line. Make sure the deletion includes the assignment and the
	// trailing comma.
	fix := NewSuggestedFix_Builder().Delete(root.GetFirstFirstChild()).Build()
	replacement := CodeReplacement_Create(4, jutil.StringLength("foo = 3, "), "")
	SuggestedFixes_AssertReplacement(fix, replacement)

	// Delete the 2nd variable.
	fix = NewSuggestedFix_Builder().Delete(root.GetFirstChild().GetSecondChild()).Build()
	replacement = CodeReplacement_Create(13, jutil.StringLength("/* blah */ bar, "), "")
	SuggestedFixes_AssertReplacement(fix, replacement)

	// Delete the last variable. Make sure it removes the leading comma.
	fix = NewSuggestedFix_Builder().Delete(root.GetFirstChild().GetLastChild()).Build()
	replacement = CodeReplacement_Create(jutil.StringIndexOf(input, ", baz"), jutil.StringLength(", baz"), "")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testDelete_multipleConstDeclaration()
func (c *SuggestedFixTest) TestDelete_multipleConstDeclaration() {
	input := "const foo = 3, bar = 4, baz = 5;"
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)

	// Delete the 1st variable on the line. Make sure the deletion includes the assignment and the
	// trailing comma.
	fix := NewSuggestedFix_Builder().Delete(root.GetFirstFirstChild()).Build()
	replacement := CodeReplacement_Create(6, jutil.StringLength("foo = 3, "), "")
	SuggestedFixes_AssertReplacement(fix, replacement)

	// Delete the 2nd variable.
	fix = NewSuggestedFix_Builder().Delete(root.GetFirstChild().GetSecondChild()).Build()
	replacement = CodeReplacement_Create(jutil.StringIndexOf(input, "bar = 4, "), jutil.StringLength("bar = 4, "), "")
	SuggestedFixes_AssertReplacement(fix, replacement)

	// Delete the last variable. Make sure it removes the leading comma.
	fix = NewSuggestedFix_Builder().Delete(root.GetFirstChild().GetLastChild()).Build()
	replacement = CodeReplacement_Create(jutil.StringIndexOf(input, ", baz = 5"), jutil.StringLength(", baz = 5"), "")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testRenameStringKey()
func (c *SuggestedFixTest) TestRenameStringKey() {
	input := "var obj = {foo: 'bar'};"
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	node := root.GetFirstFirstChild().GetFirstFirstChild()
	fix := NewSuggestedFix_Builder().Rename(node, "fooBar").Build()
	replacement := CodeReplacement_Create(11, jutil.StringLength("foo"), "fooBar")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testRenameProperty_justPropertyName()
func (c *SuggestedFixTest) TestRenameProperty_justPropertyName() {
	input := "obj.test.property"
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().Rename(root.GetFirstFirstChild(), "renamedProperty").Build()
	replacement := CodeReplacement_Create(9, jutil.StringLength("property"), "renamedProperty")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testRenameProperty_entireName()
func (c *SuggestedFixTest) TestRenameProperty_entireName() {
	input := "obj.test.property"
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().
			RenameWithReplaceNameSubtree(root.GetFirstFirstChild(), "renamedProperty", true).
			Build()
	replacement := CodeReplacement_Create(0, jutil.StringLength(input), "renamedProperty")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testRenameFunction_justFunctionName()
func (c *SuggestedFixTest) TestRenameFunction_justFunctionName() {
	input := "obj.fnCall();"
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().Rename(root.GetFirstFirstChild(), "renamedFnCall").Build()
	replacement := CodeReplacement_Create(4, jutil.StringLength("fnCall"), "renamedFnCall")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testRenameFunction_entireName()
func (c *SuggestedFixTest) TestRenameFunction_entireName() {
	fnName := "goog.dom.classes.add"
	newFnName := "goog.dom.classlist.add"
	input := fnName + "(foo, bar);"
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().RenameWithReplaceNameSubtree(root.GetFirstFirstChild(), newFnName, true).Build()
	replacement := CodeReplacement_Create(0, jutil.StringLength(fnName), newFnName)
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testRenameTaggedTemplate_justFunctionName()
func (c *SuggestedFixTest) TestRenameTaggedTemplate_justFunctionName() {
	prefix := "prt."
	fnName := "oldTaggedTemp"
	newFnName := "newTaggedTemp"
	input := prefix + fnName + "`${Foo}Bar`;"
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().Rename(root.GetFirstFirstChild(), newFnName).Build()
	replacement :=
		CodeReplacement_Create(jutil.StringLength(prefix), jutil.StringLength(fnName), newFnName)
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testRenameTaggedTemplate_entireName()
func (c *SuggestedFixTest) TestRenameTaggedTemplate_entireName() {
	fnName := "goog.dom.classes.oldTaggedTemp"
	newFnName := "goog.dom.classesList.newTaggedTemp"
	input := fnName + "`${Foo}Bar`;"
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().RenameWithReplaceNameSubtree(root.GetFirstFirstChild(), newFnName, true).Build()
	replacement := CodeReplacement_Create(0, jutil.StringLength(fnName), newFnName)
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testReplace()
func (c *SuggestedFixTest) TestReplace() {
	before := "var someRandomCode = {};\n/** some comment */\n"
	after := "goog.foo();"
	compiler := suggestedFixTest_getCompiler(before + after)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	newNode :=
		rhino.IR_ExprResult(rhino.IR_Call(rhino.IR_Getprop(rhino.IR_Name("goog2"), "get"), rhino.IR_String("service")))
	fix :=
		NewSuggestedFix_Builder().
			Replace(root.GetLastChild().GetFirstChild(), newNode, compiler).
			Build()
	replacement :=
		CodeReplacement_Create(jutil.StringLength(before), jutil.StringLength(after), "goog2.get('service');")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testReplace_functionArgument()
func (c *SuggestedFixTest) TestReplace_functionArgument() {
	before :=
		`var MyClass = function() {};
MyClass.prototype.foo = function() {};
MyClass.prototype.bar = function() {};
var clazz = new MyClass();
`
	after := "alert(clazz.foo());"
	compiler := suggestedFixTest_getCompiler(before + after)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	newNode := rhino.IR_Call(rhino.IR_Getprop(rhino.IR_Name("clazz"), "bar"))
	fix :=
		NewSuggestedFix_Builder().
			Replace(root.GetLastChild().GetFirstChild().GetLastChild(), newNode, compiler).
			Build()
	replacement :=
		CodeReplacement_Create(
			jutil.StringLength(before)+jutil.StringLength("alert("), jutil.StringLength("clazz.foo()"), "clazz.bar()")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testReplace_leftHandSideAssignment()
func (c *SuggestedFixTest) TestReplace_leftHandSideAssignment() {
	before := "var MyClass = function() {};\n"
	after := "MyClass.prototype.foo = function() {};\n"
	compiler := suggestedFixTest_getCompiler(before + after)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	newNode := rhino.IR_GetpropWithMoreProps(rhino.IR_Name("MyClass"), "prototype", "bar")
	fix :=
		NewSuggestedFix_Builder().
			Replace(root.GetLastChild().GetFirstFirstChild(), newNode, compiler).
			Build()
	replacement :=
		CodeReplacement_Create(
			jutil.StringLength(before), jutil.StringLength("MyClass.prototype.foo"), "MyClass.prototype.bar")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testReplace_lowerPrecedence()
func (c *SuggestedFixTest) TestReplace_lowerPrecedence() {
	before := "2 * "
	after := "3"
	compiler := suggestedFixTest_getCompiler(before + after + ";\n")
	root := suggestedFixTest_compileToScriptRoot(compiler)
	newNode := rhino.IR_ExprResult(rhino.IR_Sub(rhino.IR_Number(4), rhino.IR_Number(1)))
	toReplace := root.GetOnlyChild().GetOnlyChild().GetLastChild()
	fix := NewSuggestedFix_Builder().Replace(toReplace, newNode, compiler).Build()
	replacement :=
		CodeReplacement_Create(jutil.StringLength(before), jutil.StringLength(after), "(4 - 1)")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testReplace_samePrecedence()
func (c *SuggestedFixTest) TestReplace_samePrecedence() {
	before := "5 - "
	after := "1"
	compiler := suggestedFixTest_getCompiler(before + after + ";\n")
	root := suggestedFixTest_compileToScriptRoot(compiler)
	newNode := rhino.IR_ExprResult(rhino.IR_Sub(rhino.IR_Number(3), rhino.IR_Number(2)))
	toReplace := root.GetOnlyChild().GetOnlyChild().GetLastChild()
	fix := NewSuggestedFix_Builder().Replace(toReplace, newNode, compiler).Build()
	replacement :=
		CodeReplacement_Create(jutil.StringLength(before), jutil.StringLength(after), "(3 - 2)")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testReplace_higherPrecedence()
func (c *SuggestedFixTest) TestReplace_higherPrecedence() {
	before := "3 << "
	after := "2"
	compiler := suggestedFixTest_getCompiler(before + after + ";\n")
	root := suggestedFixTest_compileToScriptRoot(compiler)
	newNode := rhino.IR_ExprResult(rhino.IR_Add(rhino.IR_Number(1), rhino.IR_Number(1)))
	toReplace := root.GetOnlyChild().GetOnlyChild().GetLastChild()
	fix := NewSuggestedFix_Builder().Replace(toReplace, newNode, compiler).Build()
	replacement := CodeReplacement_Create(jutil.StringLength(before), jutil.StringLength(after), "1 + 1")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testAddCast()
func (c *SuggestedFixTest) TestAddCast() {
	input := "obj.fnCall();"
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	objNode := root.GetFirstFirstChild().GetFirstFirstChild()
	fix := NewSuggestedFix_Builder().AddCast(objNode, compiler, "FooBar").Build()
	replacement :=
		CodeReplacement_Create(0, jutil.StringLength("obj"), "/** @type {FooBar} */ (obj)")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testRemoveCast()
func (c *SuggestedFixTest) TestRemoveCast() {
	input := "var x = /** @type {string} */ (y);"
	expectedCode := "var x = y;"
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	castNode := root.GetFirstFirstChild().GetFirstChild()
	AssertThat(castNode.IsCast()).IsTrue()

	fix := NewSuggestedFix_Builder().RemoveCast(castNode, compiler).Build()
	SuggestedFixes_AssertChanges(fix, "", input, expectedCode)
}

// Java: SuggestedFixTest.testRemoveCast_complexStatement()
func (c *SuggestedFixTest) TestRemoveCast_complexStatement() {
	input :=
		`var x = /** @type {string} */ (function() {
  // Inline comment that should be preserved.
  var blah = bleh;
});
`
	expectedCode :=
		`var x = function() {
  // Inline comment that should be preserved.
  var blah = bleh;
};
`
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	castNode := root.GetFirstFirstChild().GetFirstChild()
	AssertThat(castNode.IsCast()).IsTrue()

	fix := NewSuggestedFix_Builder().RemoveCast(castNode, compiler).Build()
	SuggestedFixes_AssertChanges(fix, "", input, expectedCode)
}

// Java: SuggestedFixTest.testRemoveCast_return()
func (c *SuggestedFixTest) TestRemoveCast_return() {
	input :=
		`function f() {
  return /** @type {string} */ (
      'I am obviously a string. Why are you casting me?');
}
`
	expectedCode :=
		`function f() {
  return 'I am obviously a string. Why are you casting me?';
}
`
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	castNode := root.GetFirstChild().GetLastChild().GetFirstChild().GetLastChild()
	AssertThat(castNode.IsCast()).IsTrue()

	fix := NewSuggestedFix_Builder().RemoveCast(castNode, compiler).Build()
	SuggestedFixes_AssertChanges(fix, "", input, expectedCode)
}

// Java: SuggestedFixTest.testChangeJsDocType()
func (c *SuggestedFixTest) TestChangeJsDocType() {
	before := "/** "
	after := "@type {Foo} */\nvar foo = new Foo()"
	compiler := suggestedFixTest_getCompiler(before + after)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().
			ChangeJsDocType(root.GetFirstChild(), compiler, "Object").
			Build()
	replacement :=
		CodeReplacement_Create(jutil.StringLength(before), jutil.StringLength("@type {Foo}"), "@type {Object}")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testChangeJsDocType2()
func (c *SuggestedFixTest) TestChangeJsDocType2() {
	code := "/** @type {Foo} */\nvar foo = new Foo()"
	compiler := suggestedFixTest_getCompiler(code)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	varNode := root.GetFirstChild()
	jsdocRoot := suggestedFixTest_getOnlyElement(varNode.GetJSDocInfo().GetTypeNodes())
	fix := NewSuggestedFix_Builder().InsertBefore(jsdocRoot, "!").Build()
	replacement := CodeReplacement_Create(jutil.StringLength("/** @type {"), 0, "!")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testChangeJsDocType_packageType1()
func (c *SuggestedFixTest) TestChangeJsDocType_packageType1() {
	before := "/** "
	after := "@package {Foo} */\nvar foo = new Foo()"
	compiler := suggestedFixTest_getCompiler(before + after)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().
			ChangeJsDocType(root.GetFirstChild(), compiler, "Object").
			Build()
	replacement :=
		CodeReplacement_Create(jutil.StringLength(before), jutil.StringLength("@package {Foo}"), "@package {Object}")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testChangeJsDocType_privateType1()
func (c *SuggestedFixTest) TestChangeJsDocType_privateType1() {
	before := "/** "
	after := "@private {Foo} */\nvar foo = new Foo()"
	compiler := suggestedFixTest_getCompiler(before + after)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().
			ChangeJsDocType(root.GetFirstChild(), compiler, "Object").
			Build()
	replacement :=
		CodeReplacement_Create(jutil.StringLength(before), jutil.StringLength("@private {Foo}"), "@private {Object}")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testChangeJsDocType_privateType2()
func (c *SuggestedFixTest) TestChangeJsDocType_privateType2() {
	before := "/** @private "
	after := "@const {Foo} */\nvar foo = new Foo()"
	compiler := suggestedFixTest_getCompiler(before + after)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().
			ChangeJsDocType(root.GetFirstChild(), compiler, "Object").
			Build()
	replacement :=
		CodeReplacement_Create(jutil.StringLength(before), jutil.StringLength("@const {Foo}"), "@const {Object}")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testChangeJsDocType_privateType3()
func (c *SuggestedFixTest) TestChangeJsDocType_privateType3() {
	before := "/** @private "
	after := "@const {Foo} */\nvar foo = new Foo()"
	compiler := suggestedFixTest_getCompiler(before + after)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().
			ChangeJsDocType(root.GetFirstChild(), compiler, "Object").
			Build()
	replacement :=
		CodeReplacement_Create(jutil.StringLength(before), jutil.StringLength("@const {Foo}"), "@const {Object}")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testChangeJsDocType_privateType4()
func (c *SuggestedFixTest) TestChangeJsDocType_privateType4() {
	before := "/** "
	after := "@const {Foo} */\nvar foo = new Foo()"
	compiler := suggestedFixTest_getCompiler(before + after)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().
			ChangeJsDocType(root.GetFirstChild(), compiler, "Object").
			Build()
	replacement :=
		CodeReplacement_Create(jutil.StringLength(before), jutil.StringLength("@const {Foo}"), "@const {Object}")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testInsertArguments()
func (c *SuggestedFixTest) TestInsertArguments() {
	before := "goog.dom.classes.add("
	after := "foo, bar);"
	compiler := suggestedFixTest_getCompiler(before + after)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().InsertArguments(root.GetFirstFirstChild(), 0, "baz").Build()
	replacement := CodeReplacement_Create(jutil.StringLength(before), 0, "baz, ")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testInsertArguments_emptyArguments()
func (c *SuggestedFixTest) TestInsertArguments_emptyArguments() {
	before := "goog.dom.classes.add("
	after := ");"
	compiler := suggestedFixTest_getCompiler(before + after)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().InsertArguments(root.GetFirstFirstChild(), 0, "baz").Build()
	replacement := CodeReplacement_Create(jutil.StringLength(before), 0, "baz")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testInsertArguments_notFirstArgument()
func (c *SuggestedFixTest) TestInsertArguments_notFirstArgument() {
	before := "goog.dom.classes.add(foo, "
	after := "bar);"
	compiler := suggestedFixTest_getCompiler(before + after)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().InsertArguments(root.GetFirstFirstChild(), 1, "baz").Build()
	replacement := CodeReplacement_Create(jutil.StringLength(before), 0, "baz, ")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testInsertArguments_lastArgument()
func (c *SuggestedFixTest) TestInsertArguments_lastArgument() {
	before := "goog.dom.classes.add(foo, bar"
	after := ");"
	compiler := suggestedFixTest_getCompiler(before + after)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().InsertArguments(root.GetFirstFirstChild(), 2, "baz").Build()
	replacement := CodeReplacement_Create(jutil.StringLength(before), 0, ", baz")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testInsertArguments_beforeCastInArguments()
func (c *SuggestedFixTest) TestInsertArguments_beforeCastInArguments() {
	originalCode := "goog.dom.classes.add(foo, /** @type {String} */ (bar));"
	expectedCode := "goog.dom.classes.add(foo, baz, /** @type {String} */ (bar));"
	compiler := suggestedFixTest_getCompiler(originalCode)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().InsertArguments(root.GetFirstFirstChild(), 1, "baz").Build()
	SuggestedFixes_AssertChanges(fix, "", originalCode, expectedCode)
}

// Java: SuggestedFixTest.testInsertArguments_afterCastInArguments()
func (c *SuggestedFixTest) TestInsertArguments_afterCastInArguments() {
	originalCode := "goog.dom.classes.add(foo, /** @type {String} */ (bar));"
	expectedCode := "goog.dom.classes.add(foo, /** @type {String} */ (bar), baz);"
	compiler := suggestedFixTest_getCompiler(originalCode)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().InsertArguments(root.GetFirstFirstChild(), 2, "baz").Build()
	SuggestedFixes_AssertChanges(fix, "", originalCode, expectedCode)
}

// Java: SuggestedFixTest.testDeleteArgumentFirst()
func (c *SuggestedFixTest) TestDeleteArgumentFirst() {
	originalCode := "f(a, b, c);"
	expectedCode := "f(b, c);"

	compiler := suggestedFixTest_getCompiler(originalCode)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().DeleteArgument(root.GetFirstFirstChild(), 0).Build()

	SuggestedFixes_AssertChanges(fix, "", originalCode, expectedCode)
}

// Java: SuggestedFixTest.testDeleteArgumentFirstWithBlockComment()
func (c *SuggestedFixTest) TestDeleteArgumentFirstWithBlockComment() {
	originalCode := "f(/* blah */ a, b, c);"
	expectedCode := "f(b, c);"

	compiler := suggestedFixTest_getCompiler(originalCode)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().DeleteArgument(root.GetFirstFirstChild(), 0).Build()

	SuggestedFixes_AssertChanges(fix, "", originalCode, expectedCode)
}

// Java: SuggestedFixTest.testDeleteArgumentFirstWithLineComment()
func (c *SuggestedFixTest) TestDeleteArgumentFirstWithLineComment() {
	originalCode := "f(// blah \n a, b, c);"
	expectedCode := "f(b, c);"

	compiler := suggestedFixTest_getCompiler(originalCode)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().DeleteArgument(root.GetFirstFirstChild(), 0).Build()

	SuggestedFixes_AssertChanges(fix, "", originalCode, expectedCode)
}

// Java: SuggestedFixTest.testDeleteArgumentMiddle()
func (c *SuggestedFixTest) TestDeleteArgumentMiddle() {
	originalCode := "f(a, b, c);"
	expectedCode := "f(a, c);"

	compiler := suggestedFixTest_getCompiler(originalCode)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().DeleteArgument(root.GetFirstFirstChild(), 1).Build()

	SuggestedFixes_AssertChanges(fix, "", originalCode, expectedCode)
}

// Java: SuggestedFixTest.testDeleteArgumentMiddleWithBlockComment()
func (c *SuggestedFixTest) TestDeleteArgumentMiddleWithBlockComment() {
	originalCode := "f(a, /* blah */ b, c);"
	expectedCode := "f(a, c);"

	compiler := suggestedFixTest_getCompiler(originalCode)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().DeleteArgument(root.GetFirstFirstChild(), 1).Build()

	SuggestedFixes_AssertChanges(fix, "", originalCode, expectedCode)
}

// Java: SuggestedFixTest.testDeleteArgumentMiddleWithLineComment()
func (c *SuggestedFixTest) TestDeleteArgumentMiddleWithLineComment() {
	originalCode := "f(a, // blah \n b, c);"
	expectedCode := "f(a, c);"

	compiler := suggestedFixTest_getCompiler(originalCode)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().DeleteArgument(root.GetFirstFirstChild(), 1).Build()

	SuggestedFixes_AssertChanges(fix, "", originalCode, expectedCode)
}

// Java: SuggestedFixTest.testDeleteArgumentLast()
func (c *SuggestedFixTest) TestDeleteArgumentLast() {
	originalCode := "f(a, b, c);"
	expectedCode := "f(a, b);"

	compiler := suggestedFixTest_getCompiler(originalCode)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().DeleteArgument(root.GetFirstFirstChild(), 2).Build()

	SuggestedFixes_AssertChanges(fix, "", originalCode, expectedCode)
}

// Java: SuggestedFixTest.testDeleteArgumentLastBlockComment()
func (c *SuggestedFixTest) TestDeleteArgumentLastBlockComment() {
	originalCode := "f(a, b, /* blah */ c);"
	expectedCode := "f(a, b);"

	compiler := suggestedFixTest_getCompiler(originalCode)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().DeleteArgument(root.GetFirstFirstChild(), 2).Build()

	SuggestedFixes_AssertChanges(fix, "", originalCode, expectedCode)
}

// Java: SuggestedFixTest.testDeleteArgumentLastLineComment()
func (c *SuggestedFixTest) TestDeleteArgumentLastLineComment() {
	originalCode := "f(a, b, // blah \n c);"
	expectedCode := "f(a, b);"

	compiler := suggestedFixTest_getCompiler(originalCode)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().DeleteArgument(root.GetFirstFirstChild(), 2).Build()

	SuggestedFixes_AssertChanges(fix, "", originalCode, expectedCode)
}

// Java: SuggestedFixTest.testDeleteFirstArgumentWithPrefixComment()
func (c *SuggestedFixTest) TestDeleteFirstArgumentWithPrefixComment() {
	originalCode := "f(/** @type {number} */ (a), b, c);"
	expectedCode := "f(b, c);"

	compiler := suggestedFixTest_getCompiler(originalCode)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().DeleteArgument(root.GetFirstFirstChild(), 0).Build()

	SuggestedFixes_AssertChanges(fix, "", originalCode, expectedCode)
}

// Java: SuggestedFixTest.testDeleteFirstArgumentWithPostfixComment()
func (c *SuggestedFixTest) TestDeleteFirstArgumentWithPostfixComment() {
	originalCode := "f(a /** foo */, b, c);"
	expectedCode := "f(b, c);"

	compiler := suggestedFixTest_getCompiler(originalCode)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().DeleteArgument(root.GetFirstFirstChild(), 0).Build()

	SuggestedFixes_AssertChanges(fix, "", originalCode, expectedCode)
}

// Java: SuggestedFixTest.testDeleteArgumentWithPrefixComment()
func (c *SuggestedFixTest) TestDeleteArgumentWithPrefixComment() {
	originalCode := "f(a, /** @type {number} */ (b), c);"
	expectedCode := "f(a, c);"

	compiler := suggestedFixTest_getCompiler(originalCode)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	fix :=
		NewSuggestedFix_Builder().DeleteArgument(root.GetFirstFirstChild(), 1).Build()

	SuggestedFixes_AssertChanges(fix, "", originalCode, expectedCode)
}

// Java: SuggestedFixTest.testDeleteArgumentIndexTooLarge()
func (c *SuggestedFixTest) TestDeleteArgumentIndexTooLarge() {
	fnCall := "f(a, b, c);"
	compiler := suggestedFixTest_getCompiler(fnCall)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	AssertThrows[jutil.IIllegalArgumentException](
		func() { NewSuggestedFix_Builder().DeleteArgument(root.GetFirstFirstChild(), 3).Build() })
}

// Java: SuggestedFixTest.testDeleteArgumentIndexTooSmall()
func (c *SuggestedFixTest) TestDeleteArgumentIndexTooSmall() {
	fnCall := "f(a);"
	compiler := suggestedFixTest_getCompiler(fnCall)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	AssertThrows[jutil.IIllegalArgumentException](
		func() { NewSuggestedFix_Builder().DeleteArgument(root.GetFirstFirstChild(), -1).Build() })
}

// Java: SuggestedFixTest.testDeleteArgumentWithNoArguments()
func (c *SuggestedFixTest) TestDeleteArgumentWithNoArguments() {
	fnCall := "f();"
	compiler := suggestedFixTest_getCompiler(fnCall)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	AssertThrows[jutil.IIllegalStateException](
		func() { NewSuggestedFix_Builder().DeleteArgument(root.GetFirstFirstChild(), 0).Build() })
}

// Java: SuggestedFixTest.testAddGoogRequire_var()
func (c *SuggestedFixTest) TestAddGoogRequire_var() {
	before := "goog.provide('js.Foo');\n"
	after :=
		`goog.require('js.Bar');

var x;
/** @private */
function foo_() {};
`
	suggestedFixTest_assertAddGoogRequire(before, after, "abc.def")
}

// Java: SuggestedFixTest.testAddGoogRequire()
func (c *SuggestedFixTest) TestAddGoogRequire() {
	before := "goog.provide('js.Foo');\n"
	after :=
		`goog.require('js.Bar');

/** @private */
function foo_() {};
`
	suggestedFixTest_assertAddGoogRequire(before, after, "abc.def")
}

// Java: SuggestedFixTest.testAddGoogRequire_afterAllOtherGoogRequires()
func (c *SuggestedFixTest) TestAddGoogRequire_afterAllOtherGoogRequires() {
	before :=
		`goog.provide('js.Foo');
goog.require('js.Bar');
`
	after :=
		`/** @private */
function foo_() {};
`
	suggestedFixTest_assertAddGoogRequire(before, after, "zyx.w")
}

// Java: SuggestedFixTest.testAddGoogRequire_noGoogRequire()
func (c *SuggestedFixTest) TestAddGoogRequire_noGoogRequire() {
	before := "goog.provide('js.Foo');\n"
	after :=
		`/** @private */
function foo_() {};
`
	suggestedFixTest_assertAddGoogRequire(before, after, "abc.def")
}

// Java: SuggestedFixTest.testAddGoogRequire_noGoogRequireOrGoogProvide()
func (c *SuggestedFixTest) TestAddGoogRequire_noGoogRequireOrGoogProvide() {
	before := ""
	after :=
		`/** @private */
function foo_() {};
`
	suggestedFixTest_assertAddGoogRequire(before, after, "abc.def")
}

// Java: SuggestedFixTest.testAddGoogRequire_alreadyExists()
func (c *SuggestedFixTest) TestAddGoogRequire_alreadyExists() {
	input :=
		`goog.provide('js.Foo');
goog.require('abc.def');

/** @private */
function foo_() {};
`
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	match := NewMatch(root.GetFirstChild(), NewNodeMetadata(compiler))
	scriptMetadata := ScriptMetadata_Create(root, compiler)

	fix :=
		NewSuggestedFix_Builder().AddGoogRequire(match, "abc.def", scriptMetadata).Build()
	replacementMap := fix.GetReplacements()
	AssertThat(replacementMap.IsEmpty()).IsTrue() // MultimapSubject.isEmpty()
}

// Java: SuggestedFixTest.testAddGoogRequire_requireTypeAlreadyExists()
func (c *SuggestedFixTest) TestAddGoogRequire_requireTypeAlreadyExists() {
	input :=
		`goog.provide('js.Foo');
goog.requireType('abc.def');

/** @private */
function foo_() {};
`
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	match := NewMatch(root.GetFirstChild(), NewNodeMetadata(compiler))
	scriptMetadata := ScriptMetadata_Create(root, compiler)

	fix :=
		NewSuggestedFix_Builder().AddGoogRequire(match, "abc.def", scriptMetadata).Build()
	replacementMap := fix.GetReplacements()
	AssertThat(replacementMap.IsEmpty()).IsTrue() // MultimapSubject.isEmpty()
}

// Java: SuggestedFixTest.testAddRequireModule()
func (c *SuggestedFixTest) TestAddRequireModule() {
	input :=
		`goog.module('js.Foo');

/** @private */
function foo_() {
}
`
	expected :=
		`goog.module('js.Foo');
const safe = goog.require('goog.safe');

/** @private */
function foo_() {
}
`
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	match := NewMatch(root.GetFirstChild(), NewNodeMetadata(compiler))
	scriptMetadata := ScriptMetadata_Create(root, compiler)

	fix :=
		NewSuggestedFix_Builder().AddGoogRequire(match, "goog.safe", scriptMetadata).Build()
	SuggestedFixes_AssertChanges(fix, "", input, expected)
}

// Java: SuggestedFixTest.testAddRequireTypeModule()
func (c *SuggestedFixTest) TestAddRequireTypeModule() {
	input :=
		`goog.module('js.Foo');

/** @private */
function foo_() {
}
`
	expected :=
		`goog.module('js.Foo');
const Bar = goog.requireType('b.Bar');

/** @private */
function foo_() {
}
`
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	match := NewMatch(root.GetFirstChild(), NewNodeMetadata(compiler))
	scriptMetadata := ScriptMetadata_Create(root, compiler)

	fix :=
		NewSuggestedFix_Builder().AddGoogRequireType(match, "b.Bar", scriptMetadata).Build()
	SuggestedFixes_AssertChanges(fix, "", input, expected)
}

// Java: SuggestedFixTest.testAddRequireModule_requireTypeAlreadyExists()
func (c *SuggestedFixTest) TestAddRequireModule_requireTypeAlreadyExists() {
	input :=
		`goog.module('js.Foo');
goog.requireType('goog.safe');

/** @private */
function foo_() {
}
`
	expected :=
		`goog.module('js.Foo');
const safe = goog.requireType('goog.safe');

/** @private */
function foo_() {
}
`
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	match := NewMatch(root.GetFirstChild(), NewNodeMetadata(compiler))
	scriptMetadata := ScriptMetadata_Create(root, compiler)

	fix :=
		NewSuggestedFix_Builder().AddGoogRequire(match, "goog.safe", scriptMetadata).Build()
	SuggestedFixes_AssertChanges(fix, "", input, expected)
}

// Java: SuggestedFixTest.testAddRequireConst()
func (c *SuggestedFixTest) TestAddRequireConst() {
	input :=
		`const bar = goog.require('goog.bar');

/** @private */
function foo_() {};
`
	expected :=
		`const bar = goog.require('goog.bar');
const safe = goog.require('goog.safe');

/** @private */
function foo_() {};
`
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	match := NewMatch(root.GetFirstChild(), NewNodeMetadata(compiler))
	scriptMetadata := ScriptMetadata_Create(root, compiler)

	fix :=
		NewSuggestedFix_Builder().AddGoogRequire(match, "goog.safe", scriptMetadata).Build()
	SuggestedFixes_AssertChanges(fix, "", input, expected)
}

// Java: SuggestedFixTest.testAddRequireVar()
func (c *SuggestedFixTest) TestAddRequireVar() {
	input :=
		`var bar = goog.require('goog.bar');

/** @private */
function foo_() {};
`
	expected :=
		// We add new imports as const per the Google Style Guide
		// TODO(bangert): we could add complexity to add new imports as var if we want to.
		`var bar = goog.require('goog.bar');
const safe = goog.require('goog.safe');

/** @private */
function foo_() {};
`
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	match := NewMatch(root.GetFirstChild(), NewNodeMetadata(compiler))
	scriptMetadata := ScriptMetadata_Create(root, compiler)

	fix :=
		NewSuggestedFix_Builder().AddGoogRequire(match, "goog.safe", scriptMetadata).Build()
	SuggestedFixes_AssertChanges(fix, "", input, expected)
}

// Java: SuggestedFixTest.testAddRequireModuleUnchanged()
func (c *SuggestedFixTest) TestAddRequireModuleUnchanged() {
	input :=
		`goog.module('js.Foo');
const safe = goog.require('goog.safe');

/** @private */
function foo_() {};
`
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	match := NewMatch(root.GetFirstChild(), NewNodeMetadata(compiler))
	scriptMetadata := ScriptMetadata_Create(root, compiler)

	fix :=
		NewSuggestedFix_Builder().AddGoogRequire(match, "goog.safe", scriptMetadata).Build()
	AssertThat(fix.GetReplacements().IsEmpty()).IsTrue() // MultimapSubject.isEmpty()
}

// Java: SuggestedFixTest.testAddRequireModuleDifferentNameUnchanged()
func (c *SuggestedFixTest) TestAddRequireModuleDifferentNameUnchanged() {
	input :=
		`goog.module('js.Foo');
const googSafe = goog.require('goog.safe');

/** @private */
function foo_() {};
`
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	match := NewMatch(root.GetFirstChild(), NewNodeMetadata(compiler))
	scriptMetadata := ScriptMetadata_Create(root, compiler)

	fixBuilder :=
		NewSuggestedFix_Builder().AddGoogRequire(match, "goog.safe", scriptMetadata)
	AssertThat(fixBuilder.Build().GetReplacements().IsEmpty()).IsTrue() // MultimapSubject.isEmpty()
}

// Java: SuggestedFixTest.assertAddGoogRequire(String, String, String)
func suggestedFixTest_assertAddGoogRequire(before string, after string, namespace string) {
	compiler := suggestedFixTest_getCompiler(before + after)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	match := NewMatch(root.GetFirstChild(), NewNodeMetadata(compiler))
	scriptMetadata := ScriptMetadata_Create(root, compiler)

	fix :=
		NewSuggestedFix_Builder().AddGoogRequire(match, namespace, scriptMetadata).Build()
	replacement :=
		CodeReplacement_CreateWithSortKey(
			jutil.StringLength(before), 0, "goog.require('"+namespace+"');\n", namespace)
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testRemoveGoogRequire()
func (c *SuggestedFixTest) TestRemoveGoogRequire() {
	before :=
		`/** @fileoverview blah */

goog.provide('js.Foo');
`
	googRequire := "goog.require('abc.def');"
	input :=
		before +
			googRequire +
			`/** @private */
function foo_() {};
`
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	match := NewMatch(root.GetFirstChild(), NewNodeMetadata(compiler))
	fix := NewSuggestedFix_Builder().RemoveGoogRequire(match, "abc.def").Build()
	replacement := CodeReplacement_Create(jutil.StringLength(before), jutil.StringLength(googRequire), "")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testRemoveGoogRequireAmongSeveralGoogRequires()
func (c *SuggestedFixTest) TestRemoveGoogRequireAmongSeveralGoogRequires() {
	before :=
		`/** @fileoverview blah */

goog.provide('js.Foo');

goog.require('abc.abc');
`
	googRequire := "goog.require('abc.def');\n"
	input :=
		before +
			googRequire +
			`goog.require('def');

/** @private */
function foo_() {};
`
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	match := NewMatch(root.GetFirstChild(), NewNodeMetadata(compiler))
	fix := NewSuggestedFix_Builder().RemoveGoogRequire(match, "abc.def").Build()
	replacement := CodeReplacement_Create(jutil.StringLength(before), jutil.StringLength(googRequire), "")
	SuggestedFixes_AssertReplacement(fix, replacement)
}

// Java: SuggestedFixTest.testRemoveGoogRequire_doesNotExist()
func (c *SuggestedFixTest) TestRemoveGoogRequire_doesNotExist() {
	input :=
		`goog.require('abc.def');

/** @private */
function foo_() {};
`
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	match := NewMatch(root.GetFirstChild(), NewNodeMetadata(compiler))
	fix := NewSuggestedFix_Builder().RemoveGoogRequire(match, "fakefake").Build()
	replacementMap := fix.GetReplacements()
	AssertThat(replacementMap.IsEmpty()).IsTrue() // MultimapSubject.isEmpty()
}

// Java: SuggestedFixTest.testGenerateCode_doesNotJsEscapeStrings()
func (c *SuggestedFixTest) TestGenerateCode_doesNotJsEscapeStrings() {
	code := "someFunc('foo<>=%&\"bar');\n"
	compiler := suggestedFixTest_getCompiler(code)
	node := suggestedFixTest_compileToScriptRoot(compiler)
	generated := NewSuggestedFix_Builder().GenerateCode(compiler, node)
	AssertThat(generated).IsEqualTo(code)
}

// Java: SuggestedFixTest.testGenerateCode_var1()
func (c *SuggestedFixTest) TestGenerateCode_var1() {
	code := "var x;\nvar y;\n"
	compiler := suggestedFixTest_getCompiler(code)
	node := suggestedFixTest_compileToScriptRoot(compiler)
	generated := NewSuggestedFix_Builder().GenerateCode(compiler, node)
	AssertThat(generated).IsEqualTo(code)
}

// Java: SuggestedFixTest.testGenerateCode_var2()
func (c *SuggestedFixTest) TestGenerateCode_var2() {
	code := "var x = 0;\nvar y = 1;\n"
	compiler := suggestedFixTest_getCompiler(code)
	node := suggestedFixTest_compileToScriptRoot(compiler)
	generated := NewSuggestedFix_Builder().GenerateCode(compiler, node)
	AssertThat(generated).IsEqualTo(code)
}

// Java: SuggestedFixTest.testAttachMatchedNodeInfo()
func (c *SuggestedFixTest) TestAttachMatchedNodeInfo() {
	before :=
		`/** @fileoverview blah */

goog.provide('js.Foo');
`
	googRequire := "goog.require('abc.def');"
	input :=
		before +
			googRequire +
			`/** @private */
function foo_() {};
`
	compiler := suggestedFixTest_getCompiler(input)
	root := suggestedFixTest_compileToScriptRoot(compiler)
	match := NewMatch(root.GetFirstChild(), NewNodeMetadata(compiler))
	fix :=
		NewSuggestedFix_Builder().
			// Corresponds to the goog.require node.
			AttachMatchedNodeInfo(root.GetFirstChild().GetNext().GetFirstChild(), compiler).
			RemoveGoogRequire(match, "abc.def").
			Build()
	info := fix.GetMatchedNodeInfo()
	AssertThat(info.SourceFilename()).IsEqualTo("test")
	AssertThat(info.Lineno()).IsEqualTo(4)
	AssertThat(info.Charno()).IsEqualTo(0)
	AssertThat(info.InClosurizedFile()).IsTrue()
}

// suggestedFixTest_compileToScriptRoot returns the root script node
// produced from the compiled JS input.
//
// Java: SuggestedFixTest.compileToScriptRoot(Compiler)
func suggestedFixTest_compileToScriptRoot(compiler *Compiler) *rhino.Node {
	root := compiler.GetRoot()
	// The last child of the compiler root is a Block node, and the first
	// child of that is the Script node.
	return root.GetLastChild().GetFirstChild()
}

// Java: SuggestedFixTest.getCompiler(String)
func suggestedFixTest_getCompiler(jsInput string) *Compiler {
	compiler := NewCompiler()
	options := RefactoringDriver_GetCompilerOptions()
	compiler.Init(
		[]*SourceFile{}, // Externs
		[]*SourceFile{SourceFile_FromCode("test", jsInput)},
		options)
	compiler.Parse()
	return compiler
}

// suggestedFixTest_getOnlyElement is Iterables.getOnlyElement.
//
// Go-only: test-helper.
func suggestedFixTest_getOnlyElement(nodes []*rhino.Node) *rhino.Node {
	jutil.CheckArgumentMsg1(len(nodes) == 1, "expected one element but was: %s", len(nodes))
	return nodes[0]
}
