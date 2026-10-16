// Port of com/google/javascript/jscomp/ErrorHandler.java @ closure-compiler v20260415

package jscomp

// ErrorHandler is any generic sink for warnings and errors, after they've
// passed through any filtering WarningsGuards.
type ErrorHandler interface {
	// Report reports error at the reporting level. level may be nil (Java
	// callers pass null).
	//
	// Java: ErrorHandler.report(CheckLevel, JSError)
	Report(level *CheckLevel, error *JSError)
}
