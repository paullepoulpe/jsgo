package jscomp

import (
	"io"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func (m *customPassMap) add(time CustomPassExecutionTime, pass CompilerPass) {
	for i := range m.groups {
		if m.groups[i].time == time {
			m.groups[i].passes = append(m.groups[i].passes, pass)
			return
		}
	}
	m.groups = append(m.groups, customPassGroup{
		time:   time,
		passes: []CompilerPass{pass},
	})
}

// AddCustomPass registers a custom CompilerPass to run at executionTime.
//
// Java: CompilerOptions.addCustomPass(CustomPassExecutionTime, CompilerPass)
func (o *CompilerOptions) AddCustomPass(executionTime CustomPassExecutionTime, customPass CompilerPass) {
	jutil.CheckNotNull(customPass)
	if o.customPasses == nil {
		o.customPasses = &customPassMap{}
	}
	o.customPasses.add(executionTime, customPass)
}

// SetNameGenerator sets the shared NameGenerator used for variable and property renaming.
//
// Java: CompilerOptions.setNameGenerator(NameGenerator)
func (o *CompilerOptions) SetNameGenerator(nameGenerator NameGenerator) {
	o.nameGenerator = jutil.CheckNotNull(nameGenerator)
}

// SetCssRenamingMap sets the map used in renaming CSS class names.
//
// Java: CompilerOptions.setCssRenamingMap(CssRenamingMap)
func (o *CompilerOptions) SetCssRenamingMap(cssRenamingMap CssRenamingMap) {
	o.cssRenamingMap = cssRenamingMap
}

// SetXidHashFunction sets the hash function used for xid generation.
//
// Java: CompilerOptions.setXidHashFunction(Xid.HashFunction)
func (o *CompilerOptions) SetXidHashFunction(xidHashFunction Xid_HashFunction) {
	o.xidHashFunction = xidHashFunction
}

// SetChunkIdHashFunction sets the hash function used for generating chunk IDs.
//
// Java: CompilerOptions.setChunkIdHashFunction(Xid.HashFunction)
func (o *CompilerOptions) SetChunkIdHashFunction(chunkIdHashFunction Xid_HashFunction) {
	o.chunkIdHashFunction = chunkIdHashFunction
}

// SetStateCompressionWrapper sets an optional stream wrapper applied when
// writing saved compilation state.
//
// Java: CompilerOptions.setStateCompressionWrapper(Function<OutputStream, OutputStream>)
func (o *CompilerOptions) SetStateCompressionWrapper(wrapper func(io.Writer) io.Writer) {
	o.stateCompressionWrapper = wrapper
}

// SetStateDecompressionWrapper sets an optional stream wrapper applied when
// reading saved compilation state.
//
// Java: CompilerOptions.setStateDecompressionWrapper(Function<InputStream, InputStream>)
func (o *CompilerOptions) SetStateDecompressionWrapper(wrapper func(io.Reader) io.Reader) {
	o.stateDecompressionWrapper = wrapper
}

// SetTypedAstOutputStream configures an in-memory output stream to receive
// serialized TypedAST output during checks (wired to
// SerializeTypedAstPass_CreateFromOutputStream).
func (o *CompilerOptions) SetTypedAstOutputStream(out io.Writer) {
	o.typedAstOutputStream = out
}

// GetTypedAstOutputStream returns the configured in-memory TypedAST output stream, or nil.
func (o *CompilerOptions) GetTypedAstOutputStream() io.Writer {
	return o.typedAstOutputStream
}

// HasTypedAstOutput reports whether either a file path or an output stream is
// configured for TypedAST serialization.
func (o *CompilerOptions) HasTypedAstOutput() bool {
	return o.typedAstOutputFile != nil || o.typedAstOutputStream != nil
}

// SetLicenseFilter installs an optional predicate invoked on each file's
// @license/@preserve comment text before emission; returning false suppresses it.
func (o *CompilerOptions) SetLicenseFilter(filter func(string) bool) {
	o.licenseFilter = filter
}

// GetLicenseFilter returns the configured license filter predicate, or nil.
func (o *CompilerOptions) GetLicenseFilter() func(string) bool {
	return o.licenseFilter
}

// SetReconcileOptionsHook installs an optional callback invoked at the end of
// Compiler.ReconcileOptionsWithGuards.
func (o *CompilerOptions) SetReconcileOptionsHook(hook func(*CompilerOptions)) {
	o.reconcileOptionsHook = hook
}

// GetReconcileOptionsHook returns the configured reconcile-options callback, or nil.
func (o *CompilerOptions) GetReconcileOptionsHook() func(*CompilerOptions) {
	return o.reconcileOptionsHook
}
