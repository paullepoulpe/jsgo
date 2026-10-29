// Port of com/google/javascript/jscomp/serialization/SerializeTypedAstPassTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/serialization/SerializeTypedAstPassTest.java)

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
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// SerializeTypedAstPassTest tests SerializeTypedAstPass.
//
// The Java @Test methods not named test* get a "Test" prefix in Go so that
// RunJUnit4 runs them.
//
// Java: public final class SerializeTypedAstPassTest extends CompilerTestCase
type SerializeTypedAstPassTest struct {
	CompilerTestCaseBase

	astConsumer func(*TypedAst)
	// individual test cases may override this
	typesToForwardDeclare *jutil.LinkedHashSet[string]
}

// serializeTypedAstPassTest_protoFields names the fields that ProtoTruth's
// ignoringFieldDescriptors(...) leaves out of a comparison.
//
// Go-only: language-adapter — ProtoTruth's field descriptors: the ignored
// fields are cleared, at every nesting level, on both sides before
// comparing with equals.
type serializeTypedAstPassTest_protoFields struct {
	objectUuid, objectOwnProperty               bool
	astRelativeLine, astRelativeColumn, astType bool
}

// Proto fields commonly ignored in tests because hardcoding their values is brittle
//
// Java: SerializeTypedAstPassTest.BRITTLE_TYPE_FIELDS (ObjectTypeProto uuid, own_property)
var serializeTypedAstPassTest_BRITTLE_TYPE_FIELDS = serializeTypedAstPassTest_protoFields{
	objectUuid: true, objectOwnProperty: true,
}

// Java: SerializeTypedAstPassTest()
func NewSerializeTypedAstPassTest() *SerializeTypedAstPassTest {
	c := &SerializeTypedAstPassTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestSerializeTypedAstPassTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewSerializeTypedAstPassTest() })
}

// Java: SerializeTypedAstPassTest.setUp()
func (c *SerializeTypedAstPassTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
	c.typesToForwardDeclare = jutil.NewLinkedHashSet[string]()
	c.EnableSourceInformationAnnotator()
	c.EnableDebugLogging(true)
}

// Java: SerializeTypedAstPassTest.getProcessor(Compiler)
func (c *SerializeTypedAstPassTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewSerializeTypedAstPass(
		compiler, c.astConsumer, SerializationOptions_Builder_().SetIncludeDebugInfo(false).Build())
}

// Java: SerializeTypedAstPassTest.createCompiler()
func (c *SerializeTypedAstPassTest) CreateCompiler() *Compiler {
	compiler := c.CompilerTestCaseBase.CreateCompiler()
	for typeName := range c.typesToForwardDeclare.All() {
		compiler.ForwardDeclareType(typeName)
	}
	return compiler
}

// Java: SerializeTypedAstPassTest.testAst_constNumber()
func (c *SerializeTypedAstPassTest) TestAst_constNumber() {
	numberType := PrimitiveType_NUMBER_TYPE.GetNumber()
	result := c.compile("const x = 5;")

	serializeTypedAstPassTest_assertAstNode(result.sourceNodes[0],
		serializeTypedAstPassTest_BRITTLE_TYPE_FIELDS,
		&AstNode{
			kind_:         NodeKind_SOURCE_FILE_VALUE,
			relativeLine_: 1,
			child_: []*AstNode{{
				kind_: NodeKind_CONST_DECLARATION_VALUE,
				child_: []*AstNode{{
					kind_:                NodeKind_IDENTIFIER_VALUE,
					stringValuePointer_:  result.findInStringPool("x"),
					valueCase_:           AstNode_STRING_VALUE_POINTER_FIELD_NUMBER,
					originalNamePointer_: result.findInStringPool("x"),
					relativeColumn_:      6,
					type_:                numberType,
					bitField0_:           1,
					child_: []*AstNode{{
						kind_:           NodeKind_NUMBER_LITERAL_VALUE,
						doubleValue_:    5,
						valueCase_:      AstNode_DOUBLE_VALUE_FIELD_NUMBER,
						relativeColumn_: 4,
						type_:           numberType,
						bitField0_:      1,
					}},
				}},
			}},
		})
}

// Java: SerializeTypedAstPassTest.testAst_externs()
func (c *SerializeTypedAstPassTest) TestAst_externs() {
	// 2 externs files, the default (empty) externs file + the source file marked with @externs
	AssertThat(c.compile("/** @externs */ function foo() {}").externNodes).HasSize(2)
}

