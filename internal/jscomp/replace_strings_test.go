// Port of com/google/javascript/jscomp/ReplaceStringsTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/ReplaceStringsTest.java)

/*
 * Copyright 2010 The Closure Compiler Authors.
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

// ReplaceStringsTest tests ReplaceStrings.
//
// Java: public final class ReplaceStringsTest extends CompilerTestCase
type ReplaceStringsTest struct {
	CompilerTestCaseBase

	pass                      *ReplaceStrings
	runDisambiguateProperties bool
	rename                    bool

	functionsToInspect []string
}

// Java: ReplaceStringsTest.DEFAULT_FUNCTIONS_TO_INSPECT
var replaceStringsTest_DEFAULT_FUNCTIONS_TO_INSPECT = []string{
	"Error(?)",
	"goog.debug.Trace.startTracer(*)",
	"goog.debug.Logger.getLogger(?)",
	"goog.log.getLogger(?)",
	"goog.log.info(,?)",
	"goog.log.multiString(,?,?,)",
	"Excluded(?):!testcode",
	"NotExcluded(?):!unmatchable",
}

// Java: ReplaceStringsTest.EXTERNS
var replaceStringsTest_EXTERNS = CompilerTestCase_MINIMAL_EXTERNS +
	`var goog = {};
goog.debug = {};
/** @constructor */
goog.debug.Trace = function() {};
goog.debug.Trace.startTracer = function (var_args) {};
/** @constructor */
goog.debug.Logger = function() {};
goog.debug.Logger.prototype.info = function(msg, opt_ex) {};
/**
 * @param {?} name
 * @return {!goog.debug.Logger}
 */
goog.debug.Logger.getLogger = function(name){};
goog.log = {}
goog.log.getLogger = function(name){};
goog.log.info = function(logger, msg, opt_ex) {};
goog.log.multiString = function(logger, replace1, replace2, keep) {};
`

// Java: ReplaceStringsTest()
func NewReplaceStringsTest() *ReplaceStringsTest {
	c := &ReplaceStringsTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, replaceStringsTest_EXTERNS)
	return c
}

func TestReplaceStringsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewReplaceStringsTest() })
}

// Java: ReplaceStringsTest.getOptions()
func (c *ReplaceStringsTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_MISSING_PROPERTIES, CheckLevel_OFF)
	return options
}

// Java: ReplaceStringsTest.setUp()
func (c *ReplaceStringsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.EnableParseTypeInfo()
	c.SetGenericNameReplacements(Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS)
	c.functionsToInspect = replaceStringsTest_DEFAULT_FUNCTIONS_TO_INSPECT
	c.runDisambiguateProperties = false
	c.rename = false
}

// Java: private static class ReplaceStringsTest.Renamer extends AbstractPostOrderCallback
type replaceStringsTest_Renamer struct {
	NodeTraversal_AbstractPostOrderCallbackBase
}

// Java: ReplaceStringsTest.Renamer.visit(NodeTraversal, Node, Node)
func (r *replaceStringsTest_Renamer) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsName() || n.IsGetProp() {
		originalName := n.GetString()
		n.SetOriginalName(&originalName)
		n.SetString("renamed_" + originalName)
		t.ReportCodeChange()
	}
}

// Java: ReplaceStringsTest.getProcessor(Compiler)
func (c *ReplaceStringsTest) GetProcessor(compiler *Compiler) CompilerPass {
	c.pass = NewReplaceStrings(compiler, "`", c.functionsToInspect)

	return CompilerPassFunc(func(externs rhino.Node, js rhino.Node) {
		if c.rename {
			NodeTraversal_Traverse(compiler, js, &replaceStringsTest_Renamer{})
		}
		NewEs6NormalizeClasses(compiler).Process(externs, js)
		InlineAndCollapseProperties_Builder_(compiler).
			SetPropertyCollapseLevel(CompilerOptions_PropertyCollapseLevel_ALL).
			SetChunkOutputType(CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE).
			SetHaveModulesBeenRewritten(false).
			SetModuleResolutionMode(ModuleLoader_ResolutionMode_BROWSER).
			Build().
			Process(externs, js)
		if c.runDisambiguateProperties {
			sia := SourceInformationAnnotator_Create()
			NodeTraversal_Traverse(compiler, js, sia)

			NewDisambiguateProperties(compiler, jutil.NewLinkedHashSet[string]("foobar")).Process(externs, js)
		}
		c.pass.Process(externs, js)
	})
}

