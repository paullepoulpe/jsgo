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
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// FixingErrorManager is an error manager that finds a SuggestedFix for all
// errors if possible.
//
// Java: public class FixingErrorManager extends BasicErrorManager
type FixingErrorManager struct {
	jscomp.BasicErrorManagerBase

	fixer *ErrorToFixMapper
	// Go-only: language-adapter — a HashMap keyed by the JSError record,
	// whose hash mixes the identity hash of its Node (so Java's order is not
	// deterministic and only lookups and an any-match scan use it): an
	// insertion-ordered map with JSError.equals/hashCode (PORTING.md §6).
	sureFixes  *jutil.ObjectLinkedHashMap[*jscomp.JSError, *SuggestedFix]
	multiFixes *jutil.ObjectLinkedHashMap[*jscomp.JSError, []*SuggestedFix]
	// Go-only: language-adapter — an ImmutableSet of DiagnosticTypes, which
	// compare by key (DiagnosticType.equals): ObjectLinkedHashSet, as
	// DiagnosticGroup's types (PORTING.md §6).
	unfixableErrors *jutil.ObjectLinkedHashSet[*jscomp.DiagnosticType]
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
	return NewFixingErrorManagerWithUnfixableErrors(jutil.NewObjectLinkedHashSet[*jscomp.DiagnosticType]())
}

// Java: FixingErrorManager(ImmutableSet<DiagnosticType>)
func NewFixingErrorManagerWithUnfixableErrors(unfixableErrors *jutil.ObjectLinkedHashSet[*jscomp.DiagnosticType]) *FixingErrorManager {
	m := &FixingErrorManager{}
	m.InitBasicErrorManagerBase(m)
	m.sureFixes = jutil.NewObjectLinkedHashMap[*jscomp.JSError, *SuggestedFix]()
	m.multiFixes = jutil.NewObjectLinkedHashMap[*jscomp.JSError, []*SuggestedFix]()
	m.unfixableErrors = unfixableErrors
	return m
}

// Java: FixingErrorManager.setCompiler(AbstractCompiler)
func (m *FixingErrorManager) SetCompiler(compiler jscomp.AbstractCompiler) {
	m.fixer = NewErrorToFixMapper(compiler)
}

// Java: FixingErrorManager.report(CheckLevel, JSError)
func (m *FixingErrorManager) Report(level *jscomp.CheckLevel, error *jscomp.JSError) {
	m.BasicErrorManagerBase.Report(level, error)
	if !m.unfixableErrors.Contains(error.Type()) {
		fixes := m.fixer.GetFixesForJsError(error)
		if len(fixes) == 1 {
			m.sureFixes.Put(error, fixes[0])
		} else if len(fixes) > 0 {
			m.multiFixes.Put(error, append(m.multiFixes.Get(error), fixes...))
		}
	}
}

// Java: FixingErrorManager.containsFixableShorthandModuleWarning()
func (m *FixingErrorManager) ContainsFixableShorthandModuleWarning() bool {
	for error := range m.sureFixes.Keys() {
		if error.Type().Equals(jscomp.ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME) {
			return true
		}
	}
	for error := range m.multiFixes.Keys() {
		if error.Type().Equals(jscomp.ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME) {
			return true
		}
	}
	return false
}

// Java: FixingErrorManager.getFixesForJsError(JSError)
func (m *FixingErrorManager) GetFixesForJsError(error *jscomp.JSError) []*SuggestedFix {
	if m.sureFixes.ContainsKey(error) {
		return []*SuggestedFix{m.sureFixes.Get(error)}
	} else {
		return m.multiFixes.Get(error)
	}
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
			if fixTypes == FixingErrorManager_FixTypes_ONE_FIX && m.sureFixes.ContainsKey(error) {
				fixes = append(fixes, m.sureFixes.Get(error))
			} else if fixTypes == FixingErrorManager_FixTypes_ONE_FIX_AND_MULTI_FIXES {
				fixes = append(fixes, m.GetFixesForJsError(error)...)
			}
		}
	}
	for _, warning := range m.GetWarnings() {
		if warning.Type().Equals(jscomp.CheckExtraRequires_EXTRA_REQUIRE_WARNING) && containsFixableShorthandModuleWarning {
			// As above, don't apply the extra-require fix.
		} else {
			if fixTypes == FixingErrorManager_FixTypes_ONE_FIX && m.sureFixes.ContainsKey(warning) {
				fixes = append(fixes, m.sureFixes.Get(warning))
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
