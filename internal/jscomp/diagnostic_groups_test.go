// Port of com/google/javascript/jscomp/DiagnosticGroupsTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/DiagnosticGroupsTest.java)

/*
 * Copyright 2020 The Closure Compiler Authors.
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

import "testing"

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
