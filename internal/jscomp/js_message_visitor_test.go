// Port of com/google/javascript/jscomp/JsMessageVisitorTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/JsMessageVisitorTest.java)

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
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

// Java: private static class JsMessageVisitorTest.RenameMessagesVisitor extends AbstractPostOrderCallback
type jsMessageVisitorTest_RenameMessagesVisitor struct {
	NodeTraversal_AbstractPostOrderCallbackBase
}

// Java: JsMessageVisitorTest.RenameMessagesVisitor()
func newJsMessageVisitorTest_RenameMessagesVisitor() *jsMessageVisitorTest_RenameMessagesVisitor {
	return &jsMessageVisitorTest_RenameMessagesVisitor{}
}

// Java: JsMessageVisitorTest.RenameMessagesVisitor.visit(NodeTraversal, Node, Node)
func (v *jsMessageVisitorTest_RenameMessagesVisitor) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	// Node.getString() is never null for a NAME (the Java null check is kept as a no-op).
	if n.IsName() && strings.HasPrefix(n.GetString(), "MSG_") {
		originalName := n.GetString()
		n.SetOriginalName(jutil.StringPtr(originalName))
		n.SetString("some_prefix_" + originalName)
	} else if n.IsGetProp() && parent.IsAssign() && strings.Contains(*n.GetQualifiedName(), ".MSG_") {
		originalName := n.GetString()
		n.SetOriginalName(jutil.StringPtr(originalName))
		n.SetString("some_prefix_" + originalName)
	}
}

// JsMessageVisitorTest tests JsMessageVisitor.
//
// Java: public final class JsMessageVisitorTest
type JsMessageVisitorTest struct {
	// Go-only: test-helper — the Go test running the current test method.
	t *testing.T

	compilerOptions        *CompilerOptions
	compiler               *Compiler
	messages               []*JsMessage
	messageDefinitions     []JsMessageDefinition
	icuTemplateDefinitions []IcuTemplateDefinition
	renameMessages         bool
}

// Java: JsMessageVisitorTest()
func NewJsMessageVisitorTest() *JsMessageVisitorTest {
	return &JsMessageVisitorTest{}
}

func TestJsMessageVisitorTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewJsMessageVisitorTest() })
}

// JUnitSetT gives the instance the Go test running its test method.
func (c *JsMessageVisitorTest) JUnitSetT(t *testing.T) { c.t = t }

// JUnitBefore runs the @Before method.
func (c *JsMessageVisitorTest) JUnitBefore() { c.SetUp() }

// JUnitAfter: the class has no @After method.
func (c *JsMessageVisitorTest) JUnitAfter() {}

// Java: JsMessageVisitorTest.setUp()
func (c *JsMessageVisitorTest) SetUp() {
	c.messages = []*JsMessage{}
	c.messageDefinitions = []JsMessageDefinition{}
	c.icuTemplateDefinitions = []IcuTemplateDefinition{}
	c.compilerOptions = nil
	c.renameMessages = false
}

// Java: JsMessageVisitorTest.testIcuTemplateParsing()
func (c *JsMessageVisitorTest) TestIcuTemplateParsing() {
	icuMessageTemplateString :=
		NewJsMessageVisitor_IcuMessageTemplateString(
			// Some placeholders are ignored, because there's no additional information for them
			// (original code or example text). There's no need to generate placeholder parts for
			// those.
			"{NUM_PEOPLE, plural, offset:1\n" +
				"=0 {I see {START_BOLD}no one at all{END_BOLD} in {INTERPOLATION_2}.}\n" +
				"=1 {I see {INTERPOLATION_1} in {INTERPOLATION_2}.}\n" +
				"=2 {I see {INTERPOLATION_1} and one other person in {INTERPOLATION_2}.}\n" +
				"other {I see {INTERPOLATION_1} and # other people in {INTERPOLATION_2}.}\n" +
				"}\n")
	extractedParts :=
		icuMessageTemplateString.ExtractParts(
			jutil.NewLinkedHashSet("INTERPOLATION_1", "INTERPOLATION_2", "MISSING"))

	// The non-existent placeholder name "MISSING" should be missing from this list.
	// JsMessageVisitor will notice the name is missing from the resulting set of
	// seen placeholders and report an error, though we won't do that in this test.
	AssertThat(extractedParts.extractedPlaceholderNames).
		ContainsExactly("INTERPOLATION_1", "INTERPOLATION_2")
	parts := extractedParts.extractedParts

	AssertThat(parts[0].GetString()).
		IsEqualTo(
			"{NUM_PEOPLE, plural, offset:1\n=0 {I see {START_BOLD}no one at all{END_BOLD} in ")
	AssertThat(parts[1].GetCanonicalPlaceholderName()).IsEqualTo("INTERPOLATION_2")
	AssertThat(parts[2].GetString()).IsEqualTo(".}\n=1 {I see ")
	AssertThat(parts[3].GetCanonicalPlaceholderName()).IsEqualTo("INTERPOLATION_1")
	AssertThat(parts[4].GetString()).IsEqualTo(" in ")
	AssertThat(parts[5].GetCanonicalPlaceholderName()).IsEqualTo("INTERPOLATION_2")
	AssertThat(parts[6].GetString()).IsEqualTo(".}\n=2 {I see ")
	AssertThat(parts[7].GetCanonicalPlaceholderName()).IsEqualTo("INTERPOLATION_1")
	AssertThat(parts[8].GetString()).IsEqualTo(" and one other person in ")
	AssertThat(parts[9].GetCanonicalPlaceholderName()).IsEqualTo("INTERPOLATION_2")
	AssertThat(parts[10].GetString()).IsEqualTo(".}\nother {I see ")
	AssertThat(parts[11].GetCanonicalPlaceholderName()).IsEqualTo("INTERPOLATION_1")
	AssertThat(parts[12].GetString()).IsEqualTo(" and # other people in ")
	AssertThat(parts[13].GetCanonicalPlaceholderName()).IsEqualTo("INTERPOLATION_2")
	AssertThat(parts[14].GetString()).IsEqualTo(".}\n}\n")
}

// Java: JsMessageVisitorTest.testIcuTemplate()
func (c *JsMessageVisitorTest) TestIcuTemplate() {
	c.extractMessagesSafely(
		"const MSG_ICU_EXAMPLE = declareIcuTemplate(\n" +
			"    `{NUM_PEOPLE, plural, offset:1\n" +
			"        =0 {I see no one at all in {INTERPOLATION_2}.}\n" +
			"        =1 {I see {INTERPOLATION_1} in {INTERPOLATION_2}.}\n" +
			"        =2 {I see {INTERPOLATION_1} and one other person in {INTERPOLATION_2}.}\n" +
			"        other {I see {INTERPOLATION_1} and # other people in {INTERPOLATION_2}.}\n" +
			"    }`,\n" +
			"    {\n" +
			"      description: 'ICU example message',\n" +
			"      original_code: {\n" +
			"        'INTERPOLATION_1': '{{getPerson()}}',\n" +
			"        'INTERPOLATION_2': '{{getPlaceName()}}',\n" +
			"      },\n" +
			"      example: {\n" +
			"        'INTERPOLATION_1': 'Jane Doe',\n" +
			"        'INTERPOLATION_2': 'Paris, France',\n" +
			"      }\n" +
			"    });\n")
	AssertThat(c.icuTemplateDefinitions).HasSize(1)
	definition := c.icuTemplateDefinitions[0]

	messageNode := definition.GetMessageNode()
	AssertNode(messageNode).IsCall().HasFirstChildThat().IsName("declareIcuTemplate")

	templateTextNode := definition.GetTemplateTextNode()
	AssertNode(templateTextNode).HasToken(rhino.Token_TEMPLATELIT)

	msg := definition.GetMessage()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_ICU_EXAMPLE")
	AssertThat(msg.GetDesc()).IsEqualTo("ICU example message")
	AssertThat(msg.AsIcuMessageString()).
		IsEqualTo(
			"{NUM_PEOPLE, plural, offset:1\n" +
				"        =0 {I see no one at all in {INTERPOLATION_2}.}\n" +
				"        =1 {I see {INTERPOLATION_1} in {INTERPOLATION_2}.}\n" +
				"        =2 {I see {INTERPOLATION_1} and one other person in {INTERPOLATION_2}.}\n" +
				"        other {I see {INTERPOLATION_1} and # other people in {INTERPOLATION_2}.}\n" +
				"    }")
	// NOTE: "testcode" is the file name used by c.compiler.ParseTestCode(code)
	AssertThat(msg.GetSourceName()).IsEqualTo("testcode:1")
}

