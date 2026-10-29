// Port of com/google/javascript/jscomp/CheckConformanceTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/CheckConformanceTest.java)

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
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckConformanceTest tests for CheckConformance.
type CheckConformanceTest struct {
	CompilerTestCaseBase

	configuration string

	baseConfiguration      *string
	extendingConfiguration *string
	reportingMode          CompilerOptions_ConformanceReportingMode
}

// Java: CheckConformanceTest.EXTERNS
var checkConformanceTest_EXTERNS = CompilerTestCase_DEFAULT_EXTERNS +
	`/** @constructor */ function Window() {};
/** @type {Window} */ var window;
/** @type {Function} */ Arguments.prototype.callee;
/** @type {Function} */ Arguments.prototype.caller;
/** @type {Arguments} */ var arguments;
/** @constructor
 * @param {*=} opt_message
 * @param {*=} opt_file
 * @param {*=} opt_line
 * @return {!Error}
*/
function Error(opt_message, opt_file, opt_line) {};
function alert(y) {};
/** @constructor */ function ObjectWithNoProps() {};
function eval() {}
`

// Java: CheckConformanceTest.DEFAULT_CONFORMANCE
const checkConformanceTest_DEFAULT_CONFORMANCE = `requirement: {
  type: BANNED_NAME
  value: 'eval'
   error_message: 'eval is not allowed'
}

requirement: {
  type: BANNED_PROPERTY
  value: 'Arguments.prototype.callee'
  error_message: 'Arguments.prototype.callee is not allowed'
}
`

// Java: CheckConformanceTest()
func NewCheckConformanceTest() *CheckConformanceTest {
	c := &CheckConformanceTest{}
	c.reportingMode = CompilerOptions_ConformanceReportingMode_IGNORE_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG
	c.InitCompilerTestCaseBaseWithExterns(c, checkConformanceTest_EXTERNS)
	return c
}

func TestCheckConformanceTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckConformanceTest() })
}

// Java: CheckConformanceTest.getProcessor(Compiler)
func (c *CheckConformanceTest) GetProcessor(compiler *Compiler) CompilerPass {
	builder := ConformanceConfig_NewBuilder()
	baseBuilder := ConformanceConfig_NewBuilder()
	extendingBuilder := ConformanceConfig_NewBuilder()
	jutil.Try(func() {
		TextFormat_Merge(c.configuration, builder)
		if c.baseConfiguration != nil {
			TextFormat_Merge(*c.baseConfiguration, baseBuilder)
		}
		if c.extendingConfiguration != nil {
			TextFormat_Merge(*c.extendingConfiguration, extendingBuilder)
		}
	}, func(e *TextFormat_ParseException) {
		panic(jutil.NewRuntimeExceptionWithCause(e))
	})

	reportingMode := c.reportingMode
	return NewCheckConformanceWithConfigs(
		compiler,
		[]*ConformanceConfig{builder.Build(), baseBuilder.Build(), extendingBuilder.Build()},
		&reportingMode)
}

// TestConformance is the overloaded private helper testConformance:
// (String, String), (String, String, String) and (String, String,
// DiagnosticType).
//
// Go-only: language-adapter — the overload is chosen from the argument
// types, as javac does.
//
// Java: CheckConformanceTest.testConformance(String, String) / (String, String, String) / (String, String, DiagnosticType)
func (c *CheckConformanceTest) TestConformance(args ...any) {
	src1, src2 := args[0].(string), args[1].(string)
	switch len(args) {
	case 2:
		inputs := []*SourceFile{SourceFile_FromCode("SRC1", src1), SourceFile_FromCode("SRC2", src2)}
		c.TestNoWarning(c.Srcs(inputs))
		return
	case 3:
		switch a := args[2].(type) {
		case string:
			inputs := []*SourceFile{
				SourceFile_FromCode("SRC1", src1),
				SourceFile_FromCode("SRC2", src2),
				SourceFile_FromCode("SRC3", a),
			}
			c.TestNoWarning(c.Srcs(inputs))
			return
		case *DiagnosticType:
			inputs := []*SourceFile{SourceFile_FromCode("SRC1", src1), SourceFile_FromCode("SRC2", src2)}
			c.TestWarning(c.Srcs(inputs), c.Warning(a))
			return
		}
	}
	panic(jutil.NewIllegalArgumentExceptionMsg("testConformance: no overload"))
}

// Java: CheckConformanceTest.config(String, String, String...) (static)
func (c *CheckConformanceTest) Config(rule string, message string, fields ...string) string {
	result := strings.ReplaceAll(`requirement: {
  type: CUSTOM
  java_class: 'RULE'
`, "RULE", rule)
	for _, field := range fields {
		result += field
	}
	result += strings.ReplaceAll(`  error_message: 'MESSAGE'
}
`, "MESSAGE", message)
	return result
}

// Java: CheckConformanceTest.setUp()
func (c *CheckConformanceTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
	c.EnableClosurePass()
	c.EnableClosurePassForExpected()
	c.EnableRewriteClosureCode()
	c.configuration = checkConformanceTest_DEFAULT_CONFORMANCE
	c.IgnoreWarnings(DiagnosticGroups_MISSING_PROPERTIES)
}

// Java: CheckConformanceTest.testViolation1()
func (c *CheckConformanceTest) TestViolation1() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
}
`
	c.TestWarning("eval()", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestWarning(
		"Function.prototype.name; eval.name.length", CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testViolation_templateLiteralInvalidWithEscapeSequence()
func (c *CheckConformanceTest) TestViolation_templateLiteralInvalidWithEscapeSequence() {
	c.configuration =
		`requirement: {
  type: BANNED_STRING_REGEX
  value: '.*blah.*'
  error_message: 'Example error message'
}
`
	// invalid escape sequence is not allowed for template literals
	c.TestError("`\\000`", RhinoErrorReporter_PARSE_ERROR)
	// invalid escape sequence is allowed for tagged template literals
	c.TestNoWarning("function f(x) {  } f`\\000`;")
	c.TestWarning("`blah`", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestWarning("function f(x) {  } f`\\000blah\\000`;", CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testViolation2()
func (c *CheckConformanceTest) TestViolation2() {
	c.TestWarning("function f() { arguments.callee }", CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testConfigFile()
func (c *CheckConformanceTest) TestConfigFile() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  config_file: 'foo_conformance_proto.txt'
  error_message: 'eval is not allowed'
}
`
	c.TestWarning("eval()", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestWarning(
		"Function.prototype.name; eval.name.length",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: eval is not allowed\n  defined in foo_conformance_proto.txt")
}

// Java: CheckConformanceTest.testNotViolation1()
func (c *CheckConformanceTest) TestNotViolation1() {
	c.TestNoWarning(
		`/** @constructor */ function Foo() { this.callee = 'string'; }
/** @constructor */ function Bar() { this.callee = 1; }


function f() {
  var x;
  switch(random()) {
    case 1:
      x = new Foo();
      break;
    case 2:
      x = new Bar();
      break;
    default:
      return;
  }
  var z = x.callee;
}
`)
}

// Java: CheckConformanceTest.testNotViolation2()
func (c *CheckConformanceTest) TestNotViolation2() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'location'
  error_message: 'location is not allowed'
}
`
	c.TestNoWarning("function f() { var location = null; }")
}

// Java: CheckConformanceTest.testMaybeViolation1()
func (c *CheckConformanceTest) TestMaybeViolation1() {
	c.TestWarning("function f() { y.callee }", CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION)
	c.TestWarning(
		"function f() { new Foo().callee }", CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION)
	c.TestWarning(
		"function f() { new Object().callee }", CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION)
	c.TestWarning(
		"function f() { /** @type {*} */ var x; x.callee }",
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION)
	c.TestNoWarning("function f() {/** @const */ var x = {}; x.callee = 1; x.callee}")
}

// Java: CheckConformanceTest.testBadWhitelist1()
func (c *CheckConformanceTest) TestBadWhitelist1() {
	c.AllowSourcelessWarnings()
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'placeholder'
  whitelist_regexp: '('
}
`
	c.TestError(
		"anything;",
		CheckConformance_INVALID_REQUIREMENT_SPEC,
		`Invalid requirement. Reason: invalid regex pattern
Requirement spec:
error_message: "placeholder"
whitelist_regexp: "("
type: BANNED_NAME
value: "eval"
`)
}

// Java: CheckConformanceTest.testBadAllowlist1()
func (c *CheckConformanceTest) TestBadAllowlist1() {
	c.AllowSourcelessWarnings()
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'placeholder'
  allowlist_regexp: '('
}
`
	c.TestError(
		"anything;",
		CheckConformance_INVALID_REQUIREMENT_SPEC,
		`Invalid requirement. Reason: invalid regex pattern
Requirement spec:
error_message: "placeholder"
type: BANNED_NAME
value: "eval"
allowlist_regexp: "("
`)
}

// Java: CheckConformanceTest.testViolationWhitelisted1()
func (c *CheckConformanceTest) TestViolationWhitelisted1() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  whitelist: 'testcode'
}
`
	c.TestNoWarning("eval()")
}

// Java: CheckConformanceTest.testViolationAllowlisted1()
func (c *CheckConformanceTest) TestViolationAllowlisted1() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  allowlist: 'testcode'
}
`
	c.TestNoWarning("eval()")
}

// Java: CheckConformanceTest.testViolationWhitelistedAndAllowlistedDuplicated1()
func (c *CheckConformanceTest) TestViolationWhitelistedAndAllowlistedDuplicated1() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  whitelist: 'testcode'
  allowlist: 'testcode'
}
`
	c.TestNoWarning("eval()")
}

// Java: CheckConformanceTest.testViolationWhitelistedByWhitelistEntryPrefix()
func (c *CheckConformanceTest) TestViolationWhitelistedByWhitelistEntryPrefix() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  whitelist_entry {
    prefix: 'testcode'
  }
}
`
	c.TestNoWarning("eval()")
}

// Java: CheckConformanceTest.testViolationAllowlistedByAllowlistEntryPrefix()
func (c *CheckConformanceTest) TestViolationAllowlistedByAllowlistEntryPrefix() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  allowlist_entry {
    prefix: 'testcode'
  }
}
`
	c.TestNoWarning("eval()")
}

// Java: CheckConformanceTest.testViolationWhitelistedByWhitelistEntryRegexp()
func (c *CheckConformanceTest) TestViolationWhitelistedByWhitelistEntryRegexp() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  whitelist_entry {
    regexp: 'tes..ode'
  }
}
`
	c.TestNoWarning("eval()")
}

// Java: CheckConformanceTest.testViolationAllowlistedByAllowlistEntryRegexp()
func (c *CheckConformanceTest) TestViolationAllowlistedByAllowlistEntryRegexp() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  allowlist_entry {
    regexp: 'tes..ode'
  }
}
`
	c.TestNoWarning("eval()")
}

// Java: CheckConformanceTest.testViolationWhitelisted2()
func (c *CheckConformanceTest) TestViolationWhitelisted2() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  whitelist_regexp: 'code$'
}
`
	c.TestNoWarning("eval()")
}

// Java: CheckConformanceTest.testViolationAllowlisted2()
func (c *CheckConformanceTest) TestViolationAllowlisted2() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  allowlist_regexp: 'code$'
}
`
	c.TestNoWarning("eval()")
}

// Java: CheckConformanceTest.testViolationAllowlisted2Ts()
func (c *CheckConformanceTest) TestViolationAllowlisted2Ts() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  allowlist_regexp: 'file.ts$'
}
`
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("file.closure.js", "eval()")))
}

// Java: CheckConformanceTest.testViolationWhitelistedIgnoresRegex()
func (c *CheckConformanceTest) TestViolationWhitelistedIgnoresRegex() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  whitelist: 'file.js'
}
`
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("test/google3/file.js", "eval()")))
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("blaze-out/k8-opt/bin/file.js", "eval()")))
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("google3/blaze-out/k8-opt/bin/file.js", "eval()")))
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("bazel-out/k8-opt/bin/file.js", "eval()")))
}

// Java: CheckConformanceTest.testViolationAllowlistedIgnoresRegex()
func (c *CheckConformanceTest) TestViolationAllowlistedIgnoresRegex() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  allowlist: 'file.js'
}
`
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("test/google3/file.js", "eval()")))
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("blaze-out/k8-opt/bin/file.js", "eval()")))
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("google3/blaze-out/k8-opt/bin/file.js", "eval()")))
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("bazel-out/k8-opt/bin/file.js", "eval()")))
}

// Java: CheckConformanceTest.testViolationAllowlistedIgnoresRegexTs()
func (c *CheckConformanceTest) TestViolationAllowlistedIgnoresRegexTs() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  allowlist: 'file.ts'
}
`
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("file.closure.js", "eval()")))
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("test/google3/file.closure.js", "eval()")))
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("blaze-out/k8-opt/bin/file.closure.js", "eval()")))
	c.TestNoWarning(
		c.Srcs(SourceFile_FromCode("google3/blaze-out/k8-opt/bin/file.closure.js", "eval()")))
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("bazel-out/k8-opt/bin/file.closure.js", "eval()")))
}

// Java: CheckConformanceTest.testViolationWhitelistedAndAllowlistedIgnoresRegex()
func (c *CheckConformanceTest) TestViolationWhitelistedAndAllowlistedIgnoresRegex() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  whitelist: 'file1.js'
  allowlist: 'file2.js'
}
`
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("test/google3/file1.js", "eval()")))
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("test/google3/file2.js", "eval()")))
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("blaze-out/k8-opt/bin/file1.js", "eval()")))
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("blaze-out/k8-opt/bin/file2.js", "eval()")))
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("google3/blaze-out/k8-opt/bin/file1.js", "eval()")))
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("google3/blaze-out/k8-opt/bin/file2.js", "eval()")))
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("bazel-out/k8-opt/bin/file1.js", "eval()")))
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("bazel-out/k8-opt/bin/file2.js", "eval()")))
}

// Java: CheckConformanceTest.testViolationWhitelistedIgnoresRegex_absolutePath()
func (c *CheckConformanceTest) TestViolationWhitelistedIgnoresRegex_absolutePath() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  whitelist: '/file.js'
}
`
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("/file.js", "eval()")))
}

// Java: CheckConformanceTest.testViolationAllowlistedIgnoresRegex_absolutePath()
func (c *CheckConformanceTest) TestViolationAllowlistedIgnoresRegex_absolutePath() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  allowlist: '/file.js'
}
`
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("/file.js", "eval()")))
}

// Java: CheckConformanceTest.testViolationWhitelistedAndAllowlistedIgnoresRegex_absolutePath()
func (c *CheckConformanceTest) TestViolationWhitelistedAndAllowlistedIgnoresRegex_absolutePath() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  whitelist: '/file1.js'
  allowlist: '/file2.js'
}
`
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("/file1.js", "eval()")))
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("/file2.js", "eval()")))
}

// Java: CheckConformanceTest.testViolationWhitelistedIgnoresRegex_genfiles()
func (c *CheckConformanceTest) TestViolationWhitelistedIgnoresRegex_genfiles() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  whitelist: 'file.js'
}
`
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("blaze-out/k8-opt/genfiles/file.js", "eval()")))
}

// Java: CheckConformanceTest.testViolationAllowlistedIgnoresRegex_genfiles()
func (c *CheckConformanceTest) TestViolationAllowlistedIgnoresRegex_genfiles() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  allowlist: 'file.js'
}
`
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("blaze-out/k8-opt/genfiles/file.js", "eval()")))
}