// Java: SerializeTypedAstPassTest.testAst_typeSummaryFiles_areNotSerialized_orSearchedForTypes()
func (c *SerializeTypedAstPassTest) TestAst_typeSummaryFiles_areNotSerialized_orSearchedForTypes() {
	c.EnableCreateModuleMap()
	// SerializeTypedAstPass will drop the type summary files, which live on the externs root, to
	// avoid serializing them.
	c.AllowExternsChanges()

	closureExterns := c.Externs(NewTestExternsBuilder().AddClosureExterns().Build())
	mandatorySource := SourceFile_FromCode("mandatory.js", "/* mandatory source */")

	typeSummaryResult := c.compile(
		closureExterns,
		c.Srcs(
			mandatorySource,
			SourceFile_FromCode(
				"foo.i.js",
				`/** @fileoverview @typeSummary */

goog.loadModule(function(exports) {
  goog.module('a.Foo');
  class Foo { }
  exports = Foo;
});
`)))
	emptyResult := c.compile(
		closureExterns, //
		c.Srcs(mandatorySource))

	AssertThat(typeSummaryResult.ast.String()).IsEqualTo(emptyResult.ast.String())
}

// Java: SerializeTypedAstPassTest.testAst_letString()
func (c *SerializeTypedAstPassTest) TestAst_letString() {
	stringType := PrimitiveType_STRING_TYPE.GetNumber()
	result := c.compile("let s = 'hello';")

	ignored := serializeTypedAstPassTest_BRITTLE_TYPE_FIELDS
	ignored.astRelativeLine = true
	ignored.astRelativeColumn = true
	serializeTypedAstPassTest_assertAstNode(result.sourceNodes[0], ignored,
		&AstNode{
			kind_: NodeKind_SOURCE_FILE_VALUE,
			child_: []*AstNode{{
				kind_: NodeKind_LET_DECLARATION_VALUE,
				child_: []*AstNode{{
					kind_:                NodeKind_IDENTIFIER_VALUE,
					stringValuePointer_:  result.findInStringPool("s"),
					valueCase_:           AstNode_STRING_VALUE_POINTER_FIELD_NUMBER,
					originalNamePointer_: result.findInStringPool("s"),
					type_:                stringType,
					bitField0_:           1,
					child_: []*AstNode{{
						kind_:               NodeKind_STRING_LITERAL_VALUE,
						stringValuePointer_: result.findInStringPool("hello"),
						valueCase_:          AstNode_STRING_VALUE_POINTER_FIELD_NUMBER,
						type_:               stringType,
						bitField0_:          1,
					}},
				}},
			}},
		})
}

// Java: SerializeTypedAstPassTest.testAst_templateLit_illegalEscape()
func (c *SerializeTypedAstPassTest) TestAst_templateLit_illegalEscape() {
	// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Template_literals#es2018_revision_of_illegal_escape_sequences
	result := c.compile("latex`\\unicode`;")

	ignored := serializeTypedAstPassTest_BRITTLE_TYPE_FIELDS
	ignored.astRelativeLine = true
	ignored.astRelativeColumn = true
	ignored.astType = true
	serializeTypedAstPassTest_assertAstNode(result.sourceNodes[0], ignored,
		&AstNode{
			kind_: NodeKind_SOURCE_FILE_VALUE,
			child_: []*AstNode{{
				kind_: NodeKind_EXPRESSION_STATEMENT_VALUE,
				child_: []*AstNode{{
					kind_:              NodeKind_TAGGED_TEMPLATELIT_VALUE,
					booleanProperties_: int64(1) << rhino.NodeProperty_FREE_CALL.GetNumber(),
					child_: []*AstNode{
						{
							kind_:                NodeKind_IDENTIFIER_VALUE,
							stringValuePointer_:  result.findInStringPool("latex"),
							valueCase_:           AstNode_STRING_VALUE_POINTER_FIELD_NUMBER,
							originalNamePointer_: result.findInStringPool("latex"),
						},
						{
							kind_: NodeKind_TEMPLATELIT_VALUE,
							child_: []*AstNode{{
								kind_:      NodeKind_TEMPLATELIT_STRING_VALUE,
								valueCase_: AstNode_TEMPLATE_STRING_VALUE_FIELD_NUMBER,
								templateStringValue_: &TemplateStringValue{
									cookedStringPointer_: -1,
									rawStringPointer_:    result.findInStringPool("\\unicode"),
								},
							}},
						},
					},
				}},
			}},
		})
}

