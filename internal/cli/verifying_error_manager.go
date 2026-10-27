// Port of com/google/javascript/jscomp/VerifyingErrorManager.java @ closure-compiler v20261021

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

package cli

import (
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
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
// Its synchronized methods hold the object's monitor, which is the embedded
// ThreadSafeDelegatingErrorManager's.
type VerifyingErrorManager struct {
	*jscomp.ThreadSafeDelegatingErrorManager
	unmatchedExpectations []*jutil.Pattern
	formatter             *jscomp.LightweightMessageFormatter
}

var _ jscomp.ErrorManager = (*VerifyingErrorManager)(nil)

// Java: VerifyingErrorManager(ErrorManager, List<String>)
func NewVerifyingErrorManager(delegate jscomp.ErrorManager, expectedDiagnostics []string) *VerifyingErrorManager {
	m := &VerifyingErrorManager{
		ThreadSafeDelegatingErrorManager: jscomp.NewThreadSafeDelegatingErrorManager(delegate),
	}
	m.formatter = jscomp.LightweightMessageFormatter_WithoutSource()
	m.unmatchedExpectations = make([]*jutil.Pattern, 0, len(expectedDiagnostics))
	for _, exp := range expectedDiagnostics {
		m.unmatchedExpectations = append(m.unmatchedExpectations, jutil.Pattern_Compile(exp))
	}
	return m
}

// Java: VerifyingErrorManager.report(CheckLevel, JSError)
func (m *VerifyingErrorManager) Report(level *jscomp.CheckLevel, error *jscomp.JSError) {
	m.LockMonitor()
	defer m.UnlockMonitor()
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
				jscomp.CheckLevel_ERROR.Ptr(),
				jscomp.JSError_Make(
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
	m.LockMonitor()
	defer m.UnlockMonitor()
	m.ReportMissingExpectations()
	m.ThreadSafeDelegatingErrorManager.GenerateReport()
}

// Java: VerifyingErrorManager.reportMissingExpectations()
func (m *VerifyingErrorManager) ReportMissingExpectations() {
	m.LockMonitor()
	defer m.UnlockMonitor()
	for _, p := range m.unmatchedExpectations {
		m.ThreadSafeDelegatingErrorManager.Report(
			jscomp.CheckLevel_ERROR.Ptr(),
			jscomp.JSError_Make(AbstractCommandLineRunner_EXPECTED_DIAGNOSTIC_NOT_FOUND, p.String()))
	}
	m.unmatchedExpectations = m.unmatchedExpectations[:0]
}

// Java: VerifyingErrorManager.formatErrorForMatching(CheckLevel, JSError)
func (m *VerifyingErrorManager) FormatErrorForMatching(level *jscomp.CheckLevel, error *jscomp.JSError) string {
	if level != nil && *level == jscomp.CheckLevel_ERROR {
		return m.formatter.FormatError(error)
	}
	return m.formatter.FormatWarning(error)
}