// Java: CheckConformanceTest.testViolationWhitelistedAndAllowlistedIgnoresRegex_genfiles()
func (c *CheckConformanceTest) TestViolationWhitelistedAndAllowlistedIgnoresRegex_genfiles() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  whitelist: 'file1.js'
  allowlist: 'file2.js'
}
`
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("blaze-out/k8-opt/genfiles/file1.js", "eval()")))
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("blaze-out/k8-opt/genfiles/file2.js", "eval()")))
}

// Java: CheckConformanceTest.testFileOnOnlyApplyToIsChecked()
func (c *CheckConformanceTest) TestFileOnOnlyApplyToIsChecked() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  only_apply_to: 'foo.js'
}
`
	inputs := []*SourceFile{SourceFile_FromCode("foo.js", "eval()")}
	c.TestWarning(
		c.Srcs(inputs),
		c.Warning(CheckConformance_CONFORMANCE_VIOLATION).
			WithMessage("Violation: eval is not allowed"))
}

// Java: CheckConformanceTest.testFileOnOnlyApplyToIsCheckedTs()
func (c *CheckConformanceTest) TestFileOnOnlyApplyToIsCheckedTs() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  only_apply_to: 'foo.ts'
}
`
	inputs := []*SourceFile{SourceFile_FromCode("foo.closure.js", "eval()")}
	c.TestWarning(
		c.Srcs(inputs),
		c.Warning(CheckConformance_CONFORMANCE_VIOLATION).
			WithMessage("Violation: eval is not allowed"))
}

// Java: CheckConformanceTest.testFileNotOnOnlyApplyToIsNotChecked()
func (c *CheckConformanceTest) TestFileNotOnOnlyApplyToIsNotChecked() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  only_apply_to: 'foo.js'
}
`
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("bar.js", "eval()")))
}

// Java: CheckConformanceTest.testFileOnOnlyApplyToRegexpIsChecked()
func (c *CheckConformanceTest) TestFileOnOnlyApplyToRegexpIsChecked() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  only_apply_to_regexp: 'test.js$'
}
`
	input := []*SourceFile{SourceFile_FromCode("foo_test.js", "eval()")}
	c.TestWarning(
		c.Srcs(input),
		c.Warning(CheckConformance_CONFORMANCE_VIOLATION).
			WithMessage("Violation: eval is not allowed"))
}

// Java: CheckConformanceTest.testFileOnOnlyApplyToRegexpIsCheckedTs()
func (c *CheckConformanceTest) TestFileOnOnlyApplyToRegexpIsCheckedTs() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  only_apply_to_regexp: 'test.ts$'
}
`
	input := []*SourceFile{SourceFile_FromCode("foo_test.closure.js", "eval()")}
	c.TestWarning(
		c.Srcs(input),
		c.Warning(CheckConformance_CONFORMANCE_VIOLATION).
			WithMessage("Violation: eval is not allowed"))
}

// Java: CheckConformanceTest.testFileNotOnOnlyApplyToRegexpIsNotChecked()
func (c *CheckConformanceTest) TestFileNotOnOnlyApplyToRegexpIsNotChecked() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
  only_apply_to_regexp: 'test.js$'
}
`
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("bar.js", "eval()")))
}

// Java: CheckConformanceTest.testNoOp()
func (c *CheckConformanceTest) TestNoOp() {
	c.configuration =
		`requirement: {
  type: NO_OP,
  value: 'no matter',
  value: 'can be anything',
  error_message: 'Never happens'
}
`
	c.TestNoWarning("eval()")
}

// Java: CheckConformanceTest.testBannedEnhance()
func (c *CheckConformanceTest) TestBannedEnhance() {
	c.configuration =
		`requirement: {
  type: BANNED_ENHANCE
  value: '{some.banned.namespace}'
  value: '{another.banned.namespace}'
  error_message: 'Enhanced namespace is not allowed.'
}
`
	violationMessage := `Violation: Enhanced namespace is not allowed.
The enhanced namespace `
	ban1 := `/**
 * @enhance {some.banned.namespace}
 */
`
	c.TestWarning(
		ban1,
		CheckConformance_CONFORMANCE_VIOLATION,
		violationMessage+"\"{some.banned.namespace}\"")
	ban2 := `/**
 * @enhance {another.banned.namespace}
 */
`
	c.TestWarning(
		ban2,
		CheckConformance_CONFORMANCE_VIOLATION,
		violationMessage+"\"{another.banned.namespace}\"")
	allow1 := `/**
 * @enhance {some.allowed.namespace}
 */
`
	c.TestNoWarning(allow1)
	allow2 := `/**
 * @fileoverview no enhance annotation should always pass.
 */
`
	c.TestNoWarning(allow2)
}

// Java: CheckConformanceTest.testBannedModsRegex()
func (c *CheckConformanceTest) TestBannedModsRegex() {
	c.configuration =
		`requirement: {
  type: BANNED_MODS_REGEX
  value: '.+_bar$'
  error_message: 'Modding namespaces ending in _bar is NOT allowed.'
}
`
	violationMessage := "Violation: Modding namespaces ending in _bar is NOT allowed."
	ban := `/**
 * @mods {ns.foo_bar}
 * @modName {ns.foo_bar_baz}
 */
`
	c.TestWarning(ban, CheckConformance_CONFORMANCE_VIOLATION, violationMessage)
	allow1 := `/**
 * @mods {ns.foo}
 * @modName {ns.foo_bar}
 */
`
	c.TestNoWarning(allow1)
	allow2 := `/**
 * @fileoverview no enhance annotation should always pass.
 */
`
	c.TestNoWarning(allow2)
}

// Java: CheckConformanceTest.testBannedNameCall()
func (c *CheckConformanceTest) TestBannedNameCall() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME_CALL
  value: 'Function'
  error_message: 'Calling Function is not allowed.'
}
`
	c.TestNoWarning("f instanceof Function")
	c.TestWarning("new Function(str);", CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedName_googProvided()
func (c *CheckConformanceTest) TestBannedName_googProvided() {
	c.configuration =
		`requirement: {
  type: BANNED_NAME
  value: 'foo.bar'
  error_message: 'foo.bar is not allowed'
  allowlist: 'SRC1'
}
`
	c.TestWarning(
		c.Srcs(
			SourceFile_FromCode("SRC1", "goog.provide('foo.bar');"),
			SourceFile_FromCode("SRC2", "alert(foo.bar);")),
		CheckConformance_CONFORMANCE_VIOLATION)
	c.TestWarning(
		c.Srcs(
			SourceFile_FromCode("SRC1", "goog.provide('foo.bar'); foo.bar = {};"),
			SourceFile_FromCode("SRC2", "alert(foo.bar);")),
		CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testInferredConstCheck()
func (c *CheckConformanceTest) TestInferredConstCheck() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$InferredConstCheck'
  error_message: 'Failed to infer type of constant'
}
`
	c.TestNoWarning("/** @const */ var x = 0;")
	c.TestNoWarning("/** @const */ var x = unknown;")
	c.TestNoWarning("const x = unknown;")
	c.TestWarning(
		`/** @constructor */
function f() {
  /** @const */ this.foo = unknown;
}
var x = new f();
`,
		CheckConformance_CONFORMANCE_VIOLATION)
	c.TestWarning(
		`/** @constructor */
function f() {}
/** @this {f} */
var init_f = function() {
  /** @const */ this.foo = unknown;
};
var x = new f();
`,
		CheckConformance_CONFORMANCE_VIOLATION)
	c.TestWarning(
		`/** @constructor */
function f() {}
var init_f = function() {
  /** @const */ this.FOO = unknown;
};
var x = new f();
`,
		CheckConformance_CONFORMANCE_VIOLATION)
	c.TestWarning(
		`/** @constructor */
function f() {}
f.prototype.init_f = function() {
  /** @const */ this.FOO = unknown;
};
var x = new f();
`,
		CheckConformance_CONFORMANCE_VIOLATION)
	c.TestNoWarning(
		`/** @constructor */
function f() {}
f.prototype.init_f = function() {
  /** @const {?} */ this.FOO = unknown;
};
var x = new f();
`)
	c.TestNoWarning(
		`/** @const */
var ns = {};
/** @const */
ns.subns = ns.subns || {};
`)
	// We only check @const nodes, not @final nodes.
	c.TestNoWarning(
		`/** @constructor */
function f() {
  /** @final */ this.foo = unknown;
}
var x = new f();
`)
}

// Java: CheckConformanceTest.testBannedCodePattern1()
func (c *CheckConformanceTest) TestBannedCodePattern1() {
	c.configuration =
		`requirement: {
  type: BANNED_CODE_PATTERN
  value: '/** @param {string|String} a */function template(a) {a.blink}'
  error_message: 'blink is annoying'
}
`
	externs := checkConformanceTest_EXTERNS + "String.prototype.blink;"
	c.TestNoWarning(
		`/** @constructor */ function Foo() { this.blink = 1; }
var foo = new Foo();
foo.blink();
`)
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("'foo'.blink;"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: blink is annoying")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("'foo'.blink();"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: blink is annoying")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("String('foo').blink();"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: blink is annoying")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("foo.blink();"),
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION,
		`Possible violation: blink is annoying
The type information available for this expression is too loose to ensure conformance.
`)
}

// Java: CheckConformanceTest.testBannedDep1()
func (c *CheckConformanceTest) TestBannedDep1() {
	c.configuration =
		`requirement: {
  type: BANNED_DEPENDENCY
  value: 'testcode'
  error_message: 'testcode is not allowed'
}
`
	c.TestWarning(
		"anything;", CheckConformance_CONFORMANCE_VIOLATION, "Violation: testcode is not allowed")
}

// Java: CheckConformanceTest.testBannedDepRegexNoValue()
func (c *CheckConformanceTest) TestBannedDepRegexNoValue() {
	c.AllowSourcelessWarnings()
	c.configuration =
		`requirement: {
  type: BANNED_DEPENDENCY_REGEX
  error_message: 'testcode is not allowed'
}
`
	c.TestError(
		"anything;",
		CheckConformance_INVALID_REQUIREMENT_SPEC,
		`Invalid requirement. Reason: missing value (no banned dependency regexps)
Requirement spec:
error_message: "testcode is not allowed"
type: BANNED_DEPENDENCY_REGEX
`)
}

// Java: CheckConformanceTest.testBannedDepRegex()
func (c *CheckConformanceTest) TestBannedDepRegex() {
	c.configuration =
		`requirement: {
  type: BANNED_DEPENDENCY_REGEX
  error_message: 'testcode is not allowed'
  value: '.*test.*'
}
`
	c.TestWarning(
		"anything;", CheckConformance_CONFORMANCE_VIOLATION, "Violation: testcode is not allowed")
}

// Java: CheckConformanceTest.testBannedDepRegexWithWhitelist()
func (c *CheckConformanceTest) TestBannedDepRegexWithWhitelist() {
	c.configuration =
		`requirement: {
  type: BANNED_DEPENDENCY_REGEX
  error_message: 'testcode is not allowed'
  value: '.*test.*'
  whitelist_regexp: 'testcode'
}
`
	c.TestNoWarning("anything;")
}

// Java: CheckConformanceTest.testBannedDepRegexWithAllowlist()
func (c *CheckConformanceTest) TestBannedDepRegexWithAllowlist() {
	c.configuration =
		`requirement: {
  type: BANNED_DEPENDENCY_REGEX
  error_message: 'testcode is not allowed'
  value: '.*test.*'
  allowlist_regexp: 'testcode'
}
`
	c.TestNoWarning("anything;")
}

// Java: CheckConformanceTest.testReportLooseTypeViolations()
func (c *CheckConformanceTest) TestReportLooseTypeViolations() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY_WRITE
  value: 'HTMLScriptElement.prototype.textContent'
  error_message: 'Setting content of <script> is dangerous.'
  report_loose_type_violations: false
}
`
	externs := CompilerTestCase_DEFAULT_EXTERNS +
		`/** @constructor */ function Element() {}
/** @type {string} @implicitCast */
Element.prototype.textContent;
/** @constructor @extends {Element} */ function HTMLScriptElement() {}

`
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("(new HTMLScriptElement).textContent = 'alert(1);'"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: Setting content of <script> is dangerous.")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("HTMLScriptElement.prototype.textContent = 'alert(1);'"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: Setting content of <script> is dangerous.")
	c.TestNoWarning(c.Externs(externs), c.Srcs("(new Element).textContent = 'safe'"))
}

// Java: CheckConformanceTest.testDontCrashOnNonConstructorWithPrototype()
func (c *CheckConformanceTest) TestDontCrashOnNonConstructorWithPrototype() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY_WRITE
  value: 'Bar.prototype.method'
  error_message: 'asdf'
  report_loose_type_violations: false
}
`
	c.TestNoWarning(
		c.Externs(
			CompilerTestCase_DEFAULT_EXTERNS+
				`/** @constructor */
function Bar() {}
Bar.prototype.method = function() {};
`),
		c.Srcs(
			`function Foo() {}
Foo.prototype.method = function() {};
`))
}

// Java: CheckConformanceTest.testBannedPropertyWhitelist0()
func (c *CheckConformanceTest) TestBannedPropertyWhitelist0() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY
  value: 'C.prototype.p'
  error_message: 'C.p is not allowed'
  whitelist: 'SRC1'
}
`
	cDecl := `/** @constructor */
function C() {}
/** @type {string} */
C.prototype.p;
`
	dDecl := `/** @constructor */ function D() {}
/** @type {string} */
D.prototype.p;
`
	c.TestConformance(cDecl, dDecl)
	typedefDecl := `/** @typedef {{p: string}} */ let E;
/** @const {!E} */ const value = {p: 'foo'};
/** @type {string} */ const s = value.p;
`
	c.TestConformance(cDecl, dDecl, typedefDecl)
}

// Java: CheckConformanceTest.testBannedPropertyAllowlist0()
func (c *CheckConformanceTest) TestBannedPropertyAllowlist0() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY
  value: 'C.prototype.p'
  error_message: 'C.p is not allowed'
  allowlist: 'SRC1'
}
`
	cDecl := `/** @constructor */
function C() {}
/** @type {string} */
C.prototype.p;
`
	dDecl := `/** @constructor */ function D() {}
/** @type {string} */
D.prototype.p;
`
	c.TestConformance(cDecl, dDecl)
}

// Java: CheckConformanceTest.testBannedPropertyWhitelist1()
func (c *CheckConformanceTest) TestBannedPropertyWhitelist1() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY
  value: 'C.prototype.p'
  error_message: 'C.p is not allowed'
  whitelist: 'SRC1'
}
`
	cDecl := `/** @constructor */
function C() {
  this.p = 'str';
}
`
	dDecl := `/** @constructor */
function D() {
  this.p = 'str';
}
`
	c.TestConformance(cDecl, dDecl)
}

// Java: CheckConformanceTest.testBannedPropertyAllowlist1()
func (c *CheckConformanceTest) TestBannedPropertyAllowlist1() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY
  value: 'C.prototype.p'
  error_message: 'C.p is not allowed'
  allowlist: 'SRC1'
}
`
	cDecl := `/** @constructor */
function C() {
  this.p = 'str';
}
`
	dDecl := `/** @constructor */
