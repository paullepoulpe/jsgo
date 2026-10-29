// Port of com/google/javascript/refactoring/FixingErrorManager.java @ closure-compiler v20261021

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

package refactoring

import (
	"github.com/paullepoulpe/jsgo/internal/jscomp"
)

type fixingSureEntry struct {
	err *jscomp.JSError
	fix *SuggestedFix
}

type fixingMultiEntry struct {
	err   *jscomp.JSError
	fixes []*SuggestedFix
}

// FixingErrorManager is an error manager that finds a SuggestedFix for all
// errors if possible.
//
// Java: public class FixingErrorManager extends BasicErrorManager
type FixingErrorManager struct {
	jscomp.SortingErrorManager

	fixer                                 *ErrorToFixMapper
	sureFixes                             map[int32][]fixingSureEntry
	multiFixes                            map[int32][]fixingMultiEntry
	unfixableErrorKeys                    map[string]struct{}
	containsFixableShorthandModuleWarning bool
}

// FixingErrorManager_FixTypes selects the fixes getFixes collects.
//
// Java: private enum FixingErrorManager.FixTypes
type FixingErrorManager_FixTypes int32

const (
	FixingErrorManager_FixTypes_ONE_FIX FixingErrorManager_FixTypes = iota
	FixingErrorManager_FixTypes_ONE_FIX_AND_MULTI_FIXES
)

var fixingErrorManager_FixTypes_names = [...]string{"ONE_FIX", "ONE_FIX_AND_MULTI_FIXES"}

// Name is Enum.name().
func (t FixingErrorManager_FixTypes) Name() string { return fixingErrorManager_FixTypes_names[t] }

// String is Enum.toString().
func (t FixingErrorManager_FixTypes) String() string { return t.Name() }

// Ordinal is Enum.ordinal().
func (t FixingErrorManager_FixTypes) Ordinal() int { return int(t) }

// Java: FixingErrorManager()
func NewFixingErrorManager() *FixingErrorManager {
	return NewFixingErrorManagerWithUnfixableErrors(nil)
}

// Java: FixingErrorManager(ImmutableSet<DiagnosticType>)
func NewFixingErrorManagerWithUnfixableErrors(unfixableErrors []*jscomp.DiagnosticType) *FixingErrorManager {
	m := &FixingErrorManager{}
	m.InitSortingErrorManager(nil)
	m.sureFixes = make(map[int32][]fixingSureEntry)
	m.multiFixes = make(map[int32][]fixingMultiEntry)
	if len(unfixableErrors) > 0 {
		m.unfixableErrorKeys = make(map[string]struct{}, len(unfixableErrors))
		for _, dt := range unfixableErrors {
			if dt != nil {
				m.unfixableErrorKeys[dt.Key()] = struct{}{}
			}
		}
	}
	return m
}

// Java: FixingErrorManager.setCompiler(AbstractCompiler)
func (m *FixingErrorManager) SetCompiler(compiler jscomp.AbstractCompiler) {
	m.fixer = NewErrorToFixMapper(compiler)
}

func (m *FixingErrorManager) isUnfixable(dt *jscomp.DiagnosticType) bool {
	if dt == nil || len(m.unfixableErrorKeys) == 0 {
		return false
	}
	_, ok := m.unfixableErrorKeys[dt.Key()]
	return ok
}

func (m *FixingErrorManager) putSureFix(err *jscomp.JSError, fix *SuggestedFix) {
	if err.Type().Equals(jscomp.ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME) {
		m.containsFixableShorthandModuleWarning = true
	}
	h := err.HashCode()
	bucket := m.sureFixes[h]
	for i := range bucket {
		if bucket[i].err == err || bucket[i].err.Equals(err) {
			bucket[i].fix = fix
			return
		}
	}
	m.sureFixes[h] = append(bucket, fixingSureEntry{err: err, fix: fix})
}

func (m *FixingErrorManager) lookupSureFix(err *jscomp.JSError) (*SuggestedFix, bool) {
	if err == nil {
		return nil, false
	}
	for _, entry := range m.sureFixes[err.HashCode()] {
		if entry.err == err || entry.err.Equals(err) {
			return entry.fix, true
		}
	}
	return nil, false
}

func (m *FixingErrorManager) appendMultiFixes(err *jscomp.JSError, fixes []*SuggestedFix) {
	if err.Type().Equals(jscomp.ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME) {
		m.containsFixableShorthandModuleWarning = true
	}
	h := err.HashCode()
	bucket := m.multiFixes[h]
	for i := range bucket {
		if bucket[i].err == err || bucket[i].err.Equals(err) {
			bucket[i].fixes = append(bucket[i].fixes, fixes...)
			return
		}
	}
	cloned := make([]*SuggestedFix, len(fixes))
	copy(cloned, fixes)
	m.multiFixes[h] = append(bucket, fixingMultiEntry{err: err, fixes: cloned})
}