// Java: JsMessageVisitorTest.testIcuTemplatePlaceholderTypo()
func (c *JsMessageVisitorTest) TestIcuTemplatePlaceholderTypo() {
	c.extractMessages(
		"const MSG_ICU_EXAMPLE = declareIcuTemplate(\n" +
			"    `{NUM_PEOPLE, plural, offset:1\n" +
			"        =0 {I see no one at all in {INTERPOLATION_2}.}\n" +
			"        =1 {I see {INTERPOLATION_1} in {INTERPOLATION_2}.}\n" +
			"        =2 {I see {INTERPOLATION_1} and one other person in {INTERPOLATION_2}.}\n" +
			"        other {I see {INTERPOLATION_1} and # other people in {INTERPOLATION_2}.}\n" +
			"    }`,\n" +
			"    {\n" +
			"      description: 'ICU example message',\n" +
			"      example: {\n" +
			"        'INTERPOLATION_1': 'Jane Doe',\n" +
			"// This placeholder name doesn't match any placeholders in the message, so it should\n" +
			"// be reported as an error\n" +
			"        'INTERPOLATION_TYPO_2': 'Paris, France',\n" +
			"      }\n" +
			"    });\n")
	c.assertOneErrorWithTypeAndDescription(
		JsMessageVisitor_MESSAGE_TREE_MALFORMED,
		"Message parse tree malformed. Unknown placeholder: INTERPOLATION_TYPO_2")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
	// The malformed message is skipped.
	AssertThat(c.messages).IsEmpty()
}

// Java: JsMessageVisitorTest.testIcuTemplatePlaceholderUpperSnakeCase()
func (c *JsMessageVisitorTest) TestIcuTemplatePlaceholderUpperSnakeCase() {
	// For ICU templates you must always use UPPER_SNAKE_CASE for your placeholder names.
	// This matches the way the placeholder names are formatted in the XMB/XTB files.
	// Note that `goog.GetMsg()` requires lowerCamelCase for placeholder names, but
	// actually converts them to UPPER_SNAKE_CASE when storing them in XMB and must convert
	// the UPPER_SNAKE_CASE to lowerCamelCase when reading from XTB files. We want to avoid
	// this needless complication for `declareIcuTemplate()` message declarations.
	c.extractMessagesSafely(
		"const MSG_ICU_EXAMPLE = declareIcuTemplate(\n" +
			"    'Email: {USER_EMAIL}',\n" +
			"    {\n" +
			"      description: 'Labeled email address',\n" +
			"      example: {\n" +
			"        'USER_EMAIL': 'jane@doe.com',\n" +
			"      }\n" +
			"    });\n")
	jsMessage := c.assertOneMessage()
	AssertThat(jsMessage.CanonicalPlaceholderNames()).ContainsExactly("USER_EMAIL")
}

// Java: JsMessageVisitorTest.testIcuTemplatePlaceholderLowerCamelCase()
func (c *JsMessageVisitorTest) TestIcuTemplatePlaceholderLowerCamelCase() {
	// For ICU templates you must always use UPPER_SNAKE_CASE for your placeholder names.
	// This matches the way the placeholder names are formatted in the XMB/XTB files.
	// Note that `goog.GetMsg()` requires lowerCamelCase for placeholder names, but
	// actually converts them to UPPER_SNAKE_CASE when storing them in XMB and must convert
	// the UPPER_SNAKE_CASE to lowerCamelCase when reading from XTB files. We want to avoid
	// this needless complication for `declareIcuTemplate()` message declarations.
	c.extractMessages(
		"const MSG_ICU_EXAMPLE = declareIcuTemplate(\n" +
			"    'Email: {userEmail}', // should be USER_EMAIL\n" +
			"    {\n" +
			"      description: 'Labeled email address',\n" +
			"      example: {\n" +
			"        'userEmail': 'jane@doe.com', // should be USER_EMAIL\n" +
			"      }\n" +
			"    });\n")
	c.assertOneErrorWithTypeAndDescription(
		JsMessageVisitor_MESSAGE_TREE_MALFORMED,
		"Message parse tree malformed. Placeholder not in UPPER_SNAKE_CASE: userEmail")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
	// The malformed message is skipped.
	AssertThat(c.messages).IsEmpty()
}

// Java: JsMessageVisitorTest.testJsMessageOnVar()
func (c *JsMessageVisitorTest) TestJsMessageOnVar() {
	c.extractMessagesSafely("/** @desc Hello */ var MSG_HELLO = goog.getMsg('a')")

	msg := c.assertOneMessage()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_HELLO")
	AssertThat(msg.GetDesc()).IsEqualTo("Hello")
	// NOTE: "testcode" is the file name used by c.compiler.ParseTestCode(code)
	AssertThat(msg.GetSourceName()).IsEqualTo("testcode:1")
}

// Java: JsMessageVisitorTest.testJsMessageOnLet()
func (c *JsMessageVisitorTest) TestJsMessageOnLet() {
	c.compilerOptions = NewCompilerOptions()
	c.extractMessagesSafely("/** @desc Hello */ let MSG_HELLO = goog.getMsg('a')")

	msg := c.assertOneMessage()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_HELLO")
	AssertThat(msg.GetDesc()).IsEqualTo("Hello")
	// NOTE: "testcode" is the file name used by c.compiler.ParseTestCode(code)
	AssertThat(msg.GetSourceName()).IsEqualTo("testcode:1")
}

// Java: JsMessageVisitorTest.testJsMessageOnConst()
func (c *JsMessageVisitorTest) TestJsMessageOnConst() {
	c.compilerOptions = NewCompilerOptions()
	c.extractMessagesSafely("/** @desc Hello */ const MSG_HELLO = goog.getMsg('a')")

	msg := c.assertOneMessage()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_HELLO")
	AssertThat(msg.GetDesc()).IsEqualTo("Hello")
	// NOTE: "testcode" is the file name used by c.compiler.ParseTestCode(code)
	AssertThat(msg.GetSourceName()).IsEqualTo("testcode:1")
}

// Java: JsMessageVisitorTest.testJsMessagesWithSrcMap()
func (c *JsMessageVisitorTest) TestJsMessagesWithSrcMap() {
	sourceMap := sourcemap.NewSourceMapGeneratorV3()
	sourceMap.AddMapping(
		jutil.StringPtr("source1.html"),
		nil,
		sourcemap.NewFilePosition(10, 0),
		sourcemap.NewFilePosition(0, 0),
		sourcemap.NewFilePosition(0, 100))
	sourceMap.AddMapping(
		jutil.StringPtr("source2.html"),
		nil,
		sourcemap.NewFilePosition(10, 0),
		sourcemap.NewFilePosition(1, 0),
		sourcemap.NewFilePosition(1, 100))
	output := &strings.Builder{}
	sourceMap.AppendTo(output, jutil.StringPtr("unused.js"))

	c.compilerOptions = NewCompilerOptions()
	// NOTE: "testcode" is the file name used by c.compiler.ParseTestCode(code)
	inputSourceMaps := jutil.NewLinkedHashMap[string, *SourceMapInput]()
	inputSourceMaps.Put(
		"testcode",
		NewSourceMapInput(SourceFile_FromCode("example.srcmap", output.String())))
	c.compilerOptions.SetInputSourceMaps(inputSourceMaps)

	c.extractMessagesSafely(
		"/** @desc Hello */ var MSG_HELLO = goog.getMsg('a');\n" +
			"/** @desc Hi */ var MSG_HI = goog.getMsg('b');\n")
	AssertThat(c.messages).HasSize(2)

	msg1 := c.messages[0]
	AssertThat(msg1.GetKey()).IsEqualTo("MSG_HELLO")
	AssertThat(msg1.GetDesc()).IsEqualTo("Hello")
	AssertThat(msg1.GetSourceName()).IsEqualTo("source1.html:11")

	msg2 := c.messages[1]
	AssertThat(msg2.GetKey()).IsEqualTo("MSG_HI")
	AssertThat(msg2.GetDesc()).IsEqualTo("Hi")
	AssertThat(msg2.GetSourceName()).IsEqualTo("source2.html:11")
}

// Java: JsMessageVisitorTest.testJsMessageOnProperty()
func (c *JsMessageVisitorTest) TestJsMessageOnProperty() {
	c.extractMessagesSafely("/** @desc a */ pint.sub.MSG_MENU_MARK_AS_UNREAD = goog.getMsg('a')")

	msg := c.assertOneMessage()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_MENU_MARK_AS_UNREAD")
	AssertThat(msg.GetDesc()).IsEqualTo("a")
}

// Java: JsMessageVisitorTest.testJsMessageOnPublicField()
func (c *JsMessageVisitorTest) TestJsMessageOnPublicField() {
	c.extractMessagesSafely(
		"class Foo {\n" +
			"  /** @desc overflow menu */\n" +
			"  MSG_OVERFLOW_MENU = goog.getMsg('More options');\n" +
			"}\n")
	msg := c.assertOneMessage()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_OVERFLOW_MENU")
	AssertThat(msg.AsJsMessageString()).IsEqualTo("More options")

	// Anonymous class
	c.extractMessagesSafely(
		"foo(class {\n" +
			"  /** @desc hi */\n" +
			"  MSG_HELLO = goog.getMsg('Greetings');\n" +
			"});\n")
	AssertThat(c.messages).HasSize(2)
	msg = c.messages[1]
	AssertThat(msg.GetKey()).IsEqualTo("MSG_HELLO")
	AssertThat(msg.AsJsMessageString()).IsEqualTo("Greetings")
}

