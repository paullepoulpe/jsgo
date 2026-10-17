// Port of com/google/javascript/jscomp/SortingErrorManager.java @ closure-compiler v20260415

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SortingErrorManager is a customizable error manager that sorts all errors
// and warnings reported to it, and has customizable output through the
// ErrorReportGenerator interface.
type SortingErrorManager interface {
	ErrorManager

	// Java: SortingErrorManager.getSortedDiagnostics()
	GetSortedDiagnostics() []*SortingErrorManager_ErrorWithLevel

	GetSortingErrorManagerBase() *SortingErrorManagerBase
}

// SortingErrorManagerBase holds the fields and base implementations of
// SortingErrorManager; subclasses embed it.
type SortingErrorManagerBase struct {
	self SortingErrorManager

	messages           *jutil.TreeSet[*SortingErrorManager_ErrorWithLevel]
	originalErrorCount int
	promotedErrorCount int
	warningCount       int
	typedPercent       float64

	// Responsible for generating the report of the errors at the end of compilation
	errorReportGenerators []SortingErrorManager_ErrorReportGenerator
}

var _ SortingErrorManager = (*SortingErrorManagerBase)(nil)

// NewSortingErrorManager creates a SortingErrorManager.
//
// Java: SortingErrorManager(Set<ErrorReportGenerator>)
func NewSortingErrorManager(errorReportGenerators []SortingErrorManager_ErrorReportGenerator) *SortingErrorManagerBase {
	m := &SortingErrorManagerBase{}
	m.InitSortingErrorManagerBase(m, errorReportGenerators)
	return m
}

// InitSortingErrorManagerBase is the Go form of super(errorReportGenerators);
// self is the object being constructed.
//
// Java: SortingErrorManager(Set<ErrorReportGenerator>)
func (m *SortingErrorManagerBase) InitSortingErrorManagerBase(
	self SortingErrorManager, errorReportGenerators []SortingErrorManager_ErrorReportGenerator) {
	m.self = self
	m.messages = jutil.NewTreeSet(NewSortingErrorManager_LeveledJSErrorComparator().Compare)
	m.errorReportGenerators = sortingErrorManager_immutableSetCopyOf(errorReportGenerators)
}