function D() {
  this.p = 'str';
}
`
	c.TestConformance(cDecl, dDecl)
}

// Java: CheckConformanceTest.testBannedPropertyWhitelist2()
func (c *CheckConformanceTest) TestBannedPropertyWhitelist2() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY
  value: 'C.prototype.p'
  error_message: 'C.p is not allowed'
  whitelist: 'SRC1'
}
`
	declarations := `/** @constructor */ function SC() {}
/** @constructor @extends {SC} */
function C() {}
/** @type {string} */
C.prototype.p;
/** @constructor */ function D() {}
/** @type {string} */
D.prototype.p;
`
	c.TestConformance(declarations, "var d = new D(); d.p = 'boo';")
	c.TestConformance(
		declarations, "var c = new C(); c.p = 'boo';", CheckConformance_CONFORMANCE_VIOLATION)
	// Accessing property through a super type is possibly a violation.
	c.TestConformance(
		declarations,
		"var sc = new SC(); sc.p = 'boo';",
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION)
	c.TestConformance(
		declarations, "var c = new C(); var foo = c.p;", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestConformance(
		declarations,
		"var c = new C(); var foo = 'x' + c.p;",
		CheckConformance_CONFORMANCE_VIOLATION)
	c.TestConformance(
		declarations, "var c = new C(); c['p'] = 'boo';", CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedPropertyAllowlist2()
func (c *CheckConformanceTest) TestBannedPropertyAllowlist2() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY
  value: 'C.prototype.p'
  error_message: 'C.p is not allowed'
  allowlist: 'SRC1'
}
`
	declarations := `/** @constructor */ function SC() {}
/** @constructor @extends {SC} */
function C() {}
/** @type {string} */
C.prototype.p;
/** @constructor */ function D() {}
/** @type {string} */
D.prototype.p;
`
	c.TestConformance(declarations, "var d = new D(); d.p = 'boo';")
	c.TestConformance(
		declarations, "var c = new C(); c.p = 'boo';", CheckConformance_CONFORMANCE_VIOLATION)
	// Accessing property through a super type is possibly a violation.
	c.TestConformance(
		declarations,
		"var sc = new SC(); sc.p = 'boo';",
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION)
	c.TestConformance(
		declarations, "var c = new C(); var foo = c.p;", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestConformance(
		declarations,
		"var c = new C(); var foo = 'x' + c.p;",
		CheckConformance_CONFORMANCE_VIOLATION)
	c.TestConformance(
		declarations, "var c = new C(); c['p'] = 'boo';", CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedPropertyWhitelist3()
func (c *CheckConformanceTest) TestBannedPropertyWhitelist3() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY
  value: 'C.prototype.p'
  error_message: 'C.p is not allowed'
  whitelist: 'SRC1'
}
`
	cdecl := `/** @constructor */ function SC() {}
/** @constructor @extends {SC} */
function C() {}
/** @type {string} */
C.prototype.p;
`
	ddecl := `/** @constructor @template T */ function D() {}
/** @param {T} a */
D.prototype.method = function(a) {
  use(a.p);
};
`
	c.TestConformance(cdecl, ddecl, CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedPropertyAllowlist3()
func (c *CheckConformanceTest) TestBannedPropertyAllowlist3() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY
  value: 'C.prototype.p'
  error_message: 'C.p is not allowed'
  allowlist: 'SRC1'
}
`
	cdecl := `/** @constructor */ function SC() {}
/** @constructor @extends {SC} */
function C() {}
/** @type {string} */
C.prototype.p;
`
	ddecl := `/** @constructor @template T */ function D() {}
/** @param {T} a */
D.prototype.method = function(a) {
  use(a.p);
};
`
	c.TestConformance(cdecl, ddecl, CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedPropertyWhitelist4()
func (c *CheckConformanceTest) TestBannedPropertyWhitelist4() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY
  value: 'C.prototype.p'
  error_message: 'C.p is not allowed'
  whitelist: 'SRC1'
}
`
	cdecl := `/** @constructor */ function SC() {}
/** @constructor @extends {SC} */
function C() {}
/** @type {string} */
C.prototype.p;

/**
 * @param {K} key
 * @param {V=} opt_value
 * @constructor
 * @struct
 * @template K, V
 * @private
 */
var Entry_ = function(key, opt_value) {
  /** @const {K} */
  this.key = key;
  /** @type {V|undefined} */
  this.value = opt_value;
};
`
	ddecl := `/** @constructor @template T */ function D() {}
/** @param {T} a */
D.prototype.method = function(a) {
  var entry = new Entry('key');
  use(entry.value.p);
};
`
	c.TestConformance(cdecl, ddecl, CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedPropertyAllowlist4()
func (c *CheckConformanceTest) TestBannedPropertyAllowlist4() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY
  value: 'C.prototype.p'
  error_message: 'C.p is not allowed'
  allowlist: 'SRC1'
}
`
	cdecl := `/** @constructor */ function SC() {}
/** @constructor @extends {SC} */
function C() {}
/** @type {string} */
C.prototype.p;

/**
 * @param {K} key
 * @param {V=} opt_value
 * @constructor
 * @struct
 * @template K, V
 * @private
 */
var Entry_ = function(key, opt_value) {
  /** @const {K} */
  this.key = key;
  /** @type {V|undefined} */
  this.value = opt_value;
};
`
	ddecl := `/** @constructor @template T */ function D() {}
/** @param {T} a */
D.prototype.method = function(a) {
  var entry = new Entry('key');
  use(entry.value.p);
};
`
	c.TestConformance(cdecl, ddecl, CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedProperty5()
func (c *CheckConformanceTest) TestBannedProperty5() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY
  value: 'Array.prototype.push'
  error_message: 'banned Array.prototype.push'
}
`
	c.TestWarning("[1, 2, 3].push(4);\n", CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedPropertyWhitelist_recordType()
func (c *CheckConformanceTest) TestBannedPropertyWhitelist_recordType() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY
  value: 'Logger.prototype.config'
  error_message: 'Logger.config is not allowed'
  whitelist: 'SRC1'
}
`
	declaration := "class Logger { config() {} }"
	// Fine, because there is no explicit relationship between Logger & GoodRecord.
	c.TestConformance(
		declaration,
		`/** @record */
class GoodRecord {
  constructor() {
    /** @type {Function} */ this.config;
  }
}
`)
	// Bad, because there is a direct relationship.
	c.TestConformance(
		"/** @implements {BadRecord} */ "+declaration,
		`/** @record */
class BadRecord {
  constructor() {
    /** @type {Function} */ this.config;
  }
}
`,
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedPropertyAllowlist_recordType()
func (c *CheckConformanceTest) TestBannedPropertyAllowlist_recordType() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY
  value: 'Logger.prototype.config'
  error_message: 'Logger.config is not allowed'
  allowlist: 'SRC1'
}
`
	declaration := "class Logger { config() {} }"
	// Fine, because there is no explicit relationship between Logger & GoodRecord.
	c.TestConformance(
		declaration,
		`/** @record */
class GoodRecord {
  constructor() {
    /** @type {Function} */ this.config;
  }
}
`)
	// Bad, because there is a direct relationship.
	c.TestConformance(
		"/** @implements {BadRecord} */ "+declaration,
		`/** @record */
class BadRecord {
  constructor() {
    /** @type {Function} */ this.config;
  }
}
`,
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedPropertyWhitelist_namespacedType()
func (c *CheckConformanceTest) TestBannedPropertyWhitelist_namespacedType() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY
  value: 'ns.C.prototype.p'
  error_message: 'ns.C.p is not allowed'
  whitelist: 'SRC1'
}
`
	declarations := `/** @const */
var ns = {};
/** @constructor */ function SC() {}
/** @constructor @extends {SC} */
ns.C = function() {}
/** @type {string} */
ns.C.prototype.p;
/** @constructor */ function D() {}
/** @type {string} */
D.prototype.p;
`
	c.TestConformance(declarations, "var d = new D(); d.p = 'boo';")
	c.TestConformance(
		declarations, "var c = new ns.C(); c.p = 'boo';", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestConformance(
		declarations,
		"var c = new SC(); c.p = 'boo';",
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedPropertyAllowlist_namespacedType()
func (c *CheckConformanceTest) TestBannedPropertyAllowlist_namespacedType() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY
  value: 'ns.C.prototype.p'
  error_message: 'ns.C.p is not allowed'
  allowlist: 'SRC1'
}
`
	declarations := `/** @const */
var ns = {};
/** @constructor */ function SC() {}
/** @constructor @extends {SC} */
ns.C = function() {}
/** @type {string} */
ns.C.prototype.p;
/** @constructor */ function D() {}
/** @type {string} */
D.prototype.p;
`
	c.TestConformance(declarations, "var d = new D(); d.p = 'boo';")
	c.TestConformance(
		declarations, "var c = new ns.C(); c.p = 'boo';", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestConformance(
		declarations,
		"var c = new SC(); c.p = 'boo';",
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedPropertyWhitelist_bundledNamespacedType()
func (c *CheckConformanceTest) TestBannedPropertyWhitelist_bundledNamespacedType() {
	c.DisableRewriteClosureCode()
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY
  value: 'test.bns.BC.prototype.p'
  error_message: 'test.bns.BC.p is not allowed'
  whitelist: 'SRC1'
}
`
	declarations := `/** @fileoverview @typeSummary */
goog.loadModule(function(exports) {
goog.module('test.bns');
/** @constructor */ function SBC() {}
exports.SBC = SBC;
/** @constructor @extends {SBC} */
BC = function() {}
/** @type {string} */
BC.prototype.p;
exports.BC = BC;
/** @constructor */ function D() {}
exports.D = D;
/** @type {string} */
D.prototype.p;
return exports;
});
`
	c.TestConformance(
		declarations,
		`goog.module('test');
const {D} = goog.require('test.bns');
var d = new D();
d.p = 'boo';
`)
	// This case should be a certain violation, but the compiler cannot figure out that the imported
	// type is the same as the one found from the type registry.
	c.TestConformance(
		declarations,
		`goog.module('test');
const {BC} = goog.require('test.bns');
var bc = new BC();
bc.p = 'boo';
`,
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION)
	c.TestConformance(
		declarations,
		`goog.module('test');
const bns = goog.require('test.bns');
var bc = new bns.SBC();
bc.p = 'boo';
`,
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedPropertyAllowlist_bundledNamespacedType()
func (c *CheckConformanceTest) TestBannedPropertyAllowlist_bundledNamespacedType() {
	c.DisableRewriteClosureCode()
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY
  value: 'test.bns.BC.prototype.p'
  error_message: 'test.bns.BC.p is not allowed'
  allowlist: 'SRC1'
}
`
	declarations := `/** @fileoverview @typeSummary */
goog.loadModule(function(exports) {
goog.module('test.bns');
/** @constructor */ function SBC() {}
exports.SBC = SBC;
/** @constructor @extends {SBC} */
BC = function() {}
/** @type {string} */
BC.prototype.p;
exports.BC = BC;
/** @constructor */ function D() {}
exports.D = D;
/** @type {string} */
D.prototype.p;
return exports;
});
`
	c.TestConformance(
		declarations,
		`goog.module('test');
const {D} = goog.require('test.bns');
var d = new D();
d.p = 'boo';
`)
	// This case should be a certain violation, but the compiler cannot figure out that the imported
	// type is the same as the one found from the type registry.
	c.TestConformance(
		declarations,
		`goog.module('test');
const {BC} = goog.require('test.bns');
var bc = new BC();
bc.p = 'boo';
`,
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION)
	c.TestConformance(
		declarations,
		`goog.module('test');
const bns = goog.require('test.bns');
var bc = new bns.SBC();
bc.p = 'boo';
`,
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedPropertyWhitelist_destructuring()
func (c *CheckConformanceTest) TestBannedPropertyWhitelist_destructuring() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY
  value: 'C.prototype.p'
  error_message: 'C.p is not allowed'
  whitelist: 'SRC1'
}
`
	declarations := `/** @constructor */
var C = function() {}
/** @type {string} */
C.prototype.p;
/** @type {number} */
C.prototype.m
`
	c.TestConformance(declarations, "var {m} = new C();")
	c.TestConformance(declarations, "var {p} = new C();", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestConformance(
		declarations, "var {['p']: x} = new C();", CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedPropertyAllowlist_destructuring()
func (c *CheckConformanceTest) TestBannedPropertyAllowlist_destructuring() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY
  value: 'C.prototype.p'
  error_message: 'C.p is not allowed'
  allowlist: 'SRC1'
}
`
	declarations := `/** @constructor */
var C = function() {}
/** @type {string} */
C.prototype.p;
/** @type {number} */
C.prototype.m
`
	c.TestConformance(declarations, "var {m} = new C();")
	c.TestConformance(declarations, "var {p} = new C();", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestConformance(
		declarations, "var {['p']: x} = new C();", CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedPropertyWrite()
func (c *CheckConformanceTest) TestBannedPropertyWrite() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY_WRITE
  value: 'C.prototype.p'
  error_message: 'Assignment to C.p is not allowed'
}
`
	declarations := `/** @constructor */ function C() {}
/** @type {string} */
C.prototype.p;
/** @constructor */ function D() {}
/** @type {string} */
D.prototype.p;
`
	c.TestNoWarning(declarations + "var d = new D(); d.p = 'boo';")
	c.TestWarning(
		declarations+"var c = new C(); c.p = 'boo';", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestNoWarning(declarations + "var c = new C(); var foo = c.p;")
	c.TestNoWarning(declarations + "var c = new C(); var foo = 'x' + c.p;")
	c.TestWarning(
		declarations+"var c = new C(); c['p'] = 'boo';", CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedPropertyWriteExtern()
func (c *CheckConformanceTest) TestBannedPropertyWriteExtern() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY_WRITE
  value: 'Element.prototype.innerHTML'
  error_message: 'Assignment to Element.innerHTML is not allowed'
}
`
	externs := CompilerTestCase_DEFAULT_EXTERNS +
		`/** @constructor */ function Element() {}
/** @type {string} @implicitCast */
Element.prototype.innerHTML;
`
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("var e = new Element(); e.innerHTML = '<boo>';"),
		c.Warning(CheckConformance_CONFORMANCE_VIOLATION))
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("var e = new Element(); e.innerHTML = 'foo';"),
		c.Warning(CheckConformance_CONFORMANCE_VIOLATION))
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("var e = new Element(); e['innerHTML'] = 'foo';"),
		c.Warning(CheckConformance_CONFORMANCE_VIOLATION))
}

// Java: CheckConformanceTest.testBannedPropertyNonConstantWrite()
func (c *CheckConformanceTest) TestBannedPropertyNonConstantWrite() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY_NON_CONSTANT_WRITE
  value: 'C.prototype.p'
  error_message: 'Assignment of a non-constant value to C.p is not allowed'
}
`
	declarations := `/** @constructor */ function C() {}
/** @type {string} */
C.prototype.p;
`
	c.TestNoWarning(declarations + "var c = new C(); c.p = 'boo';")
	c.TestNoWarning(declarations + "var c = new C(); c.p = 'foo' + 'bar';")
	c.TestWarning(
		declarations+"var boo = 'boo'; var c = new C(); c.p = boo;",
		CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedPropertyRead()
func (c *CheckConformanceTest) TestBannedPropertyRead() {
	c.configuration =
		`requirement: {
  type: BANNED_PROPERTY_READ
  value: 'C.prototype.p'
  error_message: 'Use of C.p is not allowed'
}
`
	declarations := `/** @constructor */ function C() {}