// Java: SerializeTypedAstPassTest.testAst_numberInCast()
func (c *SerializeTypedAstPassTest) TestAst_numberInCast() {
	// CAST nodes in JSCompiler are a combination of a child node + JSDoc @type. Because we don't
	// serialize JSDoc @types it doesn't make sense to serialize the CAST node.
	unknownType := PrimitiveType_UNKNOWN_TYPE.GetNumber()
	ignored := serializeTypedAstPassTest_BRITTLE_TYPE_FIELDS
	ignored.astRelativeLine = true
	ignored.astRelativeColumn = true
	serializeTypedAstPassTest_assertAstNode(c.compileToAstNode("/** @type {?} */ (1);"), ignored,
		&AstNode{
			kind_: NodeKind_SOURCE_FILE_VALUE,
			child_: []*AstNode{{
				kind_: NodeKind_EXPRESSION_STATEMENT_VALUE,
				child_: []*AstNode{{
					kind_: NodeKind_NUMBER_LITERAL_VALUE,
					booleanProperties_: (int64(1) << rhino.NodeProperty_IS_PARENTHESIZED.GetNumber()) |
						(int64(1) << rhino.NodeProperty_COLOR_FROM_CAST.GetNumber()),
					doubleValue_: 1,
					valueCase_:   AstNode_DOUBLE_VALUE_FIELD_NUMBER,
					type_:        unknownType,
					bitField0_:   1,
				}},
			}},
		})
}

// Java: SerializeTypedAstPassTest.testAst_arrowFunction()
func (c *SerializeTypedAstPassTest) TestAst_arrowFunction() {
	result := c.compile("() => 'hello';")

	serializeTypedAstPassTest_assertAstNode(result.sourceNodes[0],
		serializeTypedAstPassTest_protoFields{astRelativeLine: true, astType: true, astRelativeColumn: true},
		&AstNode{
			kind_: NodeKind_SOURCE_FILE_VALUE,
			child_: []*AstNode{{
				kind_: NodeKind_EXPRESSION_STATEMENT_VALUE,
				child_: []*AstNode{{
					kind_:              NodeKind_FUNCTION_LITERAL_VALUE,
					booleanProperties_: int64(1) << rhino.NodeProperty_ARROW_FN.GetNumber(),
					child_: []*AstNode{
						{
							kind_:               NodeKind_IDENTIFIER_VALUE,
							stringValuePointer_: 0,
							valueCase_:          AstNode_STRING_VALUE_POINTER_FIELD_NUMBER,
						},
						{kind_: NodeKind_PARAMETER_LIST_VALUE},
						{
							kind_:               NodeKind_STRING_LITERAL_VALUE,
							stringValuePointer_: result.findInStringPool("hello"),
							valueCase_:          AstNode_STRING_VALUE_POINTER_FIELD_NUMBER,
						},
					},
				}},
			}},
		})
}

// Java: SerializeTypedAstPassTest.testRewrittenGoogModule_containsOriginalNames()
func (c *SerializeTypedAstPassTest) TestRewrittenGoogModule_containsOriginalNames() {
	c.EnableRewriteClosureCode()

	result := c.compile("goog.module('a.b.c'); function f(x) {}")

	serializeTypedAstPassTest_assertAstNode(result.sourceNodes[0],
		serializeTypedAstPassTest_protoFields{astRelativeLine: true, astRelativeColumn: true, astType: true},
		&AstNode{
			kind_:              NodeKind_SOURCE_FILE_VALUE,
			booleanProperties_: int64(1) << rhino.NodeProperty_GOOG_MODULE.GetNumber(),
			child_: []*AstNode{
				// `/** @const */ var module$exports$a$b$c = {};`
				{
					kind_:              NodeKind_VAR_DECLARATION_VALUE,
					booleanProperties_: int64(1) << rhino.NodeProperty_IS_NAMESPACE.GetNumber(),
					jsdoc_:             &OptimizationJsdoc{kind_: []int32{JsdocTag_JSDOC_CONST_VALUE}},
					child_: []*AstNode{{
						kind_:               NodeKind_IDENTIFIER_VALUE,
						stringValuePointer_: result.findInStringPool("module$exports$a$b$c"),
						valueCase_:          AstNode_STRING_VALUE_POINTER_FIELD_NUMBER,
						child_:              []*AstNode{{kind_: NodeKind_OBJECT_LITERAL_VALUE}},
					}},
				},
				// `function module$contents$a$b$c_f(x) {}`
				{
					kind_:                NodeKind_FUNCTION_LITERAL_VALUE,
					originalNamePointer_: result.findInStringPool("f"),
					child_: []*AstNode{
						{
							kind_:                NodeKind_IDENTIFIER_VALUE,
							stringValuePointer_:  result.findInStringPool("module$contents$a$b$c_f"),
							valueCase_:           AstNode_STRING_VALUE_POINTER_FIELD_NUMBER,
							originalNamePointer_: result.findInStringPool("f"),
						},
						{
							kind_: NodeKind_PARAMETER_LIST_VALUE,
							child_: []*AstNode{{
								kind_:                NodeKind_IDENTIFIER_VALUE,
								stringValuePointer_:  result.findInStringPool("x"),
								valueCase_:           AstNode_STRING_VALUE_POINTER_FIELD_NUMBER,
								originalNamePointer_: result.findInStringPool("x"),
							}},
						},
						{kind_: NodeKind_BLOCK_VALUE},
					},
				},
			},
		})
}

