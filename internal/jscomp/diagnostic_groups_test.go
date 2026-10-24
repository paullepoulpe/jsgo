// Port of com/google/javascript/jscomp/DiagnosticGroupsTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/DiagnosticGroupsTest.java)

/*
 * Copyright 2019 The Closure Compiler Authors.
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

// DiagnosticGroupsTest holds the unit tests of DiagnosticGroups.
type DiagnosticGroupsTest struct {
	t *testing.T
}

func (c *DiagnosticGroupsTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *DiagnosticGroupsTest) JUnitBefore()           {}
func (c *DiagnosticGroupsTest) JUnitAfter()            {}

func TestDiagnosticGroupsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return &DiagnosticGroupsTest{} })
}

// diagnosticGroupsTest_intersection is Sets.intersection(set1, set2), as a
// list in set1's order.
func diagnosticGroupsTest_intersection(set1, set2 *DiagnosticGroup) []*DiagnosticType {
	var result []*DiagnosticType
	for type_ := range set1.GetTypes().All() {
		if set2.GetTypes().Contains(type_) {
			result = append(result, type_)
		}
	}
	return result
}

// TestLintChecksGroupIsDisjointFromEveryOtherGroup tests that
// DiagnosticGroup.LINT_CHECKS does not contain DiagnosticTypes that also
// belong to other DiagnosticGroups.
//
// Because all DiagnosticGroup.LINT_CHECKS run in a single pass, enabling a
// diagnostic group containing one of the checks causes all of the other
// checks to run, which is surprising.
//
// Java: DiagnosticGroupsTest.lintChecksGroupIsDisjointFromEveryOtherGroup()
func (c *DiagnosticGroupsTest) TestLintChecksGroupIsDisjointFromEveryOtherGroup() {
	lintChecks := DiagnosticGroups_LINT_CHECKS
	for group := range DiagnosticGroups_GetRegisteredGroups().Values() {
		if group == lintChecks ||
			// TODO(lharker): stop ignoring USE_OF_GOOG_PROVIDE after migrating rules_closure
			// code to suppress useOfGoogProvide instead of lintChecks.
			group == DiagnosticGroups_USE_OF_GOOG_PROVIDE ||
			// We need a separate group for "lintVarDeclarations" so code that really cannot be
			// updated to use only `let` and `const` can suppress the warnings about `var` without
			// suppressing all lint checks.
			group == DiagnosticGroups_LINT_VAR_DECLARATIONS {
			continue
		}
		AssertWithMessage(
			"DiagnosticTypes common to DiagnosticGroups %s and %s",
			lintChecks.GetName(), group.GetName()).
			That(diagnosticGroupsTest_intersection(lintChecks, group)).
			IsEmpty()
	}
}

// Java: DiagnosticGroupsTest.conformanceErrorsCannotBeDowngraded()
func (c *DiagnosticGroupsTest) TestConformanceErrorsCannotBeDowngraded() {
	for group := range DiagnosticGroups_GetRegisteredGroups().Values() {
		AssertWithMessage("Group '%s' should not include JSC_CONFORMANCE_ERROR", group.GetName()).
			That(group.GetTypes()).
			DoesNotContain(CheckConformance_CONFORMANCE_ERROR)
	}
}

// Java: DiagnosticGroupsTest.testClosureUnawareCodeWarningsGuard()
func (c *DiagnosticGroupsTest) TestClosureUnawareCodeWarningsGuard() {
	guard := NewClosureUnawareCodeWarningsGuard()
	n1 := rhino.Node_NewString("test")
	src1 := SourceFile_FromCode("foo.js", "")
	src1.MarkAsClosureUnawareCode()
	n1.SetStaticSourceFile(src1)
	n1.SetIsInClosureUnawareSubtree(true)
	suppressibleError := JSError_MakeWithNode(n1, RhinoErrorReporter_JSDOC_MISSING_TYPE_WARNING, "")
	AssertThat(guard.Level(suppressibleError)).IsEqualTo(CheckLevel_OFF.Ptr())

	n2 := rhino.Node_NewString("test")
	nonSuppressibleError :=
		JSError_MakeWithNode(n2, RhinoErrorReporter_JSDOC_MISSING_TYPE_WARNING, "")
	AssertThat(guard.Level(nonSuppressibleError)).IsNull()

	n3 := rhino.Node_NewString("test")
	src3 := SourceFile_FromCode("foo.js", "")
	src3.MarkAsClosureUnawareCode()
	n3.SetStaticSourceFile(src3)
	n3.SetIsInClosureUnawareSubtree(true)
	otherError := JSError_MakeWithNode(n3, DiagnosticType_Error("FOO", "foo"), "")
	AssertThat(guard.Level(otherError)).IsNull()
}

// Java: DiagnosticGroupsTest.testDiagnosticGroupPathSuppressingWarningsGuard()
func (c *DiagnosticGroupsTest) TestDiagnosticGroupPathSuppressingWarningsGuard() {
	guard :=
		NewDiagnosticGroupPathSuppressingWarningsGuard(
			DiagnosticGroups_LINT_CHECKS, "suppress_me")

	var lintCheck *DiagnosticType
	for t := range DiagnosticGroups_LINT_CHECKS.GetTypes().All() {
		lintCheck = t
		break
	}
	errorInPath := JSError_MakeWithSourceLocation(jutil.StringPtr("some/path/suppress_me/file.js"), -1, -1, lintCheck)
	AssertWithMessage(
		"DiagnosticGroupPathSuppressingWarningsGuard should suppress in matching path").
		That(guard.Level(errorInPath)).
		IsEqualTo(CheckLevel_OFF.Ptr())

	errorInOtherPath := JSError_MakeWithSourceLocation(jutil.StringPtr("some/other/file.js"), -1, -1, lintCheck)
	AssertWithMessage(
		"DiagnosticGroupPathSuppressingWarningsGuard should not suppress in other path").
		That(guard.Level(errorInOtherPath)).
		IsNull()

	otherErrorInPath :=
		JSError_MakeWithSourceLocation(jutil.StringPtr("some/path/suppress_me/file.js"), -1, -1, DiagnosticType_Error("FOO", "foo"))
	AssertWithMessage(
		"DiagnosticGroupPathSuppressingWarningsGuard should not suppress unrelated errors").
		That(guard.Level(otherErrorInPath)).
		IsNull()
}