/** @type {string} */
C.prototype.p;
/** @constructor */ function D() {}
/** @type {string} */
D.prototype.p;
function use(a) {};
`
	c.TestNoWarning(declarations + "var d = new D(); d.p = 'boo';")
	c.TestNoWarning(declarations + "var c = new C(); c.p = 'boo';")
	c.TestWarning(
		declarations+"var c = new C(); use(c.p);", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestWarning(
		declarations+"var c = new C(); var foo = c.p;", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestWarning(
		declarations+"var c = new C(); var foo = 'x' + c.p;",
		CheckConformance_CONFORMANCE_VIOLATION)
	c.TestNoWarning(declarations + "var c = new C(); c['p'] = 'boo';")
	c.TestWarning(
		declarations+"var c = new C(); use(c['p']);", CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedStringRegexMissingValues()
func (c *CheckConformanceTest) TestBannedStringRegexMissingValues() {
	c.AllowSourcelessWarnings()
	c.configuration =
		`requirement: {
  type: BANNED_STRING_REGEX
  error_message: 'Empty string not allowed'
}
`
	c.TestError(
		"anything;",
		CheckConformance_INVALID_REQUIREMENT_SPEC,
		`Invalid requirement. Reason: missing value
Requirement spec:
error_message: "Empty string not allowed"
type: BANNED_STRING_REGEX
`)
}

// Java: CheckConformanceTest.testBannedStringRegexEmpty()
func (c *CheckConformanceTest) TestBannedStringRegexEmpty() {
	c.AllowSourcelessWarnings()
	c.configuration =
		`requirement: {
  type: BANNED_STRING_REGEX
  value: ' '
  error_message: 'Empty string not allowed'
}
`
	c.TestError(
		"anything;",
		CheckConformance_INVALID_REQUIREMENT_SPEC,
		`Invalid requirement. Reason: empty strings or whitespace are not allowed
Requirement spec:
error_message: "Empty string not allowed"
type: BANNED_STRING_REGEX
value: " "
`)
}

// Java: CheckConformanceTest.testBannedStringRegexMultipleValuesWithEmpty()
func (c *CheckConformanceTest) TestBannedStringRegexMultipleValuesWithEmpty() {
	c.AllowSourcelessWarnings()
	c.configuration =
		`requirement: {
  type: BANNED_STRING_REGEX
  value: 'things'
  value: ' '
  error_message: 'Empty string not allowed'
}
`
	c.TestError(
		"anything;",
		CheckConformance_INVALID_REQUIREMENT_SPEC,
		`Invalid requirement. Reason: empty strings or whitespace are not allowed
Requirement spec:
error_message: "Empty string not allowed"
type: BANNED_STRING_REGEX
value: "things"
value: " "
`)
}

// Java: CheckConformanceTest.testBannedStringRegex1()
func (c *CheckConformanceTest) TestBannedStringRegex1() {
	c.configuration =
		`requirement: {
  type: BANNED_STRING_REGEX
  value: '.*some-attr.*'
  error_message: 'Empty string not allowed'
}
`
	declarations := "let dom = '<div>sample dom template content</div>';"
	c.TestNoWarning(declarations)
	c.TestWarning(
		declarations+"let moredom = '<p some-attr=\"testval\">reflected!</p>';",
		CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedStringRegex2()
func (c *CheckConformanceTest) TestBannedStringRegex2() {
	c.configuration =
		`requirement: {
  type: BANNED_STRING_REGEX
  value: '.*things.*'
  value: 'stuff.*'
  error_message: 'Empty string not allowed'
}
`
	code := `/** @constructor */
function Base() {}; Base.prototype.m;
/** @constructor @extends {Base} */
function Sub() {}
`
	stuff := "var b = 'stuff and what not';\n"
	things := "var s = 'special things';\n"
	c.TestNoWarning(code)
	c.TestWarning(code+stuff, CheckConformance_CONFORMANCE_VIOLATION)
	c.TestWarning(code+things, CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedStringRegexExactMatch()
func (c *CheckConformanceTest) TestBannedStringRegexExactMatch() {
	c.configuration =
		`requirement: {
  type: BANNED_STRING_REGEX
  value: 'stuff'
  error_message: 'Empty string not allowed'
}
`
	code := `/** @constructor */
function Base() {}; Base.prototype.m;
/** @constructor @extends {Base} */
function Sub() {}
`
	noMatch := "var b = ' stuff ';\n"
	shouldMatch := "var s = 'stuff';\n"
	c.TestNoWarning(code + noMatch)
	c.TestWarning(code+shouldMatch, CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testBannedStringTemplateLiteral1()
func (c *CheckConformanceTest) TestBannedStringTemplateLiteral1() {
	c.configuration =
		`requirement: {
  type: BANNED_STRING_REGEX
  value: '.*good'
  error_message: 'Empty string not allowed'
}
`
	code := `/** @constructor */
function Base() {}; Base.prototype.m;
/** @constructor @extends {Base} */
function Sub() {}
`
	noMatch := "var b = `cheesy goodness`;\n"
	shouldMatch := "var b = `cheesy good`;\n"
	c.TestNoWarning(code + noMatch)
	c.TestWarning(code+shouldMatch, CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testRestrictedCall1()
func (c *CheckConformanceTest) TestRestrictedCall1() {
	c.configuration =
		`requirement: {
  type: RESTRICTED_METHOD_CALL
  value: 'C.prototype.m:function(number)'
  error_message: 'm method param must be number'
}
`
	code := `/** @constructor */ function C() {}
/** @param {*} a */
C.prototype.m = function(a){}
`
	c.TestNoWarning(code + "new C().m(1);")
	c.TestWarning(code+"new C().m('str');", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestNoWarning(code + "new C().m.call(new C(), 1);")
	c.TestWarning(code+"new C().m.call(new C(), 'str');", CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testRestrictedCall2()
func (c *CheckConformanceTest) TestRestrictedCall2() {
	c.configuration =
		`requirement: {
  type: RESTRICTED_NAME_CALL
  value: 'C.m:function(number)'
  error_message: 'C.m method param must be number'
}
`
	code := `/** @constructor */ function C() {}
/** @param {*} a */
C.m = function(a){}
`
	c.TestNoWarning(code + "C.m(1);")
	c.TestWarning(code+"C.m('str');", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestNoWarning(code + "C.m.call(C, 1);")
	c.TestWarning(code+"C.m.call(C, 'str');", CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testRestrictedCall3()
func (c *CheckConformanceTest) TestRestrictedCall3() {
	c.configuration =
		`requirement: {
  type: RESTRICTED_NAME_CALL
  value: 'C:function(number)'
  error_message: 'C method must be number'
}
`
	code := "/** @constructor @param {...*} a */ function C(a) {}\n"
	c.TestNoWarning(code + "new C(1);")
	c.TestWarning(code+"new C('str');", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestWarning(code+"new C(1, 1);", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestWarning(code+"new C();", CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testRestrictedCall4()
func (c *CheckConformanceTest) TestRestrictedCall4() {
	c.configuration =
		`requirement: {
  type: RESTRICTED_NAME_CALL
  value: 'C:function(number)'
  error_message: 'C method must be number'
}
`
	code := "/** @constructor @param {...*} a */ function C(a) {}\n"
	c.TestNoWarning(c.Externs(checkConformanceTest_EXTERNS+"goog.inherits;"), c.Srcs(code+"goog.inherits(A, C);"))
}

// Java: CheckConformanceTest.testRestrictedMethodCallThisType()
func (c *CheckConformanceTest) TestRestrictedMethodCallThisType() {
	c.configuration =
		`requirement: {
  type: RESTRICTED_METHOD_CALL
  value: 'Base.prototype.m:function(this:Sub,number)'
  error_message: 'Only call m on the subclass'
}
`
	code := `/** @constructor */
function Base() {}; Base.prototype.m;
/** @constructor @extends {Base} */
function Sub() {}
var b = new Base();
var s = new Sub();
var maybeB = cond ? new Base() : null;
var maybeS = cond ? new Sub() : null;
`
	c.TestWarning(code+"b.m(1)", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestWarning(code+"maybeB.m(1)", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestNoWarning(code + "s.m(1)")
	c.TestNoWarning(code + "maybeS.m(1)")
}

// Java: CheckConformanceTest.testRestrictedMethodCallUsingCallThisType()
func (c *CheckConformanceTest) TestRestrictedMethodCallUsingCallThisType() {
	c.configuration =
		`requirement: {
  type: RESTRICTED_METHOD_CALL
  value: 'Base.prototype.m:function(this:Sub,number)'
  error_message: 'Only call m on the subclass'
}
`
	code := `/** @constructor */
function Base() {}; Base.prototype.m;
/** @constructor @extends {Base} */
function Sub() {}
var b = new Base();
var s = new Sub();
var maybeB = cond ? new Base() : null;
var maybeS = cond ? new Sub() : null;
`
	c.TestWarning(code+"b.m.call(b, 1)", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestWarning(code+"b.m.call(maybeB, 1)", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestNoWarning(code + "b.m.call(s, 1)")
	c.TestNoWarning(code + "b.m.call(maybeS, 1)")
}

// Java: CheckConformanceTest.testRestrictedPropertyWrite()
func (c *CheckConformanceTest) TestRestrictedPropertyWrite() {
	c.configuration =
		`requirement: {
  type: RESTRICTED_PROPERTY_WRITE
  value: 'Base.prototype.x:number'
  error_message: 'Only assign number'
}
`
	code := `/** @constructor */
function Base() {}; Base.prototype.x;
var b = new Base();
`
	c.TestWarning(code+"b.x = 'a'", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestNoWarning(code + "b.x = 1")
	c.TestNoWarning(code + "var a = {}; a.x = 'a'")
	c.TestWarning(
		code+"/** @type {?} */ var a = {}; a.x = 'a'",
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION)
	c.TestWarning(
		code+"/** @type {*} */ var a = {}; a.x = 'a'",
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION)
}

// Java: CheckConformanceTest.testCustom1()
func (c *CheckConformanceTest) TestCustom1() {
	c.AllowSourcelessWarnings()
	c.configuration =
		`requirement: {
  type: CUSTOM
  error_message: 'placeholder'
}
`
	c.TestError(
		"anything;",
		CheckConformance_INVALID_REQUIREMENT_SPEC,
		`Invalid requirement. Reason: missing java_class
Requirement spec:
error_message: "placeholder"
type: CUSTOM
`)
}

// Java: CheckConformanceTest.testCustomBanForOf()
func (c *CheckConformanceTest) TestCustomBanForOf() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanForOf'
  error_message: 'BanForOf Message'
}
`
	c.TestWarning(
		"for (var x of y) { var z = x; }",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanForOf Message")
}

// Java: CheckConformanceTest.testCustomRestrictThrow1()
func (c *CheckConformanceTest) TestCustomRestrictThrow1() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanThrowOfNonErrorTypes'
  error_message: 'BanThrowOfNonErrorTypes Message'
}
`
	c.TestWarning(
		"throw 'blah';",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanThrowOfNonErrorTypes Message")
}

// Java: CheckConformanceTest.testCustomRestrictThrow2()
func (c *CheckConformanceTest) TestCustomRestrictThrow2() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanThrowOfNonErrorTypes'
  error_message: 'BanThrowOfNonErrorTypes Message'
}
`
	c.TestNoWarning("throw new Error('test');")
}

// Java: CheckConformanceTest.testCustomRestrictThrow3()
func (c *CheckConformanceTest) TestCustomRestrictThrow3() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanThrowOfNonErrorTypes'
  error_message: 'BanThrowOfNonErrorTypes Message'
}
`
	c.TestNoWarning(
		`/** @param {*} x */
function f(x) {
  throw x;
}
`)
}

// Java: CheckConformanceTest.testCustomRestrictThrow4()
func (c *CheckConformanceTest) TestCustomRestrictThrow4() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanThrowOfNonErrorTypes'
  error_message: 'BanThrowOfNonErrorTypes Message'
}
`
	c.TestNoWarning(
		`/** @constructor @extends {Error} */
function MyError() {}
/** @param {*} x */
function f(x) {
  if (x instanceof MyError) {
  } else {
    throw x;
  }
}
`)
}

// Java: CheckConformanceTest.testCustomBanUnknownThis1()
func (c *CheckConformanceTest) TestCustomBanUnknownThis1() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanUnknownThis'
  error_message: 'BanUnknownThis Message'
}
`
	c.TestWarning(
		"function f() {alert(this);}",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanUnknownThis Message")
}

// TODO(johnlenz): add a unit test for templated "this" values.
// Java: CheckConformanceTest.testCustomBanUnknownThis2()
func (c *CheckConformanceTest) TestCustomBanUnknownThis2() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanUnknownThis'
  error_message: 'BanUnknownThis Message'
}
`
	c.TestNoWarning("/** @constructor */ function C() {alert(this);}")
}

// Java: CheckConformanceTest.testCustomBanUnknownThis3()
func (c *CheckConformanceTest) TestCustomBanUnknownThis3() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanUnknownThis'
  error_message: 'BanUnknownThis Message'
}
`
	c.TestNoWarning("function f() {alert(/** @type {Error} */(this));}")
}

// Java: CheckConformanceTest.testCustomBanUnknownThis_allowsClosurePrimitiveAssert()
func (c *CheckConformanceTest) TestCustomBanUnknownThis_allowsClosurePrimitiveAssert() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanUnknownThis'
  error_message: 'BanUnknownThis Message'
}
`
	assertInstanceof := `/** @const */ var asserts = {};
/**
 * @param {?} value The value to check.
 * @param {function(new: T, ...)} type A user-defined constructor.
 * @return {T}
 * @template T
 * @closurePrimitive {asserts.matchesReturn}
 */
asserts.assertInstanceof = function(value, type) {
  return value;
};
`
	c.TestNoWarning(assertInstanceof + "function f() {asserts.assertInstanceof(this, Error);}")
}

// Java: CheckConformanceTest.testCustomBanUnknownThis_allowsGoogAssert()
func (c *CheckConformanceTest) TestCustomBanUnknownThis_allowsGoogAssert() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanUnknownThis'
  error_message: 'BanUnknownThis Message'
}
`
	c.TestNoWarning("function f() {goog.asserts.assertInstanceof(this, Error);}")
}

// Java: CheckConformanceTest.testCustomBanUnknownThis_allowsTs()
func (c *CheckConformanceTest) TestCustomBanUnknownThis_allowsTs() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanUnknownThis'
  error_message: 'BanUnknownThis Message'
}
`
	c.TestNoWarning(c.Srcs(SourceFile_FromCode("file.closure.js", "function f() {alert(this);}")))
}

// Java: CheckConformanceTest.rule(String)
// (static in Java)
func (c *CheckConformanceTest) Rule(rule string) string {
	return "com.google.javascript.jscomp.ConformanceRules$" + rule
}

// Java: CheckConformanceTest.value(String)
// (static in Java)
func (c *CheckConformanceTest) Value(value string) string {
	return "  value: '" + value + "'\n"
}

