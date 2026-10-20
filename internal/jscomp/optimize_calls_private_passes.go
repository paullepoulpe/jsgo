package jscomp

import "github.com/paullepoulpe/jsgo/internal/rhino"

// Only DefaultPassConfig constructs these adapters. Their fresh pass objects
// cannot be observed through the public AsCallGraphCompilerPass API, so their
// reference-map scope may remain private to one OptimizeCalls invocation.
type optimizeCallsPrivateReturns struct{ pass *OptimizeReturns }
type optimizeCallsPrivateParameters struct{ pass *OptimizeParameters }
type optimizeCallsPrivateConstructors struct{ pass *OptimizeConstructors }

func newOptimizeCallsPrivateReturns(compiler AbstractCompiler) OptimizeCalls_CallGraphCompilerPass {
	return optimizeCallsPrivateReturns{NewOptimizeReturns(compiler)}
}

func newOptimizeCallsPrivateParameters(compiler AbstractCompiler) OptimizeCalls_CallGraphCompilerPass {
	return optimizeCallsPrivateParameters{NewOptimizeParameters(compiler)}
}

func newOptimizeCallsPrivateConstructors(compiler AbstractCompiler) OptimizeCalls_CallGraphCompilerPass {
	return optimizeCallsPrivateConstructors{NewOptimizeConstructors(compiler)}
}

func (p optimizeCallsPrivateReturns) Process(externs, root *rhino.Node, references *OptimizeCalls_ReferenceMap) {
	p.pass.ProcessWithReferenceMap(externs, root, references)
}

func (p optimizeCallsPrivateParameters) Process(externs, root *rhino.Node, references *OptimizeCalls_ReferenceMap) {
	previousScope := p.pass.globalScope
	defer func() { p.pass.globalScope = previousScope }()
	p.pass.ProcessWithReferenceMap(externs, root, references)
}

func (p optimizeCallsPrivateConstructors) Process(externs, root *rhino.Node, references *OptimizeCalls_ReferenceMap) {
	p.pass.ProcessWithReferenceMap(externs, root, references)
}
