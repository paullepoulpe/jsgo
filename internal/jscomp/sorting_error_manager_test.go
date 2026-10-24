// Port of test/com/google/javascript/jscomp/SortingErrorManagerTest.java @ closure-compiler v20261021

/*
 * Copyright 2007 The Closure Compiler Authors.
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
	"strconv"
	"testing"
)

var (
	jsaSemNULL_SOURCE *string
	jsaSemComparator  = NewSortingErrorManager_LeveledJSErrorComparator()
	jsaSemFOO_TYPE    = DiagnosticType_Error("TEST_FOO", "Foo")
	jsaSemJOO_TYPE    = DiagnosticType_Error("TEST_JOO", "Joo")
)

func TestSortingErrorManager_OrderingBothNull(t *testing.T) {
	jsaEq(t, jsaSemComparator.Compare(nil, nil), 0)
}

func TestSortingErrorManager_OrderingSourceName1(t *testing.T) {
	e1 := JSError_MakeWithSourceLocation(jsaSemNULL_SOURCE, -1, -1, jsaSemFOO_TYPE)
	e2 := JSError_MakeWithSourceLocation(jsaStr("a"), -1, -1, jsaSemFOO_TYPE)

	jsaSemAssertSmaller(t, jsaSemError(e1), jsaSemError(e2))
}

func TestSortingErrorManager_OrderingSourceName2(t *testing.T) {
	e1 := JSError_MakeWithSourceLocation(jsaStr("a"), -1, -1, jsaSemFOO_TYPE)
	e2 := JSError_MakeWithSourceLocation(jsaStr("b"), -1, -1, jsaSemFOO_TYPE)

	jsaSemAssertSmaller(t, jsaSemError(e1), jsaSemError(e2))
}

func TestSortingErrorManager_OrderingLineno1(t *testing.T) {
	e1 := JSError_MakeWithSourceLocation(jsaSemNULL_SOURCE, -1, -1, jsaSemFOO_TYPE)
	e2 := JSError_MakeWithSourceLocation(jsaSemNULL_SOURCE, 2, -1, jsaSemFOO_TYPE)

	jsaSemAssertSmaller(t, jsaSemError(e1), jsaSemError(e2))
}

func TestSortingErrorManager_OrderingLineno2(t *testing.T) {
	e1 := JSError_MakeWithSourceLocation(jsaSemNULL_SOURCE, 8, -1, jsaSemFOO_TYPE)
	e2 := JSError_MakeWithSourceLocation(jsaSemNULL_SOURCE, 56, -1, jsaSemFOO_TYPE)
	jsaSemAssertSmaller(t, jsaSemError(e1), jsaSemError(e2))
}

func TestSortingErrorManager_OrderingCheckLevel(t *testing.T) {
	e1 := JSError_MakeWithSourceLocation(jsaSemNULL_SOURCE, -1, -1, jsaSemFOO_TYPE)
	e2 := JSError_MakeWithSourceLocation(jsaSemNULL_SOURCE, -1, -1, jsaSemFOO_TYPE)

	jsaSemAssertSmaller(t, jsaSemWarning(e1), jsaSemError(e2))
}

func TestSortingErrorManager_OrderingCharno1(t *testing.T) {
	e1 := JSError_MakeWithSourceLocation(jsaSemNULL_SOURCE, 5, -1, jsaSemFOO_TYPE)
	e2 := JSError_MakeWithSourceLocation(jsaSemNULL_SOURCE, 5, 2, jsaSemFOO_TYPE)

	jsaSemAssertSmaller(t, jsaSemError(e1), jsaSemError(e2))
	// CheckLevel preempts charno comparison
	jsaSemAssertSmaller(t, jsaSemWarning(e1), jsaSemError(e2))
}

func TestSortingErrorManager_OrderingCharno2(t *testing.T) {
	e1 := JSError_MakeWithSourceLocation(jsaSemNULL_SOURCE, 8, 7, jsaSemFOO_TYPE)
	e2 := JSError_MakeWithSourceLocation(jsaSemNULL_SOURCE, 8, 5, jsaSemFOO_TYPE)

	jsaSemAssertSmaller(t, jsaSemError(e2), jsaSemError(e1))
	// CheckLevel preempts charno comparison
	jsaSemAssertSmaller(t, jsaSemWarning(e2), jsaSemError(e1))
}

func TestSortingErrorManager_OrderingDescription(t *testing.T) {
	e1 := JSError_MakeWithSourceLocation(jsaSemNULL_SOURCE, -1, -1, jsaSemFOO_TYPE)
	e2 := JSError_MakeWithSourceLocation(jsaSemNULL_SOURCE, -1, -1, jsaSemJOO_TYPE)

	jsaSemAssertSmaller(t, jsaSemError(e1), jsaSemError(e2))
}

// jsaSemPrintingManager is the anonymous BasicErrorManager of
// testDeduplicatedErrors.
type jsaSemPrintingManager struct {
	BasicErrorManagerBase
	printedErrors []*JSError
}

func (m *jsaSemPrintingManager) Println(level *CheckLevel, error *JSError) {
	m.printedErrors = append(m.printedErrors, error)
}

func (m *jsaSemPrintingManager) PrintSummary() {}

func TestSortingErrorManager_DeduplicatedErrors(t *testing.T) {
	manager := &jsaSemPrintingManager{}
	manager.InitBasicErrorManagerBase(manager)
	e1 := JSError_MakeWithSourceLocation(jsaSemNULL_SOURCE, -1, -1, jsaSemFOO_TYPE)
	e2 := JSError_MakeWithSourceLocation(jsaSemNULL_SOURCE, -1, -1, jsaSemFOO_TYPE)
	manager.Report(CheckLevel_ERROR.Ptr(), e1)
	manager.Report(CheckLevel_ERROR.Ptr(), e2)
	manager.GenerateReport()
	jsaEq(t, len(manager.printedErrors), 1)
}

// jsaSemReportingManager is the anonymous BasicErrorManager of
// testGenerateReportCausesMoreWarnings.
type jsaSemReportingManager struct {
	BasicErrorManagerBase
	t       *testing.T
	printed int
}

func (m *jsaSemReportingManager) Println(level *CheckLevel, error *JSError) {
	if error.Type().Equals(jsaSemFOO_TYPE) {
		// TODO(b/114762232) This behavior should not be supported, and will become malformed
		// when migrating to a SortingErrorManager with an ErrorReportGenerator.
		m.Report(CheckLevel_ERROR.Ptr(), JSError_MakeWithSourceLocation(jsaSemNULL_SOURCE, -1, -1, jsaSemJOO_TYPE))
	}
	m.printed++
}

func (m *jsaSemReportingManager) PrintSummary() {
	jsaEq(m.t, m.printed, 1)
}

// This test is testing a "feature" that seems bogus and should likely be forbidden.
func TestSortingErrorManager_GenerateReportCausesMoreWarnings(t *testing.T) {
	manager := &jsaSemReportingManager{t: t}
	manager.InitBasicErrorManagerBase(manager)
	manager.Report(CheckLevel_ERROR.Ptr(), JSError_MakeWithSourceLocation(jsaSemNULL_SOURCE, -1, -1, jsaSemFOO_TYPE))
	manager.GenerateReport()
}

func jsaSemError(e *JSError) *SortingErrorManager_ErrorWithLevel {
	return NewSortingErrorManager_ErrorWithLevel(e, CheckLevel_ERROR.Ptr())
}

func jsaSemWarning(e *JSError) *SortingErrorManager_ErrorWithLevel {
	return NewSortingErrorManager_ErrorWithLevel(e, CheckLevel_WARNING.Ptr())
}

func jsaSemAssertSmaller(t *testing.T, p1, p2 *SortingErrorManager_ErrorWithLevel) {
	t.Helper()
	p1p2 := jsaSemComparator.Compare(p1, p2)
	jsaTrue(t, p1p2 < 0, strconv.Itoa(p1p2))
	p2p1 := jsaSemComparator.Compare(p2, p1)
	jsaTrue(t, p2p1 > 0, strconv.Itoa(p2p1))
}

// Go-only: extra coverage — counts, filtering and the report generators.
func TestSortingErrorManager_CountsAndGenerators(t *testing.T) {
	var reported []string
	gen := jsaSemGenerator(func(m SortingErrorManager) {
		for _, e := range m.GetSortedDiagnostics() {
			reported = append(reported, e.Error.Description())
		}
	})
	m := NewSortingErrorManager([]SortingErrorManager_ErrorReportGenerator{gen, gen})
	w := DiagnosticType_Warning("TEST_W", "w")
	m.Report(CheckLevel_ERROR.Ptr(), JSError_MakeWithSourceLocation(jsaStr("b.js"), 1, 1, jsaSemFOO_TYPE))
	m.Report(CheckLevel_ERROR.Ptr(), JSError_MakeWithSourceLocation(jsaStr("a.js"), 1, 1, w))
	m.Report(CheckLevel_WARNING.Ptr(), JSError_MakeWithSourceLocation(jsaStr("a.js"), 1, 1, w))
	m.Report(CheckLevel_OFF.Ptr(), JSError_MakeWithSourceLocation(jsaStr("a.js"), 2, 1, w))
	jsaEq(t, m.GetErrorCount(), 2)
	jsaEq(t, m.GetWarningCount(), 1)
	jsaTrue(t, m.HasHaltingErrors())
	jsaEq(t, len(m.GetErrors()), 2)
	jsaEq(t, len(m.GetWarnings()), 1)
	m.GenerateReport()
	jsaDeepEq(t, reported, []string{"w", "w", "w", "Foo"}) // OFF < WARNING < ERROR, then by file
	jsaTrue(t, ErrorManager_HasHaltingErrors(m))
	b := NewBlackHoleErrorManager()
	b.Report(CheckLevel_WARNING.Ptr(), JSError_Make(w))
	b.GenerateReport()
	jsaEq(t, b.GetWarningCount(), 1)
}

type jsaSemGenerator func(m SortingErrorManager)

func (g jsaSemGenerator) GenerateReport(m SortingErrorManager) { g(m) }
