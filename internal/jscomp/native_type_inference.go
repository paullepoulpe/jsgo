package jscomp

import "github.com/paullepoulpe/jsgo/internal/rhino"

// runNativeTypeInference uses no generic CFG or graph annotation adapters.
// All original type-inference initialization and transfer methods are shared.
// The result retains complete node/branch states and implicit-return input for
// differential fixtures; the ordinary compiler pass only needs AST/type effects.
func runNativeTypeInference(
	compiler AbstractCompiler, topology *indexedFlowTopology[*rhino.Node],
	reverseInterpreter ReverseAbstractInterpreter, syntacticScope *TypedScope,
	scopeCreator *TypedScopeCreator, assertionFunctionLookup *CodingConvention_AssertionFunctionLookup) *indexedFlowResult[*rhino.Node, FlowScope] {
	ti := &TypeInference{}
	// The pure indexed engine needs virtual dispatch, but never GetCfg or the
	// generic work queue. No placeholder graph is constructed here.
	ti.self = ti
	entry := topology.values[topology.entry]
	ti.initTypeInference(compiler, entry, entry, reverseInterpreter,
		syntacticScope, scopeCreator, assertionFunctionLookup)
	result := newIndexedFlowResult(topology, ti)
	result.compute()
	return result
}