func (m *FixingErrorManager) getMultiFixes(err *jscomp.JSError) []*SuggestedFix {
	if err == nil {
		return nil
	}
	for _, entry := range m.multiFixes[err.HashCode()] {
		if entry.err == err || entry.err.Equals(err) {
			return entry.fixes
		}
	}
	return nil
}

// Java: FixingErrorManager.report(CheckLevel, JSError)
func (m *FixingErrorManager) Report(level *jscomp.CheckLevel, error *jscomp.JSError) {
	m.SortingErrorManager.Report(level, error)
	if !m.isUnfixable(error.Type()) {
		fixes := m.fixer.GetFixesForJsError(error)
		if len(fixes) == 1 {
			m.putSureFix(error, fixes[0])
		} else if len(fixes) > 0 {
			m.appendMultiFixes(error, fixes)
		}
	}
}

// Java: FixingErrorManager.containsFixableShorthandModuleWarning()
func (m *FixingErrorManager) ContainsFixableShorthandModuleWarning() bool {
	return m.containsFixableShorthandModuleWarning
}

// Java: FixingErrorManager.getFixesForJsError(JSError)
func (m *FixingErrorManager) GetFixesForJsError(error *jscomp.JSError) []*SuggestedFix {
	if fix, ok := m.lookupSureFix(error); ok {
		return []*SuggestedFix{fix}
	}
	return m.getMultiFixes(error)
}

// Java: FixingErrorManager.getFixes(FixTypes)
func (m *FixingErrorManager) GetFixes(fixTypes FixingErrorManager_FixTypes) []*SuggestedFix {
	containsFixableShorthandModuleWarning := m.ContainsFixableShorthandModuleWarning()
	fixes := []*SuggestedFix{}
	for _, error := range m.GetErrors() {
		// Sometimes code will produce a spurious extra-require error, as well
		// as a warning about using a full namespace instead of a shorthand
		// type. In this case don't apply the extra require fix.
		if containsFixableShorthandModuleWarning && error.Type().Equals(jscomp.CheckExtraRequires_EXTRA_REQUIRE_WARNING) {
			// Don't apply this fix.
		} else {
			if fixTypes == FixingErrorManager_FixTypes_ONE_FIX {
				if fix, ok := m.lookupSureFix(error); ok {
					fixes = append(fixes, fix)
				}
			} else if fixTypes == FixingErrorManager_FixTypes_ONE_FIX_AND_MULTI_FIXES {
				fixes = append(fixes, m.GetFixesForJsError(error)...)
			}
		}
	}
	for _, warning := range m.GetWarnings() {
		if warning.Type().Equals(jscomp.CheckExtraRequires_EXTRA_REQUIRE_WARNING) && containsFixableShorthandModuleWarning {
			// As above, don't apply the extra-require fix.
		} else {
			if fixTypes == FixingErrorManager_FixTypes_ONE_FIX {
				if fix, ok := m.lookupSureFix(warning); ok {
					fixes = append(fixes, fix)
				}
			} else if fixTypes == FixingErrorManager_FixTypes_ONE_FIX_AND_MULTI_FIXES {
				fixes = append(fixes, m.GetFixesForJsError(warning)...)
			}
		}
	}
	return fixes
}

// GetSureFixes returns fixes for errors and warnings that only have one
// 'sure' guaranteed fix.
//
// Java: FixingErrorManager.getSureFixes()
func (m *FixingErrorManager) GetSureFixes() []*SuggestedFix {
	return m.GetFixes(FixingErrorManager_FixTypes_ONE_FIX)
}

// GetAllFixes returns fixes for errors first, then fixes for warnings. This
// includes 'sure' fixes with only one option, and 'multi' fixes which have
// multiple choices.
//
// Java: FixingErrorManager.getAllFixes()
func (m *FixingErrorManager) GetAllFixes() []*SuggestedFix {
	return m.GetFixes(FixingErrorManager_FixTypes_ONE_FIX_AND_MULTI_FIXES)
}

// Java: FixingErrorManager.printSummary()
func (m *FixingErrorManager) PrintSummary() {
	// No-op
}

// Java: FixingErrorManager.println(CheckLevel, JSError)
func (m *FixingErrorManager) Println(level *jscomp.CheckLevel, error *jscomp.JSError) {
	// No-op
}
