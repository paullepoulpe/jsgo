package jscomp

// newTestCompiler is `new Compiler()` followed by `initOptions(options)`, the
// set-up that many ported tests and oracles share.
//
// Go-only: test-helper
func newTestCompiler(options *CompilerOptions) *Compiler {
	c := NewCompiler()
	c.InitOptions(options)
	return c
}
