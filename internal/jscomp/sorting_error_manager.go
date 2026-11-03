// Port of com/google/javascript/jscomp/SortingErrorManager.java @ closure-compiler v20261028

/*
 * Copyright 2018 The Closure Compiler Authors.
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

type sortingErrorKey struct {
	level       CheckLevel
	hasSource   bool
	source      string
	lineno      int
	charno      int
	description string
}

func makeSortingErrorKey(err *JSError, level CheckLevel) sortingErrorKey {
	source, hasSource := err.SourceName()
	return sortingErrorKey{
		level:       level,
		hasSource:   hasSource,
		source:      source,
		lineno:      err.GetLineNumber(),
		charno:      err.Charno(),
		description: err.Description(),
	}
}

// SortingErrorManager is a customizable error manager that sorts all errors
// and warnings reported to it, and has customizable output through the
// ErrorReportGenerator interface.
type SortingErrorManager struct {
	messages           []SortingErrorManager_ErrorWithLevel
	seen               map[sortingErrorKey]struct{}
	sorted             bool
	originalErrorCount int
	promotedErrorCount int
	warningCount       int
	typedPercent       float64

	// Responsible for generating the report of the errors at the end of compilation
	errorReportGenerators []SortingErrorManager_ErrorReportGenerator
	println               func(CheckLevel, *JSError)
	printSummary          func(*SortingErrorManager)
}

var _ ErrorManager = (*SortingErrorManager)(nil)

// NewSortingErrorManager creates a SortingErrorManager.
//
// Java: SortingErrorManager(Set<ErrorReportGenerator>)
func NewSortingErrorManager(errorReportGenerators []SortingErrorManager_ErrorReportGenerator) *SortingErrorManager {
	m := &SortingErrorManager{}
	m.InitSortingErrorManager(errorReportGenerators)
	return m
}

// InitSortingErrorManager initializes an embedded SortingErrorManager.
//
// Java: SortingErrorManager(Set<ErrorReportGenerator>)
func (m *SortingErrorManager) InitSortingErrorManager(
	errorReportGenerators []SortingErrorManager_ErrorReportGenerator) {
	m.seen = make(map[sortingErrorKey]struct{})
	m.sorted = true
	if len(errorReportGenerators) > 0 {
		m.errorReportGenerators = sortingErrorManager_immutableSetCopyOf(errorReportGenerators)
	}
}

func (m *SortingErrorManager) initReportHooks(
	printlnFn func(CheckLevel, *JSError),
	printSummaryFn func(*SortingErrorManager),
) {
	m.InitSortingErrorManager(nil)
	m.println = printlnFn
	m.printSummary = printSummaryFn
}

func (m *SortingErrorManager) sortingErrorManager() *SortingErrorManager { return m }

// sortingErrorManager_immutableSetCopyOf is ImmutableSet.copyOf for the
// generators, which compare by identity.
func sortingErrorManager_immutableSetCopyOf(generators []SortingErrorManager_ErrorReportGenerator) []SortingErrorManager_ErrorReportGenerator {
	out := make([]SortingErrorManager_ErrorReportGenerator, 0, len(generators))
	for _, g := range generators {
		jutil.CheckNotNull(g)
		dup := false
		for _, o := range out {
			if jutil.ObjectsEquals(o, g) {
				dup = true
				break
			}
		}
		if !dup {
			out = append(out, g)
		}
	}
	return out
}

// Java: SortingErrorManager.report(CheckLevel, JSError)
func (m *SortingErrorManager) Report(level CheckLevel, error *JSError) {
	key := makeSortingErrorKey(error, level)
	if m.seen == nil {
		m.seen = make(map[sortingErrorKey]struct{})
	}
	if _, exists := m.seen[key]; exists {
		return
	}
	m.seen[key] = struct{}{}
	m.messages = append(m.messages, NewSortingErrorManager_ErrorWithLevel(error, level))
	m.sorted = false
	switch level {
	case CheckLevel_ERROR:
		if error.Type().level == CheckLevel_ERROR {
			m.originalErrorCount++
		} else {
			m.promotedErrorCount++
		}
	case CheckLevel_WARNING:
		m.warningCount++
	}
}

func (m *SortingErrorManager) ensureSorted() {
	if !m.sorted {
		slices.SortFunc(m.messages, sortingErrorManagerCompare)
		m.sorted = true
	}
}

// Java: SortingErrorManager.hasHaltingErrors()
func (m *SortingErrorManager) HasHaltingErrors() bool {
	return m.originalErrorCount != 0
}

// Java: SortingErrorManager.getErrorCount()
func (m *SortingErrorManager) GetErrorCount() int {
	return m.originalErrorCount + m.promotedErrorCount
}

// Java: SortingErrorManager.getWarningCount()
func (m *SortingErrorManager) GetWarningCount() int {
	return m.warningCount
}

// Java: SortingErrorManager.getErrors()
func (m *SortingErrorManager) GetErrors() []*JSError {
	return m.FilterToListOf(CheckLevel_ERROR)
}

// Java: SortingErrorManager.getWarnings()
func (m *SortingErrorManager) GetWarnings() []*JSError {
	return m.FilterToListOf(CheckLevel_WARNING)
}

// Java: SortingErrorManager.getSortedDiagnostics()
func (m *SortingErrorManager) GetSortedDiagnostics() []SortingErrorManager_ErrorWithLevel {
	// TODO(b/114762232): It should be possible to remove the copying here and switch to an
	// unmodifiable collection once we get rid of usages that add warnings during generateReport.
	m.ensureSorted()
	return slices.Clone(m.messages)
}

// Java: SortingErrorManager.setTypedPercent(double)
func (m *SortingErrorManager) SetTypedPercent(typedPercent float64) {
	m.typedPercent = typedPercent
}

// Java: SortingErrorManager.getTypedPercent()
func (m *SortingErrorManager) GetTypedPercent() float64 {
	return m.typedPercent
}

// Java: SortingErrorManager.filterToListOf(CheckLevel)
func (m *SortingErrorManager) FilterToListOf(level CheckLevel) []*JSError {
	m.ensureSorted()
	out := []*JSError{}
	for _, e := range m.messages {
		if e.Level == level {
			out = append(out, e.Error)
		}
	}
	return out
}

// Java: SortingErrorManager.generateReport()
func (m *SortingErrorManager) GenerateReport() {
	if m.println != nil || m.printSummary != nil {
		diagnostics := m.GetSortedDiagnostics()
		if m.println != nil {
			for _, message := range diagnostics {
				m.println(message.Level, message.Error)
			}
		}
		if m.printSummary != nil {
			m.printSummary(m)
		}
		return
	}
	// TODO(b/114762232): It should be invalid to report errors during the execution of this method;
	// doing so will become impossible once all subclases have migrated to an ErrorReportGenerator.
	for _, generator := range m.errorReportGenerators {
		generator.GenerateReport(m)
	}
}

// Java: ErrorManager.shouldReportConformanceViolation(Requirement, Optional<RequirementScopeEntry>, JSError, LibraryLevelNonAllowlistedConformanceViolationsBehavior, boolean)
func (m *SortingErrorManager) ShouldReportConformanceViolation(
	requirement *Requirement,
	allowlistEntry *RequirementScopeEntry,
	diagnostic *JSError,
	behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior,
	isAllowlisted bool) bool {
	return ErrorManager_ShouldReportConformanceViolation(m, requirement, allowlistEntry, diagnostic, behavior, isAllowlisted)
}

// SortingErrorManager_ErrorReportGenerator is the strategy for customizing
// the output format of the error report.
type SortingErrorManager_ErrorReportGenerator interface {
	// Java: SortingErrorManager.ErrorReportGenerator.generateReport(SortingErrorManager)
	GenerateReport(manager *SortingErrorManager)
}

func sortingErrorManagerCompare(p1, p2 SortingErrorManager_ErrorWithLevel) int {
	// null is the smallest value
	if p2.Error == nil {
		if p1.Error == nil {
			return 0
		}
		return 1
	}
	if p1.Error == nil {
		return -1
	}

	// check level
	if p1.Level != p2.Level {
		return int(p2.Level) - int(p1.Level)
	}

	// sourceName comparison
	source1, hasSource1 := p1.Error.SourceName()
	source2, hasSource2 := p2.Error.SourceName()
	if hasSource1 && hasSource2 {
		if sourceCompare := jutil.StringCompareTo(source1, source2); sourceCompare != 0 {
			return sourceCompare
		}
	} else if !hasSource1 && hasSource2 {
		return -1
	} else if hasSource1 && !hasSource2 {
		return 1
	}

	// lineno comparison
	lineno1 := p1.Error.GetLineNumber()
	lineno2 := p2.Error.GetLineNumber()
	if lineno1 != lineno2 {
		return lineno1 - lineno2
	}

	// charno comparison
	charno1 := p1.Error.Charno()
	charno2 := p2.Error.Charno()
	if charno1 != charno2 {
		return charno1 - charno2
	}

	// description
	return jutil.StringCompareTo(p1.Error.Description(), p2.Error.Description())
}

// SortingErrorManager_ErrorWithLevel is a JSError with its reported level.
//
// Go-only: language-adapter — the package-private final fields error and
// level are exported as Error and Level, since report generators read them.
type SortingErrorManager_ErrorWithLevel struct {
	Error *JSError
	Level CheckLevel
}

// Java: SortingErrorManager.ErrorWithLevel(JSError, CheckLevel)
func NewSortingErrorManager_ErrorWithLevel(error *JSError, level CheckLevel) SortingErrorManager_ErrorWithLevel {
	return SortingErrorManager_ErrorWithLevel{Error: error, Level: level}
}