// Java: SerializeTypedAstPassTest.serializesExternsFile()
func (c *SerializeTypedAstPassTest) TestSerializesExternsFile() {
	c.EnsureLibraryInjected("base")

	result := c.compileWithExterns(
		NewTestExternsBuilder().AddMath().AddObject().Build(), "let s = 'hello';")

	externAst := result.ast.GetExternAstList()[0]
	AssertThat(
		result.
			ast.
			GetSourceFilePool().
			GetSourceFileList()[externAst.GetSourceFile()-1].
			GetFilename()).
		IsEqualTo("externs")
}

// Java: SerializeTypedAstPassTest.serializesSourceOfSyntheticCode()
func (c *SerializeTypedAstPassTest) TestSerializesSourceOfSyntheticCode() {
	c.EnsureLibraryInjected("base")

	result := c.compileWithExterns(
		NewTestExternsBuilder().AddMath().AddObject().Build(), "let s = 'hello';")

	lazyAst := result.ast.GetCodeAstList()[0]
	AssertThat(
		result.
			ast.
			GetSourceFilePool().
			GetSourceFileList()[lazyAst.GetSourceFile()-1].
			GetFilename()).
		IsEqualTo("testcode")

	// 'base' is loaded at the beginning of the first script
	script := result.sourceNodes[0]
	baseLibraryStart := script.GetChild(0)
	AssertThat(
		result.
			ast.
			GetSourceFilePool().
			GetSourceFileList()[baseLibraryStart.GetSourceFile()-1].
			GetFilename()).
		EndsWith("js/base.js")
	// children of the synthetic subtree default to the root's file [synthetic:base]
	AssertThat(baseLibraryStart.GetChild(0).GetSourceFile()).IsEqualTo(int32(0))
	// "let s = 'hello'" defaults to the parent's file 'testcode'
	AssertThat(script.GetChild(script.GetChildCount() - 1).GetSourceFile()).IsEqualTo(int32(0))
}

// Java: SerializeTypedAstPassTest.serializesExternProperties()
func (c *SerializeTypedAstPassTest) TestSerializesExternProperties() {
	c.EnableGatherExternProperties()

	result := c.compileWithExterns(
		`class Foo {
// Ensure JSDoc properties are included (b/180424427)
  /** @param {{arg: string}} x */
  method(x) { }
}
`,
		"")

	AssertThat(result.ast.GetExternsSummary().GetPropNamePtrList()).
		ContainsAtLeast(result.findInStringPool("method"), result.findInStringPool("arg"))
}

// Java: SerializeTypedAstPassTest.serializeInlineSourceMappingURL()
func (c *SerializeTypedAstPassTest) TestSerializeInlineSourceMappingURL() {
	// We want TypedAST to support inline source maps (which are input source maps passed
	// by embedding a `//# sourceMappingURL=<url>` where <url> is a base64-encoded "data url").
	sourceMapTestCode := `var X = (function () {
    function X(input) {
        this.y = input;
    }
    return X;
}());
`

	base64Prefix := "data:application/json;base64,"
	encodedSourceMap := "eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZm9vLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiZm9vLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUFBO0lBR0UsV0FBWSxLQUFhO1FBQ3ZCLElBQUksQ0FBQyxDQUFDLEdBQUcsS0FBSyxDQUFDO0lBQ2pCLENBQUM7SUFDSCxRQUFDO0FBQUQsQ0FBQyxBQU5ELElBTUM7QUFFRCxPQUFPLENBQUMsR0FBRyxDQUFDLElBQUksQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMifQ=="

	code := sourceMapTestCode + "\n//# sourceMappingURL=" + base64Prefix + encodedSourceMap

	result := c.compile(code)

	lazyAst := result.ast.GetCodeAstList()[0]
	sourceMappingURL := lazyAst.GetSourceMappingUrl()

	AssertThat(sourceMappingURL).IsEqualTo(encodedSourceMap) // Do not serizile the base-64 prefix.
}

