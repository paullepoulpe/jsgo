// Port of com/google/javascript/jscomp/ThreadSafeDelegatingErrorManager.java @ closure-compiler v20260415

package jscomp

import "sync"

// ThreadSafeDelegatingErrorManager is a simple delegating ErrorManager that
// provides a thread-safe wrapper for the one being delegated.
type ThreadSafeDelegatingErrorManager struct {
	mu        sync.Mutex
	delegated ErrorManager
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
	m.mu.Lock()
	defer m.mu.Unlock()
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

// Java: ThreadSafeDelegatingErrorManager.shouldReportConformanceViolation(Requirement, Optional<RequirementScopeEntry>, JSError, LibraryLevelNonAllowlistedConformanceViolationsBehavior, boolean)
func (m *ThreadSafeDelegatingErrorManager) ShouldReportConformanceViolation(
	requirement *Requirement,
	allowlistEntry *RequirementScopeEntry,
	diagnostic *JSError,
	behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior,
	isAllowlisted bool) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.delegated.ShouldReportConformanceViolation(
		requirement, allowlistEntry, diagnostic, behavior, isAllowlisted)
}