// Java: ReplaceStringsTest.testRenameName()
func (c *ReplaceStringsTest) TestRenameName() {
	c.rename = true
	c.TestDebugStrings("Error('xyz');", "renamed_Error('a');", []string{"a", "xyz"})
}

// Java: ReplaceStringsTest.testRenameStaticProp()
func (c *ReplaceStringsTest) TestRenameStaticProp() {
	c.rename = true
	c.TestDebugStrings(
		"goog.debug.Trace.startTracer('HistoryManager.updateHistory');",
		"renamed_goog.renamed_debug.renamed_Trace.renamed_startTracer('a');",
		[]string{"a", "HistoryManager.updateHistory"})
}

// Java: ReplaceStringsTest.testThrowError2()
func (c *ReplaceStringsTest) TestThrowError2() {
	c.TestDebugStrings(
		"throw Error('x' +\n    'yz');", "throw Error('a');", []string{"a", "xyz"})
}

// Java: ReplaceStringsTest.testThrowError3()
func (c *ReplaceStringsTest) TestThrowError3() {
	c.TestDebugStrings(
		"throw Error('Unhandled mail' + ' search type ' + type);",
		"throw Error('a' + '`' + type);",
		[]string{"a", "Unhandled mail search type `"})
}

// Java: ReplaceStringsTest.testThrowError3a()
func (c *ReplaceStringsTest) TestThrowError3a() {
	c.TestDebugStrings(
		`/** @const */ var preposition = 'in';
/** @const */ var action = 'search';
/** @const */ var error = 'Unhandled ' + action;
throw Error(error + ' ' + type + ' ' + preposition + ' ' + search);
`,
		"/** @const */ var preposition = 'in';\n/** @const */ var action = 'search';\n/** @const */ var error = 'Unhandled ' + action;\nthrow Error('a' + '`' + type + '`' + search);\n",
		[]string{"a", "Unhandled search ` in `"})
}

// Java: ReplaceStringsTest.testThrowError4()
func (c *ReplaceStringsTest) TestThrowError4() {
	c.TestDebugStrings(
		`/** @constructor */
var A = function() {};
A.prototype.m = function(child) {
  if (this.haveChild(child)) {
    throw Error('Node: ' + this.getDataPath() +
                ' already has a child named ' + child);
  } else if (child.parentNode) {
    throw Error('Node: ' + child.getDataPath() +
                ' already has a parent');
  }
  child.parentNode = this;
};
`,
		"/** @constructor */\nvar A = function(){};\nA.prototype.m = function(child) {\n  if (this.haveChild(child)) {\n    throw Error('a' + '`' + this.getDataPath() + '`' + child);\n  } else if (child.parentNode) {\n    throw Error('b' + '`' + child.getDataPath());\n  }\n  child.parentNode = this;\n};\n",
		[]string{
			"a", "Node: ` already has a child named `", "b", "Node: ` already has a parent",
		})
}

// Java: ReplaceStringsTest.testThrowError_templateLiteral()
func (c *ReplaceStringsTest) TestThrowError_templateLiteral() {
	c.TestDebugStrings(
		"throw Error(`Unhandled search ${type} in ${search}`);",
		"throw Error('a' + '`' + type + '`' + search);",
		[]string{"a", "Unhandled search ` in `"})
}

// Java: ReplaceStringsTest.testThrowError_templateLiteralWithConstantTerms()
func (c *ReplaceStringsTest) TestThrowError_templateLiteralWithConstantTerms() {
	c.TestDebugStrings(
		"const preposition = 'in';\nconst action = 'search';\nconst error = `Unhandled ${action}`;\nthrow Error(`${error} ${'type ' + getType()} ${preposition} ${search}`);\n",
		"const preposition = 'in';\nconst action = 'search';\nconst error = `Unhandled ${action}`;\nthrow Error('a' + '`' + getType() + '`' + search);\n",
		[]string{"a", "Unhandled search type ` in `"})
}

