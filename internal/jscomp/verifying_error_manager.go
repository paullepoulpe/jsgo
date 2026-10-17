// Port of com/google/javascript/jscomp/VerifyingErrorManager.java @ closure-compiler v20260415

/*
 * Copyright 2026 The Closure Compiler Authors.
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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// VerifyingErrorManager is an ErrorManager that verifies expected
// diagnostics and reports missing or unexpected ones.
//
// It wraps a delegate ErrorManager. Expected diagnostics (passed as regexes)
// are swallowed if matched. Unexpected diagnostics are passed to the
// delegate. At the end of compilation (when generateReport is called), any
// unmatched expectations are reported as errors to the delegate.
//
// Go-only: language-adapter — the Java methods are synchronized on the same
// monitor as ThreadSafeDelegatingErrorManager's; PORTING.md §5 makes
// synchronized code sequential, and the embedded manager keeps its own lock
// (Go's sync.Mutex is not reentrant, so this type takes no lock itself).
type VerifyingErrorManager struct {
	*ThreadSafeDelegatingErrorManager
	unmatchedExpectations []*jutil.Pattern
	formatter             *LightweightMessageFormatter
}

var _ ErrorManager = (*VerifyingErrorManager)(nil)

// Java: VerifyingErrorManager(ErrorManager, List<String>)
func NewVerifyingErrorManager(delegate ErrorManager, expectedDiagnostics []string) *VerifyingErrorManager {
	m := &VerifyingErrorManager{
		ThreadSafeDelegatingErrorManager: NewThreadSafeDelegatingErrorManager(delegate),
	}
	m.formatter = LightweightMessageFormatter_WithoutSource()
	m.unmatchedExpectations = make([]*jutil.Pattern, 0, len(expectedDiagnostics))
	for _, exp := range expectedDiagnostics {
		m.unmatchedExpectations = append(m.unmatchedExpectations, jutil.Pattern_Compile(exp))
	}
	return m
}

// Java: VerifyingErrorManager.report(CheckLevel, JSError)
func (m *VerifyingErrorManager) Report(level *CheckLevel, error *JSError) {
	formattedError := m.FormatErrorForMatching(level, error)

	var matchedPattern *jutil.Pattern
	// Go-only: language-adapter — the Iterator with remove() is an index
	// loop over the slice; i is not advanced past a removed element.
	for i := 0; i < len(m.unmatchedExpectations); {
		p := m.unmatchedExpectations[i]
		if p.Matcher(formattedError).Find() {
			if matchedPattern == nil {
				matchedPattern = p
				m.unmatchedExpectations = slices.Delete(m.unmatchedExpectations, i, i+1)
				continue
			}

			if matchedPattern.String() == p.String() {
				// Allow duplicates that are identical.
				i++
				continue
			}

			m.ThreadSafeDelegatingErrorManager.Report(
				CheckLevel_ERROR.Ptr(),
				JSError_Make(
					AbstractCommandLineRunner_AMBIGUOUS_EXPECTATION,
					formattedError,
					matchedPattern.String(),
					p.String()))
			m.unmatchedExpectations = slices.Delete(m.unmatchedExpectations, i, i+1)
			continue
		}
		i++
	}

	if matchedPattern == nil {
		m.ThreadSafeDelegatingErrorManager.Report(level, error)
	}
}

// Java: VerifyingErrorManager.generateReport()
func (m *VerifyingErrorManager) GenerateReport() {
	m.ReportMissingExpectations()
	m.ThreadSafeDelegatingErrorManager.GenerateReport()
}

// Java: VerifyingErrorManager.reportMissingExpectations()
func (m *VerifyingErrorManager) ReportMissingExpectations() {
	for _, p := range m.unmatchedExpectations {
		m.ThreadSafeDelegatingErrorManager.Report(
			CheckLevel_ERROR.Ptr(),
			JSError_Make(AbstractCommandLineRunner_EXPECTED_DIAGNOSTIC_NOT_FOUND, p.String()))
	}
	m.unmatchedExpectations = m.unmatchedExpectations[:0]
}

// Java: VerifyingErrorManager.formatErrorForMatching(CheckLevel, JSError)
func (m *VerifyingErrorManager) FormatErrorForMatching(level *CheckLevel, error *JSError) string {
	if level != nil && *level == CheckLevel_ERROR {
		return m.formatter.FormatError(error)
	}
	return m.formatter.FormatWarning(error)
}
