package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: fast mode (--fast_mode) runs TypeCheck's traversal one script per
// worker goroutine.
//
// Type inference has finished and the TypedScopeCreator is frozen, so the
// scopes, the type registry and the JSTypes are only read (the registry is
// put in its parallel read-only mode, which bypasses its caches). Each
// worker has its own TypeCheck and TypeValidator clone whose reports and
// type mismatches are recorded, then replayed in script order: the
// diagnostics, the mismatch list and the typed-percent counts are the same
// as the sequential traversal's. The per-node writes TypeCheck makes
// (ensureTyped) stay inside the worker's script. InferJSDocInfo, which
// mutates types, runs afterwards on the main goroutine as before.

func (tc *TypeCheck) canCheckInParallel(node rhino.Node) ([]rhino.Node, int, bool) {
	opts := tc.compiler.GetOptions()
	if !opts.isFastMode() || tc.compiler.IsDebugLoggingEnabled() {
		return nil, 0, false
	}
	if tc.scopeCreator == nil || tc.topScope == nil {
		return nil, 0, false
	}
	if tc.scopeCreator.stage != TypedScopeCreator_Stage_FROZEN {
		return nil, 0, false
	}
	scripts, ok := collectRootScripts(node)
	if !ok {
		return nil, 0, false
	}
	n := parallelThreadCount(tc.compiler, len(scripts))
	if n <= 1 {
		return nil, 0, false
	}
	return scripts, n, true
}

// cloneForWorker returns a TypeCheck reporting to rc, sharing everything
// read-only with tc.
func (tc *TypeCheck) cloneForWorker(rc *recordingCompiler) *TypeCheck {
	v := *tc.validator
	v.compiler = rc
	v.mismatches = TypeMismatch_Accumulator{}
	c := &TypeCheck{
		compiler:                rc,
		validator:               &v,
		reverseInterpreter:      tc.reverseInterpreter,
		typeRegistry:            tc.typeRegistry,
		topScope:                tc.topScope,
		scopeCreator:            tc.scopeCreator,
		reportUnknownTypes:      tc.reportUnknownTypes,
		subtypingMode:           tc.subtypingMode,
		reportMissingProperties: tc.reportMissingProperties,
		inExterns:               tc.inExterns,
	}
	return c
}

// checkParallel is the traversal half of Check(node, externs) in fast mode;
// it returns false when the sequential traversal must run instead.
func (tc *TypeCheck) checkParallel(node rhino.Node) bool {
	scripts, numWorkers, ok := tc.canCheckInParallel(node)
	if !ok {
		return false
	}
	prewarmCompilerForParallelReads(tc.compiler)
	restoreReg := enableParallelTypeRegistryReads(tc.compiler)
	defer restoreReg()

	b := NodeTraversal_Builder{compiler: tc.compiler, callback: tc, scopeCreator: tc.scopeCreator}
	t := acquireNodeTraversal(&b)
	defer releaseNodeTraversal(t)
	t.InitTraversal(node)
	t.currentNode = node
	t.PushScopeOfScope(tc.topScope)
	if !tc.ShouldTraverse(t, node, rhino.NilNode) {
		t.PopScope()
		return true
	}

	type itemResult struct {
		mismatches             []TypeMismatch
		typed, null, unknown   int
		subtypingMode          rhino.JSType_SubtypingMode
		validatorSubtypingMode rhino.JSType_SubtypingMode
	}
	results := make([]itemResult, len(scripts))
	runParallelWork(tc.compiler, numWorkers, len(scripts), func(int) (func(int, *recordingCompiler), func()) {
		runItem := func(idx int, rc *recordingCompiler) {
			c := tc.cloneForWorker(rc)
			wb := NodeTraversal_Builder{compiler: rc, callback: c, scopeCreator: tc.scopeCreator}
			wt := acquireNodeTraversal(&wb)
			defer releaseNodeTraversal(wt)
			func() {
				defer func() {
					if r := recover(); r != nil {
						wt.ThrowUnexpectedException(jutil.AsThrowable(r))
					}
				}()
				wt.InitTraversal(node)
				wt.currentNode = node
				wt.PushScopeOfScope(tc.topScope)
				wt.TraverseBranch(scripts[idx], node)
				wt.PopScope()
			}()
			results[idx] = itemResult{
				mismatches:             c.validator.mismatches.mismatches,
				typed:                  c.typedCount,
				null:                   c.nullCount,
				unknown:                c.unknownCount,
				subtypingMode:          c.subtypingMode,
				validatorSubtypingMode: c.validator.subtypingMode,
			}
		}
		return runItem, nil
	}, false)

	// Replay what the sequential traversal would have left behind, in order.
	for i := range results {
		r := &results[i]
		tc.validator.mismatches.mismatches = append(tc.validator.mismatches.mismatches, r.mismatches...)
		tc.typedCount += r.typed
		tc.nullCount += r.null
		tc.unknownCount += r.unknown
		if i == len(results)-1 {
			tc.subtypingMode = r.subtypingMode
			tc.validator.SetSubtypingMode(r.validatorSubtypingMode)
		}
	}
	t.currentNode = node
	tc.Visit(t, node, rhino.NilNode)
	t.PopScope()
	return true
}