// Java: SerializeTypedAstPassTest.doesNotSerializeInputSourceMappingURL()
func (c *SerializeTypedAstPassTest) TestDoesNotSerializeInputSourceMappingURL() {
	// We do not want TypedAST to support input source maps (source maps in separate files).
	sourceMapTestCode := `var X = (function () {
    function X(input) {
        this.y = input;
    }
    return X;
}());
`

	code := sourceMapTestCode + "\n//# sourceMappingURL=foo.js.map"

	result := c.compile(code)

	lazyAst := result.ast.GetCodeAstList()[0]
	sourceMappingURL := lazyAst.GetSourceMappingUrl()

	AssertThat(sourceMappingURL).IsEmpty()
}

// Java: SerializeTypedAstPassTest.compileToAstNode(String)
func (c *SerializeTypedAstPassTest) compileToAstNode(source string) *AstNode {
	return c.compile(source).sourceNodes[0]
}

// Java: SerializeTypedAstPassTest.compile(String)
// Java: SerializeTypedAstPassTest.compile(TestPart...)
func (c *SerializeTypedAstPassTest) compile(args ...any) *serializeTypedAstPassTest_SerializationResult {
	if compilerTestCase_argsAre(args, "string") {
		return c.compileWithExterns(CompilerTestCase_DEFAULT_EXTERNS, args[0].(string))
	}
	resultAst := make([]*TypedAst, 1)
	c.astConsumer = func(ast *TypedAst) { resultAst[0] = ast }
	c.Test(args...)
	return newSerializeTypedAstPassTest_SerializationResult(resultAst[0])
}

// Java: private static class SerializeTypedAstPassTest.SerializationResult
type serializeTypedAstPassTest_SerializationResult struct {
	externNodes []*AstNode
	sourceNodes []*AstNode
	ast         *TypedAst
}

// Java: SerializeTypedAstPassTest.SerializationResult(TypedAst)
func newSerializeTypedAstPassTest_SerializationResult(ast *TypedAst) *serializeTypedAstPassTest_SerializationResult {
	r := &serializeTypedAstPassTest_SerializationResult{}
	for _, lazyAst := range ast.GetExternAstList() {
		r.externNodes = append(r.externNodes, serializeTypedAstPassTest_parseAstNode(lazyAst.GetScript()))
	}
	for _, lazyAst := range ast.GetCodeAstList() {
		r.sourceNodes = append(r.sourceNodes, serializeTypedAstPassTest_parseAstNode(lazyAst.GetScript()))
	}
	r.ast = ast
	return r
}

// findInStringPool returns the offset of the given string in this ast's
// string pool, throwing an exception if not present.
//
// Java: SerializeTypedAstPassTest.SerializationResult.findInStringPool(String)
func (r *serializeTypedAstPassTest_SerializationResult) findInStringPool(str string) int32 {
	if str == "" {
		return 0
	}

	// getStringsList().indexOf(ByteString.copyFromUtf8(str)) + 1
	offset := 0
	for i, s := range r.ast.GetStringPool().GetStringsList() {
		if bytes.Equal(s, []byte(str)) {
			offset = i + 1
			break
		}
	}
	if offset == -1 {
		jutil.CheckStateMsg2(
			false,
			"Could not find string '%s' in string pool %s",
			str,
			r.ast.GetStringPool())
	}
	return int32(offset)
}

// Java: SerializeTypedAstPassTest.parseAstNode(ByteString)
func serializeTypedAstPassTest_parseAstNode(byteString []byte) *AstNode {
	return AstNode_ParseFrom(byteString)
}

// Java: SerializeTypedAstPassTest.compileWithExterns(String, String)
func (c *SerializeTypedAstPassTest) compileWithExterns(externs string, source string) *serializeTypedAstPassTest_SerializationResult {
	return c.compile(c.Externs(externs), c.Srcs(source))
}