// Java: ReplaceStringsTest.testThrowError_templateLiteralWithNonconstantTerms()
func (c *ReplaceStringsTest) TestThrowError_templateLiteralWithNonconstantTerms() {
	// NOTE: We can only inline error if it is *transitively* constant. Otherwise it is left alone,
	// since any transitive strings may have changed, and we certainly don't want to call functions
	// a second time.
	c.TestDebugStrings(
		"const error = `Unhandled ${action}`;\nthrow Error(`${error} ${type} in ${search}`);\n",
		"const error = `Unhandled ${action}`;\nthrow Error('a' + '`' + error + '`' + type + '`' + search);\n",
		[]string{"a", "` ` in `"})
}

// Java: ReplaceStringsTest.testThrowError_templateLiteralConcatenation()
func (c *ReplaceStringsTest) TestThrowError_templateLiteralConcatenation() {
	c.TestDebugStrings(
		"throw Error(`Unhandled mail` + ` search type ${type}`);",
		"throw Error('a' + '`' + type);",
		[]string{"a", "Unhandled mail search type `"})
}

// Java: ReplaceStringsTest.testThrowNonStringError()
func (c *ReplaceStringsTest) TestThrowNonStringError() {
	// No replacement is done when an error is neither a string literal nor
	// a string concatenation expression.
	c.TestDebugStrings("throw Error(x('abc'));", "throw Error(x('abc'));", []string{})
}

// Java: ReplaceStringsTest.testThrowError_taggedTemplateLiteral()
func (c *ReplaceStringsTest) TestThrowError_taggedTemplateLiteral() {
	// No replacement is done when there is a tag function.
	c.TestDebugStrings(
		"throw Error(x`abc`);", //
		"throw Error(x`abc`);",
		[]string{})
}

// Java: ReplaceStringsTest.testThrowConstStringError()
func (c *ReplaceStringsTest) TestThrowConstStringError() {
	c.TestDebugStrings(
		"var AA = 'uvw', AB = AA + 'xyz'; throw Error(AB);",
		"var AA = 'uvw', AB = AA + 'xyz'; throw Error('a');",
		[]string{"a", "uvwxyz"})
}

// Java: ReplaceStringsTest.testThrowConstStringError_templateLiteral()
func (c *ReplaceStringsTest) TestThrowConstStringError_templateLiteral() {
	c.TestDebugStrings(
		"var AA = 'uvw', AB = `${AA}xyz`; throw Error(AB);",
		"var AA = 'uvw', AB = `${AA}xyz`; throw Error('a');",
		[]string{"a", "uvwxyz"})
}

// Java: ReplaceStringsTest.testThrowNewError1()
func (c *ReplaceStringsTest) TestThrowNewError1() {
	c.TestDebugStrings(
		"throw new Error('abc');", "throw new Error('a');", []string{"a", "abc"})
}

// Java: ReplaceStringsTest.testThrowNewError2()
func (c *ReplaceStringsTest) TestThrowNewError2() {
	c.TestDebugStrings("throw new Error();", "throw new Error();", []string{})
}

// Java: ReplaceStringsTest.testStartTracer1()
func (c *ReplaceStringsTest) TestStartTracer1() {
	c.TestDebugStrings(
		"goog.debug.Trace.startTracer('HistoryManager.updateHistory');",
		"goog.debug.Trace.startTracer('a');",
		[]string{"a", "HistoryManager.updateHistory"})
}

// Java: ReplaceStringsTest.testStartTracer2()
func (c *ReplaceStringsTest) TestStartTracer2() {
	c.TestDebugStrings(
		"goog$debug$Trace.startTracer('HistoryManager', 'updateHistory');",
		"goog$debug$Trace.startTracer('a', 'b');",
		[]string{
			"a", "HistoryManager",
			"b", "updateHistory",
		})
}

