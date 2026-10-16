// Port of com/google/javascript/jscomp/ijs/IjsErrors.java @ closure-compiler v20260415

package jscomp

// IjsErrors holds IJS-related errors occurring outside the
// ConvertToTypedInterface pass.
type IjsErrors struct{}

func newIjsErrors() *IjsErrors { return &IjsErrors{} }

var (
	IjsErrors_BAD_IJS_FILE_NAME = DiagnosticType_Error("JSC_BAD_IJS_FILE_NAME", "--ijs file name must end in .i.js: {0}")

	IjsErrors_CONFLICTING_IJS_FILE = DiagnosticType_Warning(
		"JSC_CONFLICTING_IJS_FILE",
		"Found --ijs file in a conflict with --js or --weakdep file. This usually means the same"+
			" file appears both in the srcs and the transitive deps of a build target.\n"+
			"File 1: {0}\n"+
			"File 2: {1}")
)