// Java: JsMessageVisitorTest.testJsMessageOnPublicField_directAlias()
func (c *JsMessageVisitorTest) TestJsMessageOnPublicField_directAlias() {
	c.extractMessagesSafely(
		"/** @desc Something */\n" +
			"const MSG_SOMETHING = goog.getMsg('Something');\n" +
			"class Foo {\n" +
			"  MSG_SOMETHING = MSG_SOMETHING;\n" +
			"}\n")
	AssertThat(c.assertOneMessage().GetKey()).IsEqualTo("MSG_SOMETHING")
}

// Note: This may be undesirable but represents the current behavior.
// Ideally the assignment to this MSG in the constructor should be sufficient to avoid the error.
// Java: JsMessageVisitorTest.testJsMessageOnPublicField_indirectAliasPresentlyErrors()
func (c *JsMessageVisitorTest) TestJsMessageOnPublicField_indirectAliasPresentlyErrors() {
	c.extractMessages(
		"/** @desc Anything */\n" +
			"const MSG_ANYTHING = goog.getMsg('Anything');\n" +
			"class Baz {\n" +
			"  MSG_ANYTHING;\n" +
			"  constructor() {\n" +
			"    this.MSG_ANYTHING = MSG_ANYTHING;\n" +
			"  }\n" +
			"}\n")
	c.assertOneError(JsMessageVisitor_MESSAGE_HAS_NO_VALUE)
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testErrorOnPublicFields()
func (c *JsMessageVisitorTest) TestErrorOnPublicFields() {
	c.extractMessages(
		"class Foo {\n" +
			"  /** @desc */\n" +
			"  MSG_WITH_NO_RHS;\n" +
			"}\n")
	c.assertOneError(JsMessageVisitor_MESSAGE_HAS_NO_VALUE)
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testErrorOnStaticField()
func (c *JsMessageVisitorTest) TestErrorOnStaticField() {
	c.extractMessages(
		"class Bar {\n" +
			"  /** @desc */\n" +
			"  static MSG_STATIC_FIELD_WITH_NO_RHS;\n" +
			"}\n")
	c.assertOneError(JsMessageVisitor_MESSAGE_HAS_NO_VALUE)
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testJsMessageOnPublicStaticField()
func (c *JsMessageVisitorTest) TestJsMessageOnPublicStaticField() {
	c.extractMessagesSafely(
		"class Bar {\n" +
			"  /** @desc menu */\n" +
			"  static MSG_MENU = goog.getMsg('Options');\n" +
			"}\n")
	msg := c.assertOneMessage()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_MENU")
	AssertThat(msg.AsJsMessageString()).IsEqualTo("Options")

	c.extractMessagesSafely(
		"let G = class {\n" +
			"  /** @desc apples */\n" +
			"  static MSG_FRUIT = goog.getMsg('Apples');\n" +
			"}\n")
	AssertThat(c.messages).HasSize(2)
	msg = c.messages[1]
	AssertThat(msg.GetKey()).IsEqualTo("MSG_FRUIT")
	AssertThat(msg.AsJsMessageString()).IsEqualTo("Apples")
}

// Java: JsMessageVisitorTest.testStaticInheritance()
func (c *JsMessageVisitorTest) TestStaticInheritance() {
	c.extractMessagesSafely(
		"/** @desc a */\n" +
			"foo.bar.BaseClass.MSG_MENU = goog.getMsg('hi');\n" +
			"/**\n" +
			" * @desc a\n" +
			" * @suppress {visibility}\n" +
			" */\n" +
			"foo.bar.Subclass.MSG_MENU = foo.bar.BaseClass.MSG_MENU;\n")

	msg := c.assertOneMessage()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_MENU")
	AssertThat(msg.GetDesc()).IsEqualTo("a")
}

// Java: JsMessageVisitorTest.testMsgInEnum()
func (c *JsMessageVisitorTest) TestMsgInEnum() {
	c.extractMessages(
		"/**\n" +
			" * @enum {number}\n" +
			" */\n" +
			"var MyEnum = {\n" +
			"  MSG_ONE: 0\n" +
			"};\n")
	AssertThat(c.compiler.GetErrors()).IsEmpty()
	c.assertOneWarning(JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY)
}

// Java: JsMessageVisitorTest.testMsgInEnumWithSuppression()
func (c *JsMessageVisitorTest) TestMsgInEnumWithSuppression() {
	c.extractMessagesSafely(
		"/** @fileoverview\n" +
			" * @suppress {messageConventions}\n" +
			" */\n" +
			"\n" +
			"/**\n" +
			" * @enum {number}\n" +
			" */\n" +
			"var MyEnum = {\n" +
			"  MSG_ONE: 0\n" +
			"};\n")
}

// Java: JsMessageVisitorTest.testJsMessageOnObjLit()
func (c *JsMessageVisitorTest) TestJsMessageOnObjLit() {
	c.extractMessagesSafely(
		"pint.sub = {\n" +
			"/** @desc a */ MSG_MENU_MARK_AS_UNREAD: goog.getMsg('a')}\n")

	msg := c.assertOneMessage()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_MENU_MARK_AS_UNREAD")
	AssertThat(msg.GetDesc()).IsEqualTo("a")
}

// Java: JsMessageVisitorTest.testInvalidJsMessageOnObjLit()
func (c *JsMessageVisitorTest) TestInvalidJsMessageOnObjLit() {
	c.extractMessages(
		"pint.sub = {\n" +
			"  /** @desc a */ MSG_MENU_MARK_AS_UNREAD: undefined\n" +
			"}\n")
	AssertThat(c.compiler.GetErrors()).IsEmpty()
	c.assertOneWarning(JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY)
}

// Java: JsMessageVisitorTest.testJsMessageAliasOnObjLit()
func (c *JsMessageVisitorTest) TestJsMessageAliasOnObjLit() {
	c.extractMessagesSafely(
		"pint.sub = {\n" +
			"  MSG_MENU_MARK_AS_UNREAD: another.namespace.MSG_MENU_MARK_AS_UNREAD\n" +
			"}\n")
	AssertThat(c.messages).IsEmpty()
}

// Java: JsMessageVisitorTest.testMessageAliasedToObject()
func (c *JsMessageVisitorTest) TestMessageAliasedToObject() {
	c.extractMessagesSafely("a.b.MSG_FOO = MSG_FOO;")
	AssertThat(c.messages).IsEmpty()
}

// Java: JsMessageVisitorTest.testMsgPropertyAliasesMsgVariable_mismatchedMSGNameIsAllowed()
func (c *JsMessageVisitorTest) TestMsgPropertyAliasesMsgVariable_mismatchedMSGNameIsAllowed() {
	c.extractMessagesSafely("a.b.MSG_FOO_ALIAS = MSG_FOO;")
	AssertThat(c.messages).IsEmpty()
}

// Java: JsMessageVisitorTest.testMsgPropertyAliasesMsgProperty_mismatchedMSGNameIsAllowed()
func (c *JsMessageVisitorTest) TestMsgPropertyAliasesMsgProperty_mismatchedMSGNameIsAllowed() {
	c.extractMessagesSafely("a.b.MSG_FOO_ALIAS = c.d.MSG_FOO;")
	AssertThat(c.messages).IsEmpty()
}

// Java: JsMessageVisitorTest.testMessageAliasedToObject_nonMSGNameIsNotAllowed()
func (c *JsMessageVisitorTest) TestMessageAliasedToObject_nonMSGNameIsNotAllowed() {
	c.extractMessages("a.b.MSG_FOO_ALIAS = someVarName;")
	AssertThat(c.messages).IsEmpty()
	AssertThat(c.compiler.GetErrors()).IsEmpty()
	c.assertOneWarning(JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY)
}

// Java: JsMessageVisitorTest.testMessageExport_shortHand()
func (c *JsMessageVisitorTest) TestMessageExport_shortHand() {
	c.extractMessagesSafely("exports = {MSG_FOO};")
	AssertThat(c.messages).IsEmpty()
}

// Java: JsMessageVisitorTest.testMessageExport_longHand()
func (c *JsMessageVisitorTest) TestMessageExport_longHand() {
	c.extractMessagesSafely("exports = {MSG_FOO: MSG_FOO};")
	AssertThat(c.messages).IsEmpty()
}

// Java: JsMessageVisitorTest.testMessageDefinedInExportsIsNotOrphaned()
func (c *JsMessageVisitorTest) TestMessageDefinedInExportsIsNotOrphaned() {
	c.extractMessagesSafely(
		"exports = {\n" +
			"  /** @desc Description. */\n" +
			"  MSG_FOO: goog.getMsg('Foo'),\n" +
			"};\n")
}

// Java: JsMessageVisitorTest.testJsMessageAlias_fromObjectDestrucuturing_longhand()
func (c *JsMessageVisitorTest) TestJsMessageAlias_fromObjectDestrucuturing_longhand() {
	c.extractMessagesSafely("({MSG_MENU_MARK_AS_UNREAD: MSG_MENU_MARK_AS_UNREAD} = x);")
	AssertThat(c.messages).IsEmpty()
}

// Java: JsMessageVisitorTest.testJsMessageAlias_fromObjectDestrucuturing_MSGlonghand_allowed()
func (c *JsMessageVisitorTest) TestJsMessageAlias_fromObjectDestrucuturing_MSGlonghand_allowed() {
	c.extractMessagesSafely(
		"({MSG_FOO: MSG_FOO_ALIAS} = {/** @desc Foo */ MSG_FOO: goog.getMsg('Foo')});")
	c.assertOneMessage()
}

// Java: JsMessageVisitorTest.testJsMessageAlias_fromObjectDestrucuturing_shorthand()
func (c *JsMessageVisitorTest) TestJsMessageAlias_fromObjectDestrucuturing_shorthand() {
	c.extractMessagesSafely("({MSG_MENU_MARK_AS_UNREAD} = x);")
	AssertThat(c.messages).IsEmpty()
}

// Java: JsMessageVisitorTest.testJsMessageOnRHSOfVar()
func (c *JsMessageVisitorTest) TestJsMessageOnRHSOfVar() {
	c.extractMessagesSafely("var MSG_MENU_MARK_AS_UNREAD = a.name.space.MSG_MENU_MARK_AS_UNREAD;")
	AssertThat(c.messages).IsEmpty()
}

// Java: JsMessageVisitorTest.testOrphanedJsMessage()
func (c *JsMessageVisitorTest) TestOrphanedJsMessage() {
	c.extractMessages("goog.getMsg('a')")
	AssertThat(c.messages).IsEmpty()

	c.assertOneError(JsMessageVisitor_MESSAGE_NODE_IS_ORPHANED)
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testOrphanedIcuTemplate()
func (c *JsMessageVisitorTest) TestOrphanedIcuTemplate() {
	c.extractMessages(
		"const {declareIcuTemplate} = goog.require('goog.i18n.messages');\n" +
			"\n" +
			"declareIcuTemplate('a')\n")
	AssertThat(c.messages).IsEmpty()

	c.assertOneError(JsMessageVisitor_MESSAGE_NODE_IS_ORPHANED)
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testMessageWithoutDescription()
func (c *JsMessageVisitorTest) TestMessageWithoutDescription() {
	c.extractMessages("var MSG_HELLO = goog.getMsg('a')")

	msg := c.assertOneMessage()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_HELLO")

	AssertThat(c.compiler.GetErrors()).IsEmpty()
	c.assertOneWarning(JsMessageVisitor_MESSAGE_HAS_NO_DESCRIPTION)
}

// Java: JsMessageVisitorTest.testIncorrectMessageReporting()
func (c *JsMessageVisitorTest) TestIncorrectMessageReporting() {
	c.extractMessages("var MSG_HELLO = goog.getMsg('a' + + 'b')")
	c.assertOneErrorWithTypeAndDescription(
		JsMessageVisitor_MESSAGE_TREE_MALFORMED,
		"Message parse tree malformed. literal string or concatenation expected")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
	AssertThat(c.messages).IsEmpty()
}

// Java: JsMessageVisitorTest.testTemplateLiteral()
func (c *JsMessageVisitorTest) TestTemplateLiteral() {
	c.compilerOptions = NewCompilerOptions()

	c.extractMessagesSafely("/** @desc Hello */ var MSG_HELLO = goog.getMsg(`hello`);")

	msg := c.assertOneMessage()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_HELLO")
	AssertThat(msg.AsJsMessageString()).IsEqualTo("hello")
}

// Java: JsMessageVisitorTest.testTemplateLiteralWithSubstitution()
func (c *JsMessageVisitorTest) TestTemplateLiteralWithSubstitution() {
	c.compilerOptions = NewCompilerOptions()

	c.extractMessages("/** @desc Hello */ var MSG_HELLO = goog.getMsg(`hello ${name}`);")
	c.assertOneErrorWithTypeAndDescription(
		JsMessageVisitor_MESSAGE_TREE_MALFORMED,
		"Message parse tree malformed. Template literals with substitutions are not allowed.")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
	AssertThat(c.messages).IsEmpty()
}

// Java: JsMessageVisitorTest.testClosureMessageWithHelpPostfix()
func (c *JsMessageVisitorTest) TestClosureMessageWithHelpPostfix() {
	c.extractMessagesSafely(
		"/** @desc help text */\n" +
			"var MSG_FOO_HELP = goog.getMsg('Help!');\n")

	msg := c.assertOneMessage()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_FOO_HELP")
	AssertThat(msg.GetDesc()).IsEqualTo("help text")
	AssertThat(msg.AsJsMessageString()).IsEqualTo("Help!")
}

// Java: JsMessageVisitorTest.testClosureMessageWithoutGoogGetmsg()
func (c *JsMessageVisitorTest) TestClosureMessageWithoutGoogGetmsg() {
	c.extractMessages("var MSG_FOO_HELP = 'I am a bad message';")

	AssertThat(c.messages).IsEmpty()
	AssertThat(c.compiler.GetErrors()).IsEmpty()
	c.assertOneWarning(JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY)
}

// Java: JsMessageVisitorTest.testAllowOneMSGtoAliasAnotherMSG()
func (c *JsMessageVisitorTest) TestAllowOneMSGtoAliasAnotherMSG() {
	// NOTE: tsickle code generation can end up creating new MSG_* variables that are temporary
	// aliases of existing ones that were defined correctly using goog.GetMsg(). Don't complain
	// about them.
	c.extractMessagesSafely(
		"/** @desc A foo message */\n" +
			"var MSG_FOO = goog.getMsg('Foo message');\n" +
			"var MSG_FOO_1 = MSG_FOO;\n")

	msg := c.assertOneMessage()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_FOO")
	AssertThat(msg.AsJsMessageString()).IsEqualTo("Foo message")
}

// Java: JsMessageVisitorTest.testDisallowOneMSGtoAliasNONMSG()
func (c *JsMessageVisitorTest) TestDisallowOneMSGtoAliasNONMSG() {
	// NOTE: tsickle code generation can end up creating new MSG_* variables that are temporary
	// aliases of existing ones that were defined correctly using goog.GetMsg(). Don't complain
	// about them.
	c.extractMessages(
		"/** @desc A foo message */\n" +
			"var mymsg = 'Foo message';\n" +
			"var MSG_FOO_1 = mymsg;\n")

	AssertThat(c.messages).IsEmpty()
	AssertThat(c.compiler.GetErrors()).IsEmpty()
	c.assertOneWarning(JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY)
}

// Java: JsMessageVisitorTest.testClosureFormatParametizedFunction()
func (c *JsMessageVisitorTest) TestClosureFormatParametizedFunction() {
	c.extractMessagesSafely(
		"/** @desc help text */\n" +
			"var MSG_SILLY = goog.getMsg('{$adjective} ' + 'message',\n" +
			"{'adjective': 'silly'});\n")

	AssertThat(c.messageDefinitions).HasSize(1)
	jsMessageDefinition := c.messageDefinitions[0]

	// `goog.GetMsg(...)`
	messageNode := jsMessageDefinition.GetMessageNode()
	AssertNode(messageNode).IsCall().HasFirstChildThat().MatchesQualifiedName("goog.getMsg")

	// `'{$adjective} ' + 'message'`
	templateTextNode := jsMessageDefinition.GetTemplateTextNode()
	AssertNode(templateTextNode).HasToken(rhino.Token_ADD).HasFirstChildThat().IsString("{$adjective} ")

	// `{'adjective': ...}`
	placeholderValuesNode := jsMessageDefinition.GetPlaceholderValuesNode()
	AssertNode(placeholderValuesNode).
		IsObjectLit().
		HasFirstChildThat().
		HasToken(rhino.Token_STRING_KEY).
		HasStringThat().
		IsEqualTo("adjective")

	// Map containing "adjective" -> string node containing "silly"
	placeholderValueMap :=
		jsMessageDefinition.GetPlaceholderValueMap()
	AssertThat(placeholderValueMap).HasSize(1)
	AssertNode(placeholderValueMap.Get("adjective")).IsString("silly")

	msg := jsMessageDefinition.GetMessage()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_SILLY")
	AssertThat(msg.GetDesc()).IsEqualTo("help text")
	AssertThat(msg.AsJsMessageString()).IsEqualTo("{$adjective} message")
}

// Java: JsMessageVisitorTest.testHugeMessage()
func (c *JsMessageVisitorTest) TestHugeMessage() {
	c.extractMessagesSafely(
		"/**\n" +
			" * @desc A message with lots of stuff.\n" +
			" */\n" +
			"var MSG_HUGE = goog.getMsg(\n" +
			"    '{$startLink_1}Google{$endLink}' +\n" +
			"    '{$startLink_2}blah{$endLink}{$boo}{$foo_001}{$boo}' +\n" +
			"    '{$foo_002}{$xxx_001}{$image}{$image_001}{$xxx_002}',\n" +
			"    {'startLink_1': '<a href=http://www.google.com/>',\n" +
			"     'endLink': '</a>',\n" +
			"     'startLink_2': '<a href=\"' + opt_data.url + '\">',\n" +
			"     'boo': opt_data.boo,\n" +
			"     'foo_001': opt_data.foo,\n" +
			"     'foo_002': opt_data.boo.foo,\n" +
			"     'xxx_001': opt_data.boo + opt_data.foo,\n" +
			"     'image': htmlTag7,\n" +
			"     'image_001': opt_data.image,\n" +
			"     'xxx_002': foo.callWithOnlyTopLevelKeys(\n" +
			"         bogusFn, opt_data, null, 'bogusKey1',\n" +
			"         opt_data.moo, 'bogusKey2', param10)});\n")

	msg := c.assertOneMessage()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_HUGE")
	AssertThat(msg.GetDesc()).IsEqualTo("A message with lots of stuff.")
	AssertThat(msg.AsJsMessageString()).
		IsEqualTo(
			"{$startLink_1}Google{$endLink}{$startLink_2}blah{$endLink}{$boo}{$foo_001}{$boo}{$foo_002}{$xxx_001}{$image}{$image_001}{$xxx_002}")
}

// Java: JsMessageVisitorTest.testUnnamedGoogleMessage()
func (c *JsMessageVisitorTest) TestUnnamedGoogleMessage() {
	c.extractMessages("var MSG_UNNAMED = goog.getMsg('Hullo');")

	msg := c.assertOneMessage()
	AssertThat(msg.GetDesc()).IsNull()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_16LJMYKCXT84X")
	AssertThat(msg.GetId()).IsEqualTo("MSG_16LJMYKCXT84X")
	AssertThat(c.compiler.GetErrors()).IsEmpty()
	c.assertOneWarning(JsMessageVisitor_MESSAGE_HAS_NO_DESCRIPTION)
}

// Java: JsMessageVisitorTest.testUnnamedIcuTemplate()
func (c *JsMessageVisitorTest) TestUnnamedIcuTemplate() {
	// Unlike `goog.GetMsg()` we do require a description for anonymous ICU templates.
	// Requiring the description is less complicated, and there doesn't seem to be any reason
	// not to require them.
	c.extractMessagesSafely(
		"var MSG_UNNAMED = declareIcuTemplate('Hullo', {description: 'description'});")

	AssertThat(c.icuTemplateDefinitions).HasSize(1)
	msg := c.icuTemplateDefinitions[0].GetMessage()
	AssertThat(msg.GetDesc()).IsEqualTo("description")
	AssertThat(msg.GetKey()).IsEqualTo("MSG_16LJMYKCXT84X")
	AssertThat(msg.GetId()).IsEqualTo("MSG_16LJMYKCXT84X")
}

// Java: JsMessageVisitorTest.testMeaningGetsUsedAsIdIfTheresNoGenerator()
func (c *JsMessageVisitorTest) TestMeaningGetsUsedAsIdIfTheresNoGenerator() {
	c.extractMessagesSafely(
		"/**\n" +
			" * @desc some description\n" +
			" * @meaning some meaning\n" +
			" */\n" +
			"var MSG_HULLO = goog.getMsg('Hullo');\n")

	msg := c.assertOneMessage()
	AssertThat(msg.GetDesc()).IsEqualTo("some description")
	AssertThat(msg.GetKey()).IsEqualTo("MSG_HULLO")
	AssertThat(msg.GetMeaning()).IsEqualTo("some meaning")
	AssertThat(msg.GetId()).IsEqualTo("some meaning")
}

// Java: JsMessageVisitorTest.testEmptyTextMessage()
func (c *JsMessageVisitorTest) TestEmptyTextMessage() {
	c.extractMessages("/** @desc text */ var MSG_FOO = goog.getMsg('');")

	c.assertOneMessage()
	AssertThat(c.compiler.GetErrors()).IsEmpty()
	c.assertOneWarningWithDescription(
		"Message value of MSG_FOO is just an empty string. Empty messages are forbidden.")
}

// Java: JsMessageVisitorTest.testEmptyTextComplexMessage()
func (c *JsMessageVisitorTest) TestEmptyTextComplexMessage() {
	c.extractMessages(
		"/** @desc text */ var MSG_BAR = goog.getMsg(\n" +
			"'' + '' + ''     + '' +'');\n")

	c.assertOneMessage()
	AssertThat(c.compiler.GetErrors()).IsEmpty()
	c.assertOneWarningWithDescription(
		"Message value of MSG_BAR is just an empty string. Empty messages are forbidden.")
}

// Java: JsMessageVisitorTest.testMsgVarWithoutAssignment()
func (c *JsMessageVisitorTest) TestMsgVarWithoutAssignment() {
	c.extractMessages("var MSG_SILLY;")

	c.assertOneError(JsMessageVisitor_MESSAGE_HAS_NO_VALUE)
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testRegularVarWithoutAssignment()
func (c *JsMessageVisitorTest) TestRegularVarWithoutAssignment() {
	c.extractMessagesSafely("var SILLY;")

	AssertThat(c.messages).IsEmpty()
}

// Java: JsMessageVisitorTest.testMsgPropertyWithoutAssignment()
func (c *JsMessageVisitorTest) TestMsgPropertyWithoutAssignment() {
	c.t.Skip("@Ignore // Currently unimplemented.")
	c.extractMessages("goog.message.MSG_SILLY_PROP;")

	c.assertOneErrorWithDescription("Message MSG_SILLY_PROP has no value")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testMsgVarWithIncorrectRightSide()
func (c *JsMessageVisitorTest) TestMsgVarWithIncorrectRightSide() {
	c.extractMessages("var MSG_SILLY = 0;")

	AssertThat(c.compiler.GetErrors()).IsEmpty()
	c.assertOneWarningWithTypeAndDescription(
		JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY,
		"Message must be initialized using a call to goog.getMsg or"+
			" goog.i18n.messages.declareIcuTemplate")
}

// Java: JsMessageVisitorTest.testIncorrectMessage()
func (c *JsMessageVisitorTest) TestIncorrectMessage() {
	c.extractMessages("DP_DatePicker.MSG_DATE_SELECTION = {};")

	AssertThat(c.messages).IsEmpty()
	AssertThat(c.compiler.GetErrors()).IsEmpty()
	c.assertOneWarningWithTypeAndDescription(
		JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY,
		"Message must be initialized using a call to goog.getMsg or"+
			" goog.i18n.messages.declareIcuTemplate")
}

// Java: JsMessageVisitorTest.testUnrecognizedFunction()
func (c *JsMessageVisitorTest) TestUnrecognizedFunction() {
	c.extractMessages("DP_DatePicker.MSG_DATE_SELECTION = somefunc('a')")

	AssertThat(c.messages).IsEmpty()
	c.assertOneErrorWithDescription(
		"Message parse tree malformed. Message must be initialized using a call to goog.getMsg" +
			" or declareIcuTemplate (from goog.i18n.messages).")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testExtractPropertyMessage()
func (c *JsMessageVisitorTest) TestExtractPropertyMessage() {
	c.extractMessagesSafely(
		"/**\n" +
			" * @desc A message that demonstrates placeholders\n" +
			" */\n" +
			"a.b.MSG_SILLY = goog.getMsg(\n" +
			"    '{$adjective} ' + '{$someNoun}',\n" +
			"    {'adjective': adj, 'someNoun': noun});\n")

	msg := c.assertOneMessage()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_SILLY")
	AssertThat(msg.AsJsMessageString()).IsEqualTo("{$adjective} {$someNoun}")
	AssertThat(msg.GetDesc()).IsEqualTo("A message that demonstrates placeholders")
}

// Java: JsMessageVisitorTest.testExtractPropertyMessageInFunction()
func (c *JsMessageVisitorTest) TestExtractPropertyMessageInFunction() {
	c.extractMessagesSafely(
		"function f() {\n" +
			"  /**\n" +
			"   * @desc A message that demonstrates placeholders\n" +
			"   */\n" +
			"  a.b.MSG_SILLY = goog.getMsg(\n" +
			"      '{$adjective} ' + '{$someNoun}',\n" +
			"      {'adjective': adj, 'someNoun': noun});\n" +
			"}\n")

	msg := c.assertOneMessage()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_SILLY")
	AssertThat(msg.AsJsMessageString()).IsEqualTo("{$adjective} {$someNoun}")
	AssertThat(msg.GetDesc()).IsEqualTo("A message that demonstrates placeholders")
}

// Java: JsMessageVisitorTest.testAlmostButNotExternalMessage()
func (c *JsMessageVisitorTest) TestAlmostButNotExternalMessage() {
	c.extractMessagesSafely("/** @desc External */ var MSG_EXTERNAL = goog.getMsg('External');")

	msg := c.assertOneMessage()
	AssertThat(msg.IsExternal()).IsFalse()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_EXTERNAL")
}

// Java: JsMessageVisitorTest.testExternalMessage()
func (c *JsMessageVisitorTest) TestExternalMessage() {
	c.extractMessagesSafely("var MSG_EXTERNAL_111 = goog.getMsg('Hello World');")

	msg := c.assertOneMessage()
	AssertThat(msg.IsExternal()).IsTrue()
	AssertThat(msg.GetId()).IsEqualTo("111")
}

// Java: JsMessageVisitorTest.testExternalMessage_customSuffix()
func (c *JsMessageVisitorTest) TestExternalMessage_customSuffix() {
	c.extractMessagesSafely("var MSG_EXTERNAL_111$$1 = goog.getMsg('Hello World');")

	msg := c.assertOneMessage()
	AssertThat(msg.IsExternal()).IsTrue()
	AssertThat(msg.GetId()).IsEqualTo("111")
}

// Java: JsMessageVisitorTest.testIsValidMessageNameStrict()
func (c *JsMessageVisitorTest) TestIsValidMessageNameStrict() {
	var visitor JsMessageVisitor = newJsMessageVisitorTest_DummyJsVisitor()

	AssertThat(visitor.IsMessageName("MSG_HELLO")).IsTrue()
	AssertThat(visitor.IsMessageName("MSG_")).IsTrue()
	AssertThat(visitor.IsMessageName("MSG_HELP")).IsTrue()
	AssertThat(visitor.IsMessageName("MSG_FOO_HELP")).IsTrue()

	AssertThat(visitor.IsMessageName("_FOO_HELP")).IsFalse()
	AssertThat(visitor.IsMessageName("MSGFOOP")).IsFalse()
}

// Java: JsMessageVisitorTest.testUnexistedPlaceholders()
func (c *JsMessageVisitorTest) TestUnexistedPlaceholders() {
	c.extractMessages("var MSG_FOO = goog.getMsg('{$foo}:', {});")

	AssertThat(c.messages).IsEmpty()
	c.assertOneErrorWithTypeAndDescription(
		JsMessageVisitor_MESSAGE_TREE_MALFORMED,
		"Message parse tree malformed. Unrecognized message placeholder referenced: foo")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testUnusedReferenesAreNotOK()
func (c *JsMessageVisitorTest) TestUnusedReferenesAreNotOK() {
	c.extractMessages(
		"/** @desc AA */\n" +
			"var MSG_FOO = goog.getMsg('lalala:', {foo:1});\n")
	AssertThat(c.messages).IsEmpty()
	c.assertOneErrorWithTypeAndDescription(
		JsMessageVisitor_MESSAGE_TREE_MALFORMED, "Message parse tree malformed. Unused message placeholder: foo")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testDuplicatePlaceHoldersAreBad()
func (c *JsMessageVisitorTest) TestDuplicatePlaceHoldersAreBad() {
	c.extractMessages(
		"var MSG_FOO = goog.getMsg(\n" +
			"'{$foo}:', {'foo': 1, 'foo' : 2});\n")

	AssertThat(c.messages).IsEmpty()
	c.assertOneErrorWithTypeAndDescription(
		JsMessageVisitor_MESSAGE_TREE_MALFORMED, "Message parse tree malformed. duplicate string key: foo")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testDuplicatePlaceholderReferencesAreOk()
func (c *JsMessageVisitorTest) TestDuplicatePlaceholderReferencesAreOk() {
	c.extractMessagesSafely(
		"/** @desc Sample description */\n" +
			"var MSG_FOO = goog.getMsg('{$foo}:, {$foo}', {'foo': 1});\n")

	msg := c.assertOneMessage()
	AssertThat(msg.AsJsMessageString()).IsEqualTo("{$foo}:, {$foo}")
}

// Java: JsMessageVisitorTest.testCamelcasePlaceholderNamesAreOk()
func (c *JsMessageVisitorTest) TestCamelcasePlaceholderNamesAreOk() {
	c.extractMessagesSafely(
		"/** @desc Hello */\n" +
			"var MSG_WITH_CAMELCASE = goog.getMsg('Slide {$slideNumber}:', {'slideNumber': opt_index + 1});\n")

	msg := c.assertOneMessage()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_WITH_CAMELCASE")
	AssertThat(msg.AsJsMessageString()).IsEqualTo("Slide {$slideNumber}:")
	parts := msg.GetParts()
	AssertThat(parts).HasSize(3)
	AssertThat(parts[1].GetJsPlaceholderName()).IsEqualTo("slideNumber")
	AssertThat(parts[1].GetCanonicalPlaceholderName()).IsEqualTo("SLIDE_NUMBER")
}

// Java: JsMessageVisitorTest.testNonCamelcasePlaceholderNamesAreNotOkInMsgText()
func (c *JsMessageVisitorTest) TestNonCamelcasePlaceholderNamesAreNotOkInMsgText() {
	c.extractMessages(
		"var MSG_WITH_CAMELCASE = goog.getMsg(\n" +
			"    'Slide {$SLIDE_NUMBER}:',\n" +
			"    {'slideNumber': opt_index + 1});\n")

	AssertThat(c.messages).IsEmpty()
	c.assertOneErrorWithTypeAndDescription(
		JsMessageVisitor_MESSAGE_TREE_MALFORMED,
		"Message parse tree malformed. Placeholder name not in lowerCamelCase: SLIDE_NUMBER")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testNonCamelcasePlaceholderNamesAreNotOkInPlaceholderObject()
func (c *JsMessageVisitorTest) TestNonCamelcasePlaceholderNamesAreNotOkInPlaceholderObject() {
	c.extractMessages(
		"var MSG_WITH_CAMELCASE = goog.getMsg(\n" +
			"    'Slide {$slideNumber}:',\n" +
			"    {'SLIDE_NUMBER': opt_index + 1});\n")

	AssertThat(c.messages).IsEmpty()
	c.assertOneErrorWithTypeAndDescription(
		JsMessageVisitor_MESSAGE_TREE_MALFORMED,
		"Message parse tree malformed. Unrecognized message placeholder referenced: slideNumber")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testUnquotedPlaceholdersAreOk()
func (c *JsMessageVisitorTest) TestUnquotedPlaceholdersAreOk() {
	c.extractMessagesSafely(
		"/** @desc Hello */\n" +
			"var MSG_FOO = goog.getMsg('foo {$unquoted}:', {unquoted: 12});\n")

	c.assertOneMessage()
}

// Java: JsMessageVisitorTest.testDuplicateMessageError()
func (c *JsMessageVisitorTest) TestDuplicateMessageError() {
	c.extractMessages(
		"(function () {/** @desc Hello */ var MSG_HELLO = goog.getMsg('a')})\n" +
			"(function () {/** @desc Hello2 */ var MSG_HELLO = goog.getMsg('a')})\n")

	c.assertOneError(JsMessageVisitor_MESSAGE_DUPLICATE_KEY)
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testNoDuplicateErrorOnExternMessage()
func (c *JsMessageVisitorTest) TestNoDuplicateErrorOnExternMessage() {
	c.extractMessagesSafely(
		"(function () {/** @desc Hello */\n" +
			"var MSG_EXTERNAL_2 = goog.getMsg('a')})\n" +
			"(function () {/** @desc Hello2 */\n" +
			"var MSG_EXTERNAL_2 = goog.getMsg('a')})\n")
}

// Java: JsMessageVisitorTest.testUsingMsgPrefixWithFallback()
func (c *JsMessageVisitorTest) TestUsingMsgPrefixWithFallback() {
	c.extractMessagesSafely(
		"function f() {\n" +
			"/** @desc Hello */ var MSG_UNNAMED_1 = goog.getMsg('hello');\n" +
			"/** @desc Hello */ var MSG_UNNAMED_2 = goog.getMsg('hello');\n" +
			"var x = goog.getMsgWithFallback(\n" +
			"    MSG_UNNAMED_1, MSG_UNNAMED_2);\n" +
			"}\n")
}

// Java: JsMessageVisitorTest.testUsingMsgPrefixWithFallback_rename()
func (c *JsMessageVisitorTest) TestUsingMsgPrefixWithFallback_rename() {
	c.renameMessages = true
	c.extractMessagesSafely(
		"function f() {\n" +
			"/** @desc Hello */ var MSG_A = goog.getMsg('hello');\n" +
			"/** @desc Hello */ var MSG_B = goog.getMsg('hello!');\n" +
			"var x = goog.getMsgWithFallback(MSG_A, MSG_B);\n" +
			"}\n")
}

// Java: JsMessageVisitorTest.testUsingMsgPrefixWithFallback_duplicateUnnamedKeys_rename()
func (c *JsMessageVisitorTest) TestUsingMsgPrefixWithFallback_duplicateUnnamedKeys_rename() {
	c.renameMessages = true
	c.extractMessagesSafely(
		"function f() {\n" +
			"  /** @desc Hello */ var MSG_UNNAMED_1 = goog.getMsg('hello');\n" +
			"  /** @desc Hello */ var MSG_UNNAMED_2 = goog.getMsg('hello');\n" +
			"  var x = goog.getMsgWithFallback(\n" +
			"      MSG_UNNAMED_1, MSG_UNNAMED_2);\n" +
			"}\n" +
			"function g() {\n" +
			"  /** @desc Hello */ var MSG_UNNAMED_1 = goog.getMsg('hello');\n" +
			"  /** @desc Hello */ var MSG_UNNAMED_2 = goog.getMsg('hello');\n" +
			"  var x = goog.getMsgWithFallback(\n" +
			"      MSG_UNNAMED_1, MSG_UNNAMED_2);\n" +
			"}\n")
}

// Java: JsMessageVisitorTest.testUsingMsgPrefixWithFallback_module()
func (c *JsMessageVisitorTest) TestUsingMsgPrefixWithFallback_module() {
	c.extractMessagesSafely(
		"/** @desc Hello */ var MSG_A = goog.getMsg('hello');\n" +
			"/** @desc Hello */ var MSG_B = goog.getMsg('hello!');\n" +
			"var x = goog.getMsgWithFallback(messages.MSG_A, MSG_B);\n")
}

// Java: JsMessageVisitorTest.testUsingMsgPrefixWithFallback_moduleRenamed()
func (c *JsMessageVisitorTest) TestUsingMsgPrefixWithFallback_moduleRenamed() {
	c.extractMessagesSafely(
		"/** @desc Hello */ var MSG_A = goog.getMsg('hello');\n" +
			"/** @desc Hello */ var MSG_B = goog.getMsg('hello!');\n" +
			"var x = goog.getMsgWithFallback(module$exports$messages$MSG_A, MSG_B);\n")
}

// Java: JsMessageVisitorTest.testErrorWhenUsingMsgPrefixWithFallback()
func (c *JsMessageVisitorTest) TestErrorWhenUsingMsgPrefixWithFallback() {
	c.extractMessages(
		"/** @desc Hello */ var MSG_HELLO_1 = goog.getMsg('hello');\n" +
			"/** @desc Hello */ var MSG_HELLO_2 = goog.getMsg('hello');\n" +
			"/** @desc Hello */\n" +
			"var MSG_HELLO_3 = goog.getMsgWithFallback(MSG_HELLO_1, MSG_HELLO_2);\n")
	c.assertOneError(JsMessageVisitor_MESSAGE_TREE_MALFORMED)
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testRenamedMessages_var()
func (c *JsMessageVisitorTest) TestRenamedMessages_var() {
	c.renameMessages = true

	c.extractMessagesSafely("/** @desc Hello */ var MSG_HELLO = goog.getMsg('a')")

	msg := c.assertOneMessage()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_HELLO")
	AssertThat(msg.GetDesc()).IsEqualTo("Hello")
	// NOTE: "testcode" is the file name used by c.compiler.ParseTestCode(code)
	AssertThat(msg.GetSourceName()).IsEqualTo("testcode:1")
}

// Java: JsMessageVisitorTest.testRenamedMessages_getprop()
func (c *JsMessageVisitorTest) TestRenamedMessages_getprop() {
	c.renameMessages = true

	c.extractMessagesSafely("/** @desc a */ pint.sub.MSG_MENU_MARK_AS_UNREAD = goog.getMsg('a')")

	msg := c.assertOneMessage()
	AssertThat(msg.GetKey()).IsEqualTo("MSG_MENU_MARK_AS_UNREAD")
	AssertThat(msg.GetDesc()).IsEqualTo("a")
}

// Java: JsMessageVisitorTest.testGetMsgWithOptions()
func (c *JsMessageVisitorTest) TestGetMsgWithOptions() {
	c.extractMessagesSafely(
		"/** @desc Hello */\n" +
			"var MSG_HELLO =\n" +
			"    goog.getMsg(\n" +
			"        'Hello, {$name}',\n" +
			"        {\n" +
			"          'name': getName(),\n" +
			"        },\n" +
			"        {\n" +
			"          html: false,\n" +
			"          unescapeHtmlEntities: true,\n" +
			"          example: {\n" +
			"            'name': 'George',\n" +
			"          },\n" +
			"          original_code: {\n" +
			"            'name': 'getName()',\n" +
			"          },\n" +
			"        })\n")

	msg := c.assertOneMessage()
	// MSG_HELLO =
	AssertThat(msg.GetKey()).IsEqualTo("MSG_HELLO")
	// * @desc Hello
	AssertThat(msg.GetDesc()).IsEqualTo("Hello")

	msgDefinition := c.messageDefinitions[0]

	// goog.GetMsg(...)
	callNode := msgDefinition.GetMessageNode()
	AssertNode(callNode).IsCall().HasFirstChildThat().MatchesQualifiedName("goog.getMsg")

	AssertNode(msgDefinition.GetTemplateTextNode()).IsString("Hello, {$name}")

	// `{ 'name': getName() }`
	placeholderValuesNode := msgDefinition.GetPlaceholderValuesNode()
	AssertNode(placeholderValuesNode).IsObjectLit()
	AssertThat(callNode.GetChildAtIndex(2)).IsSameInstanceAs(placeholderValuesNode)

	// placeholder name 'name' maps to the `getName()` call in the values map
	placeholderValueMap := msgDefinition.GetPlaceholderValueMap()
	AssertThat(placeholderValueMap.KeySlice()).ContainsExactly("name")
	nameValueNode := placeholderValueMap.Get("name")
	AssertNode(nameValueNode).IsCall().HasOneChildThat().IsName("getName")
	AssertThat(nameValueNode.GetGrandparent()).IsSameInstanceAs(placeholderValuesNode)

	// `html: false`
	AssertThat(msgDefinition.ShouldEscapeLessThan()).IsFalse()
	// `unescapeHtmlEntities: true`
	AssertThat(msgDefinition.ShouldUnescapeHtmlEntities()).IsTrue()

	// MapSubject.containsExactly(k, v) is containsExactlyEntriesIn of the one-entry map.
	// `example: { 'name': 'George' }`
	AssertThat(msg.GetPlaceholderNameToExampleMap()).ContainsExactlyEntriesIn(map[string]string{"name": "George"})
	// `original_code: {'name': 'getName()' }`
	AssertThat(msg.GetPlaceholderNameToOriginalCodeMap()).
		ContainsExactlyEntriesIn(map[string]string{"name": "getName()"})
}

// Java: JsMessageVisitorTest.testGoogGetMsgWithNoArgs()
func (c *JsMessageVisitorTest) TestGoogGetMsgWithNoArgs() {
	c.extractMessages(
		"/** @desc something */\n" +
			"const MSG_X = goog.getMsg(); // no arguments\n")
	c.assertOneErrorWithDescription("Message parse tree malformed. Message string literal expected")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testGoogGetMsgWithBadValuesArg()
func (c *JsMessageVisitorTest) TestGoogGetMsgWithBadValuesArg() {
	c.extractMessages(
		"/** @desc something */\n" +
			"const MSG_X =\n" +
			"    goog.getMsg(\n" +
			"        'Hello, {$name}',\n" +
			"        getName()); // should be an object literal\n")
	c.assertOneErrorWithDescription("Message parse tree malformed. object literal expected")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testGoogGetMsgWithBadValuesKey()
func (c *JsMessageVisitorTest) TestGoogGetMsgWithBadValuesKey() {
	c.extractMessages(
		"/** @desc something */\n" +
			"const MSG_X =\n" +
			"    goog.getMsg(\n" +
			"        'Hello, {$name}',\n" +
			"        {\n" +
			"          [name]: getName() // a computed key is not allowed\n" +
			"        });\n")
	c.assertOneErrorWithDescription("Message parse tree malformed. string key expected")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testGoogGetMsgWithBadOptionsArg()
func (c *JsMessageVisitorTest) TestGoogGetMsgWithBadOptionsArg() {
	c.extractMessages(
		"/** @desc something */\n" +
			"const MSG_X =\n" +
			"    goog.getMsg(\n" +
			"        'Hello, {$name}',\n" +
			"        { 'name': getName() },\n" +
			"        options); // options bag must be an object literal\n")
	c.assertOneErrorWithDescription("Message parse tree malformed. object literal expected")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testGoogGetMsgWithComputedKeyInOptions()
func (c *JsMessageVisitorTest) TestGoogGetMsgWithComputedKeyInOptions() {
	c.extractMessages(
		"/** @desc something */\n" +
			"const MSG_X =\n" +
			"    goog.getMsg(\n" +
			"        'Hello, {$name}',\n" +
			"        { 'name': getName() },\n" +
			"        {\n" +
			"          [options]: true // option names cannot be computed keys\n" +
			"        });\n")
	c.assertOneErrorWithDescription("Message parse tree malformed. string key expected")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testGoogGetMsgWithUnknownOption()
func (c *JsMessageVisitorTest) TestGoogGetMsgWithUnknownOption() {
	c.extractMessages(
		"/** @desc something */\n" +
			"const MSG_X =\n" +
			"    goog.getMsg(\n" +
			"        'Hello, {$name}',\n" +
			"        { 'name': getName() },\n" +
			"        {\n" +
			"          unknownOption: true // not a valid option name\n" +
			"        });\n")
	c.assertOneErrorWithDescription("Message parse tree malformed. Unknown option: unknownOption")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testGoogGetMsgWithInvalidBooleanOption()
func (c *JsMessageVisitorTest) TestGoogGetMsgWithInvalidBooleanOption() {
	c.extractMessages(
		"/** @desc something */\n" +
			"const MSG_X =\n" +
			"    goog.getMsg(\n" +
			"        'Hello, {$name}',\n" +
			"        { 'name': getName() },\n" +
			"        {\n" +
			"          html: 'true' // boolean option value must be a boolean literal\n" +
			"        });\n")
	c.assertOneErrorWithDescription("Message parse tree malformed. html: Literal true or false expected")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testGoogGetMsgWithOriginalCodeForInvalidExample()
func (c *JsMessageVisitorTest) TestGoogGetMsgWithOriginalCodeForInvalidExample() {
	c.extractMessages(
		"/** @desc something */\n" +
			"const MSG_X =\n" +
			"    goog.getMsg(\n" +
			"        'Hello, {$name}',\n" +
			"        { 'name': getName() },\n" +
			"        {\n" +
			"          example: 'name: something' // not an object literal\n" +
			"        });\n")
	c.assertOneErrorWithDescription("Message parse tree malformed. object literal expected")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testGoogGetMsgWithInvalidOriginalCodeValue()
func (c *JsMessageVisitorTest) TestGoogGetMsgWithInvalidOriginalCodeValue() {
	c.extractMessages(
		"/** @desc something */\n" +
			"const MSG_X =\n" +
			"    goog.getMsg(\n" +
			"        'Hello, {$name}',\n" +
			"        { 'name': getName() },\n" +
			"        {\n" +
			"          original_code: {\n" +
			"            'name': getName() // value is not a string\n" +
			"          }\n" +
			"        });\n")
	c.assertOneErrorWithDescription("Message parse tree malformed. literal string or concatenation expected")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testGoogGetMsgWithOriginalCodeForUnknownPlaceholder()
func (c *JsMessageVisitorTest) TestGoogGetMsgWithOriginalCodeForUnknownPlaceholder() {
	c.extractMessages(
		"/** @desc something */\n" +
			"const MSG_X =\n" +
			"    goog.getMsg(\n" +
			"        'Hello, {$name}',\n" +
			"        { 'name': getName() },\n" +
			"        {\n" +
			"          original_code: {\n" +
			"            'unknownPlaceholder': 'something' // not a valid placeholder name\n" +
			"          }\n" +
			"        });\n")
	c.assertOneErrorWithDescription("Message parse tree malformed. Unknown placeholder: unknownPlaceholder")
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
}

// Java: JsMessageVisitorTest.testGetMsgWithGoogScope()
func (c *JsMessageVisitorTest) TestGetMsgWithGoogScope() {
	c.extractMessagesSafely(
		"/** @desc Suggestion Code found outside of <head> tag. */\n" +
			"var $jscomp$scope$12345$0$MSG_CONSUMER_SURVEY_CODE_OUTSIDE_BODY_TAG =\n" +
			"goog.getMsg('Code should be added to <body> tag.');\n")

	msg := c.assertOneMessage()
	AssertThat(msg.GetId()).IsEqualTo("MSG_CONSUMER_SURVEY_CODE_OUTSIDE_BODY_TAG")
	AssertThat(msg.GetKey()).
		IsEqualTo("$jscomp$scope$12345$0$MSG_CONSUMER_SURVEY_CODE_OUTSIDE_BODY_TAG")
}

// Java: JsMessageVisitorTest.assertOneMessage()
func (c *JsMessageVisitorTest) assertOneMessage() *JsMessage {
	AssertThat(c.messages).HasSize(1)
	return c.messages[0]
}

// Java: JsMessageVisitorTest.assertOneError(DiagnosticType, String)
func (c *JsMessageVisitorTest) assertOneErrorWithTypeAndDescription(type_ *DiagnosticType, description string) {
	c.assertOneError(type_)
	c.assertOneErrorWithDescription(description)
}

// Java: JsMessageVisitorTest.assertOneError(DiagnosticType)
func (c *JsMessageVisitorTest) assertOneError(type_ *DiagnosticType) {
	AssertThat(c.compiler.GetErrors()).
		ComparingElementsUsing(JSCompCorrespondences_DIAGNOSTIC_EQUALITY).
		ContainsExactly(type_)
}

// Java: JsMessageVisitorTest.assertOneError(String)
func (c *JsMessageVisitorTest) assertOneErrorWithDescription(description string) {
	AssertThat(c.compiler.GetErrors()).
		ComparingElementsUsing(JSCompCorrespondences_DESCRIPTION_EQUALITY).
		ContainsExactly(description)
}

// Java: JsMessageVisitorTest.assertOneWarning(DiagnosticType, String)
func (c *JsMessageVisitorTest) assertOneWarningWithTypeAndDescription(type_ *DiagnosticType, description string) {
	c.assertOneWarning(type_)
	c.assertOneWarningWithDescription(description)
}

// Java: JsMessageVisitorTest.assertOneWarning(DiagnosticType)
func (c *JsMessageVisitorTest) assertOneWarning(type_ *DiagnosticType) {
	AssertThat(c.compiler.GetWarnings()).
		ComparingElementsUsing(JSCompCorrespondences_DIAGNOSTIC_EQUALITY).
		ContainsExactly(type_)
}

// Java: JsMessageVisitorTest.assertOneWarning(String)
func (c *JsMessageVisitorTest) assertOneWarningWithDescription(description string) {
	AssertThat(c.compiler.GetWarnings()).
		ComparingElementsUsing(JSCompCorrespondences_DESCRIPTION_EQUALITY).
		ContainsExactly(description)
}

// Java: JsMessageVisitorTest.extractMessagesSafely(String)
func (c *JsMessageVisitorTest) extractMessagesSafely(input string) {
	c.extractMessages(input)
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
	AssertThat(c.compiler.GetErrors()).IsEmpty()
}

// Java: JsMessageVisitorTest.extractMessages(String)
func (c *JsMessageVisitorTest) extractMessages(input string) {
	c.compiler = NewCompiler()
	if c.compilerOptions != nil {
		c.compiler.InitOptions(c.compilerOptions)
	}
	root := c.compiler.ParseTestCode(input)
	var visitor JsMessageVisitor = newJsMessageVisitorTest_CollectMessages(c, c.compiler)
	if c.renameMessages {
		renameMessagesVisitor := newJsMessageVisitorTest_RenameMessagesVisitor()
		NodeTraversal_Traverse(c.compiler, root, renameMessagesVisitor)
	}
	visitor.Process(nil, root)
}

// Java: private class JsMessageVisitorTest.CollectMessages extends JsMessageVisitor
type jsMessageVisitorTest_CollectMessages struct {
	JsMessageVisitorBase
	outer *JsMessageVisitorTest
}

// Java: JsMessageVisitorTest.CollectMessages(Compiler)
func newJsMessageVisitorTest_CollectMessages(outer *JsMessageVisitorTest, compiler *Compiler) *jsMessageVisitorTest_CollectMessages {
	v := &jsMessageVisitorTest_CollectMessages{outer: outer}
	v.InitJsMessageVisitorBase(v, compiler, nil)
	return v
}

// Java: JsMessageVisitorTest.CollectMessages.processJsMessageDefinition(JsMessageDefinition)
func (v *jsMessageVisitorTest_CollectMessages) ProcessJsMessageDefinition(definition JsMessageDefinition) {
	v.outer.messages = append(v.outer.messages, definition.GetMessage())
	v.outer.messageDefinitions = append(v.outer.messageDefinitions, definition)
}

// Java: JsMessageVisitorTest.CollectMessages.processIcuTemplateDefinition(IcuTemplateDefinition)
func (v *jsMessageVisitorTest_CollectMessages) ProcessIcuTemplateDefinition(definition IcuTemplateDefinition) {
	v.outer.icuTemplateDefinitions = append(v.outer.icuTemplateDefinitions, definition)
	v.outer.messages = append(v.outer.messages, definition.GetMessage())
}

// Java: private static class JsMessageVisitorTest.DummyJsVisitor extends JsMessageVisitor
type jsMessageVisitorTest_DummyJsVisitor struct {
	JsMessageVisitorBase
}

// Java: JsMessageVisitorTest.DummyJsVisitor()
func newJsMessageVisitorTest_DummyJsVisitor() *jsMessageVisitorTest_DummyJsVisitor {
	v := &jsMessageVisitorTest_DummyJsVisitor{}
	v.InitJsMessageVisitorBase(v, nil, nil)
	return v
}

// Java: JsMessageVisitorTest.DummyJsVisitor.processJsMessageDefinition(JsMessageDefinition)
func (v *jsMessageVisitorTest_DummyJsVisitor) ProcessJsMessageDefinition(definition JsMessageDefinition) {
	// no-op
}

// Java: JsMessageVisitorTest.DummyJsVisitor.processIcuTemplateDefinition(IcuTemplateDefinition)
func (v *jsMessageVisitorTest_DummyJsVisitor) ProcessIcuTemplateDefinition(definition IcuTemplateDefinition) {
	// no-op
}