// Java: ReplaceStringsTest.testStartTracer3()
func (c *ReplaceStringsTest) TestStartTracer3() {
	c.TestDebugStrings(
		`goog$debug$Trace.startTracer('ThreadlistView',
                             'Updating ' + array.length + ' rows');
`,
		"goog$debug$Trace.startTracer('a', 'b' + '`' + array.length);",
		[]string{"a", "ThreadlistView", "b", "Updating ` rows"})
}

// Java: ReplaceStringsTest.testStartTracer4()
func (c *ReplaceStringsTest) TestStartTracer4() {
	c.TestDebugStrings(
		"goog.debug.Trace.startTracer(s, 'HistoryManager.updateHistory');",
		"goog.debug.Trace.startTracer(s, 'a');",
		[]string{"a", "HistoryManager.updateHistory"})
}

// Java: ReplaceStringsTest.testLoggerInitialization()
func (c *ReplaceStringsTest) TestLoggerInitialization() {
	c.TestDebugStrings(
		"goog$debug$Logger$getLogger('my.app.Application');",
		"goog$debug$Logger$getLogger('a');",
		[]string{"a", "my.app.Application"})
}

// Java: ReplaceStringsTest.testLoggerOnVar()
func (c *ReplaceStringsTest) TestLoggerOnVar() {
	c.TestDebugStrings(
		`var logger = goog.debug.Logger.getLogger('foo');
logger.info('Some message');
`,
		`var logger = goog.debug.Logger.getLogger('a');
logger.info('Some message');
`,
		[]string{"a", "foo"})
}

// Java: ReplaceStringsTest.testRepeatedErrorString1()
func (c *ReplaceStringsTest) TestRepeatedErrorString1() {
	c.TestDebugStrings(
		"Error('abc');Error('def');Error('abc');",
		"Error('a');Error('b');Error('a');",
		[]string{"a", "abc", "b", "def"})
}

// Java: ReplaceStringsTest.testRepeatedErrorString2()
func (c *ReplaceStringsTest) TestRepeatedErrorString2() {
	c.TestDebugStrings(
		"Error('a:' + u + ', b:' + v); Error('a:' + x + ', b:' + y);",
		"Error('a' + '`' + u + '`' + v); Error('a' + '`' + x + '`' + y);",
		[]string{"a", "a:`, b:`"})
}

// Java: ReplaceStringsTest.testRepeatedErrorString3()
func (c *ReplaceStringsTest) TestRepeatedErrorString3() {
	c.TestDebugStrings(
		"var AB = 'b'; throw Error(AB); throw Error(AB);",
		"var AB = 'b'; throw Error('a'); throw Error('a');",
		[]string{"a", "b"})
}

// Java: ReplaceStringsTest.testRepeatedTracerString()
func (c *ReplaceStringsTest) TestRepeatedTracerString() {
	c.TestDebugStrings(
		"goog$debug$Trace.startTracer('A', 'B', 'A');",
		"goog$debug$Trace.startTracer('a', 'b', 'a');",
		[]string{"a", "A", "b", "B"})
}

// Java: ReplaceStringsTest.testRepeatedLoggerString()
func (c *ReplaceStringsTest) TestRepeatedLoggerString() {
	c.TestDebugStrings(
		`goog$debug$Logger$getLogger('goog.net.XhrTransport');
goog$debug$Logger$getLogger('my.app.Application');
goog$debug$Logger$getLogger('my.app.Application');
`,
		`goog$debug$Logger$getLogger('a');
goog$debug$Logger$getLogger('b');
goog$debug$Logger$getLogger('b');
`,
		[]string{"a", "goog.net.XhrTransport", "b", "my.app.Application"})
}

