// Port of com/google/javascript/rhino/StaticSourceFile.java
// @ closure-compiler v20260415

package rhino

// SourceKind mirrors StaticSourceFile.SourceKind.
type SourceKind int

const (
	// SourceKindStrong is a file whose contents are necessary both for type
	// checking and emitting code.
	SourceKindStrong SourceKind = iota
	// SourceKindWeak is a file whose contents are necessary for type checking only.
	SourceKindWeak
	// SourceKindExtern is a file whose contents are extern declarations.
	SourceKindExtern
	// SourceKindNonCode is a file that is not part of the compilation unit and is
	// not strong, weak, or an extern.
	SourceKindNonCode
)

// StaticSourceFile contains information about a compiler input.
// Mirrors com.google.javascript.rhino.StaticSourceFile.
type StaticSourceFile interface {
	// GetName returns the name of the file. Must be unique across all files in
	// the compilation.
	GetName() string

	// GetKind returns the source kind.
	GetKind() SourceKind

	// IsClosureUnawareCode reports whether this source file contains code that
	// should be optimized using a minimal set of known-safe optimizations.
	IsClosureUnawareCode() bool

	// GetLineOffset returns the offset of the given line number relative to the
	// file start. Line number should be 1-based. Returns math.MinInt32 when the
	// source file has no line information.
	GetLineOffset(lineNumber int) int

	// GetLineOfOffset returns the 1-based line number of the given source offset.
	GetLineOfOffset(offset int) int

	// GetColumnOfOffset returns the 0-based column number of the given source
	// offset.
	GetColumnOfOffset(offset int) int
}

// IsStrong reports whether sf has kind SourceKindStrong.
func IsStrong(sf StaticSourceFile) bool { return sf.GetKind() == SourceKindStrong }

// IsWeak reports whether sf has kind SourceKindWeak.
func IsWeak(sf StaticSourceFile) bool { return sf.GetKind() == SourceKindWeak }

// IsExtern reports whether sf has kind SourceKindExtern.
func IsExtern(sf StaticSourceFile) bool { return sf.GetKind() == SourceKindExtern }

// IsNonCode reports whether sf has kind SourceKindNonCode.
func IsNonCode(sf StaticSourceFile) bool { return sf.GetKind() == SourceKindNonCode }

// IsTypeScriptSource reports whether sf is a TypeScript-derived source file,
// matching StaticSourceFile.isTypeScriptSource().
func IsTypeScriptSource(sf StaticSourceFile) bool {
	name := sf.GetName()
	return len(name) > len(".closure.js") && name[len(name)-len(".closure.js"):] == ".closure.js" ||
		len(name) > len(".tsx.cl.js") && name[len(name)-len(".tsx.cl.js"):] == ".tsx.cl.js"
}
