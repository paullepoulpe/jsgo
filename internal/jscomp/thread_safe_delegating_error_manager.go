// Port of com/google/javascript/jscomp/ThreadSafeDelegatingErrorManager.java @ closure-compiler v20261021

/*
 * Copyright 2017 The Closure Compiler Authors.
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

import "sync"

// ConformanceViolationHandler handles conformance violations reported to a
// ThreadSafeDelegatingErrorManager and returns whether the violation should be
// surfaced as a compiler diagnostic.
type ConformanceViolationHandler func(
	requirement *Requirement,
	allowlistEntry *RequirementScopeEntry,
	diagnostic *JSError,
	behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior,
	isAllowlisted bool,
) bool

// ThreadSafeDelegatingErrorManager is a simple delegating ErrorManager that
// provides a thread-safe wrapper for the one being delegated.
type ThreadSafeDelegatingErrorManager struct {
	mu                          sync.Mutex
	delegated                   ErrorManager
	conformanceViolationHandler ConformanceViolationHandler
}

var _ ErrorManager = (*ThreadSafeDelegatingErrorManager)(nil)

// Java: ThreadSafeDelegatingErrorManager(ErrorManager)
func NewThreadSafeDelegatingErrorManager(delegated ErrorManager) *ThreadSafeDelegatingErrorManager {
	return &ThreadSafeDelegatingErrorManager{delegated: delegated}
}

// Java: ThreadSafeDelegatingErrorManager.report(CheckLevel, JSError)
func (m *ThreadSafeDelegatingErrorManager) Report(level *CheckLevel, error *JSError) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.delegated.Report(level, error)
}

// Java: ThreadSafeDelegatingErrorManager.generateReport()
func (m *ThreadSafeDelegatingErrorManager) GenerateReport() {
	m.delegated.GenerateReport()
}

// Java: ThreadSafeDelegatingErrorManager.hasHaltingErrors()
func (m *ThreadSafeDelegatingErrorManager) HasHaltingErrors() bool {
	return m.delegated.HasHaltingErrors()
}

// Java: ThreadSafeDelegatingErrorManager.getErrorCount()
func (m *ThreadSafeDelegatingErrorManager) GetErrorCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.delegated.GetErrorCount()
}

// Java: ThreadSafeDelegatingErrorManager.getWarningCount()
func (m *ThreadSafeDelegatingErrorManager) GetWarningCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.delegated.GetWarningCount()
}

// Java: ThreadSafeDelegatingErrorManager.getErrors()
func (m *ThreadSafeDelegatingErrorManager) GetErrors() []*JSError {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.delegated.GetErrors()
}

// Java: ThreadSafeDelegatingErrorManager.getWarnings()
func (m *ThreadSafeDelegatingErrorManager) GetWarnings() []*JSError {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.delegated.GetWarnings()
}

// Java: ThreadSafeDelegatingErrorManager.setTypedPercent(double)
func (m *ThreadSafeDelegatingErrorManager) SetTypedPercent(typedPercent float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.delegated.SetTypedPercent(typedPercent)
}

// Java: ThreadSafeDelegatingErrorManager.getTypedPercent()
func (m *ThreadSafeDelegatingErrorManager) GetTypedPercent() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.delegated.GetTypedPercent()
}

// SetConformanceViolationHandler installs an optional callback invoked by
// ShouldReportConformanceViolation under the manager's lock. When nil, calls
// delegate to the wrapped ErrorManager.
func (m *ThreadSafeDelegatingErrorManager) SetConformanceViolationHandler(h ConformanceViolationHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.conformanceViolationHandler = h
}

// Java: ThreadSafeDelegatingErrorManager.shouldReportConformanceViolation(Requirement, Optional<RequirementScopeEntry>, JSError, LibraryLevelNonAllowlistedConformanceViolationsBehavior, boolean)
func (m *ThreadSafeDelegatingErrorManager) ShouldReportConformanceViolation(
	requirement *Requirement,
	allowlistEntry *RequirementScopeEntry,
	diagnostic *JSError,
	behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior,
	isAllowlisted bool) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.conformanceViolationHandler != nil {
		return m.conformanceViolationHandler(
			requirement, allowlistEntry, diagnostic, behavior, isAllowlisted)
	}
	return m.delegated.ShouldReportConformanceViolation(
		requirement, allowlistEntry, diagnostic, behavior, isAllowlisted)
}