// Java: ReplaceStringsTest.testRepeatedStringsWithDifferentMethods()
func (c *ReplaceStringsTest) TestRepeatedStringsWithDifferentMethods() {
	c.Test(
		`throw Error('A');
goog$debug$Trace.startTracer('B', 'A');
goog$debug$Logger$getLogger('C');
goog$debug$Logger$getLogger('B');
goog$debug$Logger$getLogger('A');
throw Error('D');
throw Error('C');
throw Error('B');
throw Error('A');
`,
		`throw Error('a');
goog$debug$Trace.startTracer('b', 'a');
goog$debug$Logger$getLogger('c');
goog$debug$Logger$getLogger('b');
goog$debug$Logger$getLogger('a');
throw Error('d');
throw Error('c');
throw Error('b');
throw Error('a');
`)
}

// Java: ReplaceStringsTest.testLoggerWithNoReplacedParam()
func (c *ReplaceStringsTest) TestLoggerWithNoReplacedParam() {
	c.TestDebugStrings(
		`var x = {};
x.logger_ = goog.log.getLogger('foo');
goog.log.info(x.logger_, 'Some message');
`,
		`var x$logger_ = goog.log.getLogger('a');
goog.log.info(x$logger_, 'b');
`,
		[]string{
			"a", "foo",
			"b", "Some message",
		})
}

// Java: ReplaceStringsTest.testLoggerWithSomeParametersNotReplaced()
func (c *ReplaceStringsTest) TestLoggerWithSomeParametersNotReplaced() {
	c.TestDebugStrings(
		`var x = {};
x.logger_ = goog.log.getLogger('foo');
goog.log.multiString(x.logger_, 'Some message', 'Some message2',
'Do not replace');
`,
		`var x$logger_ = goog.log.getLogger('a');
goog.log.multiString(x$logger_, 'b', 'c', 'Do not replace');
`,
		[]string{
			"a", "foo",
			"b", "Some message",
			"c", "Some message2",
		})
}

// Java: ReplaceStringsTest.testWarningForTaggedTemplates_unqualifiedName()
func (c *ReplaceStringsTest) TestWarningForTaggedTemplates_unqualifiedName() {
	c.TestWarning(
		"throw Error`Unhandled mail search type ${type}`;",
		ReplaceStrings_STRING_REPLACEMENT_TAGGED_TEMPLATE)
}

// Java: ReplaceStringsTest.testWarningForTaggedTemplates_qualifiedName()
func (c *ReplaceStringsTest) TestWarningForTaggedTemplates_qualifiedName() {
	c.TestWarning(
		"goog.debug.Logger.getLogger`foo`;", //
		ReplaceStrings_STRING_REPLACEMENT_TAGGED_TEMPLATE)
}

// Java: ReplaceStringsTest.testWarnsIfPassingPrototypeMethod()
func (c *ReplaceStringsTest) TestWarnsIfPassingPrototypeMethod() {
	// ReplaceStrings supported this configuration until November 2020, so make sure users don't
	// pass it thinking it is still supported.
	c.AllowSourcelessWarnings()

	builder := []string{}
	builder = append(builder, "A.prototype.f(?)")
	c.functionsToInspect = builder
	c.TestError("", ReplaceStrings_BAD_REPLACEMENT_CONFIGURATION)
}

// Java: ReplaceStringsTest.testExcludedFile()
func (c *ReplaceStringsTest) TestExcludedFile() {
	c.TestDebugStrings("Excluded('xyz');", "Excluded('xyz');", []string{})
	c.TestDebugStrings("NotExcluded('xyz');", "NotExcluded('a');", []string{"a", "xyz"})
}

// Java: ReplaceStringsTest.testDebugStrings(String, String, String[])
func (c *ReplaceStringsTest) TestDebugStrings(js string, expected string, substitutedStrings []string) {
	// Verify that the strings are substituted correctly in the JS code.
	c.Test(js, expected)

	results := c.pass.GetResult()
	AssertThat(len(substitutedStrings) % 2).IsEqualTo(0)
	AssertThat(results).HasSize(len(substitutedStrings) / 2)

	// Verify that substituted strings are decoded correctly.
	for i := 0; i < len(substitutedStrings); i += 2 {
		result := results[i/2]
		original := substitutedStrings[i+1]
		AssertThat(result.original).IsEqualTo(original)

		replacement := substitutedStrings[i]
		AssertThat(result.replacement).IsEqualTo(replacement)
	}
}