// Java: CheckConformanceTest.testBanUnknownDirectThisPropsReferences_implicitUnknownOnEs5Constructor_warn()
func (c *CheckConformanceTest) TestBanUnknownDirectThisPropsReferences_implicitUnknownOnEs5Constructor_warn() {
	c.configuration = c.Config(c.Rule("BanUnknownDirectThisPropsReferences"), "My rule message")
	c.TestWarning(
		`/** @constructor */
function f() {}
f.prototype.prop;
f.prototype.method = function() { alert(this.prop); };
`,
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: My rule message")
}

// Java: CheckConformanceTest.testBanUnknownDirectThisPropsReferences_explicitUnknownOnEs5Constructor_ok()
func (c *CheckConformanceTest) TestBanUnknownDirectThisPropsReferences_explicitUnknownOnEs5Constructor_ok() {
	c.configuration = c.Config(c.Rule("BanUnknownDirectThisPropsReferences"), "My rule message")
	c.TestNoWarning(
		`/** @constructor */
function f() {};
/** @type {?} */
f.prototype.prop;
f.prototype.method = function() { alert(this.prop); }
`)
}

// Java: CheckConformanceTest.testBanUnknownDirectThisPropsReferences_implicitUnknownOnEs6Class_warn()
func (c *CheckConformanceTest) TestBanUnknownDirectThisPropsReferences_implicitUnknownOnEs6Class_warn() {
	c.configuration = c.Config(c.Rule("BanUnknownDirectThisPropsReferences"), "My rule message")
	c.TestWarning(
		`class F {
  constructor() {
    this.prop;
  }
  method() {
    alert(this.prop);
  }
}
`,
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: My rule message")
}

// Java: CheckConformanceTest.testBanUnknownDirectThisPropsReferences_explicitUnknownOnEs6Class_ok()
func (c *CheckConformanceTest) TestBanUnknownDirectThisPropsReferences_explicitUnknownOnEs6Class_ok() {
	c.configuration = c.Config(c.Rule("BanUnknownDirectThisPropsReferences"), "My rule message")
	c.TestNoWarning(
		`class F {
  constructor() {
    /** @type {?} */
    this.prop;
  }
  method() {
    alert(this.prop);
  }
}
`)
}

// Java: CheckConformanceTest.testBanUnknownDirectThisPropsReferences_implicitUnknownOnClassField_warn()
func (c *CheckConformanceTest) TestBanUnknownDirectThisPropsReferences_implicitUnknownOnClassField_warn() {
	c.configuration = c.Config(c.Rule("BanUnknownDirectThisPropsReferences"), "My rule message")
	// TODO(b/192088118): need to fix so test gives warning for implicit field reference
	c.TestNoWarning(
		`class F {
  prop;
  method() {
    alert(this.prop);
  }
}
`)
}

// Java: CheckConformanceTest.testBanUnknownDirectThisPropsReferences_explicitUnknownOnClassField_ok()
func (c *CheckConformanceTest) TestBanUnknownDirectThisPropsReferences_explicitUnknownOnClassField_ok() {
	c.configuration = c.Config(c.Rule("BanUnknownDirectThisPropsReferences"), "My rule message")
	c.TestNoWarning(
		`class F {
  /** @type {?} */
  prop = 2;
  method() {
    alert(this.prop);
  }
}
`)
}

// Java: CheckConformanceTest.testBanUnknownDirectThisPropsReferences_inferredNotUnknown_ok()
func (c *CheckConformanceTest) TestBanUnknownDirectThisPropsReferences_inferredNotUnknown_ok() {
	c.configuration = c.Config(c.Rule("BanUnknownDirectThisPropsReferences"), "My rule message")
	c.TestNoWarning(
		`class F {
  constructor() {
    this.prop = 42;
  }
  method() {
    alert(this.prop);
  }
}
`)
}

// Java: CheckConformanceTest.testBanUnknownDirectThisPropsReferences_implicitUnknownAssignedButNotUsed_ok()
func (c *CheckConformanceTest) TestBanUnknownDirectThisPropsReferences_implicitUnknownAssignedButNotUsed_ok() {
	c.configuration = c.Config(c.Rule("BanUnknownDirectThisPropsReferences"), "My rule message")
	c.TestNoWarning(
		`/** @constructor */
function f() {}
f.prototype.prop;
f.prototype.method = function() { this.prop = foo; };
`)
}

// Java: CheckConformanceTest.testCustomBanUnknownProp1()
func (c *CheckConformanceTest) TestCustomBanUnknownProp1() {
	c.configuration =
		c.Config(c.Rule("BanUnknownTypedClassPropsReferences"), "My rule message", c.Value("String"))
	c.TestWarning(
		`/** @constructor */ function f() {}; f.prototype.prop;
f.prototype.method = function() { alert(this.prop); }
`,
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: My rule message\nThe property \"prop\" on type \"f\"")
}

// Java: CheckConformanceTest.testCustomBanUnknownProp_templateUnion()
func (c *CheckConformanceTest) TestCustomBanUnknownProp_templateUnion() {
	c.configuration =
		c.Config(c.Rule("BanUnknownTypedClassPropsReferences"), "My rule message", c.Value("String"))
	c.TestSame(
		`/** @record @template T */
class X {
  constructor() {
    /** @type {T|!Array<T>} */ this.x;
    f(this.x);
  }
}

/**
 * @param {T|!Array<T>} value
 * @template T
 */
function f(value) {}
`)
}

// Java: CheckConformanceTest.testCustomBanUnknownProp1_es6Class()
func (c *CheckConformanceTest) TestCustomBanUnknownProp1_es6Class() {
	c.configuration =
		c.Config(c.Rule("BanUnknownTypedClassPropsReferences"), "My rule message", c.Value("String"))
	c.TestWarning(
		`class F {
  constructor() {
    this.prop;
  }
}

alert(new F().prop);
`,
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: My rule message\nThe property \"prop\" on type \"F\"")
}

// Java: CheckConformanceTest.testCustomBanUnknownProp2()
func (c *CheckConformanceTest) TestCustomBanUnknownProp2() {
	c.configuration =
		c.Config(c.Rule("BanUnknownTypedClassPropsReferences"), "My rule message", c.Value("String"))
	js := `Object.prototype.foobar;
 /** @param {ObjectWithNoProps} a */
function f(a) { alert(a.foobar); };
`
	c.TestWarning(
		js,
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: My rule message\nThe property \"foobar\" on type \"(ObjectWithNoProps|null)\"")
}

// Java: CheckConformanceTest.testCustomBanUnknownProp3()
func (c *CheckConformanceTest) TestCustomBanUnknownProp3() {
	c.configuration =
		c.Config(c.Rule("BanUnknownTypedClassPropsReferences"), "My rule message", c.Value("String"))
	c.TestNoWarning(
		`/** @constructor */ function f() {}
f.prototype.method = function() { this.prop = foo; };
`)
}

// Java: CheckConformanceTest.testCustomBanUnknownProp4()
func (c *CheckConformanceTest) TestCustomBanUnknownProp4() {
	c.configuration =
		c.Config(c.Rule("BanUnknownTypedClassPropsReferences"), "My rule message", c.Value("String"))
	c.TestNoWarning(
		`/** @constructor */ function f() { /** @type {?} */ this.prop = null; };
f.prototype.method = function() { alert(this.prop); }
`)
}

// Java: CheckConformanceTest.testCustomBanUnknownProp5()
func (c *CheckConformanceTest) TestCustomBanUnknownProp5() {
	c.configuration =
		c.Config(c.Rule("BanUnknownTypedClassPropsReferences"), "My rule message", c.Value("String"))
	c.TestWarning(
		`/** @typedef {?} */ var Unk;
/** @constructor */ function f() { /** @type {?Unk} */ this.prop = null; };
f.prototype.method = function() { alert(this.prop); }
`,
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: My rule message\nThe property \"prop\" on type \"f\"")
}

// Java: CheckConformanceTest.testCustomBanUnknownProp6()
func (c *CheckConformanceTest) TestCustomBanUnknownProp6() {
	c.configuration =
		c.Config(c.Rule("BanUnknownTypedClassPropsReferences"), "My rule message", c.Value("String"))
	c.TestWarning(
		`goog.module('example');
/** @constructor */ function f() { this.prop; };
f.prototype.method = function() { alert(this.prop); }
`,
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: My rule message\nThe property \"prop\" on type \"f\"")
}

// Java: CheckConformanceTest.testCustomBanUnknownProp7()
func (c *CheckConformanceTest) TestCustomBanUnknownProp7() {
	c.configuration =
		c.Config(c.Rule("BanUnknownTypedClassPropsReferences"), "My rule message", c.Value("String"))
	c.TestNoWarning(
		`/** @constructor */
function Foo() {
  /** @type {!Object<number, number>} */
  this.prop;
}
function f(/** !Foo */ x) {
  return x.prop[1] + 123;
}
`)
}

// Java: CheckConformanceTest.testCustomBanUnknownProp_getPropInVoidOperatorDoesntCauseSpuriousWarning()
func (c *CheckConformanceTest) TestCustomBanUnknownProp_getPropInVoidOperatorDoesntCauseSpuriousWarning() {
	// See b/112072360
	c.configuration =
		c.Config(c.Rule("BanUnknownTypedClassPropsReferences"), "My rule message", c.Value("String"))
	c.TestNoWarning(
		`/** @constructor */
function Foo() {
  /** @type {!Object<number, number>} */
  this.prop;
}
const foo = new Foo();
const f = () => void foo.prop;
`)
}

// Java: CheckConformanceTest.testCustomBanUnknownProp_getPropInDestructuringDoesntCauseSpuriousWarning()
func (c *CheckConformanceTest) TestCustomBanUnknownProp_getPropInDestructuringDoesntCauseSpuriousWarning() {
	c.configuration =
		c.Config(c.Rule("BanUnknownTypedClassPropsReferences"), "My rule message", c.Value("String"))
	c.Test(
		c.Externs("var unknownItem;"),
		c.Srcs(
			"/** @constructor */\n"+
				"function Foo() {}\n"+
				"const foo = new Foo();\n"+
				"// note that `foo.prop` is unknown here, but we don't warn because it's not being\n"+
				"// used\n"+
				"[foo.prop] = unknownItem;\n"))
}

// Java: CheckConformanceTest.testCustomBanUnknownProp_unionUndefined()
func (c *CheckConformanceTest) TestCustomBanUnknownProp_unionUndefined() {
	c.configuration = c.Config(c.Rule("BanUnknownTypedClassPropsReferences"), "My rule message")
	c.TestNoWarning(
		`/** @constructor */
function Foo() {}
if (false) {
  /** @type {(null|?)} */
  Foo.prototype.prop;
}
function f() {
  return new Foo().prop;
}
`)
}

// Java: CheckConformanceTest.testSimpleExtends()
func (c *CheckConformanceTest) TestSimpleExtends() {
	c.baseConfiguration =
		jutil.StringPtr(`requirement: {
  rule_id: "gws:goog.Promise.X"
  type: BANNED_NAME
  error_message: "Prefer using native Promise equivalents. See go/gws-js-conformance#goog-promise"

  value: "goog.Promise.all"
}
`)
	c.extendingConfiguration =
		jutil.StringPtr(`requirement: {
  extends: 'gws:goog.Promise.X'
}
`)
	c.TestWarning(
		`goog.Promise.all();
`,
		CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testSimpleExtends_withAllowlist()
func (c *CheckConformanceTest) TestSimpleExtends_withAllowlist() {
	c.baseConfiguration =
		jutil.StringPtr(`requirement: {
  rule_id: "gws:goog.Promise.X"
  type: BANNED_NAME
  error_message: "Prefer using native Promise equivalents. See go/gws-js-conformance#goog-promise"

  value: "goog.Promise.all"
}
`)
	c.extendingConfiguration =
		jutil.StringPtr(`requirement: {
  extends: 'gws:goog.Promise.X'
  allowlist: "testcode"
}
`)
	// respects the allowlist in the extending config and does not warn
	c.TestNoWarning(
		`goog.Promise.all();
`)
}

// Java: CheckConformanceTest.testSimpleExtends_withAllowlistRegexp()
func (c *CheckConformanceTest) TestSimpleExtends_withAllowlistRegexp() {
	c.baseConfiguration =
		jutil.StringPtr(`requirement: {
  rule_id: "gws:goog.Promise.X"
  type: BANNED_NAME
  error_message: "Prefer using native Promise equivalents. See go/gws-js-conformance#goog-promise"

  value: "goog.Promise.all"
}
`)
	c.extendingConfiguration =
		jutil.StringPtr(`requirement: {
  extends: 'gws:goog.Promise.X'
  allowlist_regexp: "code$"
}
`)
	// respects the allowlist_regexp in the extending config and does not warn
	c.TestNoWarning(
		`goog.Promise.all();
`)
}

// Java: CheckConformanceTest.testRequirementHasBothRuleIdAndExtends_reportsInvalidRequirementSpec()
func (c *CheckConformanceTest) TestRequirementHasBothRuleIdAndExtends_reportsInvalidRequirementSpec() {
	c.AllowSourcelessWarnings()
	c.baseConfiguration =
		jutil.StringPtr(`requirement: {
  rule_id: "gws:goog.Promise.X"
  type: BANNED_NAME
  error_message: "Prefer using native Promise equivalents. See go/gws-js-conformance#goog-promise"

  value: "goog.Promise.all"
}
`)
	// has both rule_id and extends fields set
	c.extendingConfiguration =
		jutil.StringPtr(`requirement: {
  rule_id: "extends:gws:goog.Promise.X"
  extends: 'gws:goog.Promise.X'
}
`)
	c.TestError(
		`goog.Promise.all();
`,
		CheckConformance_INVALID_REQUIREMENT_SPEC)
}

// Java: CheckConformanceTest.testMissingBehaviorInLibraryLevelReportingMode()
func (c *CheckConformanceTest) TestMissingBehaviorInLibraryLevelReportingMode() {
	// the conformance reporting mode is RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG, i.e.
	// as if the CheckJS action is being run.
	c.reportingMode = CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG
	c.AllowSourcelessWarnings()
	c.baseConfiguration =
		jutil.StringPtr(`requirement: {
  rule_id: "gws:goog.Promise.X"
  type: BANNED_NAME
  error_message: "Prefer using native Promise equivalents. See go/gws-js-conformance#goog-promise"
  value: "goog.Promise.all"
}
`)
	c.TestWarning(
		`goog.Promise.all();
`,
		CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testSimpleExtends_disallowsOverridingLibraryLevelBehaviorToADifferentValue_inLibraryLevelConformanceMode()
func (c *CheckConformanceTest) TestSimpleExtends_disallowsOverridingLibraryLevelBehaviorToADifferentValue_inLibraryLevelConformanceMode() {
	// the conformance reporting mode is RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG, i.e.
	// as if the CheckJS action is being run.
	c.reportingMode = CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG
	c.AllowSourcelessWarnings()
	c.baseConfiguration =
		jutil.StringPtr(`library_level_non_allowlisted_conformance_violations_behavior: REPORT_AS_BUILD_ERROR
requirement: {
  rule_id: "gws:goog.Promise.X"
  type: BANNED_NAME
  error_message: "Prefer using native Promise equivalents. See go/gws-js-conformance#goog-promise"

  value: "goog.Promise.all"
}
`)
	c.extendingConfiguration =
		jutil.StringPtr(`# tries to override the library-level behavior, not allowed
library_level_non_allowlisted_conformance_violations_behavior: RECORD_ONLY
requirement: {
  extends: 'gws:goog.Promise.X'
}
`)
	c.TestError(
		`goog.Promise.all();
`,
		CheckConformance_INVALID_REQUIREMENT_SPEC,
		`Invalid requirement. Reason: extending rule's config may not specify a different value of 'library_level_non_allowlisted_conformance_violations_behavior' than the base rule's config. Skipping all conformance checks.
Requirement spec:
extends: "gws:goog.Promise.X"
`)
}

// Java: CheckConformanceTest.testSimpleExtends_allowsOverridingLibraryLevelBehaviorWhenBaseRuleHasNoBehavior()
func (c *CheckConformanceTest) TestSimpleExtends_allowsOverridingLibraryLevelBehaviorWhenBaseRuleHasNoBehavior() {
	// the conformance reporting mode is RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG, i.e.
	// as if the CheckJS action is being run.
	c.reportingMode = CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG
	c.AllowSourcelessWarnings()
	c.baseConfiguration =
		jutil.StringPtr(`requirement: {
  rule_id: "gws:goog.Promise.X"
  type: BANNED_NAME
  error_message: "Prefer using native Promise equivalents. See go/gws-js-conformance#goog-promise"

  value: "goog.Promise.all"
}
`)
	c.extendingConfiguration =
		jutil.StringPtr(`library_level_non_allowlisted_conformance_violations_behavior: REPORT_AS_BUILD_ERROR
requirement: {
  extends: 'gws:goog.Promise.X'
}
`)
	c.TestWarning(
		`goog.Promise.all();
`,
		CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testSimpleExtends_appliesLibraryLevelBehaviorThatIsSpecified()
func (c *CheckConformanceTest) TestSimpleExtends_appliesLibraryLevelBehaviorThatIsSpecified() {
	c.reportingMode = CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG
	c.AllowSourcelessWarnings()
	c.baseConfiguration =
		jutil.StringPtr(`requirement: {
  rule_id: "gws:goog.Promise.X"
  type: BANNED_NAME
  error_message: "Prefer using native Promise equivalents. See go/gws-js-conformance#goog-promise"

  value: "goog.Promise.all"
}
`)
	c.extendingConfiguration =
		jutil.StringPtr(`library_level_non_allowlisted_conformance_violations_behavior: RECORD_ONLY
requirement: {
  extends: 'gws:goog.Promise.X'
}
`)
	// respects the library-level behavior in the extending config and does not warn
	c.TestNoWarning(
		`goog.Promise.all();
`)
}

// Java: CheckConformanceTest.testSimpleExtends_appliesLibraryLevelBehaviorThatIsSpecified_recordOnlyInBaseConfig()
func (c *CheckConformanceTest) TestSimpleExtends_appliesLibraryLevelBehaviorThatIsSpecified_recordOnlyInBaseConfig() {
	c.reportingMode = CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG
	c.AllowSourcelessWarnings()
	c.baseConfiguration =
		jutil.StringPtr(`library_level_non_allowlisted_conformance_violations_behavior: RECORD_ONLY
requirement: {
  rule_id: "gws:goog.Promise.X"
  type: BANNED_NAME
  error_message: "Prefer using native Promise equivalents. See go/gws-js-conformance#goog-promise"

  value: "goog.Promise.all"
}
`)
	c.extendingConfiguration =
		jutil.StringPtr(`requirement: {
  extends: 'gws:goog.Promise.X'
}
`)
	// respects the library-level behavior in the base config and does not warn
	c.TestNoWarning(
		`goog.Promise.all();
`)
}

// Java: CheckConformanceTest.testLibraryLevelConformance_requiresRuleIdOrExtends_failsWhenBothMissing()
func (c *CheckConformanceTest) TestLibraryLevelConformance_requiresRuleIdOrExtends_failsWhenBothMissing() {
	c.reportingMode = CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG
	c.AllowSourcelessWarnings()
	c.configuration =
		`library_level_non_allowlisted_conformance_violations_behavior: RECORD_ONLY
requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
}
`
	c.TestError(
		"eval();",
		CheckConformance_INVALID_REQUIREMENT_SPEC,
		"Invalid requirement. Reason: Library-level conformance requirements require their"+
			" rule_id or extends field set.\n"+
			"Requirement spec:\n"+
			"error_message: \"eval is not allowed\"\n"+
			"type: BANNED_NAME\n"+
			"value: \"eval\"\n")
}

// Java: CheckConformanceTest.testLibraryLevelConformance_requiresRuleIdOrExtends_passesWhenRuleIdPresent()
func (c *CheckConformanceTest) TestLibraryLevelConformance_requiresRuleIdOrExtends_passesWhenRuleIdPresent() {
	c.reportingMode = CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG
	c.AllowSourcelessWarnings()
	c.configuration =
		`library_level_non_allowlisted_conformance_violations_behavior: RECORD_ONLY
requirement: {
  rule_id: 'eval-rule'
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
}
`
	c.TestNoWarning("eval();")
}

// Java: CheckConformanceTest.testLibraryLevelConformance_requiresRuleIdOrExtends_passesWhenExtendsPresent()
func (c *CheckConformanceTest) TestLibraryLevelConformance_requiresRuleIdOrExtends_passesWhenExtendsPresent() {
	c.reportingMode = CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG
	c.AllowSourcelessWarnings()
	c.configuration = ""
	c.baseConfiguration =
		jutil.StringPtr(`requirement: {
  rule_id: 'base-rule'
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
}
`)
	c.extendingConfiguration =
		jutil.StringPtr(`library_level_non_allowlisted_conformance_violations_behavior: RECORD_ONLY
requirement: {
  extends: 'base-rule'
}
`)
	c.TestNoWarning("eval();")
}

// Java: CheckConformanceTest.testLibraryLevelConformance_requiresRuleIdOrExtends_passesWhenBehaviorIsNotRecordOnly()
func (c *CheckConformanceTest) TestLibraryLevelConformance_requiresRuleIdOrExtends_passesWhenBehaviorIsNotRecordOnly() {
	c.reportingMode = CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG
	c.AllowSourcelessWarnings()
	c.configuration =
		`library_level_non_allowlisted_conformance_violations_behavior: REPORT_AS_BUILD_ERROR
requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
}
`
	c.TestWarning("eval();", CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testLibraryLevelConformance_requiresRuleIdOrExtends_passesWhenNotLibraryLevelMode()
func (c *CheckConformanceTest) TestLibraryLevelConformance_requiresRuleIdOrExtends_passesWhenNotLibraryLevelMode() {
	c.reportingMode = CompilerOptions_ConformanceReportingMode_IGNORE_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG
	c.AllowSourcelessWarnings()
	c.configuration =
		`library_level_non_allowlisted_conformance_violations_behavior: RECORD_ONLY
requirement: {
  type: BANNED_NAME
  value: 'eval'
  error_message: 'eval is not allowed'
}
`
	c.TestWarning("eval();", CheckConformance_CONFORMANCE_VIOLATION)
}

// Java: CheckConformanceTest.testCustomBanUnknownProp_mergeConfigWithValue()
func (c *CheckConformanceTest) TestCustomBanUnknownProp_mergeConfigWithValue() {
	c.configuration =
		c.Config(
			c.Rule("BanUnknownTypedClassPropsReferences"),
			"My message",
			"  rule_id: 'x'",
			"  allow_extending_value: true") +
			`requirement: {
  extends: 'x'
  value: 'F'
}

`
	c.TestNoWarning(
		`/** @typedef {?} */ var Unk;
/** @constructor */ function F() { /** @type {?Unk} */ this.prop = null; };
F.prototype.method = function() { alert(this.prop); }
`)
}

// Java: CheckConformanceTest.testCustomBanUnknownInterfaceProp1()
func (c *CheckConformanceTest) TestCustomBanUnknownInterfaceProp1() {
	c.configuration =
		c.Config(c.Rule("BanUnknownTypedClassPropsReferences"), "My rule message", c.Value("String"))
	js := `/** @interface */ function I() {}
I.prototype.method = function() {};
I.prototype.gak;
/** @param {!I} a */
function f(a) {
  a.gak();
}
`
	c.TestWarning(
		js,
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: My rule message\nThe property \"gak\" on type \"I\"")
}

// Java: CheckConformanceTest.testCustomBanUnknownInterfaceProp2()
func (c *CheckConformanceTest) TestCustomBanUnknownInterfaceProp2() {
	c.configuration =
		c.Config(c.Rule("BanUnknownTypedClassPropsReferences"), "My rule message", c.Value("String"))
	c.TestNoWarning(
		`/** @interface */ function I() {}
I.prototype.method = function() {};
/** @param {I} a */ function f(a) {
  a.method();
}
`)
}

// Java: CheckConformanceTest.testCustomBanGlobalVars1()
func (c *CheckConformanceTest) TestCustomBanGlobalVars1() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanGlobalVars'
  error_message: 'BanGlobalVars Message'
}
`
	c.TestWarning(
		"var x;", CheckConformance_CONFORMANCE_VIOLATION, "Violation: BanGlobalVars Message")
	c.TestWarning(
		"function fn() {}",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanGlobalVars Message")
	c.TestNoWarning("goog.provide('x');")
	c.TestNoWarning("/** @externs */ var x;")
	// TODO(johnlenz): This might be overly conservative but doing otherwise is more complicated
	// so let see if we can get away with this.
	c.TestWarning(
		"goog.provide('x'); var x;",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanGlobalVars Message")
}

// Java: CheckConformanceTest.testCustomBanGlobalVars2()
func (c *CheckConformanceTest) TestCustomBanGlobalVars2() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanGlobalVars'
  error_message: 'BanGlobalVars Message'
}
`
	c.TestNoWarning(
		`goog.scope(function() {
  var x = {y: 'y'}
  var z = {
     [x.y]: 2
  }
});
`)
	// Test with let and const
	c.TestNoWarning(
		`goog.scope(function() {
  let x = {y: 'y'}
  let z = {
     [x.y]: 2
  }
});
`)
	c.TestNoWarning(
		`goog.scope(function() {
  const x = {y: 'y'}
  const z = {
     [x.y]: 2
  }
});
`)
}

// Java: CheckConformanceTest.testCustomBanGlobalVarsWithDestructuring()
func (c *CheckConformanceTest) TestCustomBanGlobalVarsWithDestructuring() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanGlobalVars'
  error_message: 'BanGlobalVars Message'
}
`
	c.TestWarning(
		"var [x] = [];",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanGlobalVars Message")
	c.TestNoWarning("/** @externs */ var [x] = [];")
}

// Java: CheckConformanceTest.testCustomBanGlobalVarsWithAllowlist()
func (c *CheckConformanceTest) TestCustomBanGlobalVarsWithAllowlist() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanGlobalVars'
  value: 'foo'
  value: 'bar'
  error_message: 'BanGlobalVars Message'
}
`
	c.TestWarning(
		"var baz;", CheckConformance_CONFORMANCE_VIOLATION, "Violation: BanGlobalVars Message")
	c.TestNoWarning("var foo; var bar;")
}

// Java: CheckConformanceTest.testBanGlobalVarsInEs6Module()
func (c *CheckConformanceTest) TestBanGlobalVarsInEs6Module() {
	// ES6 modules cannot be type checked yet
	c.DisableTypeCheck()
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanGlobalVars'
  error_message: 'BanGlobalVars Message'
}
`
	c.TestNoWarning("export function foo() {}")
	c.TestNoWarning("var s; export {x}")
	c.TestNoWarning("export var s;")
}

// Java: CheckConformanceTest.testCustomBanUnresolvedType()
func (c *CheckConformanceTest) TestCustomBanUnresolvedType() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanUnresolvedType'
  error_message: 'BanUnresolvedType Message'
}
`
	c.TestWarning(
		"goog.forwardDeclare('Foo'); /** @param {Foo} a */ function f(a) {a.foo()};",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanUnresolvedType Message\nReference to type 'Foo' never resolved.")
	c.TestWarning(
		"goog.forwardDeclare('Foo'); /** @type {Foo} */ var f = makeFoo(); f.foo();",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanUnresolvedType Message\nReference to type 'Foo' never resolved.")
	c.Test(
		c.Srcs(
			`goog.forwardDeclare('Foo');
goog.forwardDeclare('Bar');
/** @param {?Foo|Bar} foobar */
function f(foobar) {
  return foobar.m();
}
`),
		c.Warning(CheckConformance_CONFORMANCE_VIOLATION).
			WithMessage(
				`Violation: BanUnresolvedType Message
Reference to type 'Foo|Bar' never resolved.
`))
	c.TestNoWarning(
		`/**
 *  @param {!Object<string, ?>} data
 */
function foo(data) {
  data['bar'].baz();
}
`)
}

// Java: CheckConformanceTest.testCustomStrictBanUnresolvedType()
func (c *CheckConformanceTest) TestCustomStrictBanUnresolvedType() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$StrictBanUnresolvedType'
  error_message: 'StrictBanUnresolvedType Message'
}
`
	c.TestWarning(
		"goog.forwardDeclare('Foo'); /** @param {Foo} a */ var f = function(a) {}",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: StrictBanUnresolvedType Message\nReference to type 'Foo' never resolved.")
	c.TestWarning(
		c.Srcs("goog.forwardDeclare('Foo'); /** @param {!Foo} a */ var f;", "f(5);"),
		TypeValidator_TYPE_MISMATCH_WARNING)
	c.TestWarning(
		c.Srcs(
			"goog.forwardDeclare('Foo'); /** @return {!Foo} */ var f;", //
			"f();"),
		CheckConformance_CONFORMANCE_VIOLATION)
	c.Test(
		c.Srcs(
			`goog.forwardDeclare('Foo');
goog.forwardDeclare('Bar');
/**
 * @param {?Foo} foo
 * @param {?Bar} bar
 */
function f(foo, bar) {
  return foo || bar;
}
`),
		c.Warning(CheckConformance_CONFORMANCE_VIOLATION).
			WithMessage(
				`Violation: StrictBanUnresolvedType Message
Reference to type 'Foo' never resolved.
`),
		c.Warning(CheckConformance_CONFORMANCE_VIOLATION).
			WithMessage(
				`Violation: StrictBanUnresolvedType Message
Reference to type 'Bar' never resolved.
`),
		c.Warning(CheckConformance_CONFORMANCE_VIOLATION).
			WithMessage(
				`Violation: StrictBanUnresolvedType Message
Reference to type 'Foo' never resolved.
`),
		c.Warning(CheckConformance_CONFORMANCE_VIOLATION).
			WithMessage(
				`Violation: StrictBanUnresolvedType Message
Reference to type 'Bar|Foo' never resolved.
`),
		c.Warning(CheckConformance_CONFORMANCE_VIOLATION).
			WithMessage(
				`Violation: StrictBanUnresolvedType Message
Reference to type 'Bar' never resolved.
`))
}

// Java: CheckConformanceTest.testMergeRequirementsWhitelist()
func (c *CheckConformanceTest) TestMergeRequirementsWhitelist() {
	compiler := c.CreateCompiler()
	builder := ConformanceConfig_NewBuilder()
	builder.
		AddRequirementBuilder().
		SetRuleId("a").
		AddWhitelist("x").
		AddWhitelistRegexp("m").
		AddWhitelistEntry(
			RequirementScopeEntry_NewBuilder().AddPrefix("x2").AddRegexp("m2").Build())
	builder.
		AddRequirementBuilder().
		SetExtends("a").
		AddWhitelist("y").
		AddWhitelistRegexp("n").
		AddWhitelistEntry(
			RequirementScopeEntry_NewBuilder().AddPrefix("a2").AddRegexp("y2").Build())
	checkConformance := NewCheckConformance(compiler)
	requirements := checkConformance.MergeRequirements(
		compiler,
		[]*ConformanceConfig{builder.Build()},
		/* isLibraryLevelConformanceReportingMode= */ false)
	AssertThat(requirements).HasSize(1)
	requirement := requirements[0]
	AssertThat(len(requirement.GetWhitelistList())).IsEqualTo(2)
	AssertThat(len(requirement.GetWhitelistRegexpList())).IsEqualTo(2)
	AssertThat(len(requirement.GetWhitelistEntryList())).IsEqualTo(2)
}

// Java: CheckConformanceTest.testMergeRequirementsAllowlist()
func (c *CheckConformanceTest) TestMergeRequirementsAllowlist() {
	compiler := c.CreateCompiler()
	builder := ConformanceConfig_NewBuilder()
	builder.
		AddRequirementBuilder().
		SetRuleId("a").
		AddAllowlist("x").
		AddAllowlistRegexp("m").
		AddAllowlistEntry(
			RequirementScopeEntry_NewBuilder().AddPrefix("x2").AddRegexp("m2").Build())
	builder.
		AddRequirementBuilder().
		SetExtends("a").
		AddAllowlist("y").
		AddAllowlistRegexp("n").
		AddAllowlistEntry(
			RequirementScopeEntry_NewBuilder().AddPrefix("a2").AddRegexp("y2").Build())
	checkConformance := NewCheckConformance(compiler)
	requirements := checkConformance.MergeRequirements(
		compiler,
		[]*ConformanceConfig{builder.Build()},
		/* isLibraryLevelConformanceReportingMode= */ false)
	AssertThat(requirements).HasSize(1)
	requirement := requirements[0]
	AssertThat(len(requirement.GetAllowlistList())).IsEqualTo(2)
	AssertThat(len(requirement.GetAllowlistRegexpList())).IsEqualTo(2)
	AssertThat(len(requirement.GetAllowlistEntryList())).IsEqualTo(2)
}

// Java: CheckConformanceTest.testMergeRequirementsWhitelist_findsDuplicates()
func (c *CheckConformanceTest) TestMergeRequirementsWhitelist_findsDuplicates() {
	compiler := c.CreateCompiler()
	errorManager := NewBlackHoleErrorManager()
	compiler.SetErrorManager(errorManager)
	builder := ConformanceConfig_NewBuilder()
	builder.AddRequirementBuilder().AddWhitelist("x").AddWhitelist("x")
	config := builder.Build()
	checkConformance := NewCheckConformance(compiler)
	requirements := checkConformance.MergeRequirements(
		compiler,
		[]*ConformanceConfig{config},
		/* isLibraryLevelConformanceReportingMode= */ false)
	AssertThat(len(requirements[0].GetWhitelistList())).IsEqualTo(1)
	AssertThat(errorManager.GetErrorCount()).IsEqualTo(0)
}

// Java: CheckConformanceTest.testMergeRequirementsAllowlist_findsDuplicates()
func (c *CheckConformanceTest) TestMergeRequirementsAllowlist_findsDuplicates() {
	compiler := c.CreateCompiler()
	errorManager := NewBlackHoleErrorManager()
	compiler.SetErrorManager(errorManager)
	builder := ConformanceConfig_NewBuilder()
	builder.AddRequirementBuilder().AddAllowlist("x").AddAllowlist("x")
	config := builder.Build()
	checkConformance := NewCheckConformance(compiler)
	requirements := checkConformance.MergeRequirements(
		compiler,
		[]*ConformanceConfig{config},
		/* isLibraryLevelConformanceReportingMode= */ false)
	AssertThat(len(requirements[0].GetAllowlistList())).IsEqualTo(1)
	AssertThat(errorManager.GetErrorCount()).IsEqualTo(0)
}

// Java: CheckConformanceTest.testCustomBanNullDeref1()
func (c *CheckConformanceTest) TestCustomBanNullDeref1() {
	c.configuration = c.Config(c.Rule("BanNullDeref"), "My rule message")
	externs := checkConformanceTest_EXTERNS + "String.prototype.prop;"
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("/** @param {string|null} n */ function f(n) { alert(n.prop); }"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: My rule message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("/** @param {string|null} n */ function f(n) { alert(n['prop']); }"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: My rule message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs(
			`/** @param {string|null} n */
function f(n) { alert('prop' in n); }
`),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: My rule message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("/** @param {string|undefined} n */ function f(n) { alert(n.prop); }"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: My rule message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("/** @param {?Function} fnOrNull */ function f(fnOrNull) { fnOrNull(); }"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: My rule message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("/** @param {?Function} fnOrNull */ function f(fnOrNull) { new fnOrNull(); }"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: My rule message")
	c.TestNoWarning(
		c.Externs(externs), c.Srcs("/** @param {string} n */ function f(n) { alert(n.prop); }"))
	c.TestNoWarning(c.Externs(externs), c.Srcs("/** @param {?} n */ function f(n) { alert(n.prop); }"))
}

// Java: CheckConformanceTest.testCustomBanNullDeref2()
func (c *CheckConformanceTest) TestCustomBanNullDeref2() {
	c.configuration = c.Config(c.Rule("BanNullDeref"), "My rule message")
	externs := checkConformanceTest_EXTERNS + "String.prototype.prop;"
	code := "/** @param {?String} n */ function f(n) { alert(n.prop); }"
	c.TestWarning(
		c.Externs(externs),
		c.Srcs(code),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: My rule message")
	c.configuration = c.Config(c.Rule("BanNullDeref"), "My rule message", c.Value("String"))
	c.TestNoWarning(c.Externs(externs), c.Srcs(code))
}

// Java: CheckConformanceTest.testCustomBanNullDeref3()
func (c *CheckConformanceTest) TestCustomBanNullDeref3() {
	c.configuration = c.Config(c.Rule("BanNullDeref"), "My rule message")
	typedefExterns := checkConformanceTest_EXTERNS +
		`/** @fileoverview */
/** @const */ var ns = {};
/** @enum {number} */ ns.Type.State = {OPEN: 0};
/** @typedef {{a:string}} */ ns.Type;
`
	code := `/** @return {void} n */
function f() { alert(ns.Type.State.OPEN); }
`
	c.TestNoWarning(c.Externs(typedefExterns), c.Srcs(code))
}

// Java: CheckConformanceTest.testCustomBanNullDeref4()
func (c *CheckConformanceTest) TestCustomBanNullDeref4() {
	c.configuration = c.Config(c.Rule("BanNullDeref"), "My rule message")
	c.TestNoWarning(
		`/** @param {*} x */
function f(x) {
  return x.toString();
}
`)
}

// Java: CheckConformanceTest.testRequireUseStrict0()
func (c *CheckConformanceTest) TestRequireUseStrict0() {
	c.configuration = c.Config(c.Rule("RequireUseStrict"), "My rule message")
	c.TestWarning("anything;", CheckConformance_CONFORMANCE_VIOLATION, "Violation: My rule message")
}

// Java: CheckConformanceTest.testRequireUseStrictScript()
func (c *CheckConformanceTest) TestRequireUseStrictScript() {
	c.configuration = c.Config(c.Rule("RequireUseStrict"), "My rule message")
	c.TestNoWarning("'use strict';")
}

// Java: CheckConformanceTest.testRequireUseStrictGoogModule()
func (c *CheckConformanceTest) TestRequireUseStrictGoogModule() {
	c.configuration = c.Config(c.Rule("RequireUseStrict"), "My rule message")
	c.TestNoWarning("goog.module('foo');")
}

// Java: CheckConformanceTest.testRequireUseStrictEs6Module()
func (c *CheckConformanceTest) TestRequireUseStrictEs6Module() {
	c.configuration = c.Config(c.Rule("RequireUseStrict"), "My rule message")
	c.EnableRewriteEsModules()
	c.TestNoWarning("export var x = 2;")
}

// Java: CheckConformanceTest.testBanCreateElement()
func (c *CheckConformanceTest) TestBanCreateElement() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanCreateElement'
  error_message: 'BanCreateElement Message'
  value: 'script'
}
`
	c.TestWarning(
		"goog.dom.createElement('script');",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanCreateElement Message")
	c.TestWarning(
		"goog.dom.createDom('script', {});",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanCreateElement Message")
	externs := CompilerTestCase_DEFAULT_EXTERNS +
		`/** @constructor */ function Document() {}
/** @const {!Document} */ var document;
/** @const */ goog.dom = {};
/** @constructor */ goog.dom.DomHelper = function() {};
`
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("document.createElement('script');"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanCreateElement Message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("new goog.dom.DomHelper().createElement('script');"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanCreateElement Message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("function f(/** ?Document */ doc) { doc.createElement('script'); }"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanCreateElement Message")
	c.TestNoWarning("goog.dom.createElement('iframe');")
	c.TestNoWarning("goog.dom.createElement(goog.dom.TagName.SCRIPT);")
}

// Java: CheckConformanceTest.testBanCreateDom()
func (c *CheckConformanceTest) TestBanCreateDom() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanCreateDom'
  error_message: 'BanCreateDom Message'
  value: 'iframe.src'
  value: 'div.class'
}
`
	c.TestWarning(
		"goog.dom.createDom('iframe', {'src': src});",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanCreateDom Message")
	c.TestWarning(
		"goog.dom.createDom('iframe', {'src': src, 'name': ''}, '');",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanCreateDom Message")
	c.TestWarning(
		"goog.dom.createDom(goog.dom.TagName.IFRAME, {'src': src});",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanCreateDom Message")
	c.TestWarning(
		"goog.dom.createDom('div', 'red');",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanCreateDom Message")
	c.TestWarning(
		"goog.dom.createDom('div', ['red']);",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanCreateDom Message")
	c.TestWarning(
		"goog.dom.createDom('DIV', ['red']);",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanCreateDom Message")
	externs := CompilerTestCase_DEFAULT_EXTERNS +
		`/** @const */ goog.dom = {};
/** @constructor */ goog.dom.DomHelper = function() {};
`
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("new goog.dom.DomHelper().createDom('iframe', {'src': src});"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanCreateDom Message")
	c.TestWarning(
		"goog.dom.createDom(tag, {'src': src});",
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION,
		"Possible violation: BanCreateDom Message")
	c.TestWarning(
		"goog.dom.createDom('iframe', attrs);",
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION,
		"Possible violation: BanCreateDom Message")
	c.TestWarning(
		"goog.dom.createDom(tag, attrs);",
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION,
		"Possible violation: BanCreateDom Message")
	c.TestWarning(
		"goog.dom.createDom('iframe', x ? {'src': src} : 'class');",
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION,
		"Possible violation: BanCreateDom Message")
	// Give a possible violation message if there are computed properties.
	c.TestWarning(
		"goog.dom.createDom('iframe', {['not_src']: src});",
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION,
		`Possible violation: BanCreateDom Message
The type information available for this expression is too loose to ensure conformance.
`)
	c.TestNoWarning("goog.dom.createDom('iframe');")
	c.TestNoWarning("goog.dom.createDom('iframe', {'src': ''});")
	c.TestNoWarning("goog.dom.createDom('iframe', {'name': name});")
	c.TestNoWarning("goog.dom.createDom('iframe', 'red' + '');")
	c.TestNoWarning("goog.dom.createDom('iframe', ['red']);")
	c.TestNoWarning("goog.dom.createDom('iframe', undefined);")
	c.TestNoWarning("goog.dom.createDom('iframe', null);")
	c.TestNoWarning("goog.dom.createDom('img', {'src': src});")
	c.TestNoWarning("goog.dom.createDom('img', attrs);")
	c.TestNoWarning(
		`goog.dom.createDom(
'iframe', /** @type {?string|!Array|undefined} */ (className));
`)
	c.TestNoWarning("goog.dom.createDom(tag, {});")
	c.TestNoWarning(
		`/** @enum {string} */ var Classes = {A: ''};
goog.dom.createDom('iframe', Classes.A);
`)
}

// Java: CheckConformanceTest.testBanCreateDomIgnoreLooseType()
func (c *CheckConformanceTest) TestBanCreateDomIgnoreLooseType() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanCreateDom'
  error_message: 'BanCreateDom Message'
  report_loose_type_violations: false
  value: 'iframe.src'
}
`
	c.TestWarning(
		"goog.dom.createDom('iframe', {'src': src});",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanCreateDom Message")
	c.TestNoWarning("goog.dom.createDom('iframe', attrs);")
	c.TestNoWarning("goog.dom.createDom(tag, {'src': src});")
}

// Java: CheckConformanceTest.testBanCreateDomTagNameType()
func (c *CheckConformanceTest) TestBanCreateDomTagNameType() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanCreateDom'
  error_message: 'BanCreateDom Message'
  value: 'div.class'
}
`
	externs := CompilerTestCase_DEFAULT_EXTERNS +
		`/** @const */ goog.dom = {};
/** @constructor @template T */ goog.dom.TagName = function() {};
/** @type {!goog.dom.TagName<!HTMLDivElement>} */
goog.dom.TagName.DIV = new goog.dom.TagName();
/** @constructor */ function HTMLDivElement() {}

`
	c.TestWarning(
		c.Externs(externs),
		c.Srcs(
			`function f(/** !goog.dom.TagName<!HTMLDivElement> */ div) {
  goog.dom.createDom(div, 'red');
}
`),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanCreateDom Message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs(
			`const TagName = goog.dom.TagName;
goog.dom.createDom(TagName.DIV, 'red');
`),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanCreateDom Message")
}

// Java: CheckConformanceTest.testBanCreateDomMultiType()
func (c *CheckConformanceTest) TestBanCreateDomMultiType() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanCreateDom'
  error_message: 'BanCreateDom Message'
  value: 'h2.class'
}
`
	externs := CompilerTestCase_DEFAULT_EXTERNS +
		`/** @const */ goog.dom = {};
/** @constructor @template T */ goog.dom.TagName = function() {}
/** @constructor */ function HTMLHeadingElement() {}

`
	c.TestWarning(
		c.Externs(externs),
		c.Srcs(
			`function f(/** !goog.dom.TagName<!HTMLHeadingElement> */ heading) {
  goog.dom.createDom(heading, 'red');
}
`),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanCreateDom Message")
}

// Java: CheckConformanceTest.testBanCreateDomAnyTagName()
func (c *CheckConformanceTest) TestBanCreateDomAnyTagName() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanCreateDom'
  error_message: 'BanCreateDom Message'
  value: '*.innerHTML'
}
`
	c.TestWarning(
		"goog.dom.createDom('span', {'innerHTML': html});",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanCreateDom Message")
	c.TestWarning(
		"goog.dom.createDom(tag, {'innerHTML': html});",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanCreateDom Message")
	c.TestWarning(
		"goog.dom.createDom('span', attrs);",
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION,
		"Possible violation: BanCreateDom Message")
	c.TestNoWarning("goog.dom.createDom('span', {'innerHTML': ''});")
	c.TestNoWarning("goog.dom.createDom('span', {'innerhtml': html});")
}

// Java: CheckConformanceTest.testBanCreateDomTextContent()
func (c *CheckConformanceTest) TestBanCreateDomTextContent() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanCreateDom'
  error_message: 'BanCreateDom Message'
  value: 'script.textContent'
}
`
	c.TestWarning(
		"goog.dom.createDom('script', {}, source);",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanCreateDom Message")
	c.TestWarning(
		"goog.dom.createDom('script', {'textContent': source});",
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanCreateDom Message")
	c.TestWarning(
		"goog.dom.createDom(script, {}, source);",
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION,
		"Possible violation: BanCreateDom Message")
	c.TestNoWarning("goog.dom.createDom('script', {});")
	c.TestNoWarning("goog.dom.createDom('span', {'textContent': text});")
	c.TestNoWarning("goog.dom.createDom('span', {}, text);")
}

// Java: CheckConformanceTest.testBanElementSetAttribute()
func (c *CheckConformanceTest) TestBanElementSetAttribute() {
	c.configuration =
		`requirement: {

  type: CUSTOM

  value: 'src'

  value: 'HREF'

  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanElementSetAttribute'

  report_loose_type_violations: false

  error_message: 'BanSetAttribute Message'

}
`
	externs := CompilerTestCase_DEFAULT_EXTERNS +
		`/** @constructor */ function Element() {}
/** @constructor @extends {Element} */ function HTMLScriptElement() {}

`
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("(new HTMLScriptElement).setAttribute('SRc', 'xxx')"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanSetAttribute Message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("(new HTMLScriptElement).setAttribute('href', 'xxx')"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanSetAttribute Message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("(new HTMLScriptElement).setAttributeNS(null, 'SRc', 'xxx')"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanSetAttribute Message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("(new HTMLScriptElement)['SRc'] = 'xxx';"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanSetAttribute Message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("(new HTMLScriptElement)['href'] = 'xxx';"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanSetAttribute Message")
	c.TestNoWarning(
		c.Externs(externs),
		c.Srcs("(new HTMLScriptElement).setAttributeNS('ns', 'data-random', 'xxx')"))
	c.TestNoWarning(
		c.Externs(externs), c.Srcs("(new HTMLScriptElement).setAttribute('data-random', 'xxx')"))
	c.TestNoWarning(
		c.Externs(externs),
		c.Srcs("(new HTMLScriptElement).setAttributeNS(null, 'data-random', 'xxx')"))
	c.TestNoWarning(c.Externs(externs), c.Srcs("(new HTMLScriptElement)['innerHTML'] = 'xxx';"))
	c.TestNoWarning(
		c.Externs(externs), c.Srcs("var attr = 'unknown'; (new HTMLScriptElement)[attr] =  'xxx';"))
	c.TestNoWarning(
		c.Externs(externs),
		c.Srcs("var attr = 'unknown'; (new HTMLScriptElement).setAttributeNS(null, attr, 'xxx')"))
	c.TestNoWarning(
		c.Externs(externs),
		c.Srcs(
			`/** @param {string|null} attr */
function foo(attr) { (new HTMLScriptElement)[attr] =  'xxx'; }
`))
	c.TestNoWarning(
		c.Externs(externs),
		c.Srcs(
			`const foo = 'safe';
var bar = foo;
(new HTMLScriptElement).setAttribute(bar, 'xxx');
`))
	c.TestNoWarning(c.Externs(externs), c.Srcs("(new HTMLScriptElement)['data-random'] = 'xxx';"))
	c.TestNoWarning(
		c.Externs(externs),
		c.Srcs(
			`const foo = 'safe';
const bar = foo;
(new HTMLScriptElement).setAttribute(bar, 'xxx');
`))
}

// Java: CheckConformanceTest.testBanElementSetAttributeLoose()
func (c *CheckConformanceTest) TestBanElementSetAttributeLoose() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  value: 'src'
  value: 'HREF'
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanElementSetAttribute'
  error_message: 'BanSetAttribute Message'
  report_loose_type_violations: true
}
`
	externs := CompilerTestCase_DEFAULT_EXTERNS +
		`/** @constructor */ function Element() {}
/** @constructor @extends {Element} */ function HTMLScriptElement() {}
`
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("(new HTMLScriptElement).setAttribute('SRc', 'xxx')"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanSetAttribute Message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("(new HTMLScriptElement).setAttribute('href', 'xxx')"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanSetAttribute Message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("(new HTMLScriptElement).setAttributeNS(null, 'SRc', 'xxx')"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanSetAttribute Message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("(new HTMLScriptElement).setAttributeNS('ns', 'data-random', 'xxx')"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanSetAttribute Message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("var attr = 'unknown'; (new HTMLScriptElement).setAttributeNS(null, attr, 'xxx')"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanSetAttribute Message")
	c.TestNoWarning(
		c.Externs(externs), c.Srcs("(new HTMLScriptElement).setAttribute('data-random', 'xxx')"))
	c.TestNoWarning(
		c.Externs(externs),
		c.Srcs("(new HTMLScriptElement).setAttributeNS(null, 'data-random', 'xxx')"))
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("(new HTMLScriptElement)['SRc'] = 'xxx';"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanSetAttribute Message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("(new HTMLScriptElement)['src'] = 'xxx';"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanSetAttribute Message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("(new HTMLScriptElement)['href'] = 'xxx';"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanSetAttribute Message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("(new HTMLScriptElement)['innerHTML'] = 'xxx';"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanSetAttribute Message")
	c.TestNoWarning(
		c.Externs(externs), c.Srcs("const attr = 'unknown'; (new HTMLScriptElement)[attr] =  'xxx';"))
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("var attr = 'unknown'; (new HTMLScriptElement)[attr] =  'xxx';"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanSetAttribute Message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs(
			`/** @param {string|null} attr */
function foo(attr) { (new HTMLScriptElement)[attr] =  'xxx'; }
`),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanSetAttribute Message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs(
			`const foo = 'safe';
var bar = foo;
(new HTMLScriptElement).setAttribute(bar, 'xxx');
`),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanSetAttribute Message")
	c.TestNoWarning(c.Externs(externs), c.Srcs("(new HTMLScriptElement)['data-random'] = 'xxx';"))
	c.TestNoWarning(
		c.Externs(externs),
		c.Srcs(
			`const foo = 'safe';
const bar = foo;
(new HTMLScriptElement).setAttribute(bar, 'xxx');
`))
	c.TestNoWarning(
		c.Externs(externs),
		c.Srcs(
			`/** @const */
var foo = 'safe';
(new HTMLScriptElement).setAttribute(foo, 'xxx');
`))
	c.TestNoWarning(
		c.Externs(externs),
		c.Srcs(
			`goog.provide('test.Attribute');
/** @const */
test.Attribute.foo = 'safe';
(new HTMLScriptElement).setAttribute(test.Attribute.foo, 'xxx');
`))
	c.TestNoWarning(
		c.Externs(externs),
		c.Srcs(
			`goog.provide('test.Attribute');

/** @enum {string} */
test.Attribute = {SRC: 'src', HREF: 'href', SAFE: 'safe'};
`,
			`goog.module('test.setAttribute');

const Attribute = goog.require('test.Attribute');

const attr = Attribute.SAFE;
(new HTMLScriptElement).setAttribute(attr, 'xxx');
`))
	c.TestNoWarning(
		c.Externs(externs),
		c.Srcs(
			`goog.provide('xid');
goog.provide('xid.String');
/** @enum {string} */ xid.String = {DO_NOT_USE: ''};
/**
 * @param {string} id
 * @return {xid.String}
 */
xid = function(id) {return /** @type {xid.String} */ (id);};
const attr = xid('src');
(new HTMLScriptElement).setAttribute(attr, 'xxx');
(new HTMLScriptElement)[attr] = 'xxx';
`))
}

// Java: CheckConformanceTest.testBanSettingAttributes()
func (c *CheckConformanceTest) TestBanSettingAttributes() {
	c.configuration =
		`requirement: {

  type: CUSTOM

  value: 'Element.prototype.attr'

  value: 'Foo.prototype.attrib'

  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanSettingAttributes'

  error_message: 'BanSettingAttributes Message'

}
`
	externs := CompilerTestCase_DEFAULT_EXTERNS +
		`/** @constructor */ function Foo() {}
/** @constructor */ function Bar() {}
/** @constructor */ function Element() {}
/** @constructor @extends {Element} */ function HTMLScriptElement() {}
`
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("(new Foo).attrib('SRc', 'xxx')"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanSettingAttributes Message")
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("(new HTMLScriptElement).attr('href', 'xxx')"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanSettingAttributes Message")
	c.TestNoWarning(c.Externs(externs), c.Srcs("(new Bar).attr('src', 'xxx')"))
	c.TestNoWarning(c.Externs(externs), c.Srcs("(new Foo).attrib('src')"))
	c.TestNoWarning(c.Externs(externs), c.Srcs("(new HTMLScriptElement).attrib('src', 'xxx')"))
	c.TestNoWarning(c.Externs(externs), c.Srcs("(new HTMLScriptElement).attr('data-random', 'xxx')"))
	c.TestWarning(
		c.Externs(externs),
		c.Srcs(
			`const foo = 'safe';
var bar = foo;
(new HTMLScriptElement).attr(bar, 'xxx');
`),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanSettingAttributes Message")
	c.TestNoWarning(
		c.Externs(externs),
		c.Srcs(
			`const foo = 'safe';
const bar = foo;
(new HTMLScriptElement).attr(bar, 'xxx');
`))
	c.TestNoWarning(
		c.Externs(externs),
		c.Srcs(
			`goog.provide('test.Attribute');

/** @enum {string} */
test.Attribute = {SRC: 'src', HREF: 'href', SAFE: 'safe'};
`,
			`goog.module('test.attr');

const Attribute = goog.require('test.Attribute');

const attr = Attribute.SAFE;
(new HTMLScriptElement).attr(attr, 'xxx');
`))
	c.TestNoWarning(
		c.Externs(externs),
		c.Srcs(
			`goog.provide('xid');
goog.provide('xid.String');
/** @enum {string} */ xid.String = {DO_NOT_USE: ''};
/**
 * @param {string} id
 * @return {xid.String}
 */
xid = function(id) {return /** @type {xid.String} */ (id);};
const attr = xid('src');
(new HTMLScriptElement).attr(attr, 'xxx');
`))
}

// Java: CheckConformanceTest.testBanExecCommand()
func (c *CheckConformanceTest) TestBanExecCommand() {
	c.configuration =
		`requirement: {

  type: CUSTOM

  value: 'insertHTML'

  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanExecCommand'

  error_message: 'BanExecCommand message'

}
`
	externs := CompilerTestCase_DEFAULT_EXTERNS + "/** @constructor */ function Document() {}"
	c.TestWarning(
		c.Externs(externs),
		c.Srcs("(new Document).execCommand('insertHtml', false, 'xxx')"),
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanExecCommand message")
	c.TestNoWarning(c.Externs(externs), c.Srcs("(new Document).execCommand('bold')"))
}

// Java: CheckConformanceTest.testBanStaticThis()
func (c *CheckConformanceTest) TestBanStaticThis() {
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanStaticThis'
  error_message: 'BanStaticThis Message'
}
`
	c.TestWarning(
		`class Foo {
  static bar() {
    this;
  }
}
`,
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanStaticThis Message")
	c.TestWarning(
		`class Foo {
  static bar() {
    this.buzz();
  }
  static buzz() {}
}
`,
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanStaticThis Message")
	c.TestWarning(
		`class Foo {
  static bar() {
    let fn = () => this;
  }
}
`,
		CheckConformance_CONFORMANCE_VIOLATION,
		"Violation: BanStaticThis Message")
	c.TestNoWarning("let fn = function() { this.buzz(); };")
	c.TestNoWarning(
		`class Foo {
  bar() {
    this.buzz();
  }
  buzz() {}
}
`)
	c.TestNoWarning(
		`class Foo {
  static buzz() {}
}
Foo.bar = function() {
  this.buzz();
}
`)
	c.TestNoWarning(
		`class Foo {
  static bar() {
    let fn = function() { this; };
  }
}
`)
}

// Java: CheckConformanceTest.testInvalidCustomRuleJavaClass()
func (c *CheckConformanceTest) TestInvalidCustomRuleJavaClass() {
	c.AllowSourcelessWarnings()
	c.configuration =
		`requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$NoSuchClass'
  error_message: 'Should fail'
}
`
	c.TestError(
		"var x;",
		CheckConformance_INVALID_REQUIREMENT_SPEC,
		`Invalid requirement. Reason: Unknown conformance rule: com.google.javascript.jscomp.ConformanceRules$NoSuchClass
Requirement spec:
error_message: "Should fail"
type: CUSTOM
java_class: "com.google.javascript.jscomp.ConformanceRules$NoSuchClass"
`)
}

// Java: CheckConformanceTest.testAllCustomRulesRegistered()
func (c *CheckConformanceTest) TestAllCustomRulesRegistered() {
	c.AllowSourcelessWarnings()
	rules := []string{
		"NoOp",
		"InferredConstCheck",
		"BannedDependency",
		"BannedDependencyRegex",
		"BannedName",
		"BannedProperty",
		"BannedStringRegex",
		"RestrictedNameCall",
		"RestrictedMethodCall",
		"RestrictedPropertyWrite",
		"BannedCodePattern",
		"BanForOf",
		"RequireUseStrict",
		"BanThrowOfNonErrorTypes",
		"BanNullDeref",
		"BanUnknownThis",
		"BanUnknownDirectThisPropsReferences",
		"BanUnknownTypedClassPropsReferences",
		"BanUnresolvedType",
		"StrictBanUnresolvedType",
		"BanGlobalVars",
		"BannedEnhance",
		"BannedModsRegex",
		"BanCreateElement",
		"BanCreateDom",
		"BanElementSetAttribute",
		"BanSettingAttributes",
		"BanExecCommand",
		"BanStaticThis",
	}
	for _, rule := range rules {
		c.configuration =
			"requirement: {\n" +
				"  type: CUSTOM\n" +
				"  java_class: 'com.google.javascript.jscomp.ConformanceRules$" +
				rule +
				"'\n" +
				"  error_message: 'msg'\n" +
				"  value: 'value'\n" +
				"}\n"
		// Go-only: language-adapter — the two catch clauses of one try: a
		// Throwable is caught once and dispatched on its class, so that a
		// failure in the first handler is not caught by the second.
		var thrown jutil.Throwable
		jutil.Try(func() {
			c.TestNoWarning("var x;")
		}, func(e jutil.Throwable) {
			thrown = e
		})
		switch e := thrown.(type) {
		case nil:
		case jutil.IRuntimeException:
			AssertThat(e.GetMessage()).DoesNotContain("Unknown conformance rule")
		case *jutil.AssertionError:
			// Assertions are expected if the rule requires specific values
		default:
			panic(e)
		}
	}
}

func (c *CheckConformanceTest) TestReportSourcelessNodeWithOnlyApplyTo() {
	compiler := c.CreateCompiler()
	compiler.InitOptions(NewCompilerOptions())
	reqWithOnlyApplyTo := Requirement_NewBuilder().
		SetType(Requirement_Type_BANNED_NAME).
		SetErrorMessage("eval is banned").
		AddValue("eval").
		AddOnlyApplyTo("foo/").
		Build()
	ruleWithOnlyApplyTo, err := newConformanceRules_BannedName(compiler, reqWithOnlyApplyTo)
	if err != nil {
		panic(err)
	}
	sourcelessNode := rhino.IR_Name("eval")
	ruleWithOnlyApplyTo.Report(
		sourcelessNode,
		ConformanceRules_ConformanceResult_VIOLATION,
		ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_UNSPECIFIED)
	AssertThat(compiler.GetErrorManager().GetWarningCount()).IsEqualTo(0)

	reqUnscoped := Requirement_NewBuilder().
		SetType(Requirement_Type_BANNED_NAME).
		SetErrorMessage("eval is banned").
		AddValue("eval").
		Build()
	ruleUnscoped, err := newConformanceRules_BannedName(compiler, reqUnscoped)
	if err != nil {
		panic(err)
	}
	ruleUnscoped.Report(
		sourcelessNode,
		ConformanceRules_ConformanceResult_VIOLATION,
		ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_UNSPECIFIED)
	AssertThat(compiler.GetErrorManager().GetWarningCount()).IsEqualTo(1)
}
