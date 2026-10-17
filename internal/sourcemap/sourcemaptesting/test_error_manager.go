// Port of com/google/javascript/jscomp/testing/TestErrorManager.java @ closure-compiler v20260415

/*
 * Copyright 2015 The Closure Compiler Authors.
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

package sourcemaptesting

// Go-only: test-helper — TestErrorManager belongs to
// com.google.javascript.jscomp.testing, whose other classes are ported as
// jscomp test files. Its only user, SourceMapTestCase, must be importable by
// the tests of two Go packages (sourcemap and jscomp), so both live here.

import (
	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// TestErrorManager is an ErrorManager that compares errors with a list of
// expected errors.
//
// Java: public final class TestErrorManager extends BasicErrorManager
type TestErrorManager struct {
	jscomp.BasicErrorManagerBase

	errorIndex   int
	errors       []string
	warningIndex int
	warnings     []string
}

// Java: TestErrorManager()
func NewTestErrorManager() *TestErrorManager {
	m := &TestErrorManager{errors: []string{}, warnings: []string{}}
	m.InitBasicErrorManagerBase(m)
	return m
}

// Java: TestErrorManager.expectErrors(String...)
func (m *TestErrorManager) ExpectErrors(expectedErrors ...string) {
	m.errorIndex = 0
	m.errors = expectedErrors
}

// Java: TestErrorManager.expectWarnings(String...)
func (m *TestErrorManager) ExpectWarnings(expectedWarnings ...string) {
	m.warningIndex = 0
	m.warnings = expectedWarnings
}

// Java: TestErrorManager.report(CheckLevel, JSError)
func (m *TestErrorManager) Report(level *jscomp.CheckLevel, error *jscomp.JSError) {
	m.BasicErrorManagerBase.Report(level, error)

	switch *level {
	case jscomp.CheckLevel_ERROR:
		if m.errorIndex >= len(m.errors) {
			panic(jutil.NewAssertionError("Unexpected error: " + error.String()))
		}
		testErrorManager_assertEquals(m.errors[m.errorIndex], error.Description())
		m.errorIndex++
	case jscomp.CheckLevel_WARNING:
		if m.errorIndex >= len(m.warnings) {
			panic(jutil.NewAssertionError("Unexpected warning: " + error.String()))
		}
		testErrorManager_assertEquals(m.warnings[m.warningIndex], error.Description())
		m.warningIndex++
	case jscomp.CheckLevel_OFF:
		// no-op
	}
}

// testErrorManager_assertEquals is org.junit.Assert.assertEquals(Object,
// Object) on two Strings.
func testErrorManager_assertEquals(expected string, actual string) {
	if expected != actual {
		panic(jutil.NewAssertionError("expected:<" + expected + "> but was:<" + actual + ">"))
	}
}

// Java: TestErrorManager.println(CheckLevel, JSError)
func (m *TestErrorManager) Println(level *jscomp.CheckLevel, error *jscomp.JSError) {
	// no-op
}

// Java: TestErrorManager.printSummary()
func (m *TestErrorManager) PrintSummary() {
	// no-op
}

// Java: TestErrorManager.hasEncounteredAllErrors()
func (m *TestErrorManager) HasEncounteredAllErrors() bool {
	return m.errorIndex == len(m.errors)
}

// Java: TestErrorManager.hasEncounteredAllWarnings()
func (m *TestErrorManager) HasEncounteredAllWarnings() bool {
	return m.warningIndex == len(m.warnings)
}
