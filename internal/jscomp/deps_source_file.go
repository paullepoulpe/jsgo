// Port of com/google/javascript/jscomp/deps/SourceFile.java @ closure-compiler v20260415

package jscomp

// Deps_SourceFile is an interface for representing source files for the
// Compilation Service.
//
// Go-only: language-adapter — com.google.javascript.jscomp.deps.SourceFile
// and com.google.javascript.jscomp.SourceFile are both in the Go package
// jscomp, so this one carries its subpackage name.
type Deps_SourceFile interface {
	GetName() string
	GetContent() string
	WasModified() bool
}