// sortingErrorManager_immutableSetCopyOf is ImmutableSet.copyOf for the
// generators, which compare by identity.
func sortingErrorManager_immutableSetCopyOf(generators []SortingErrorManager_ErrorReportGenerator) []SortingErrorManager_ErrorReportGenerator {
	out := []SortingErrorManager_ErrorReportGenerator{}
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

// GetSortingErrorManagerBase gives access to the base fields.
func (m *SortingErrorManagerBase) GetSortingErrorManagerBase() *SortingErrorManagerBase { return m }

// Java: SortingErrorManager.report(CheckLevel, JSError)
func (m *SortingErrorManagerBase) Report(level *CheckLevel, error *JSError) {
	e := NewSortingErrorManager_ErrorWithLevel(error, level)
	if m.messages.Add(e) {
		if level != nil && *level == CheckLevel_ERROR {
			if error.Type().level == CheckLevel_ERROR {
				m.originalErrorCount++
			} else {
				m.promotedErrorCount++
			}
		} else if level != nil && *level == CheckLevel_WARNING {
			m.warningCount++
		}
	}
}

// Java: SortingErrorManager.hasHaltingErrors()
func (m *SortingErrorManagerBase) HasHaltingErrors() bool {
	return m.originalErrorCount != 0
}

// Java: SortingErrorManager.getErrorCount()
func (m *SortingErrorManagerBase) GetErrorCount() int {
	return m.originalErrorCount + m.promotedErrorCount
}

// Java: SortingErrorManager.getWarningCount()
func (m *SortingErrorManagerBase) GetWarningCount() int {
	return m.warningCount
}

// Java: SortingErrorManager.getErrors()
func (m *SortingErrorManagerBase) GetErrors() []*JSError {
	return m.FilterToListOf(CheckLevel_ERROR)
}

// Java: SortingErrorManager.getWarnings()
func (m *SortingErrorManagerBase) GetWarnings() []*JSError {
	return m.FilterToListOf(CheckLevel_WARNING)
}

// Java: SortingErrorManager.getSortedDiagnostics()
func (m *SortingErrorManagerBase) GetSortedDiagnostics() []*SortingErrorManager_ErrorWithLevel {
	// TODO(b/114762232): It should be possible to remove the copying here and switch to an
	// unmodifiable collection once we get rid of usages that add warnings during generateReport.
	return m.messages.ToSlice()
}

// Java: SortingErrorManager.setTypedPercent(double)
func (m *SortingErrorManagerBase) SetTypedPercent(typedPercent float64) {
	m.typedPercent = typedPercent
}

// Java: SortingErrorManager.getTypedPercent()
func (m *SortingErrorManagerBase) GetTypedPercent() float64 {
	return m.typedPercent
}

// Java: SortingErrorManager.filterToListOf(CheckLevel)
func (m *SortingErrorManagerBase) FilterToListOf(level CheckLevel) []*JSError {
	out := []*JSError{}
	for e := range m.messages.All() {
		if e.Level != nil && *e.Level == level {
			out = append(out, e.Error)
		}
	}
	return out
}

// Java: SortingErrorManager.generateReport()
func (m *SortingErrorManagerBase) GenerateReport() {
	// TODO(b/114762232): It should be invalid to report errors during the execution of this method;
	// doing so will become impossible once all subclases have migrated to an ErrorReportGenerator.
	for _, generator := range m.errorReportGenerators {
		generator.GenerateReport(m.self)
	}
}

// Java: ErrorManager.shouldReportConformanceViolation(Requirement, Optional<RequirementScopeEntry>, JSError, LibraryLevelNonAllowlistedConformanceViolationsBehavior, boolean)
func (m *SortingErrorManagerBase) ShouldReportConformanceViolation(
	requirement *Requirement,
	allowlistEntry *RequirementScopeEntry,
	diagnostic *JSError,
	behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior,
	isAllowlisted bool) bool {
	return ErrorManager_ShouldReportConformanceViolation(m.self, requirement, allowlistEntry, diagnostic, behavior, isAllowlisted)
}

// SortingErrorManager_ErrorReportGenerator is the strategy for customizing
// the output format of the error report.
type SortingErrorManager_ErrorReportGenerator interface {
	// Java: SortingErrorManager.ErrorReportGenerator.generateReport(SortingErrorManager)
	GenerateReport(manager SortingErrorManager)
}

// SortingErrorManager_LeveledJSErrorComparator is the comparator of JSError
// with an associated CheckLevel. The ordering is the standard lexical
// ordering on the quintuple (file name, line number, CheckLevel, character
// number, description).
//
// Note: this comparator imposes orderings that are inconsistent with
// JSError.equals(Object).
type SortingErrorManager_LeveledJSErrorComparator struct{}

const (
	sortingErrorManager_LeveledJSErrorComparator_P1_LT_P2 = -1
	sortingErrorManager_LeveledJSErrorComparator_P1_GT_P2 = 1
)

// Java: SortingErrorManager.LeveledJSErrorComparator()
func NewSortingErrorManager_LeveledJSErrorComparator() *SortingErrorManager_LeveledJSErrorComparator {
	return &SortingErrorManager_LeveledJSErrorComparator{}
}

// Java: SortingErrorManager.LeveledJSErrorComparator.compare(ErrorWithLevel, ErrorWithLevel)
func (c *SortingErrorManager_LeveledJSErrorComparator) Compare(p1, p2 *SortingErrorManager_ErrorWithLevel) int {
	// null is the smallest value
	if p2 == nil {
		if p1 == nil {
			return 0
		} else {
			return sortingErrorManager_LeveledJSErrorComparator_P1_GT_P2
		}
	}

	// check level
	if p1.Level != p2.Level {
		if p1.Level == nil || p2.Level == nil {
			// p2.level.compareTo(p1.level) with a null level
			panic(jutil.NewNullPointerException())
		}
		return p2.Level.CompareTo(*p1.Level)
	}

	// sourceName comparison
	source1 := p1.Error.SourceName()
	source2 := p2.Error.SourceName()
	if source1 != nil && source2 != nil {
		sourceCompare := jutil.StringCompareTo(*source1, *source2)
		if sourceCompare != 0 {
			return sourceCompare
		}
	} else if source1 == nil && source2 != nil {
		return sortingErrorManager_LeveledJSErrorComparator_P1_LT_P2
	} else if source1 != nil && source2 == nil {
		return sortingErrorManager_LeveledJSErrorComparator_P1_GT_P2
	}

	// lineno comparison
	lineno1 := p1.Error.GetLineNumber()
	lineno2 := p2.Error.GetLineNumber()
	if lineno1 != lineno2 {
		return lineno1 - lineno2
	} else if lineno1 < 0 && 0 <= lineno2 {
		return sortingErrorManager_LeveledJSErrorComparator_P1_LT_P2
	} else if 0 <= lineno1 && lineno2 < 0 {
		return sortingErrorManager_LeveledJSErrorComparator_P1_GT_P2
	}

	// charno comparison
	charno1 := p1.Error.Charno()
	charno2 := p2.Error.Charno()
	if charno1 != charno2 {
		return charno1 - charno2
	} else if charno1 < 0 && 0 <= charno2 {
		return sortingErrorManager_LeveledJSErrorComparator_P1_LT_P2
	} else if 0 <= charno1 && charno2 < 0 {
		return sortingErrorManager_LeveledJSErrorComparator_P1_GT_P2
	}

	// description
	return jutil.StringCompareTo(p1.Error.Description(), p2.Error.Description())
}

// SortingErrorManager_ErrorWithLevel is a JSError with its reported level
// (which may be nil).
//
// Go-only: language-adapter — the package-private final fields error and
// level are exported as Error and Level, since report generators read them.
type SortingErrorManager_ErrorWithLevel struct {
	Error *JSError
	Level *CheckLevel
}

// Java: SortingErrorManager.ErrorWithLevel(JSError, CheckLevel)
func NewSortingErrorManager_ErrorWithLevel(error *JSError, level *CheckLevel) *SortingErrorManager_ErrorWithLevel {
	if level != nil {
		level = level.Ptr()
	}
	return &SortingErrorManager_ErrorWithLevel{Error: error, Level: level}
}

// HashCode is Objects.hash(level, description, sourceName, lineNumber,
// charno).
//
// Go-only: language-adapter — Java hashes the CheckLevel by identity; its
// ordinal is used here.
//
// Java: SortingErrorManager.ErrorWithLevel.hashCode()
func (e *SortingErrorManager_ErrorWithLevel) HashCode() int32 {
	h := int32(1)
	if e.Level != nil {
		h = 31*h + int32(e.Level.Ordinal())
	} else {
		h = 31 * h
	}
	h = 31*h + jutil.StringHashCode(e.Error.Description())
	if s := e.Error.SourceName(); s != nil {
		h = 31*h + jutil.StringHashCode(*s)
	} else {
		h = 31 * h
	}
	h = 31*h + int32(e.Error.GetLineNumber())
	h = 31*h + int32(e.Error.Charno())
	return h
}

// Java: SortingErrorManager.ErrorWithLevel.equals(Object)
func (e *SortingErrorManager_ErrorWithLevel) Equals(obj any) bool {
	o, ok := obj.(*SortingErrorManager_ErrorWithLevel)
	if !ok || o == nil {
		return false
	}
	return ((e.Level == nil && o.Level == nil) || (e.Level != nil && o.Level != nil && *e.Level == *o.Level)) &&
		e.Error.Description() == o.Error.Description() &&
		sortingErrorManager_nullableStringEquals(e.Error.SourceName(), o.Error.SourceName()) &&
		e.Error.GetLineNumber() == o.Error.GetLineNumber() &&
		e.Error.Charno() == o.Error.Charno()
}

// sortingErrorManager_nullableStringEquals is Objects.equals for nullable
// Strings.
func sortingErrorManager_nullableStringEquals(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