// Java: SerializeTypedAstPassTest.generateDiagnosticFiles()
func (c *SerializeTypedAstPassTest) generateDiagnosticFiles() {
	c.compile(
		`class Foo0 {
  w() { }
}
class Foo1 extends Foo0 {
  y() { }
}
const ns = {Foo0, Foo1};
/** @suppress {checkTypes} */
const /** !Foo1 */ typeMismatch = new Foo0();
`)
}

// Java: SerializeTypedAstPassTest.loadFile(Path)
func serializeTypedAstPassTest_loadFile(path string) string {
	// Files.lines(path).collect(joining("\n"))
	data, err := os.ReadFile(path)
	if err != nil {
		e := jutil.NewAssertionError(err.Error())
		e.InitCause(jutil.NewIOException(err.Error()))
		panic(e)
	}
	text := strings.ReplaceAll(strings.ReplaceAll(string(data), "\r\n", "\n"), "\r", "\n")
	return strings.TrimSuffix(text, "\n")
}

// Java: SerializeTypedAstPassTest.debugLogFiles()
func (c *SerializeTypedAstPassTest) debugLogFiles() []string {
	dir := filepath.Join(
		c.GetLastCompiler().GetOptions().GetDebugLogDirectory().String(),
		"SerializeTypesToPointers")
	entries, err := os.ReadDir(dir)
	if err != nil {
		e := jutil.NewAssertionError(err.Error())
		e.InitCause(jutil.NewIOException(err.Error()))
		panic(e)
	}
	// Go-only: language-adapter — Files.list's order is unspecified; sorted.
	var files []string
	for _, e := range entries {
		files = append(files, filepath.Join(dir, e.Name()))
	}
	sort.Strings(files)
	return files
}

// Java: SerializeTypedAstPassTest.assertValidJson(String)
func serializeTypedAstPassTest_assertValidJson(src string) {
	AssertThat(src).IsNotEmpty()

	// Go-only: language-adapter — new Gson().fromJson(src, LinkedHashMap.class / ArrayList.class)
	// is parsed as a JsonObject / JsonArray; both throw if invalid.
	if jutil.StringCharAt(src, 0) == '{' {
		jutil.NewGson().FromJsonObject(src)
	} else {
		jutil.NewGson().FromJsonArray(src)
	}
}

// serializeTypedAstPassTest_assertAstNode is
// assertThat(actual).ignoringFieldDescriptors(ignored).isEqualTo(expected).
//
// Go-only: language-adapter — ProtoTruth (see
// serializeTypedAstPassTest_protoFields).
func serializeTypedAstPassTest_assertAstNode(actual *AstNode, ignored serializeTypedAstPassTest_protoFields, expected *AstNode) {
	AssertThat(serializeTypedAstPassTest_clearAstNode(actual, ignored).String()).
		IsEqualTo(serializeTypedAstPassTest_clearAstNode(expected, ignored).String())
}

// serializeTypedAstPassTest_clearAstNode returns node with the ignored
// fields cleared at every level. (AstNode holds no ObjectTypeProto, so the
// BRITTLE_TYPE_FIELDS do not apply to it.)
func serializeTypedAstPassTest_clearAstNode(node *AstNode, ignored serializeTypedAstPassTest_protoFields) *AstNode {
	cp := *node
	cp.cachedSize = 0
	if ignored.astRelativeLine {
		cp.relativeLine_ = 0
	}
	if ignored.astRelativeColumn {
		cp.relativeColumn_ = 0
	}
	if ignored.astType {
		cp.type_ = 0
		cp.bitField0_ &^= 1
	}
	if len(node.child_) > 0 {
		children := make([]*AstNode, len(node.child_))
		for i, child := range node.child_ {
			children[i] = serializeTypedAstPassTest_clearAstNode(child, ignored)
		}
		cp.child_ = children
	}
	return &cp
}

// serializeTypedAstPassTest_clearTypeProto returns type_ with the ignored
// ObjectTypeProto fields cleared.
func serializeTypedAstPassTest_clearTypeProto(type_ *TypeProto, ignored serializeTypedAstPassTest_protoFields) *TypeProto {
	if !type_.HasObject() {
		return type_
	}
	obj := *type_.GetObject()
	obj.cachedSize = 0
	if ignored.objectUuid {
		obj.uuid_ = nil
	}
	if ignored.objectOwnProperty {
		obj.ownProperty_ = nil
	}
	tp := *type_
	tp.cachedSize = 0
	tp.object_ = &obj
	return &tp
}
